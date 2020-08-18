package automation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetNumberOfSubtasks(t *testing.T) {
	var task Task
	for numSubtasks := 0; numSubtasks < 5; numSubtasks++ {
		task.SetNumberOfSubtasks(numSubtasks)

		task.mutex.Lock()
		assert.Equal(t, numSubtasks, len(task.subtasks), "wrong number of subtasks created")
		task.mutex.Unlock()
	}
}

func TestGetNextSubtask(t *testing.T) {
	var task Task
	numSubtasks := 10
	task.SetNumberOfSubtasks(numSubtasks)
	for subtaskIndex := 0; subtaskIndex < numSubtasks; subtaskIndex++ {
		subtask := task.GetNextSubtask()

		task.mutex.Lock()
		assert.Equal(t, &task.subtasks[subtaskIndex], subtask, "wrong next subtask")
		task.mutex.Unlock()

		task.IncrementDone()
	}
}

func TestGetProgress(t *testing.T) {
	var task Task
	numSubtasks := 10
	task.SetNumberOfSubtasks(numSubtasks)
	var progressFractions []float64
	for subtaskIndex := 0; subtaskIndex < numSubtasks; subtaskIndex++ {
		done, all := task.GetProgress()
		progressFractions = append(progressFractions, float64(done)/float64(all))

		task.IncrementDone()
	}
	task.SetAllDone()
	done, all := task.GetProgress()
	assert.True(t, done == all, "everything should be done already")

	progressFractions = append(progressFractions, float64(done)/float64(all))

	assert.IsIncreasing(t, progressFractions, "progress should increase over time")
}

func TestSubSubtasks(t *testing.T) {
	var task Task
	for numSubtasks := 2; numSubtasks <= 5; numSubtasks++ {
		task.SetNumberOfSubtasks(numSubtasks)

		subtask := task.GetNextSubtask()
		subtask.SetNumberOfSubtasks(numSubtasks)

		subtask.IncrementDone()

		done, all := subtask.GetProgress()
		fraction := float64(done) / float64(all)

		epsilon := 0.0001
		assert.InEpsilon(t, 1.0/float64(numSubtasks), fraction, epsilon, "incorrect progress for subtask")

		done, all = task.GetProgress()
		fraction = float64(done) / float64(all)
		assert.InEpsilon(t, 1.0/float64(numSubtasks)/float64(numSubtasks), fraction, epsilon, "incorrect progress for task")
	}
}
