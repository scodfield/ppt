package dbuffer

import (
	"context"
	"fmt"
	"ppt/login/db"
	"sync"
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
	Step  int64
}

type FunctionUUID struct {
	Type        int
	Buffer1     BasicID
	Buffer2     BasicID
	UseBuffer1  bool
	Offset      int64
	IsNewIDSync bool
	sync.Mutex
	syncMutex sync.Mutex
}

type DoubleBuffer struct {
	Functions []*FunctionUUID
}

func InitDBuffer() {
	DBuffer = &DoubleBuffer{
		Functions: []*FunctionUUID{},
	}
	for idType := UUID_TYPE_ITEMID; idType < UUID_TYPE_MAX; idType++ {
		DBuffer.Functions = append(DBuffer.Functions, initFunctionID(idType))
	}
}

func initFunctionID(id int) *FunctionUUID {
	idType := formatFunctionType(id)
	maxID, step := db.GetFunctionMaxID(idType)
	return &FunctionUUID{
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

func GetUUIDByType(id int) (newID int64) {
	var useFunction *FunctionUUID
	var curStep int64
	var isSwitch bool
	for _, function := range DBuffer.Functions {
		if function.Type != id {
			continue
		}

		function.Lock()
		if function.UseBuffer1 {
			newID = function.Buffer1.MaxID + function.Offset
			curStep = function.Buffer1.Step
			if newID >= function.Buffer1.MaxID+function.Buffer1.Step-1 {
				isSwitch = true
				curStep = function.Buffer2.Step
			}
		} else {
			newID = function.Buffer2.MaxID + function.Offset
			curStep = function.Buffer2.Step
			if newID >= function.Buffer2.MaxID+function.Buffer2.Step-1 {
				isSwitch = true
				curStep = function.Buffer1.Step
			}
		}
		if isSwitch {
			function.UseBuffer1 = !function.UseBuffer1
			function.Offset = 0
			function.IsNewIDSync = false
		} else {
			function.Offset++
		}
		function.Unlock()

		useFunction = function
		break
	}
	if useFunction.Offset >= (curStep/2) && !useFunction.IsNewIDSync {
		go SyncNewID(useFunction)
	}
	return
}

func SyncNewID(function *FunctionUUID) {
	function.syncMutex.Lock()
	defer function.syncMutex.Unlock()
	if !function.IsNewIDSync {
		idType := formatFunctionType(function.Type)
		maxID, step := db.GetFunctionMaxID(idType)
		if function.UseBuffer1 {
			function.Buffer2.MaxID = maxID
			function.Buffer2.Step = step
		} else {
			function.Buffer1.MaxID = maxID
			function.Buffer1.Step = step
		}
		function.IsNewIDSync = true
	}
}
