package main

import (
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
)


func GetConnection()(db *gorm.DB,err error) {
  dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
  db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if err != nil {
    return nil,err
  }
  return db,nil

}
