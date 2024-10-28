package actor

const MSG_QUEUE_SIZE = 20

type MsgNode struct {
	Msg       Msg
	pre, next *MsgNode
}

type MsgQueue struct {
	head, tail *MsgNode
	size       int
}

func NewMsgQueue() MsgQueue {
	return MsgQueue{}
}

func (mq *MsgQueue) Enqueue(msg Msg) {
	if mq.size == 0 || mq.head == nil || mq.tail == nil {
		mq.head, mq.tail = &MsgNode{}, &MsgNode{}
		mq.head.next = mq.tail
		mq.tail.pre = mq.head
	}
	newNode := &MsgNode{Msg: msg}
	newNode.pre, newNode.next = mq.tail.pre, mq.tail
	mq.tail.pre.next, mq.tail.pre = newNode, newNode
	mq.size++
}

func (mq *MsgQueue) Dequeue() Msg {
	if mq.size == 0 {
		return nil
	}
	ret := mq.head.next
	mq.head.next = mq.head.next.next
	mq.head.next.pre = mq.head
	mq.size--
	return ret.Msg
}

func (mq *MsgQueue) Size() int {
	return mq.size
}
