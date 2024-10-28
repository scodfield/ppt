package actor

import (
	"time"
)

type RoleActor struct {
	roleID         uint64
	msgQueue       MsgQueue
	msgChan        chan Msg
	closedChan     chan bool
	procClosedChan chan bool
	logout         bool
}

func (role *RoleActor) GetRoleID() uint64 {
	return role.roleID
}

func (role *RoleActor) Logout() bool {
	return role.logout
}

func (role *RoleActor) GetMsgChan() chan Msg {
	return role.msgChan
}

func (role *RoleActor) Start() {
	go role.MsgReceiveLoop()
	go role.MsgProcLoop()
}

func (role *RoleActor) Stop() {
	// 退出消息收发
	close(role.msgChan)
	<-role.closedChan
	close(role.closedChan)
	// 退出消息处理
	role.procClosedChan <- true
	close(role.procClosedChan)
}

func (role *RoleActor) MsgReceiveLoop() {
	for msg := range role.msgChan {
		role.msgQueue.Enqueue(msg)
	}
	role.closedChan <- true
}

func (role *RoleActor) MsgProcLoop() {
	ticker := time.NewTicker(time.Millisecond * 10)
	for {
		if msg := role.msgQueue.Dequeue(); msg != nil {
			msg.Proc()
			continue
		}
		select {
		case <-role.procClosedChan:
			goto END
		case <-ticker.C:
			//TODO
		}
	}
END:
	//  处理剩余消息
	if role.msgQueue.Size() > 0 {
		for i := 0; i < role.msgQueue.Size(); i++ {
			if msg := role.msgQueue.Dequeue(); msg != nil {
				msg.Proc()
			}
		}
	}
	role.logout = true
}

func NewRoleActor(roleID uint64) Actor {
	role := &RoleActor{}
	role.roleID = roleID
	role.msgQueue = NewMsgQueue()
	role.msgChan = make(chan Msg, MSG_QUEUE_SIZE)
	role.closedChan = make(chan bool)
	role.procClosedChan = make(chan bool)
	return role
}

func SendToMsg(actor RoleActor, msg Msg) {
	actor.GetMsgChan() <- msg
}
