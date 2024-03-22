package s3

import (
	"context"
	"os"
	"testing"
)

func initSvc() Service {
	return New("", "", "", "", "default")
}

func TestService_Upload(t *testing.T) {
	pwd, _ := os.Getwd()
	f, _ := os.Open(pwd + "/service_test.go")
	//f, _ := os.ReadFile(pwd + "/service_test.go")
	//osFileToMultipartFile(f)
	err := initSvc().Upload(context.Background(), "", "test/test.txt", f, "application/json; charset=utf-8")
	if err != nil {
		t.Error(err)
		return
	}
}
