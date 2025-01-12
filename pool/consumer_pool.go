package pool

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
)

var consumerPool *ants.Pool

func InitConsumerPool(num int) error {
	var err error
	consumerPool, err = ants.NewPool(num)
	if err != nil {
		fmt.Println("init consumer pool err:", err)
		return err
	}
	return nil
}

func CloseConsumerPool() {
	if consumerPool != nil {
		consumerPool.Release()
	}
}

func GetConsumerPool() *ants.Pool {
	return consumerPool
}
