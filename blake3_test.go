package goblake3

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Logf("%v", err)
		}
	}()

	b3h := New()

	start := time.Now()
	buf := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		buf[i] = byte(i)
	}
	b3h.Update(buf)
	result := b3h.Finalize()
	b3hStr := base64.StdEncoding.EncodeToString(result)
	t.Logf("Result:%s\nTime:%d\n", b3hStr, time.Since(start).Microseconds())

	start = time.Now()
	var output []byte
	for i := 0; i < 1024; i++ {
		for j := 0; j < 1024; j++ {
			buf[j] = byte(j)
		}
		b3h.Update(buf)
		output = b3h.FinalizeSeek()
	}
	b3hStr = base64.StdEncoding.EncodeToString(output)
	t.Logf("Result:%s\nTime:%d\n", b3hStr, time.Since(start).Microseconds())
}

func TestNew2(t *testing.T) {
	fmt.Println("seekHash---------------")
	do(seekHash)
	fmt.Println("hash-------------------")
	do(hash)
}

func hash(hasher *Blake3Hasher) {
	buf := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		buf[i] = byte(i)
	}
	hasher.Update(buf)
	hasher.Finalize()
}

func seekHash(hasher *Blake3Hasher) {
	buf := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		buf[i] = byte(i)
	}
	hasher.Update(buf)
	hasher.FinalizeSeek()
}

func do(cb func(hasher *Blake3Hasher)) {
	b := New()
	p := context.TODO()
	c, _ := context.WithTimeout(p, 10*time.Second)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	counter := 0
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-c.Done():
				return
			default:
				counter++
				cb(b)
			}
		}
	}(c)
	wg.Wait()
	fmt.Printf("Times:%d\n", counter)
}