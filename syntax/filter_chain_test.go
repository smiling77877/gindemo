package syntax

import (
	"log"
	"testing"
)

// 函数式的 Filter
type Filter func() error

type FilterChain func(next Filter) Filter

type MyServer struct {
	root Filter
}

func NewMyServer(flts ...FilterChain) *MyServer {
	var root Filter = func() error {
		log.Println("这是最后一个 filter", "正常来说也是框架核心")
		return nil
	}
	// 从后往前组装
	for i := len(flts) - 1; i >= 0; i-- {
		root = flts[i](root)
	}
	return &MyServer{root: root}
}

func (m *MyServer) Serve() error {
	return m.root()
}

func TestMyServer(t *testing.T) {
	var first FilterChain = func(next Filter) Filter {
		return func() error {
			log.Println("第一个执行前")
			err := next()
			log.Println("第一个执行后")
			return err
		}
	}
	var second FilterChain = func(next Filter) Filter {
		return func() error {
			log.Println("第二个执行前")
			err := next()
			log.Println("第二个执行后")
			return err
		}
	}

	server := NewMyServer(first, second)
	server.Serve()
}
