package encode

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"reflect"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
)

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"-"`
	Message string      `json:"message"`
	TraceId string      `json:"traceId"`
	Stream  bool        `json:"-"`
}

type Failure interface {
	Failed() error
}

type Errorer interface {
	Error() error
}

func Error(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

func JsonError(ctx context.Context, err error, w http.ResponseWriter) {
	headers, ok := ctx.Value("response-headers").(map[string]string)
	if ok {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}

	var errDefined bool
	for k := range ResponseMessage {
		if strings.Contains(err.Error(), k.Error().Error()) {
			errDefined = true
			break
		}
	}

	if !errDefined {
		err = ErrSystem.Wrap(err)
	}
	if err == nil {
		err = errors.Wrap(err, ErrSystem.Error().Error())
	}
	traceId, _ := ctx.Value("traceId").(string)
	var code int
	code = ResponseMessage[ResStatus(strings.Split(err.Error(), ":")[0])]
	if code == 0 {
		code = 500
	}
	w.Header().Set("TraceId", traceId)
	_ = kithttp.EncodeJSONResponse(ctx, w, map[string]interface{}{
		"message": err.Error(),
		"code":    code,
		"success": false,
		"traceId": traceId,
	})
}

func JsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(Failure); ok && f.Failed() != nil {
		JsonError(ctx, f.Failed(), w)
		return nil
	}
	resp, ok := response.(Response)
	if !ok {
		log.Println("response is not Response type", reflect.TypeOf(response))
	}
	if resp.Error == nil {
		resp.Code = 200
		resp.Success = true
	} else {
		var code int
		code = ResponseMessage[ResStatus(strings.Split(resp.Error.Error(), ":")[0])]
		if code == 0 {
			code = 500
		}
		resp.Code = code
		resp.Message = resp.Error.Error()
	}

	headers, ok := ctx.Value("response-headers").(map[string]string)
	if ok {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}
	traceId, _ := ctx.Value("traceId").(string)
	resp.TraceId = traceId
	w.Header().Set("TraceId", traceId)
	return kithttp.EncodeJSONResponse(ctx, w, resp)
}
