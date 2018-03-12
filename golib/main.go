package main

import (
	"fmt"
	"golib/goicon"
	"golib/gopw"
)

func main() {
	testgopw()
	testgoicon()
}

func testgopw() {
	pwd, _ := gopw.Create()
	fmt.Println(pwd)
}

func testgoicon() {
	ico, _ := goicon.Get("http://www.baidu.com/")
	fmt.Println(ico)
}
