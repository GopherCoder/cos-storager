package database

import (
	"cos-storager/config"
	"fmt"

	"qiniupkg.com/x/log.v7"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var POSTGRES *gorm.DB

func Init() *gorm.DB {
	databaseConnectInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s password=%s",
		config.DATABASEHOST,
		config.DATABASEUSER,
		config.DATABASENAME,
		config.DATABASEPORT,
		config.DATABASESSLMODE,
		config.DATABASEPASSWORD)
	log.Println(databaseConnectInfo)
	connect, err := gorm.Open("postgres", databaseConnectInfo)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	connect.LogMode(true)
	POSTGRES = connect
	return POSTGRES

}
