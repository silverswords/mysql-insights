package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Id   int64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Name string
	Age  int `gorm:"default:18"`
}

func main() {
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/masterSlaveDB?parseTime=true")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	result := db.HasTable(&User{})
	log.Println("result:", result)
	if result != true {
		db.CreateTable(new(User))
	}

	// db.NewRecord(user)
	// db.Create(&user)
	// start := time.Now()
	// for i := 0; i < 10000; i++ {
	// user := User{Name: "aaa", Age: 20}
	// 	db.NewRecord(user)
	// 	db.Create(&user)
	// 	log.Println("[current count]", i)
	// 	log.Println("[current insert time]", time.Now().Sub(start).Seconds())
	// }

	// start := time.Now()
	// db.Find(&user, 10000)
	// log.Println("[query time]", time.Now().Sub(start).Seconds())
	// log.Println(db.Find(&user, 10000).Value)

	// start := time.Now()
	// db.Table("users").Where("id IN (?)", 9999).Updates(map[string]interface{}{"name": "bbb", "age": 18})
	// log.Println("[update time]", time.Now().Sub(start).Seconds())

	start := time.Now()
	db.Where("id = ?", 10000).Delete(&User{})
	log.Println("[delete time]", time.Now().Sub(start).Seconds())
}
