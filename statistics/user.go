package statistics

// UseBase 用户基础信息
type UseBase struct {
	UserID  uint64 `json:"user_id"`
	BrandID int32  `json:"brand_id"`
	Channel string `json:"channel"`
	VipLv   int32  `json:"vip_lv"`
}

// UserLoginStatistics 用户登录统计
type UserLoginStatistics struct {
	Base               *UseBase `json:"base"`
	LastLoginTime      int64    `json:"last_login_time"`      // 上次登录时间
	ConsecutiveDays    int32    `json:"consecutive_days"`     // 最近连续登录天数
	MaxConsecutiveDays int32    `json:"max_consecutive_days"` // 最长连续登录天数
}

// UserOnlineStatistics 用户在线统计
type UserOnlineStatistics struct {
	Base               *UseBase `json:"base"`
	LastLogoutTime     int64    `json:"last_logout_time"`     // 上次离线时间
	ConsecutiveDays    int32    `json:"consecutive_days"`     // 最近连续在线天数
	MaxConsecutiveDays int32    `json:"max_consecutive_days"` // 最长连续在线天数
}

// UserRechargeStatistics 用户充值统计
type UserRechargeStatistics struct {
	Base                 *UseBase `json:"base"`
	TodayRecharge        int64    `json:"today_recharge"`          // 今日累充
	YesterdayRecharge    int64    `json:"yesterday_recharge"`      // 昨日累充
	LastThreeDayRecharge int64    `json:"last_three_day_recharge"` // 过去3日累充
	LastSevenDayRecharge int64    `json:"last_seven_day_recharge"` // 过去7日累充
	TotalRecharge        int64    `json:"total_recharge"`          // 历史累充
	FirstRecharge        int64    `json:"first_recharge"`          // 首充金额
	FirstRechargeTime    int64    `json:"first_recharge_time"`     // 首充时间
	LastRechargeTime     int64    `json:"last_recharge_time"`      // 上次充值时间
}

// UserVipLvStatistics 用户VIP升级统计
type UserVipLvStatistics struct {
	Base       *UseBase              `json:"base"`
	ChangeList []UserVipLvChangeItem `json:"change_list"`
}

type UserVipLvChangeItem struct {
	VipLv  int32 `json:"vip_lv"`
	UpTime int64 `json:"up_time"` // 升级时间
}
