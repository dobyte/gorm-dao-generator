package dao

import (
	"github.com/dobyte/gorm-dao-generator/example/dao/internal"
	"gorm.io/gorm"
)

type (
	UserColumns = internal.UserColumns
	UserOrderBy = internal.UserOrderBy
	UserFilterFunc = internal.UserFilterFunc
	UserUpdateFunc = internal.UserUpdateFunc
	UserColumnFunc = internal.UserColumnFunc
	UserOrderFunc = internal.UserOrderFunc
)

type User struct {
	*internal.User
}

func NewUser(db *gorm.DB) *User {
	return &User{User: internal.NewUser(db)}
}
