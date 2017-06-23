package main

import (
	"BlackWidow/src/Scheduler"
	"time"
)

func main() {
	s := new(Scheduler.SpiderDispther)
	s.Start_urls = make([]string, 3)
	s.Start_urls[0] = "http://www.99.com.cn"
	s.Start_urls[1] = "http://www.qq.com"
	s.Start_urls[2] = "http://www.google.com"
	Scheduler.InitData(s)
	/*var ch chan int
	ch <- 1*/
	time.Sleep(time.Second * 10)
}

