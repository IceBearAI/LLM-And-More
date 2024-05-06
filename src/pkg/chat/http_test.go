package chat

import (
	"testing"
)

func TestMakeHTTPHandler(t *testing.T) {
	t.Skip("Not implemented")
}

func TestHTTP_ChatCompletionStream(t *testing.T) {
	svc := initSvc()
	httpHandler := MakeHTTPHandler(svc, nil, nil)

}
