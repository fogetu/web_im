package gotest

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithDeadline(t *testing.T) {
	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(10)*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	chiHanBao2(ctx)
	defer cancel() // 防止意外
}

func chiHanBao2(ctx context.Context) {
	n := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("stop \n")
			return
		default:
			incr := getRand(5) + 1
			n += incr
			fmt.Printf("我吃了 %d 个汉堡\n", n)
		}
		time.Sleep(time.Second)
	}
}
