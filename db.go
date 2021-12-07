package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type dbconfig struct {
	Host     string
	Port     string
	UserName string
	Password string
	Database string
}

func initDB(cfg *dbconfig) (*gorm.DB, error) {
	//url := "topnews:topnews2016@tcp(master.mysql.hurybuy.xbnet.com:3306)/dropship?parseTime=True&loc=Local&multiStatements=true&charset=utf8mb4"
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local&multiStatements=true&charset=utf8mb4", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := gorm.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}
	if err := db.DB().Ping(); err != nil {
		return db, err
	}
	return db, nil
}
