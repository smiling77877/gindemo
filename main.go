package main

import "fmt"

func testv1() {
	for i := 0; i < 10; i++ {
		j := i
		fmt.Printf("循环 %p, %p \n", &i, &j)
		defer func() {
			fmt.Printf("%p \n", &j)
			println(j)
		}()
	}
	println("跳出循环")
}

func testv2() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("%p \n", &i)
			println(i)
		}()
	}
}

func main() {
	testv2()
}
