package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

var mongoCollection *mgo.Collection

func main() {

	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	mongoCollection = session.DB("test").C("data")

	router := httprouter.New()
	router.GET("/user/:name/profile", profile)
	router.POST("/user/:name/profile", store)

	http.Handle("/", router)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func profile(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	var result *interface{}
	err := mongoCollection.Find(bson.M{"name": vars["name"]}).One(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		writeJson(w, result)
	}
}

func store(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	// This is a map with string keys and any type as value.
	// "When mgo marshals a struct, it lowercases the fields by default" so
	// a struct wouldn't work.
	var m map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	m["name"] = vars["name"]
	err = mongoCollection.Insert(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func writeJson(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
