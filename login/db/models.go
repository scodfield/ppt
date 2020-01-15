package db

import (
	"time"
)

type User struct {
	Id int
	AccId  int
	Name string 
	Password string 
}


type LoginLog struct {
	Id int `orm:pk;auto`
	AccId int `orm:index`
	Name string
	LoginTime time.Time
}