package encode

import (
	"github.com/pkg/errors"
)

type ResStatus string

var ResponseMessage = map[ResStatus]int{
	Invalid:        400,
	InvalidParams:  400,
	ErrParamsPhone: 401,
	ErrBadRoute:    401,
	ErrSystem:      500,
	ErrNotfound:    404,
	ErrLimiter:     429,

	ErrAccountNotFound:    404,
	ErrAccountLogin:       1002,
	ErrAccountLoginIsNull: 1003,
	ErrAccountNotLogin:    501,
	ErrAccountASD:         1004,
	ErrAccountLocked:      1005,
	ErrAuthNotLogin:       501,
	ErrAuthLogin:          501,

	ErrParamsNamespace: 1100,
	ErrParamsService:   1101,
	ErrAuthTimeout:     1102,

	ErrConversationIdNotFound:        1200,
	ErrMessageIdNotFound:             1201,
	ErrChatChannelNotFound:           1202,
	ErrChatChannelApiKey:             1203,
	ErrChatChannelModelId:            1204,
	ErrChatChannelModelIdNotAllow:    1205,
	ErrFastChatTokensNumFromMessages: 1206,
	ErrChatCreateOpenAI:              1207,
	ErrAudioTranslationTextEmpty:     1208,
	ErrTenantNotFound:                1209,

	ErrDatasetNotFound: 1300,

	ErrToolNotFound: 1350,

	ErrAssistantNotFound: 1400,
}

const (
	// 公共错误信息
	Invalid                 ResStatus = "invalid"
	InvalidParams           ResStatus = "请求参数错误"
	ErrNotfound             ResStatus = "不存在"
	ErrBadRoute             ResStatus = "请求路由错误"
	ErrParamsPhone          ResStatus = "手机格式不正确"
	ErrParamsNamespace      ResStatus = "项目空间不能为空"
	ErrParamsService        ResStatus = "服务不能为空"
	ErrLimiter              ResStatus = "请求过于频繁，请稍后再试"
	ErrServerStartDbConnect ResStatus = "数据库连接失败"
	ErrAuthTimeout          ResStatus = "授权已过期"

	// 中间件错误信息
	ErrSystem                  ResStatus = "系统错误"
	ErrAccountNotLogin         ResStatus = "用户没登录"
	ErrAuthNotLogin            ResStatus = "请先登录"
	ErrAccountLoginIsNull      ResStatus = "用户名和密码不能为空"
	ErrAccountLogin            ResStatus = "用户名或密码错误"
	ErrAccountNotFound         ResStatus = "账号不存在, 请联系管理员"
	ErrAccountASD              ResStatus = "权限验证失败"
	ErrAccountLocked           ResStatus = "用户已被锁定"
	ErrAccountUpdate           ResStatus = "用户更新失败"
	ErrAuthLogin               ResStatus = "登录失败"
	ErrAuthCheckCaptchaCode    ResStatus = "图形验证码错误"
	ErrAuthCheckCaptchaNotnull ResStatus = "图形验证码不能为空"
	ErrTenantNotFound          ResStatus = "租户不存在"
	// ErrConversationIdNotFound Chat
	ErrConversationIdNotFound        ResStatus = "会话会ID不能为空"
	ErrMessageIdNotFound             ResStatus = "会话消息ID不能为空"
	ErrChatChannelNotFound           ResStatus = "获取渠道不存在或已失效，请联系管理员"
	ErrChatChannelApiKey             ResStatus = "渠道APIKEY不能为空"
	ErrChatChannelModelId            ResStatus = "渠道模型不被允许"
	ErrChatChannelModelIdNotAllow    ResStatus = "模型不被允许"
	ErrFastChatTokensNumFromMessages ResStatus = "请求Tokens数量超过了此模型所支持的最大Tokens数量"
	ErrChatCreateMessage             ResStatus = "创建消息失败"
	ErrChatCreateOpenAI              ResStatus = "对话生成失败"
	ErrAudioTranslationTextEmpty     ResStatus = "音频翻译文本为空"

	ErrDatasetNotFound ResStatus = "数据集不存在"

	ErrToolNotFound ResStatus = "工具不存在"

	ErrAssistantNotFound ResStatus = "助手不存在"
)

func (c ResStatus) String() string {
	return string(c)
}

func (c ResStatus) Error() error {
	return errors.New(string(c))
}

func (c ResStatus) Wrap(err error) error {
	return errors.Wrap(err, string(c))
}
