package main

import (
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
)


func GetConnection() {
  dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // Migrate the schema
  db.AutoMigrate(&Product{})

  // Create
  db.Create(&Product{Code: "D43", Price: 153})

  // Read
  var product Product
  db.First(&product, 2) // find product with integer primary key
  db.First(&product, "code = ?", "D43") // find product with code D42

  // Update - update product's price to 200
  db.Model(&product).Update("Price", 303)

  // Delete - delete product
  db.Delete(&product, 1)
  

}
