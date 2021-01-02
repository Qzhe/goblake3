package goblake3

import (
	"encoding/base64"
	"testing"
)

func TestNew(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Logf("%v", err)
		}
	}()

	b3h := New()

	buf := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		buf[i] = byte(i)
	}
	result := b3h.Finalize()
	b3hStr := base64.StdEncoding.EncodeToString(result)
	t.Logf("Result:%s", b3hStr)
}
