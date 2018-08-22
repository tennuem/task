package main

import (
	"sync"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type task struct {
	Method string `json:"method" binding:"required"`
	URL    string `json:"url" binding:"required"`
}

type taskDelete struct {
	ID string `json:"id" binding:"required"`
}

var tasks = make(map[uuid.UUID]task)
var mu sync.Mutex

func addTask(c *gin.Context) {
	var t task
	if err := c.Bind(&t); err != nil {
		c.Status(400)
		return
	}

	ch <- t

	id, err := uuid.NewV4()
	if err != nil {
		c.Status(500)
		return
	}

	mu.Lock()
	tasks[id] = t
	mu.Unlock()
	c.Status(201)
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

	id, err := uuid.FromString(t.ID)
	if err != nil {
		c.Status(500)
	}

	mu.Lock()
	if _, ok := tasks[id]; !ok {
		c.Status(404)
		return
	}

	delete(tasks, id)
	mu.Unlock()
	c.Status(200)
}

func getResponses(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(200, responses)
}
