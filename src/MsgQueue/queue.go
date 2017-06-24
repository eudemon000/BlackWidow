package MsgQueue

import (
	"container/list"
	//"fmt"
)


//消息队列

type QueenManager struct {
	list	*list.List
	ch	chan int
	size	int

}

type QueueHandler func(data interface{})

var h QueueHandler

var q *QueenManager

func InitManager() *QueenManager {
	q = new(QueenManager)
	q.list = list.New()
	q.size = q.list.Len()
	q.ch = make(chan int)
	return q
}


/*
//读取队列
func (q *QueenManager)Pull() {
		var n *list.Element
		for e := q.list.Front(); e != nil; e = n {
			switch e.Value.(type) {
			case string:
				fmt.Println("pull", e.Value.(string))
			}

		}
		q.ch <- 1
}
*/

//读取队列
func (q *QueenManager)Pull() *list.List {
	l := list.New()
	l.PushBackList(q.list)
	return l
}

func (q * QueenManager)Push(data interface{}) {
	q.list.PushBack(data)
	q.size = q.list.Len()
}


/*func (q * QueenManager)Push(data interface{}) {
	q.list.PushBack(data)
	q.size = q.list.Len()
	<- q.ch
}*/

