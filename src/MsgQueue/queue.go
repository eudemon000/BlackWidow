package MsgQueue

import (
	"container/list"
	"fmt"
)


//消息队列

type QueenManager struct {
	list	*list.List
	ch	chan int
	size	int

}

type QueueHandler func(data interface{})

var h QueueHandler

func InitManager() *QueenManager {
	q := new(QueenManager)
	q.list = list.New()
	q.size = q.list.Len()
	q.ch = make(chan int)
	go q.pull()
	return q
}


//读取队列
func (q *QueenManager)pull() {
	for {
		var n *list.Element
		for e := q.list.Front(); e != nil; e = n {
			switch e.Value.(type) {
			case string:
				//h(e.Value.(string))
				fmt.Println("pull", e.Value.(string))
			}
			n = e.Next()
			q.list.Remove(e)
		}
		q.ch <- 1
	}
}

func (q * QueenManager)Push(data interface{}) {
	q.list.PushBack(data)
	q.size = q.list.Len()
	<- q.ch
}

