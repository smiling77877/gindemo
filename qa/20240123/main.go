package main

import "fmt"

func main() {
	Defer()
}

func Defer() {
	i := 0
	fmt.Printf("print(i): %v\n", i)
	defer fmt.Printf("defer 直接形参调用 print(i): %v\n", i)

	defer func() {
		fmt.Printf("defer 闭包函数 print(i): %v\n", i)
	}()

	i++
	fmt.Printf("print(i): %v\n", i)
}
