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
	"sync"
)

// Task is the structure that helps get the percentage of work done
type Task struct {
	subtasks     []Task
	subtasksDone int
	taskDone     bool
	mutex        sync.Mutex
}

// AddSubtasks is used to add subtasks to the task
func (t *Task) AddSubtasks(num int) {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.subtasks = make([]Task, num)
	t.mutex.Unlock()
}

// IncrementDone increments subtasksDone
func (t *Task) IncrementDone() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.subtasksDone++
	t.mutex.Unlock()
}

// GetNextSubtask returns the pointer to next not done subtask
func (t *Task) GetNextSubtask() *Task {
	if t == nil {
		return nil
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.subtasksDone >= len(t.subtasks) {
		return nil
	}
	return &t.subtasks[t.subtasksDone]
}

// SetAllDone sets the task done
func (t *Task) SetAllDone() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.taskDone = true
	t.mutex.Unlock()
}

func floatToInts(x float64) (int32, int32) {
	all := 1000000
	return int32(x * float64(all)), int32(all)
}

// GetProgress returns the fraction of work done
func (t *Task) GetProgress() (int32, int32) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.taskDone {
		return 1, 1
	}
	answer := 0.0
	if len(t.subtasks) != 0 {
		answer += float64(t.subtasksDone) / float64(len(t.subtasks))
		if t.subtasksDone < len(t.subtasks) {
			done, all := t.subtasks[t.subtasksDone].GetProgress()
			answer += float64(done) / float64(all) / float64(len(t.subtasks))
		}
	}
	done, all := floatToInts(answer)
	if done >= all {
		done = all - 1
	}
	return done, all
}
