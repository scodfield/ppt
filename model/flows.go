package model

import "gorm.io/gorm"

// UserEventFlows 用户事件流水
type UserEventFlows struct {
	UserID       uint64 `gorm:"user_id" json:"user_id"`
	UserName     string `gorm:"user_name" json:"user_name"`
	ServerID     int32  `gorm:"server_id;comment:服务器ID" json:"server_id"`
	IP           string `gorm:"ip;comment:登录IP" json:"ip"`
	Level        int32  `gorm:"level" json:"level"`
	VipLevel     int32  `gorm:"vip_level" json:"vip_level"`
	Balance      int64  `gorm:"balance;comment:充值余额" json:"balance"`
	TotalBalance int64  `gorm:"total_balance;comment:充值总额" json:"total_balance"`
	EventType    int32  `gorm:"event_type;comment:事件类型" json:"event_type"`
	ExtraArgs    string `gorm:"extra_args" json:"extra_args"`
	EventTime    int64  `gorm:"event_time;type:Int64;autoCreateTime:milli" json:"event_time"`
}

func (UserEventFlows) TableName() string {
	return "user_event_flows"
}

func MigrateUserEventFlows(sqlSession *gorm.DB) error {
	return sqlSession.Set("gorm:table_options", "ENGINE=ReplicatedReplacingMergeTree() PARTITION BY (event_time) ORDER BY (event_time, user_id, event_type) SETTINGS index_granularity = 8192").AutoMigrate(&UserEventFlows{})
}

// UserEventFlowsDistributed 用户事件流水-分布式表
type UserEventFlowsDistributed struct {
	UserID       uint64 `gorm:"user_id" json:"user_id"`
	UserName     string `gorm:"user_name" json:"user_name"`
	ServerID     int32  `gorm:"server_id;comment:服务器ID" json:"server_id"`
	IP           string `gorm:"ip;comment:登录IP" json:"ip"`
	Level        int32  `gorm:"level" json:"level"`
	VipLevel     int32  `gorm:"vip_level" json:"vip_level"`
	Balance      int64  `gorm:"balance;comment:充值余额" json:"balance"`
	TotalBalance int64  `gorm:"total_balance;comment:充值总额" json:"total_balance"`
	EventType    int32  `gorm:"event_type;comment:事件类型" json:"event_type"`
	ExtraArgs    string `gorm:"extra_args" json:"extra_args"`
	EventTime    int64  `gorm:"event_time;type:Int64;autoCreateTime:milli" json:"event_time"`
}

func MigrateUserEventFlowsDistributed(sqlSession *gorm.DB) error {
	return sqlSession.Set("gorm:table_options", "ENGINE=Distributed(default, user_flows, user_event_flows, rand())").AutoMigrate(&UserEventFlowsDistributed{})
}
