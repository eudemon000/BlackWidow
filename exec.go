package main

import (
	"BlackWidow/src/Scheduler"
	"time"
)

func main() {
	s := new(Scheduler.SpiderDispther)
	s.Start_urls = make([] string, 1)
	s.Start_urls[0] = "http://www.99.com.cn"
	Scheduler.InitData(s)
	/*var ch chan int
	ch <- 1*/
	time.Sleep(time.Second * 10)
}

