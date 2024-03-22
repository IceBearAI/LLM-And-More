package tests

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
	"net/http"
)

type ToolOptions struct {
	Name         string
	ActionID     string
	Params       map[string]string
	APIKey       string
	AccessToken  string
	UserAgent    string
	Client       *http.Client
	FunctionCall func(ctx context.Context, input string) (string, error)
	Description  string
	Header       map[string]string
	Host         string
}

type dynamicTool struct {
	CallbacksHandler callbacks.Handler
	client           *http.Client
	functionCall     func(ctx context.Context, input string) (string, error)
	name             string
	description      string
	actionID         string
	params           map[string]string
}

func (d dynamicTool) Name() string {
	return d.name
}

func (d dynamicTool) Description() string {
	return d.description
}

func (d dynamicTool) Call(ctx context.Context, input string) (string, error) {
	fmt.Println(fmt.Sprintf("调用工具: %s, 输入: %s", d.Name(), input))
	if d.CallbacksHandler != nil {
		d.CallbacksHandler.HandleToolStart(ctx, input)
	}
	var result string
	// TODO 这里应该替换成http.Client 请求
	result, err = d.functionCall(ctx, input)
	if err != nil {
		if d.CallbacksHandler != nil {
			d.CallbacksHandler.HandleToolError(ctx, err)
		}
		return "", err
	}

	if d.CallbacksHandler != nil {
		d.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}

func NewDynamicTool(opts ToolOptions) tools.Tool {
	t := &dynamicTool{
		functionCall: opts.FunctionCall,
		name:         opts.Name,
		actionID:     opts.ActionID,
		params:       opts.Params,
		description:  opts.Description,
		client:       opts.Client,
	}

	return t
}
