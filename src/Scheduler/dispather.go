package Scheduler

import (
	dbbase "BlackWidow/src/Database"
	"runtime"
	_"fmt"
	_"time"
	mq "BlackWidow/src/MsgQueue"
)

type SpiderDispther struct {
	Start_urls	[]string
}

var sDispther *SpiderDispther

var qm *mq.QueenManager

//初始化分发器
func InitData (s *SpiderDispther) {
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)
	sDispther = s
	/*for _, item := range  sDispther.Start_urls {
		checkSpiderExist(item)
	}*/
	qm = mq.InitManager()

	go getWaitUrls()
	go checkSpiderExist(s.Start_urls)

}

var checkChan chan int = make(chan int)

/*
1、检查该url是否爬取过，如果爬取过就返回
2、检查待爬取列表里是否有该url，如果没有就添加到列表
3、通知获取待爬url的函数开始取数据
 */
func checkSpiderExist (urls []string) {
	for _, url := range urls {
		isExist := dbbase.GetUrl(url)
		if isExist {
			continue
		}
		waitIsExist := dbbase.GetWaitUrl(url)
		if !waitIsExist {
			dbbase.InsertWaitUrl(url)
		}
	}
	<- checkChan
}

//从数据库查找待爬取的url
func getWaitUrls() {
	for {
		checkChan <- 1
		urls := dbbase.GetWaitUrls()
		for _, item := range urls {
			//fmt.Println(index, item.Url)
			qm.Push(item.Url)
		}
	}
}



