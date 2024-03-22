package util

import (
	"bytes"
	"testing"
)

func TestKrand(t *testing.T) {
	t.Log(string(bytes.ToLower(Krand(32, KC_RAND_KIND_ALL))))
	t.Log(string(Krand(8, KC_RAND_KIND_ALL)))
	t.Log(string(Krand(24, KC_RAND_KIND_ALL)))
	t.Log(string(Krand(12, KC_RAND_KIND_LOWER)))
	t.Log(string(Krand(18, KC_RAND_KIND_UPPER)))
	t.Log(string(Krand(24, KC_RAND_KIND_NUM)))
}

func TestFormatBytes(t *testing.T) {
	t.Log(FormatBytes(1234567890))
}
