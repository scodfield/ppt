package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
	"ppt/dao"
)

type User struct {
	ID        string `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uint64 `gorm:"not null;uniqueIndex;column:user_id;comment:玩家UserID" json:"user_id"`
	Username  string `gorm:"not null;size:255;column:user_name;comment:玩家名" json:"user_name"`
	Password  string `gorm:"not null;size:255;column:password;comment:密码" json:"password"`
	Email     string `gorm:"not null;uniqueIndex;size:255;column:email;comment:注册邮箱" json:"email"`
	BrandID   int32  `gorm:"not null;column:brand_id;comment:品牌" json:"brand_id"`
	Channel   string `gorm:"not null;column:channel;comment:渠道" json:"channel"`
	Lang      string `gorm:"not null;column:lang;comment:语言包" json:"lang"`
	CreatedAt int64  `gorm:"autoCreateTime:milli;column:create_at;comment:创建时间" json:"created_at"`
	UpdateAt  int64  `gorm:"autoUpdateTime:milli;column:update_at;comment:最后更新" json:"update_at"`
}

func (User) TableName() string {
	return "user"
}

func MigrateUserModel(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

// CreateUser 创建用户
func CreateUser(userID uint64, userName, password, email string) (*User, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       uuidV7.String(),
		UserID:   userID,
		Username: userName,
		Password: password,
		Email:    email,
	}, nil
}

// InsertUsersByPgxPool 批量插入用户数据
func InsertUsersByPgxPool(pgxPool *pgxpool.Pool, users []*User) error {
	userModel := User{}
	batchSize := 1000
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}
		usersBatch := users[i:end]
		batch := &pgx.Batch{}
		for _, user := range usersBatch {
			batch.Queue(fmt.Sprintf(`insert into %s (id, user_id, user_name, password, emal) values ($1,$2,$3,$4,$5)`, userModel.TableName()), user.ID, user.UserID, user.Username, user.Password, user.Email)
		}
		err := pgxPool.SendBatch(dao.Ctx, batch).Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// InsertUsersByBatch 批量创建用户
func InsertUsersByBatch(pgDB *gorm.DB, users []*User) error {
	return pgDB.CreateInBatches(&users, 200).Error
}

// GetUserByID 获取用户
func GetUserByID(pgDB *gorm.DB, userID uint64) (*User, error) {
	var user User
	if err := pgDB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserSpecifyFieldsByID(pgDB *gorm.DB, userID uint64, fields []string) (*User, error) {
	var user User
	if err := pgDB.Where("user_id = ?", userID).Select(fields).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
