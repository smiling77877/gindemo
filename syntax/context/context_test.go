package context

import (
	"context"
	"testing"
	"time"
)

type key struct {
}

func TestContextValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), key{}, "value1")
	val, ok := ctx.Value(key{}).(string)
	t.Log(val, ok)
}

func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		t.Log("准备调用 cancel 了")
		cancel()
	}()
	<-ctx.Done()
	t.Log("已经 cancel 了")
	t.Log(ctx.Err())
}

func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	<-ctx.Done()
	t.Log("超时了")
	t.Log(ctx.Err())
}

func TestContextParentCancel(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second, func() {
		cancel()
	})
	son, sonCancel := context.WithCancel(parent)
	<-son.Done()
	t.Log("son 已经过来了")
	sonCancel()
}

func TestContextParentCancel1(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())

	_, sonCancel := context.WithCancel(parent)
	time.AfterFunc(time.Second, func() {
		sonCancel()
	})
	<-parent.Done()
	t.Log("parent 过来了")
	cancel()
}
