package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type task struct {
	Method string `json:"method" binding:"required"`
	URL    string `json:"url" binding:"required"`
}

type taskDelete struct {
	ID int `json:"id" binding:"required"`
}

var tasks = make(map[int]task)
var mu sync.Mutex

func addTask(c *gin.Context) {
	var t task
	if err := c.Bind(&t); err != nil {
		c.Status(400)
		return
	}
	//a
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

	mu.Lock()
	if _, ok := tasks[t.ID]; !ok {
		c.Status(404)
		return
	}

	delete(tasks, t.ID)
	mu.Unlock()
	c.Status(200)
}

func getResponses(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(200, responses)
}
