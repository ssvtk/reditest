package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"os"
)

type User struct {
	Name  string `redis:"name" json:"name"`
	Age   int    `redis:"age" json:"age"`
	Email string `redis:"email" json:"email"`
	Posts int    `redis:"posts" json:"posts"`
}

type Users struct {
	Users []User `json:"users"`
}

//String is a Stringer of User
func (u User) String() string {
	return fmt.Sprintf("Name: %s\nAge: %d\nEmail: %s\nPosts: %d\n", u.Name, u.Age, u.Email, u.Posts)
}

func main() {
	conn := Connect()
	//addUsers(conn)
	defer conn.Close()
	for i := 1; i < getUsersAmount(conn); i++ {
		fmt.Println(getUser(i, conn))

	}

}

//Connect perform the connection to redis
func Connect() redis.Conn {
	conn, err := redis.Dial("tcp", "localhost:6379")

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to REDIS is DONE")
	return conn
}

//getUser unmarshalls user from redis incoming data to User struct.
func getUser(userId int, conn redis.Conn) *User {
	var User User
	values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("user:%d", userId)))
	err = redis.ScanStruct(values, &User)
	if err != nil {
		log.Fatal(err)
	}
	return &User
}

//getUsersAmount returns an amount of users stored in redis, by managing the parallel list of added users "mylist"
func getUsersAmount(conn redis.Conn) int {
	num, err := redis.Int(conn.Do("LLEN", "mylist"))
	if err != nil {
		log.Fatal(err)
	}
	return num
}

//addUsers adds Users from json data parallel adding new users to mylist.
func addUsers(conn redis.Conn) {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	startPoint := getUsersAmount(conn)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var users Users
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		log.Fatal(err)
	}
	for i, user := range users.Users {
		_, err := conn.Do("HMSET", fmt.Sprintf("user:%d", startPoint+i), "name", user.Name, "age", user.Age, "email", user.Email, "posts", user.Posts)
		_, err = conn.Do("LPUSH", "mylist", fmt.Sprintf("user:%d", i))
		if err != nil {
			log.Print(err)
		}
	}
}
