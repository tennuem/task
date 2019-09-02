package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tennuem/task/pkg/task"
)

func TestCreateGetTask(t *testing.T) {
	testData := []*task.Task{
		{
			Title:       "Task title",
			Description: "Task description",
		},
		{
			Title:       "Task title updated",
			Description: "Task description updated",
		},
	}
	repo := NewRepository()

	for _, v := range testData {
		createResp, err := repo.Create(context.Background(), v)
		getResp, err := repo.GetByID(context.Background(), createResp.ID)

		assert.NoError(t, err)
		assert.NotNil(t, getResp.ID)
		assert.Equal(t, v.Title, getResp.Title)
		assert.Equal(t, v.Description, getResp.Description)
	}
}

func TestGetTasksList(t *testing.T) {
	testData := []*task.Task{
		{
			Title:       "Task title",
			Description: "Task description",
		},
		{
			Title:       "Task title updated",
			Description: "Task description updated",
		},
	}
	repo := NewRepository()

	for i := 0; i < len(testData); i++ {
		v := testData[i]

		_, err := repo.Create(context.Background(), v)
		resp := repo.GetList(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, v.Title, resp[i].Title)
		assert.Equal(t, v.Description, resp[i].Description)
	}
}

func TestUpdateTask(t *testing.T) {
	testData := []*task.Task{
		{
			Title:       "Task title",
			Description: "Task description",
		},
		{
			Title:       "Task title updated",
			Description: "Task description updated",
		},
	}
	repo := NewRepository()

	for _, v := range testData {
		_, err := repo.Create(context.Background(), v)

		updatedTitle := "Updated title"
		updatedDescription := "Updated description"
		v.Title = updatedTitle
		v.Description = updatedDescription
		resp, err := repo.Update(context.Background(), v)

		assert.NoError(t, err)
		assert.NotNil(t, resp.ID)
		assert.Equal(t, updatedTitle, resp.Title)
		assert.Equal(t, updatedDescription, resp.Description)
	}
}

func TestDeleteTask(t *testing.T) {
	testData := []*task.Task{
		{
			Title:       "Task title",
			Description: "Task description",
		},
		{
			Title:       "Task title updated",
			Description: "Task description updated",
		},
	}
	repo := NewRepository()

	for _, v := range testData {
		createResp, err := repo.Create(context.Background(), v)

		err = repo.Delete(context.Background(), createResp.ID)

		_, err = repo.GetByID(context.Background(), createResp.ID)

		assert.EqualError(t, err, ErrTaskNotExist.Error())
	}
}
