package assistants

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Schema struct {
	Openapi string `json:"openapi"`
	Info    struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Version     string `json:"version"`
	} `json:"info"`
	Servers []struct {
		Url string `json:"url"`
	} `json:"servers"`
	Paths struct {
		Location struct {
			Get struct {
				Description string `json:"description"`
				OperationId string `json:"operationId"`
				Parameters  []struct {
					Name        string `json:"name"`
					In          string `json:"in"`
					Description string `json:"description"`
					Required    bool   `json:"required"`
					Schema      struct {
						Type string `json:"type"`
					} `json:"schema"`
				} `json:"parameters"`
				Deprecated bool `json:"deprecated"`
			} `json:"get"`
		} `json:"/location"`
	} `json:"paths"`
	Components struct {
		Schemas struct {
		} `json:"schemas"`
	} `json:"components"`
}

// MetadataHttp 动态工具元数据
type MetadataHttp struct {
	// Url 请求地址
	Url string `json:"url"`
	// Method 请求方法
	Method string `json:"method"`
	// Headers 请求头
	Headers map[string]string `json:"header"`
	// Body 请求体
	Body interface{} `json:"body"`
	// UserAgent 请求头
	UserAgent string `json:"userAgent"`
}

// ToolOptions 动态工具选项
type ToolOptions struct {
	// Name 工具名称
	Name string
	// ActionID 动作ID
	ActionID string
	// ToolType 工具类型
	ToolType types.ToolType
	// Description 工具描述
	Description string
	// Metadata 工具元数据
	Metadata string
	logger   log.Logger
	TraceId  string
	Opts     []kithttp.ClientOption
}

type dynamicTool struct {
	CallbacksHandler callbacks.Handler
	name             string
	description      string
	actionID         string
	params           map[string]string
	toolType         types.ToolType
	logger           log.Logger
	traceId          string
	metadata         string
	opts             []kithttp.ClientOption
}

func (s dynamicTool) Name() string {
	return s.name
}

func (s dynamicTool) Description() string {
	return s.description
}

func (s dynamicTool) Call(ctx context.Context, input string) (result string, err error) {
	logger := log.With(s.logger, "actionId", s.actionID, "tool", s.Name(), "toolType", s.toolType)
	_ = level.Info(logger).Log("msg", fmt.Sprintf("调用工具: %s, 输入: %s", s.Name(), input))
	if s.CallbacksHandler != nil {
		s.CallbacksHandler.HandleToolStart(ctx, input)
	}
	switch s.toolType {
	case types.ToolTypeTypeFunction:
		result, err = s.httpSend(ctx, input)
	case types.ToolTypeCodeInterpreter:
		// TODO: 代码解释器，调用paas创建job执行代码
	case types.ToolTypeRetrieval:
		// TODO: 知识库检索
	}

	if err != nil {
		if s.CallbacksHandler != nil {
			s.CallbacksHandler.HandleToolError(ctx, err)
		}
		_ = level.Error(logger).Log("msg", "调用工具失败", "err", err)
		return "", err
	}

	if s.CallbacksHandler != nil {
		s.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}

func (s dynamicTool) getHttpClient(ctx context.Context, host string) *http.Client {
	// 如果是com的域名，使用代理
	if strings.HasSuffix(host, ".com") {
		return &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	return http.DefaultClient
}

func (s dynamicTool) httpSend(ctx context.Context, input string) (res string, err error) {
	var metaData MetadataHttp
	if err = json.Unmarshal([]byte(s.metadata), &metaData); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("解析元数据失败: %s", s.metadata))
		return "", err
	}
	tgt, err := url.Parse(metaData.Url)
	if err != nil {
		err = errors.Wrap(err, "解析URL失败")
		return "", err
	}
	opts := s.opts
	opts = append(opts, kithttp.SetClient(s.getHttpClient(ctx, metaData.Url)))
	ep := kithttp.NewClient(strings.ToUpper(metaData.Method), tgt, func(ctx context.Context, r *http.Request, request interface{}) error {
		r.Header.Set("User-Agent", metaData.UserAgent)
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		r.Header.Set("Action-Id", s.actionID)
		r.Header.Set("Action-Name", s.Name())
		if metaData.Headers != nil {
			for k, v := range metaData.Headers {
				r.Header.Set(k, v)
			}
		}
		var b bytes.Buffer
		r.Body = io.NopCloser(&b)
		return json.NewEncoder(&b).Encode(request)
	}, func(ctx context.Context, response2 *http.Response) (response interface{}, err error) {
		if response2.StatusCode != http.StatusOK {
			if response2.Body != nil {
				return io.ReadAll(response2.Body)
			}
			err = errors.New(fmt.Sprintf("调用HTTP失败: %s", response2.Status))
			return nil, err
		}
		return io.ReadAll(response2.Body)
	}, opts...).Endpoint()

	var data interface{}

	if strings.TrimSpace(input) != "" && strings.TrimSpace(input) != "None" && strings.TrimSpace(input) != "{}" {
		data = map[string]interface{}{}
		_ = level.Info(s.logger).Log("msg", fmt.Sprintf("解析输入: %s", input))
		// 解码JSON到map
		err = json.Unmarshal([]byte(input), &data)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("解析输入失败: %s", input))
			//return "", err
		}
	}

	result, err := ep(ctx, data)
	if err != nil {
		err = errors.Wrap(err, "调用HTTP失败")
		return "", err
	}
	resByte, _ := result.([]byte)
	return string(resByte), err
}

func NewDynamicTool(opts ToolOptions) tools.Tool {
	logger := log.With(opts.logger, "actionId", opts.ActionID, "tool", opts.Name)
	t := &dynamicTool{
		name:        opts.Name,
		actionID:    opts.ActionID,
		description: opts.Description,
		toolType:    opts.ToolType,
		logger:      logger,
		traceId:     opts.TraceId,
		metadata:    opts.Metadata,
		opts:        opts.Opts,
	}
	return t
}
