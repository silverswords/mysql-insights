package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	hang := make(chan struct{})
	// const count = 100000
	// var wg sync.WaitGroup
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3307)/masterSlaveDB?parseTime=true")
	if err != nil {
		log.Println(err)
	}
	// // defer db.Close()

	result := db.HasTable(&User{})
	log.Println("result:", result)
	if result != true {
		db.CreateTable(new(User))
	}

	// // insert
	// start := time.Now()
	// wg.Add(100)
	// for i := 0; i < 100; i++ {
	// 	user := User{Name: "aaa", Age: 20}
	// 	go func() {
	// 		db.Create(&user)
	// 		wg.Done()
	// 	}()
	// }

	// wg.Wait()

	// // query
	// user := User{Name: "aaa", Age: 20}
	// start = time.Now()
	// db.Find(&user, 50000)
	// log.Println("[query time]", time.Now().Sub(start).Seconds())
	// log.Println(db.Find(&user, count).Value)

	// // update
	// start = time.Now()
	// db.Table("users").Where("id = ?", 100000).Updates(map[string]interface{}{"name": "bbb", "age": 18})
	// log.Println("[update time]", time.Now().Sub(start).Seconds())

	// start = time.Now()
	// db.Table("users").Updates(map[string]interface{}{"name": "bbb", "age": 18})
	// log.Println("[update time]", time.Now().Sub(start).Seconds())

	// // delete
	// start = time.Now()
	// db.Where("id = ?", 100000).Delete(&User{})
	// log.Println("[delete time]", time.Now().Sub(start).Seconds())

	// start = time.Now()
	// db.Delete(&User{})
	// log.Println("[delete time]", time.Now().Sub(start).Seconds())

	go CreateAnimals(db)
	<-hang
}

type User struct {
	Id   int64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Name string
	// Age  int `gorm:"default:18"`
}

func CreateAnimals(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&User{Name: "Giraffe"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&User{Name: "Lion"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
