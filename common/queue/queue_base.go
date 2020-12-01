package queue

import (
	"fmt"
)

type QueueNode struct {
	Data interface{}
	Next *QueueNode
}

//创建链列（数据）
func (queue *QueueNode) Create(Data ...interface{}) {
	if queue == nil {
		return
	}
	if len(Data) == 0 {
		return
	}

	//创建链列
	for _, v := range Data {
		newNode := new(QueueNode)
		newNode.Data = v

		queue.Next = newNode
		queue = queue.Next
	}

}

//打印链列
func (queue *QueueNode) Print() {
	if queue == nil {
		return
	}
	for queue != nil {
		if queue.Data != nil {
			fmt.Print(queue.Data, " ")
		}
		queue = queue.Next
	}
	fmt.Println()
}

func (queue *QueueNode) Get() (data interface{}) {
	data = nil
	if queue == nil {
		return
	}
	if queue.Next.Data != nil {
		data = queue.Next.Data
	}
	queue.Next = queue.Next.Next
	return
}

//链列个数
func (queue *QueueNode) Length() int {
	if queue == nil {
		return -1
	}

	i := 0
	for queue.Next != nil {
		i++
		queue = queue.Next
	}
	return i
}

//入列(insert)
func (queue *QueueNode) Push(Data interface{}) {
	//放在队列的末尾

	if queue == nil {
		return
	}
	if Data == nil {
		return
	}

	//找到队列末尾
	for queue.Next != nil {
		queue = queue.Next
	}

	//创建新节点 将新节点加入队列末尾
	newNode := new(QueueNode)
	newNode.Data = Data

	queue.Next = newNode
}

//出队(delete)
func (queue *QueueNode) Pop() {
	//队头出列
	if queue == nil {
		return
	}
	//记录列队第一个的节点
	//node:=queue.Next
	//queue.Next=node.Next

	queue.Next = queue.Next.Next
}
