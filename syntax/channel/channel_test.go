package channel

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	// 声明
	//var ch chan struct{}
	// 声明并创建
	//ch1 := make(chan int)
	// 这种是带buffer的
	ch2 := make(chan int, 3)
	// 把 123 发送到 ch2 里面
	ch2 <- 123
	data := <-ch2
	t.Log(data)
	// 这个是关闭 channel
	close(ch2)
}

func TestChannel2(t *testing.T) {
	c := make(chan struct{})
	close(c)
	assert.Panics(t, func() {
		// 往已经关闭的channel写数据，会 Panic
		c <- struct{}{}
	})
}

func TestChannelClose(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 0
	// ok 代表没有读到
	val, ok := <-ch
	t.Log("读到了数据吗？", ok, val)
	close(ch)
	// 这个操作会引起 panic
	//ch <- 123
	val, ok = <-ch
	t.Log("读到了数据吗？", ok, val)

	// 也会 panic
	//close(ch)
}

func TestChannelLoop(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
			time.Sleep(time.Second)
		}
		close(ch)
	}()

	for val := range ch {
		t.Log(val)
	}
}

func TestChannelBlocking(t *testing.T) {
	ch := make(chan int)
	b1 := BigStruct{}
	go func() {
		var b BigStruct
		// 这个就是goroutine泄露
		ch <- 123
		t.Log(b, b1)
	}()
}

func TestChannelBlocking2(t *testing.T) {
	// 没有初始化，c == nil
	var c chan struct{}
	// 这两个都会导致阻塞

	go func() {
		<-c
		t.Log("111不会输出这一句")
	}()
	var b1 BigStruct
	go func() {
		var b2 BigStruct
		c <- struct{}{}
		t.Log("222不会输出这一句", b1, b2)
	}()
}

type BigStruct struct{}

func TestChannelSelect(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 2)
	go func() {
		time.Sleep(time.Second)
		ch1 <- 123
	}()

	go func() {
		time.Sleep(time.Second)
		ch2 <- 123
	}()
	select {
	case val := <-ch1:
		t.Log("进来了 ch1 这里", val)
	case val := <-ch2:
		t.Log("进来了 ch2 这里", val)
	}
}

func Close[T io.Closer](t T) {
	t.Close()
}
