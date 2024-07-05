package util

import "testing"

func Test_GetFileList(t *testing.T) {
	dirPath := "/Users/cong/.cache/huggingface/hub/models--Qwen--Qwen2-0.5B-Instruct/snapshots/c291d6fce4804a1d39305f388dd32897d1f7acc4"
	files, err := GetFileList(dirPath)
	if err != nil {
		t.Error(err)
		return
	}
	for _, file := range files {
		t.Logf("%+v", file)
	}
}
