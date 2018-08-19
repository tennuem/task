package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	log.Println("Starting server")

	http.HandleFunc("/add", addTask)
	http.HandleFunc("/get", getTask)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func fetchTask(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var req request
	var resp response

	if err := json.Unmarshal(b, &req); err != nil {
		log.Fatal(err)
	}

	if err := fetchLink(req, &resp); err != nil {
		log.Fatal(err)
	}

	result, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(result)
}

func fetchLink(r request, w *response) error {

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	w.Status = resp.Status
	w.Body = len(b)
	w.Header = resp.Header

	return nil
}
