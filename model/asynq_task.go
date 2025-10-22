package model

type PptAsynqTaskType int32

const (
	AsynqTaskTypeDefault PptAsynqTaskType = iota * 100
	AsynqTaskTypeMail
	AsynqTaskTypeNotice
)

type PptAsynqTask struct {
	TaskID   string           `json:"task_id"`
	TaskType PptAsynqTaskType `json:"task_type"`
}
