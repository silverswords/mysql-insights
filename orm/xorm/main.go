package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func main() {
	const count = 100000
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3307)/masterSlaveDB?parseTime=true")
	if err != nil {
		log.Println(err, "err")
	}

	err = engine.Sync2(new(Userinfo))
	if err != nil {
		log.Println(err, "[CreateTable err]")
	}

	// insert
	user := User{Name: "aaa", Address: "bbb"}
	start := time.Now()
	for i := 0; i < count; i++ {
		user := User{Name: "aaa", Address: "bbb"}
		_, err := engine.Insert(user)
		log.Println("[current count]", i)
		log.Println("[current insert time]", time.Now().Sub(start).Seconds())
		if err != nil {
			log.Println(err, "[insert err]")
		}
	}

	// query
	start = time.Now()
	result, err := engine.Id(count).Get(&user)
	log.Println("[query time]", time.Now().Sub(start).Seconds())
	if err != nil {
		log.Println(err)
	}
	log.Println(result, "[query result]")

	// update
	start = time.Now()
	affected, err := engine.Update(&User{Name: "ccc", Address: "ddd"})
	log.Println("[update time]", time.Now().Sub(start).Seconds())
	if err != nil {
		log.Println(err)
	}
	log.Println(affected)

	// delete
	start = time.Now()
	_, err = engine.Exec("delete from user")
	log.Println("[delete time]", time.Now().Sub(start).Seconds())
	if err != nil {
		log.Println(err)
	}

	//transaction
	res, err := engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		log.Println("111")
		user1 := Userinfo{Username: "xixi", Departname: "dev", Alias: "lunny", Created: time.Now()}
		if _, err := session.Insert(&user1); err != nil {
			return nil, err
		}
		return nil, nil
	})
	log.Println("result", res)
}

type User struct {
	Id      int64
	Name    string `xorm:"VARCHAR(64) 'name'"`
	Address string `xorm:"VARCHAR(256) 'address'"`
}

type Userinfo struct {
	Id         int64
	Username   string
	Departname string
	Alias      string
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}
