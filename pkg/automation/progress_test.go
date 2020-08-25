/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package automation

import (
	"sync/atomic"
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
	for numSubtasks := 2; numSubtasks <= 5; numSubtasks++ {
		var task Task
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

func TestThreadsTask(t *testing.T) {
	for numGoroutines := 1; numGoroutines < 20; numGoroutines++ {
		for numSubtasks := 0; numSubtasks < 5; numSubtasks++ {
			var task Task
			task.SetNumberOfSubtasks(numSubtasks)
			ch := make(chan bool)
			var subtasksDone int32
			for i := 0; i < numGoroutines; i++ {
				go func() {
					for task.GetNextSubtask() != nil {
						task.IncrementDone()
						atomic.AddInt32(&subtasksDone, 1)
					}
					ch <- true
				}()
			}
			for i := 0; i < numGoroutines; i++ {
				<-ch
			}
			task.SetAllDone()
			done, all := task.GetProgress()
			assert.True(t, done == all, "Task should be finished")
			assert.Equal(t, int32(numSubtasks), subtasksDone, "Wrong number of done subtasks")
		}
	}
}

func TestThreadsGetProgress(t *testing.T) {
	for numGoroutines := 1; numGoroutines < 10; numGoroutines++ {
		for numSubtasks := 0; numSubtasks < 10; numSubtasks++ {
			task := &Task{}
			task.SetNumberOfSubtasks(numSubtasks)
			for i := 0; i < numGoroutines; i++ {
				go func() {
					var progress []float64
					done, all := 0, 1
					for done < all {
						done, all := task.GetProgress()
						progress = append(progress, float64(done)/float64(all))
					}
					assert.IsNonDecreasing(t, progress, "Progress should not decrease")
				}()
			}
			for i := 0; i < numSubtasks; i++ {
				task.IncrementDone()
			}
		}

	}
}
