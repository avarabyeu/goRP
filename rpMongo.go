package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	//"log"
	"fmt"
	"github.com/avarabyeu/goRP/rplog"

)

func main1() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var logs []rplog.Log
	session.Clone().DB("server").C("log").Find(nil).All(&logs)
	fmt.Println("Results All: ", logs)


	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	//
	//var all []bson.M
	//
	//session.DB("local").C("oplog.$main").Find(bson.M{}).All(&all)
	//
	//for key, value := range all {
	//	fmt.Println(key, value)
	//}
	//log.Println(fmt.Sprint(all))

	//results := make(chan map[string]interface{})
	//go watch(session, results)
	//for r := range results {
	//	log.Println(fmt.Sprint(r))
	//}

}

func watch(session *mgo.Session, results chan map[string]interface{}) {
	iter := session.DB("local").C("oplog.$main").Find(bson.M{}).Tail(1 * time.Minute)
	for {
		obs := map[string]interface{}{}
		for iter.Next(&obs) {
			results <- obs
		}
	}
}
