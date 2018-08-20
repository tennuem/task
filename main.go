package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Status string              `json:"status"`
	Header map[string][]string `json:"headers"`
	Body   int                 `json:"body"`
}

var ch = make(chan task)
var responses []response

func main() {
	go func() {
		tokens := make(chan struct{}, 5)
		for v := range ch {
			tokens <- struct{}{}
			go fetchLink(v)
			<-tokens
		}
		defer close(tokens)
	}()

	r := gin.Default()
	routes(r)
	r.Run()
}

func fetchLink(t task) {
	client := &http.Client{}
	req, err := http.NewRequest(t.Method, t.URL, nil)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	mu.Lock()
	responses = append(responses, response{
		Status: resp.Status,
		Body:   len(b),
		Header: resp.Header,
	})
	mu.Unlock()
}
