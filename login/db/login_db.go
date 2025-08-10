package db

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	redic "ppt/cache"
	"ppt/dao"
	"ppt/logger"
	"ppt/login/utils"
	"ppt/model"
	"strconv"
	"sync"
	"time"
)

var (
	redisCache redic.Cache
	o          orm.Ormer
	count      int
	countMutex sync.Mutex
)

const regHash = "user_reg"

func init() {
	// cache conn
	redisCache, _ = redic.NewCache(`{"conn":"127.0.0.1:6379","password":"foobared"}`)

	// init mysql
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(LoginLog))
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/thd_login?charset=utf8&loc=Local", 30)
	orm.RunSyncdb("default", false, true)

	// create orm instance
	o = orm.NewOrm()

	// init acc_id
	curReg, _ := o.QueryTable("user").Count()
	count = int(curReg)
	fmt.Println("current count: ", count)
}

// SetAccountInfo 设置用户账号信息
func SetAccountInfo(name string, user model.User) error {
	//user.AccId = get_account_id()
	// set cache info & mark registered
	//SetLoginCache(user)
	//redisCache.Command("HSET", regHash, name, 1)
	//_, insertErr := o.Insert(&user)
	//if insertErr != nil {
	//	fmt.Println("mys insert err, ", insertErr)
	//}

	var err error
	now := time.Now().UnixMilli()
	exists := false
	if exists, err = dao.RedisDB.HSetNX(dao.Ctx, dao.UserNameRegisterKey, name, now).Result(); err != nil {
		logger.Error("SetAccountInfo HSetNX user name error", zap.Uint64("user_id", user.UserID), zap.String("user_name", name))
		return err
	}
	if !exists {
		logger.Warn("SetAccountInfo HSetNX user name already exists", zap.Uint64("user_id", user.UserID), zap.String("user_name", name))
		return errors.New("user name already exists")
	}

	return nil
}

func GetAccountInfo(name string) interface{} {
	userCache := GetLoginCache(name)
	if userCache != nil {
		return userCache
	}

	user := User{Name: name}
	readErr := o.Read(&user, "Name")
	if readErr != nil {
		fmt.Println("get account info err: ", readErr)
		return nil
	}
	return user
}

func get_account_id() int {
	countMutex.Lock()
	defer countMutex.Unlock()
	count++
	return count
}

// whether registered
func WhetherRegistered(name string) bool {
	regReply, regErr := redisCache.Command("HGET", regHash, name)
	if regErr != nil {
		fmt.Println("user register error: ", regErr)
		return false
	} else if regReply == nil {
		return false
	}
	regArray := regReply.([]byte)
	regStr := string(regArray)
	if regStr == "1" {
		return true
	}
	return false
}

// user login cache
func SetLoginCache(user User) error {
	cacheDuration := 3600 * time.Second
	name := user.Name
	userMarshal, _ := json.Marshal(user)
	cacheErr := redisCache.SetEx(name+":cache", userMarshal, cacheDuration)
	if cacheErr != nil {
		fmt.Println("set account cache error: ", cacheErr)
		return cacheErr
	}
	return nil
}

// get login cache
func GetLoginCache(name string) interface{} {
	var userCache User
	cacheRes := redisCache.Get(name + ":cache")
	if cacheRes != nil {
		cacheStr := cacheRes.([]byte)
		json.Unmarshal(cacheStr, &userCache)
		return userCache
	}
	return nil
}

// login log
func SetLoginLog(loginLog LoginLog) error {
	if _, logErr := o.Insert(&loginLog); logErr != nil {
		fmt.Println("loginLog err, ", logErr)
		return logErr
	}
	// update sign info
	time := loginLog.LoginTime
	month := time.Month()
	day := time.Day()
	monthInd := utils.MonthMapping(month.String())
	signKey := loginLog.Name + "_" + strconv.Itoa(monthInd)
	redisCache.Command("SETBIT", signKey, day, 1)
	// monthSign := redisCache.Get(signKey)
	// daySign,_ := redisCache.Command("GETBIT", signKey, day)
	// fmt.Println("month sign info: ", monthSign, ", daySign: ", daySign)
	return nil
}
