package main

import (
	"context"
	"log"
	"time"

	"github.com/dobyte/gorm-dao-generator/example/dao"
	"github.com/dobyte/gorm-dao-generator/example/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/game?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}))
	if err != nil {
		log.Fatalf("connect mysql server failed: %v", err)
	}

	mailDao := dao.NewMail(db)
	baseCtx := context.Background()

	_, err = mailDao.Insert(baseCtx, &model.Mail{
		Title:    "gorm-dao-generator introduction",
		Content:  "The gorm-dao-generator is a tool for automatically generating Mysql Data Access Object.",
		Sender:   1,
		Receiver: 2,
		Status:   1,
		SendTime: time.Now(),
	})
	if err != nil {
		log.Fatalf("failed to insert into mongo database: %v", err)
	}

	mail, err := mailDao.FindOne(baseCtx, func(cols *dao.MailColumns) interface{} {
		return map[string]interface{}{
			cols.Receiver: 2,
		}
	})
	if err != nil {
		log.Fatalf("failed to find a row of data from mongo database: %v", err)
	}

	log.Printf("%+v", mail)
}
