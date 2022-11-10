package global

import (
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
)

var (
	// DATABASE 连接
	DB *mgo.Database

	// CONNECT  连接
	Session *mgo.Session
)

// NewDB
func NewDB() {
	var err error
	Session, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     []string{fmt.Sprintf("%s:%d", Config().DataBase.DBHost, Config().DataBase.DBPort)},
		Direct:    false,
		Timeout:   time.Second * 3,
		Database:  Config().DataBase.DBName,
		Username:  Config().DataBase.DBUser,
		Password:  Config().DataBase.DBPass,
		PoolLimit: 4096,
	})
	if err != nil {
		log.Fatalln(fmt.Sprintf("open db fail:%s", err.Error()))
	}
	if err := Session.Ping(); err != nil {
		log.Fatalln(fmt.Sprintf("ping db fail:%s", err.Error()))
	}
	DB = Session.DB(Config().DataBase.DBName)
}
