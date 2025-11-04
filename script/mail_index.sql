-- 创建用户邮件过期时间索引
CREATE index idx_user_mail_expired_time on user_mail(expire_time);