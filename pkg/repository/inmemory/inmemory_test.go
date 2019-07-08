package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tennuem/task/pkg/task"
)

func TestInmemory(t *testing.T) {
	var (
		testData1 = task.Task{
			ID:          1,
			Title:       "Task title",
			Description: "Task description",
		}
		testData2 = task.Task{
			ID:          1,
			Title:       "Task title updated",
			Description: "Task description updated",
		}
		repo = NewRepository()
	)

	err := repo.Create(context.Background(), &testData1)
	tasks := repo.GetList(context.Background())
	taskResp, err := repo.GetByID(context.Background(), testData1.ID)
	err = repo.Update(context.Background(), &testData2)
	err = repo.Delete(context.Background(), testData2.ID)

	assert.NoError(t, err)
	assert.Equal(t, tasks[0].ID, testData1.ID)
	assert.Equal(t, tasks[0].Title, testData1.Title)
	assert.Equal(t, tasks[0].Description, testData1.Description)
	assert.Equal(t, taskResp.ID, testData1.ID)
	assert.Equal(t, taskResp.Title, testData1.Title)
	assert.Equal(t, taskResp.Description, testData1.Description)
}
