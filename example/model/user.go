package model

import (
	"time"
)

type Gender int

const (
	GenderUnknown Gender = iota // 未知
	GenderMale                  // 男性
	GenderFemale                // 女性
)

// Type 用户类型
type Type int

const (
	TypeRobot   Type = 0 // 机器人用户
	TypeGuest   Type = 1 // 游客用户
	TypeGeneral Type = 2 // 普通用户
	TypeSystem  Type = 3 // 系统用户
)

// Status 用户状态
type Status int

const (
	StatusNormal    Status = iota // 正常
	StatusForbidden               // 封禁
)

//go:generate gorm-dao-generator -model-dir=. -model-names=User:user -dao-dir=../dao/
type User struct {
	ID             int64          `gorm:"column:id"`
	UID            int32          `gorm:"column:uid"`                            // 用户ID
	Account        string         `gorm:"column:account"`                        // 用户账号
	Password       string         `gorm:"column:password"`                       // 用户密码
	Salt           string         `gorm:"column:salt"`                           // 密码
	Mobile         string         `gorm:"column:mobile"`                         // 用户手机
	Email          string         `gorm:"column:email"`                          // 用户邮箱
	Nickname       string         `gorm:"column:nickname"`                       // 用户昵称
	Signature      string         `gorm:"column:signature"`                      // 用户签名
	Gender         Gender         `gorm:"column:gender"`                         // 用户性别
	Level          int            `gorm:"column:level"`                          // 用户等级
	Experience     int            `gorm:"column:experience"`                     // 用户经验
	Coin           int            `gorm:"column:coin"`                           // 用户金币
	Type           Type           `gorm:"column:type"`                           // 用户类型
	Status         Status         `gorm:"column:status"`                         // 用户状态
	DeviceID       string         `gorm:"column:device_id"`                      // 设备ID
	ThirdPlatforms ThirdPlatforms `gorm:"column:third_platforms"`                // 第三方平台
	RegisterIP     string         `gorm:"column:register_ip"`                    // 注册IP
	RegisterTime   time.Time      `gorm:"column:register_time" gen:"autoFill"`   // 注册时间
	LastLoginIP    string         `gorm:"column:last_login_ip"`                  // 最近登录IP
	LastLoginTime  time.Time      `gorm:"column:last_login_time" gen:"autoFill"` // 最近登录时间
}

// ThirdPlatforms 第三方平台
type ThirdPlatforms struct {
	Wechat   string `bson:"wechat"`   // 微信登录openid
	Google   string `bson:"google"`   // 谷歌登录userid
	Facebook string `bson:"facebook"` // 脸书登录userid
}
