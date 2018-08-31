package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

var session *mgo.Session

func init() {
	session, _ = mgo.Dial("172.16.14.52:27017")
}

type User struct {
	Name string
	Age  int
	Msgs []UserMessage
}

type UserMessage struct {
	Id         int
	Msg        string
	CreateTime time.Time
}

func Test_Mongo_Insert(t *testing.T) {
	defer session.Close()
	collection := session.DB("test").C("demo")
	collection.Insert(&User{Name: "Admin", Age: 20, Msgs: []UserMessage{
		{Id: 1, Msg: "一个例子", CreateTime: time.Now()},
		{Id: 2, Msg: "第二个例子", CreateTime: time.Now()},
		{Id: 3, Msg: "第三个例子", CreateTime: time.Now()},
	}})
}

func Test_Mongo_Update(t *testing.T) {
	defer session.Close()
	collection := session.DB("test").C("users")
	if e := collection.Update(bson.M{"msgs.id": 2}, bson.M{"$set": bson.M{"msgs.$.msg": "修改成功"}}); e != nil {
		t.Fatal(e.Error())
	}
}

func Test_Mongo_Select(t *testing.T) {
	defer session.Close()
	collection := session.DB("test").C("users")
	query := collection.Find(bson.M{"msgs.id": 2})
	//if n, err := query.Count(); err != nil {
	//	t.Log(err)
	//} else {
	//	t.Log(n)
	//}

	iter := query.Iter()
	defer iter.Close()
	user := User{}
	for iter.Next(&user) {
		fmt.Printf("%v", user)
	}
}
