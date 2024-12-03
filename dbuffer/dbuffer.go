package dbuffer

import (
	"context"
	"fmt"
	"ppt/login/db"
)

var ctx = context.Background()

const (
	UUID_TYPE_ITEMID = iota
	UUID_TYPE_PETID
	UUID_TYPE_MAX
)

var DBuffer *DoubleBuffer

type BasicID struct {
	MaxID int64
	Step  int32
}

type FunctionUUID struct {
	Type       int
	Buffer1    BasicID
	Buffer2    BasicID
	UseBuffer1 bool
}

type DoubleBuffer struct {
	Functions []FunctionUUID
}

func InitDBuffer() {
	DBuffer = &DoubleBuffer{
		Functions: []FunctionUUID{},
	}
	for idType := UUID_TYPE_ITEMID; idType < UUID_TYPE_MAX; idType++ {
		DBuffer.Functions = append(DBuffer.Functions, initFunctionID(idType))
	}
}

func initFunctionID(id int) FunctionUUID {
	idType := formatFunctionType(id)
	maxID, step := db.GetFunctionMaxID(idType)
	return FunctionUUID{
		Type: id,
		Buffer1: BasicID{
			MaxID: maxID,
			Step:  step,
		},
		UseBuffer1: true,
	}
}

func formatFunctionType(idType int) string {
	return fmt.Sprintf("global_uuid_%d", idType)
}
