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
		}
		job.TrainLog = jobLog
		if err = store.FineTuning().UpdateFineTuningJob(ctx, &job); err != nil {
			_ = level.Warn(logger).Log("msg", "update job log failed", "err", err)
		}
	}
	return
}
