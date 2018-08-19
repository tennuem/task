package main

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

type task struct {
	Method string `json:"method"`
	URL    string `json:"url"`
}

type taskDelete struct {
	ID int `json:"id"`
}

var tasks = make(map[int]task)
var mu sync.Mutex

func addTask(c *gin.Context) {
	var t task
	if err := c.Bind(&t); err != nil {
		c.Status(400)
		return
	}

	ch <- t

	mu.Lock()
	tasks[len(tasks)] = t
	mu.Unlock()
	c.Status(200)
}

func getTask(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(200, tasks)
}

func deleteTask(c *gin.Context) {
	var t taskDelete
	if err := c.Bind(&t); err != nil {
		c.Status(400)
		return
	}

	log.Println(tasks[t.ID])
	if _, ok := tasks[t.ID]; !ok {
		c.Status(404)
		return
	}

	mu.Lock()
	delete(tasks, t.ID)
	mu.Unlock()
	c.Status(200)
}

func getResponses(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(200, responses)
}
