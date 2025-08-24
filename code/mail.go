package code

var (
	MailDefaultValidDays int32 = 30 // 邮件默认有效天数
)

const (
	MailDefaultOperator = "system"
)

const (
	MailTypeSystem = 100 // 系统邮件
)

const (
	MailTemplateStatusInUse   = 0  // 使用中-默认
	MailTemplateStatusDeleted = 10 // 已删除
)

const (
	MailReadStatusUnread = 0  // 未读
	MailReadStatusRead   = 10 // 已读
)

const (
	MailAccessoryStatusUnDraw = 0  // 未领取
	MailAccessoryStatusDraw   = 10 // 已领取
)

const (
	MailVisibleStatusVisible = 0  // 可见
	MailVisibleStatusDeleted = 10 // 已删除
	MailVisibleStatusExpired = 20 // 已过期
)
