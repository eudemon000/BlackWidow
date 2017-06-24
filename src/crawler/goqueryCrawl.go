package crawler

import (
	"fmt"
	"container/list"
	"github.com/PuerkitoBio/goquery"
	logUtil "BlackWidow/src/logPackage"
)

type Spider interface {
	SpiderWeb(string) list.List
}

type Test struct {

}

func (t *Test)SpiderWeb(url string) list.List {
	fmt.Println("spiderWeb", url)
	l := list.List{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		logUtil.Msg(logUtil.Error, err)
		return l
	}
	aTags := doc.Find("body").Find("a")
	aTags.Each(func(index int, aSele *goquery.Selection) {
		result, ok := aSele.Attr("href")
		if ok {
			l.PushBack(result)
		}

	})

	return l
}
