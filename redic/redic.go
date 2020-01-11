// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// * * * * * * * * * * * * * * * * * * * * * * * * * *
// This package is based on github.com/astaxie/beego/cache, all
// the licenses are follow the prior permissions and limitations. 
// Since we only focus on the redis for the moment, so we just
// extend the interfaces to make it easier to use.
// If there is any question, please contact me by 953836918@qq.com
// More docs http://beego.me/docs/module/cache.md.
// * * * * * * * * * * * * * * * * * * * * * * * * * *


package redic

import (
	"fmt"
	"time"
	"encoding/json"
	"errors"
	"strconv"
	"github.com/gomodule/redigo/redis"
	"strings"
)

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey = "redic"
)


// Cache interface contains all behaviors for cache adapter.
// usage:
//	c,err := cache.NewCache("{....}")
//	c.Set("key",value)
//	c.SetEX("key",value, 3600 * time.Second)
//	v := c.Get("key")
//
//	c.Incr("counter")  // now is 1
//	c.Incr("counter")  // now is 2
//	count := c.Get("counter").(int)
type Cache interface {
	// get cached value by key.
	Get(key string) interface{}
	// GetMulti is a batch version of Get.
	GetMulti(keys []string) []interface{}
	// set cached value with key.
	Set(key string, val interface{}) error
	// set cached value with key and expire time.
	SetEx(key string, val interface{}, timeout time.Duration) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// custom command
	Command(commandName string, args ...interface{}) (interface{}, error)
	// clear all cache.
	ClearAll() error
	// start gc routine based on config string settings.
	StartAndGC(config string) error
}


// RedisC is Redis cache adapter.
type RedisC struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
	password string
	maxIdle  int
	idleTimeout int
}

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() Cache {
	return &RedisC{key: DefaultKey}
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *RedisC) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *RedisC) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

// Get cache from redis.
func (rc *RedisC) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// GetMulti get cache from redis.
func (rc *RedisC) GetMulti(keys []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Set set cache to redis.
func (rc *RedisC) Set(key string, val interface{}) error {
	_, err := rc.do("SET", key, val)
	return err
}

// SetEx set cache to redis with expire time.
func (rc *RedisC) SetEx(key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Delete delete cache in redis.
func (rc *RedisC) Delete(key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *RedisC) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Command execute custom commands
func (rc *RedisC) Command(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// Incr increase counter in redis.
func (rc *RedisC) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (rc *RedisC) Decr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *RedisC) ClearAll() error {
	c := rc.p.Get()
	defer c.Close()
	cachedKeys, err := redis.Strings(c.Do("KEYS", rc.key+":*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *RedisC) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}

	// Format redis://<password>@<host>:<port>
	cf["conn"] = strings.Replace(cf["conn"], "redis://", "", 1)
	if i := strings.Index(cf["conn"], "@"); i > -1 {
		cf["password"] = cf["conn"][0:i]
		cf["conn"] = cf["conn"][i+1:]
	}

	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	if _, ok := cf["maxIdle"]; !ok {
		cf["maxIdle"] = "3"
	}
	if _, ok := cf["idleTimeout"]; !ok {
		cf["idleTimeout"] = "180"
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]
	rc.maxIdle, _ = strconv.Atoi(cf["maxIdle"])
	rc.idleTimeout, _ = strconv.Atoi(cf["idleTimeout"])

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *RedisC) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     rc.maxIdle,
		IdleTimeout: time.Duration(rc.idleTimeout) * time.Second,
		Dial:        dialFunc,
	}
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(config string) (adapter Cache, err error) {
	adapter = NewRedisCache()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}

