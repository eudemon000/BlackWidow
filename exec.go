package main

import (
	"BlackWidow/src/Scheduler"
	"time"
	"fmt"
)

func main() {
	s := new(Scheduler.SpiderDispther)
	s.Start_urls = make([]string, 1)
	s.Start_urls[0] = "http://www.99.com.cn"
	s.Allowed_domains = make([]string, 1)
	s.Allowed_domains[0] = "99.com.cn"
	//s.Start_urls[0] = "http://www.qq.com"
	Scheduler.InitData(s)
	var ch chan int
	ch <- 1
	time.Sleep(time.Second * 10)
	fmt.Println("正常执行完毕")
}

