package dao

import (
	"github.com/dobyte/gorm-dao-generator/example/dao/internal"
	"gorm.io/gorm"
)

type (
	MailColumns = internal.MailColumns
	MailOrderBy = internal.MailOrderBy
)

type Mail struct {
	*internal.Mail
}

func NewMail(db *gorm.DB) *Mail {
	return &Mail{Mail: internal.NewMail(db)}
}
