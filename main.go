package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

type User struct {
	Name  string `redis:"name"`
	Age   int    `redis:"age"`
	Email string `redis:"email"`
	Posts int    `redis:"posts"`
}

func (u User) String() string {
	return fmt.Sprintf("Name: %s\nAge: %d\nEmail: %s\nPosts: %d\n", u.Name, u.Age, u.Email, u.Posts)
}

func main() {
	conn := Connect()
	defer conn.Close()
	for i := 1; i < getUserNumber(conn); i++ {
		fmt.Println(getUser(i, conn))

	}

}

func Connect() redis.Conn {
	conn, err := redis.Dial("tcp", "localhost:6379")

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to REDIS is DONE")
	return conn
}

func getUser(userId int, conn redis.Conn) *User {
	var User User
	values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("user:%d", userId)))
	err = redis.ScanStruct(values, &User)
	if err != nil {
		log.Fatal(err)
	}
	return &User
}

func getUserNumber(conn redis.Conn) int {
	num, err := redis.Int(conn.Do("LLEN", "mylist"))
	if err != nil {
		log.Fatal(err)
	}
	return num
}
