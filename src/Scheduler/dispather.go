package Scheduler

import (
	dbbase "BlackWidow/src/Database"
	"runtime"
	_"fmt"
	_"time"
	mq "BlackWidow/src/MsgQueue"
	"BlackWidow/src/crawler"
	"container/list"
)

type SpiderDispther struct {
	Start_urls	[]string
}

var sDispther *SpiderDispther

var qm *mq.QueenManager

var test *crawler.Test

//初始化分发器
func InitData (s *SpiderDispther) {
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)
	sDispther = s
	/*for _, item := range  sDispther.Start_urls {
		checkSpiderExist(item)
	}*/
	qm = mq.InitManager()
	test = new(crawler.Test)

	go getWaitUrls()
	go checkSpiderExist(s.Start_urls)
	go disTask()

}

var checkChan chan int = make(chan int)

var taskChan chan int = make(chan int)

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
		if len(urls) > 0 {
			for _, item := range urls {
				//fmt.Println(index, item.Url)
				qm.Push(item.Url)
			}
			taskChan <- 1
		}
	}
}

func disTask() {
	for {
		<- taskChan
		l := qm.Pull()
		var n *list.Element

		for e := l.Front(); e != nil; e = n {
			n = e.Next()
			l.Remove(e)
			switch e.Value.(type) {
			case string:
				nextList := test.SpiderWeb(e.Value.(string))
				start_urls := make([]string, nextList.Len())

				var nn *list.Element
				index := 0
				for ne := nextList.Front(); ne != nil; ne = nn {
					nn = ne.Next()
					nextList.Remove(ne)
					start_urls[index] = ne.Value.(string)
					index++
				}
				go checkSpiderExist(start_urls)
			}

		}

	}


}



