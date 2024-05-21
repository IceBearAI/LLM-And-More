package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/pkg/files"
	"github.com/IceBearAI/aigc/src/pkg/finetuning"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/util"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
	"time"
)

var (
	jobFineTuningCmd = &cobra.Command{
		Use:               "finetuning command <args> [flags]",
		Short:             "微调任务命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
可用的配置类型：
[run-waiting-train, running-log]

aigc-server job -h
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err = prepare(cmd.Context()); err != nil {
				return errors.Wrap(err, "prepare")
			}
			fileSvc = files.NewService(logger, traceId, store, apiSvc, []files.CreationOption{
				files.WithLocalDataPath(serverStoragePath),
				files.WithServerUrl(fmt.Sprintf("%s/storage", serverDomain)),
				files.WithStorageType("local"),
			}...)
			fineTuningSvc = finetuning.New(traceId, logger, store, fileSvc, apiSvc,
				finetuning.WithVolumeName(runtimeK8sVolumeName),
				finetuning.WithGpuTolerationValue(datasetsGpuToleration),
				finetuning.WithCallbackHost(serverDomain),
			)
			return nil
		},
	}

	jobFineTuningJobRunWaitingTrainCmd = &cobra.Command{
		Use:               `run-waiting-train [flags]`,
		Short:             "微调任务等待训练",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
aigc-server job finetuning run-waiting-train
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fineTuningRunWaitingTrain(cmd.Context())
		},
	}

	jobFineTuningJobRunningJobLogCmd = &cobra.Command{
		Use:               `running-log [flags]`,
		Short:             "同步正在训练脚本日志",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
aigc-server job finetuning running-log
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			runningJobs, err := store.FineTuning().FindFineTuningJobRunning(cmd.Context())
			if err != nil {
				_ = level.Error(logger).Log("msg", "find running job failed", "err", err.Error())
				return err
			}
			return fineTuningRunningJobLog(cmd.Context(), runningJobs)
		},
	}
)

func fineTuningRunWaitingTrain(ctx context.Context) (err error) {
	return fineTuningSvc.RunWaitingTrain(ctx)
}

type logEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	Loss         float64   `json:"loss"`
	LearningRate float64   `json:"learning_rate"`
	Epoch        float64   `json:"epoch"`
	GradNorm     float64   `json:"grad_norm"`
	EvalLoss     float64   `json:"eval_loss"`
	Step         int       `json:"step"`
	Rank         int       `json:"rank"`
}

func fineTuningRunningJobLog(ctx context.Context, jobs []types.FineTuningTrainJob) (err error) {
	for _, job := range jobs {
		var jobLog string
		jobLog, err = apiSvc.Runtime().GetJobLogs(ctx, job.PaasJobName)
		if err != nil {
			_ = level.Warn(logger).Log("msg", "get job pods log failed", "err", err.Error())
			continue
		}
		if len(strings.TrimSpace(jobLog)) == 0 {
			continue
		}

		lineArr := strings.Split(jobLog, "\n")
		re := regexp.MustCompile(`\{[^}]*\}`)

		var logEntryList []logEntry

		for _, log := range lineArr {
			log = strings.TrimSpace(log)
			matches := re.FindAllString(log, -1)
			for _, match := range matches {
				if len(match) > 0 {
					// 将单引号替换为双引号以符合JSON格式
					jsonStr := strings.Replace(match, "'", "\"", -1)         // 将单引号替换为双引号
					jsonStr = strings.Replace(jsonStr, "False", "false", -1) // 将 False 替换为 false
					jsonStr = strings.Replace(jsonStr, "True", "true", -1)   // 将 True 替换为 true

					var entry logEntry
					err := json.Unmarshal([]byte(jsonStr), &entry)
					if err != nil {
						level.Info(logger).Log("json.Unmarshal", "unmarshal json failed", "err", err.Error())
						continue
					}
					if entry.Epoch == 0 {
						continue
					}
					logEntryList = append(logEntryList, entry)
				}
			}
		}
		if len(logEntryList) > 1 {
			lastLine := logEntryList[len(logEntryList)-1]
			job.ProgressLoss = lastLine.Loss
			job.ProgressLearningRate = lastLine.LearningRate
			job.ProgressEpochs = lastLine.Epoch
			job.Progress = util.RoundToFourDecimalPlaces(lastLine.Epoch / float64(job.TrainEpoch))
			if lastLine.EvalLoss > 0 {
				b, _ := json.Marshal(fineTuningDiagnosisLog(logEntryList))
				job.Diagnosis = string(b)
			}
		}
		job.TrainLog = jobLog
		if err = store.FineTuning().UpdateFineTuningJob(ctx, &job); err != nil {
			_ = level.Warn(logger).Log("msg", "update job log failed", "err", err)
		}
	}
	return
}

type diagnosis struct {
	// 过拟合
	Overfitting string `json:"overfitting,omitempty"`
	// 欠拟合
	Underfitting string `json:"underfitting,omitempty"`
	// 灾难性遗忘
	CatastrophicForgetting string `json:"catastrophicForgetting,omitempty"`
}

func fineTuningDiagnosisLog(lines []logEntry) (res diagnosis) {
	var (
		// 设置默认警报阈值
		//overfitting            = 0.02 // 过拟合阈值：验证损失的增加率
		underfitting           = 0.5 // 欠拟合阈值：训练损失高于此值
		catastrophicForgetting = 1.5 // 灾难性遗忘阈值：训练损失的突然增加率

		// 计算平均变化率
		totalTrainLossChange float64 = 0 // 训练损失变化总和
		totalEvalLossChange  float64 = 0 // 验证损失变化总和
	)
	if len(lines) == 1 {
		if lines[0].Loss < 0.009 {
			res.Overfitting = "High"
		} else if lines[0].Loss > underfitting {
			res.Underfitting = "High"
		}
	} else {
		for i := 0; i < len(lines)-1; i++ {
			totalTrainLossChange += lines[i+1].Loss - lines[i].Loss
			totalEvalLossChange += lines[i+1].EvalLoss - lines[i].EvalLoss
		}
		avgTrainLossChange := totalTrainLossChange / float64(len(lines)-1)
		//avgEvalLossChange := totalEvalLossChange / float64(len(lines)-1)
		// 检查过拟合
		if lines[len(lines)-1].EvalLoss > lines[len(lines)-2].EvalLoss {
			// 当验证损失增加时，表示过拟合
			res.Overfitting = "High"
		}
		// 检查欠拟合
		if lines[len(lines)-1].Loss > underfitting {
			// 当训练损失高于阈值时，表示欠拟合
			res.Underfitting = "High"
		}
		// 检查灾难性遗忘
		if avgTrainLossChange > catastrophicForgetting {
			// 当训练损失的平均变化率高于阈值时，表示灾难性遗忘
			res.CatastrophicForgetting = "High"
		}
	}
	return
}
