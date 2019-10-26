package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	Name string
	Age  int `gorm:"default:18"`
}

func main() {
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/masterSlaveDB?parseTime=true")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	db.CreateTable(new(User))

	user := User{Name: "lala", Age: 18}
	db.NewRecord(user)
	db.Create(&user)
	db.NewRecord(user)

	db.First(&user, 0)

	// db.Delete(&user, 0)
}
