package Scheduler

import (
dbbase "BlackWidow/src/Database"
"runtime"
_"fmt"
_"time"
mq "BlackWidow/src/MsgQueue"
"fmt"
_"BlackWidow/src/Downloader"
"sync"
"BlackWidow/src/Downloader"
"strings"
	"time"
)

type SpiderDispther struct {
	Start_urls		[]string
	Allowed_domains	[]string
	lock			*sync.Mutex
}

var sDispther *SpiderDispther

var qm *mq.QueenManager

//初始化分发器
func InitData (s *SpiderDispther) {
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)
	sDispther = s
	sDispther.lock = new(sync.Mutex)
	/*for _, item := range  sDispther.Start_urls {
		checkSpiderExist(item)
	}*/
	qm = mq.InitManager(func(data interface{}){
		fmt.Println("回调", data)
		results := Downloader.Parser(data.(string))
		fmt.Println(len(results))
		sDispther.pull(results)
	})

	go getWaitUrls()
	//checkSpiderExist(s.Start_urls)
	go checkSpiderExist()

	go recordTime()

}

var checkChan chan int = make(chan int)

/*
1、检查该url是否爬取过，如果爬取过就返回
2、检查待爬取列表里是否有该url，如果没有就添加到列表
3、通知获取待爬url的函数开始取数据
 */
/*func checkSpiderExist (urls []string) {
	for _, url := range urls {
		fmt.Println(url)
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
}*/

func checkSpiderExist () {
	for {
		for _, url := range sDispther.Start_urls {
			fmt.Println("检查URL===>", url)
			isExist := dbbase.GetUrl(url)
			if isExist {
				continue
			}
			//fmt.Println("检查通过===>", url)
			waitIsExist := dbbase.GetWaitUrl(url)
			if !waitIsExist {
				dbbase.InsertWaitUrl(url)
			}
		}
		sDispther.Start_urls = []string{}
		<-checkChan
	}
}

//从数据库查找待爬取的url
func getWaitUrls() {
	for {
		checkChan <- 1
		urls := dbbase.GetWaitUrls()
		for _, item := range urls {
			//fmt.Println(index, item.Url)
			//判断URL是否为Allowed_domains允许的URL
			for _, domain := range sDispther.Allowed_domains {
				if strings.Contains(item.Url, domain) {
					qm.Push(item.Url)
					fmt.Println(item.Url, "已放入队列")
				}
			}

		}
	}
}

func(s *SpiderDispther) pull(urls []string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, item := range urls {
		s.Start_urls = append(s.Start_urls, item)
	}
}

func recordTime() {
	for {
		/*year := time.Now().Year()
		month := time.Now().Month()
		day := time.Now().Day()*/
		fmt.Println(time.Now())
		time.Sleep(10 * time.Second)
	}
}



