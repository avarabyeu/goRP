package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"testing"
	"fmt"
	"log"
)

func TTestLog(t *testing.T) {
	session, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	var logRepo = NewLogRepo(session, "reportportal")
	logs, _ := logRepo.FindAll()
	for l := range logs {
		log.Println(fmt.Sprint(l))

	}

	// Optional. Switch the session to a monotonic behavior.
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
