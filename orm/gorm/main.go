package main

import (
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	const count = 100000
	var wg sync.WaitGroup
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3307)/masterSlaveDB?parseTime=true")
	if err != nil {
		log.Println(err)
	}
	// defer db.Close()

	result := db.HasTable(&User{})
	log.Println("result:", result)
	if result != true {
		db.CreateTable(new(User))
	}

	// db.NewRecord(user)
	// db.Create(&user)
	// start := time.Now()
	wg.Add(100)

	for i := 0; i < 100; i++ {
		user := User{Name: "aaa", Age: 20}
		go func() {
			db.Create(&user)
			// log.Println("[current count]", i)
			wg.Done()
			// log.Println("[current insert time]", time.Now().Sub(start).Seconds())
		}()
	}

	wg.Wait()

	user := User{Name: "aaa", Age: 20}
	start := time.Now()
	db.Find(&user, 50000)
	log.Println("[query time]", time.Now().Sub(start).Seconds())
	log.Println(db.Find(&user, count).Value)

	// start := time.Now()
	// db.Table("users").Where("id IN (?)", 100000).Updates(map[string]interface{}{"name": "bbb", "age": 18})
	// log.Println("[update time]", time.Now().Sub(start).Seconds())

	// start := time.Now()
	// db.Table("users").Updates(map[string]interface{}{"name": "bbb", "age": 18})
	// log.Println("[update time]", time.Now().Sub(start).Seconds())

	// start := time.Now()
	// db.Where("id = ?", 100000).Delete(&User{})
	// log.Println("[delete time]", time.Now().Sub(start).Seconds())

	// start := time.Now()
	// db.Delete(&User{})
	// log.Println("[delete time]", time.Now().Sub(start).Seconds())
}

type User struct {
	Id   int64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Name string
	Age  int `gorm:"default:18"`
}
