package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type request struct {
	Method string
	URL    string
}

type response struct {
	ID     int                 `json:"id"`
	Status string              `json:"status"`
	Header map[string][]string `json:"headers"`
	Body   int                 `json:"body"`
}

type pagination struct {
	Count int
	Page  int
}

var tasks = make(map[int]request)

func addTask(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var req request

	if err := json.Unmarshal(b, &req); err != nil {
		log.Fatal(err)
	}

	tasks[len(tasks)] = req

	log.Println("Add task", req.Method, req.URL)

	w.WriteHeader(http.StatusOK)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	// b, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer r.Body.Close()

	result, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(result)
}
