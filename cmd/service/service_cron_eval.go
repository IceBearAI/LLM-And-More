package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/IceBearAI/aigc/src/pkg/files"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	openai2 "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strings"
)

type llmWaitingEvalCronJob struct {
	logger   log.Logger
	Name     string
	ctx      context.Context
	store    repository.Repository
	apiSvc   services.Service
	running  bool
	filesSvc files.Service
}

type messageLine struct {
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type alpacaData struct {
	ID            string                `json:"id"`
	Conversations []alpacaConversations `json:"conversations"`
}

type alpacaConversations struct {
	From  string `json:"from"`
	Value string `json:"value"`
}

func (s *llmWaitingEvalCronJob) Run() {
	var runErr error
	if s.running {
		return
	}
	s.running = true
	defer func() {
		s.running = false
	}()
	// 如果有正在运行的评估任务，则不执行
	_, runErr = s.store.LLMEval().GetRunningEvalTask(s.ctx)
	if !errors.Is(runErr, gorm.ErrRecordNotFound) {
		_ = level.Info(s.logger).Log("msg", "get running eval task failed", "err", runErr)
		return
	}
	task, runErr := s.store.LLMEval().GetEarliestWaitingEvalTask(s.ctx, "DatasetFile")
	if runErr != nil {
		_ = level.Info(s.logger).Log("msg", "get earliest waiting eval task failed", "err", runErr)
		return
	}

	defer func() {
		if runErr != nil {
			if updateErr := s.store.LLMEval().UpdateEvalStatus(s.ctx, task.ID, types.EvalStatusFailed, runErr.Error()); updateErr != nil {
				_ = level.Error(s.logger).Log("msg", "update eval status failed", "err", updateErr)
			}
		}
	}()

	// 判断模型是否可用
	modelInfo, runErr := s.store.Model().FindByModelId(s.ctx, task.ModelName)
	if runErr != nil {
		_ = level.Warn(s.logger).Log("msg", "get model info failed", "err", runErr)
		return
	}
	if !modelInfo.Enabled {
		_ = level.Warn(s.logger).Log("msg", "model not enabled")
		return
	}

	datasetFileInfo := task.DatasetFile
	//if task.DatasetFile.ID == 0 {
	//	datasetFileInfo = task.DatasetSampleFile
	//}

	// 更新评估状态
	if runErr = s.store.LLMEval().UpdateEvalStatus(s.ctx, task.ID, types.EvalStatusRunning, ""); err != nil {
		_ = level.Error(s.logger).Log("msg", "update eval status failed", "err", runErr)
		return
	}
	// 更新开始时间
	if runErr = s.store.LLMEval().UpdateEvalStartTime(s.ctx, task.ID); runErr != nil {
		_ = level.Warn(logger).Log("msg", "update eval start time failed", "err", runErr)
	}

	if task.MetricName == "equal" {
		if runErr = s.completionMetricEqual(s.ctx, datasetFileInfo.S3Url, task); runErr != nil {
			_ = level.Error(s.logger).Log("msg", "completion metric equal failed", "err", runErr)
			runErr = errors.Wrap(err, "completionMetricEqual")
		}
		return
	}
}

func (s *llmWaitingEvalCronJob) completionMetricEqual(ctx context.Context, httpUrl string, evalTask types.LLMEvalResults) (err error) {
	body, err := getHttpFileBody(httpUrl)
	if err != nil {
		err = errors.Wrap(err, "getHttpFileBody")
		return
	}
	_ = level.Info(s.logger).Log("msg", "开始评估", "httpUrl", httpUrl, "evalTaskModelName¬", evalTask.ModelName)
	var progressedNum int
	var successNum int
	var rate float64
	var notEqualLines []byte
	dataList := bytes.Split(body, []byte("\n"))
	totalLine := len(dataList)
	for k, line := range dataList {
		var inputMsg messageLine
		if err = json.Unmarshal(line, &inputMsg); err != nil {
			continue
		}
		var chatMessages []openai2.ChatCompletionMessage
		chatMessages = append(chatMessages, openai2.ChatCompletionMessage{
			Role:    "system",
			Content: "",
		})
		isEqual := true
		for _, msg := range inputMsg.Messages {
			if !util.StringInArray([]string{"user", "assistant"}, msg.Role) {
				isEqual = false
				continue
			}
			if strings.EqualFold(msg.Role, "user") {
				chatMessages = append(chatMessages, openai2.ChatCompletionMessage{
					Role:    "user",
					Content: msg.Content,
				})
			}
			if strings.EqualFold(msg.Role, "assistant") {
				chatMessages = append(chatMessages, openai2.ChatCompletionMessage{
					Role:    "assistant",
					Content: msg.Content,
				})
			}

			completion, _, err := s.apiSvc.FastChat().ChatCompletion(ctx, evalTask.ModelName, chatMessages, 0, 0, 0.5, 0.5, 2048, 1, nil, "", nil, nil)
			if err != nil {
				_ = level.Warn(s.logger).Log("msg", "create completion failed", "err", err)
				isEqual = false
				continue
			}
			if !strings.EqualFold(strings.TrimSpace(completion.Choices[0].Message.Content), strings.TrimSpace(msg.Content)) {
				_ = level.Info(s.logger).Log("msg", "completion not equal", "err", err)
				isEqual = false
				continue
			}
		}
		chatMessages = nil
		isEqual = true
		if !isEqual {
			line = append(line, []byte("\n")...)
			notEqualLines = append(notEqualLines, line...)
			_ = level.Info(s.logger).Log("msg", "completion not equal", "line", k, "err", err)
			continue
		}
		successNum++
		progressedNum++
		if progressedNum == 100 {
			rate = float64(successNum) / float64(k)
			progressed := float64(k) / float64(evalTask.EvalTotal)
			// 更新数据库进度
			if err = s.store.LLMEval().UpdateEvalProgress(s.ctx, evalTask.ID, rate*10, progressed, types.EvalStatusRunning); err != nil {
				_ = level.Error(s.logger).Log("msg", "update eval progress failed", "err", err)
			}
			progressedNum = 0
		}
		if k >= (evalTask.EvalTotal - 1) {
			_ = level.Info(s.logger).Log("msg", "completion rate", "score", rate, "successNum", successNum, "totalLine", totalLine)
			break
		}
	}
	// 计算百分比
	rate = float64(successNum) / float64(evalTask.EvalTotal)
	// 更新评估状态
	if err = s.store.LLMEval().UpdateEvalProgress(s.ctx, evalTask.ID, rate*10, 1, types.EvalStatusSuccess); err != nil {
		_ = level.Error(s.logger).Log("msg", "update eval status failed", "err", err)
		return
	}
	// 将失败的传到s3或者其他地方
	if len(notEqualLines) > 0 {
		if fileUrl, uploadErr := s.filesSvc.UploadLocal(ctx, util.NewFile(body), "json"); uploadErr == nil {
			_ = level.Info(logger).Log("msg", "upload to s3 success", "panUrl", fileUrl)
			detail := map[string]interface{}{
				"evalResultNotEqualUrl": fileUrl,
			}
			b, _ := json.Marshal(detail)
			if evalErr := s.store.LLMEval().UpdateEvalDetail(ctx, evalTask.ID, string(b)); evalErr != nil {
				_ = level.Warn(logger).Log("msg", "update eval detail failed", "err", evalErr)
			}
		}
	}
	_ = level.Info(s.logger).Log("msg", "completion rate", "score", rate, "successNum", successNum, "evalTask.EvalTotal", evalTask.EvalTotal, "totalLine", totalLine)
	return nil
}

func getHttpFileBody(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	//resp, err := http.Get(url)
	if err != nil {
		err = errors.Wrap(err, "http.Get")
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "io.ReadAll")
		return
	}
	return
}
