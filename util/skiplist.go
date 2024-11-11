package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	MAX_LEVEL int = 32
)

type SkipNode struct {
	right *SkipNode
	down  *SkipNode
	val   int
}

type SkipList struct {
	head     *SkipNode
	level    int
	pathList []*SkipNode
}

func (skip *SkipList) clearPathList() {
	if skip == nil {
		return
	}
	skip.pathList = []*SkipNode{}
}

func (skip *SkipList) pushPath(node *SkipNode) {
	if skip == nil || node == nil {
		return
	}
	skip.pathList = append(skip.pathList, node)
}

func (skip *SkipList) popPath() *SkipNode {
	if skip == nil || len(skip.pathList) == 0 {
		return nil
	}
	node := skip.pathList[len(skip.pathList)-1]
	skip.pathList = skip.pathList[:len(skip.pathList)-1]
	return node
}

func (skip *SkipList) ToString() string {
	if skip == nil || skip.head == nil {
		return ""
	}
	var result = []string{}
	head := skip.head
	for head.down != nil {
		head = head.down
	}
	for head.right != nil {
		result = append(result, fmt.Sprintf("%v", head.right.val))
		head = head.right
	}
	return fmt.Sprintf("Skiplist:%s\n", strings.Join(result, ","))
}

func newNode(right *SkipNode, down *SkipNode, val int) *SkipNode {
	return &SkipNode{
		right: right,
		down:  down,
		val:   val,
	}
}

func (skip *SkipList) Insert(val int) {
	if skip == nil {
		return
	}
	skip.clearPathList()
	cur := skip.head
	if cur.val == val {
		return
	}
	for cur != nil {
		for cur.right != nil && cur.right.val < val {
			cur = cur.right
		}
		if cur.right != nil && cur.val == val {
			return
		}
		skip.pathList = append(skip.pathList, cur)
		cur = cur.down
	}
	insertFlag := true
	var downNode *SkipNode
	for insertFlag && len(skip.pathList) > 0 {
		cur = skip.popPath()
		cur.right = newNode(cur.right, downNode, val)
		downNode = cur.right
		insertFlag = rand.Intn(100) < 30
	}
	if insertFlag && (skip.level < MAX_LEVEL) {
		skip.head = newNode(newNode(nil, downNode, val), skip.head, skip.head.val)
		skip.level++
	}
}

func (skip *SkipList) Delete(val int) bool {
	if skip == nil || skip.head == nil {
		return false
	}
	found := false
	cur := skip.head
	for cur != nil {
		for cur.right != nil && cur.right.val < val {
			cur = cur.right
		}
		if cur.right != nil && cur.right.val == val {
			cur.right = cur.right.right
			found = true
			cur = cur.down
		} else if cur.right != nil && cur.right.val > val {
			cur = cur.down
		}
	}
	return found
}

func (skip *SkipList) Search(val int) bool {
	if skip == nil || skip.head == nil {
		return false
	}
	cur := skip.head
	for cur != nil {
		for cur.right != nil && cur.right.val < val {
			cur = cur.right
		}
		if cur.right != nil && cur.right.val == val {
			return true
		} else if cur.right != nil && cur.right.val > val {
			cur = cur.down
		} else {
			return false
		}
	}
	return false
}

func NewSkipList() *SkipList {
	return &SkipList{
		head: &SkipNode{
			val: -1,
		},
		level:    1,
		pathList: make([]*SkipNode, 0),
	}
}
