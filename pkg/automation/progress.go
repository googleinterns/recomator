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

// Task is the structure that helps get the percentage of work done.
// GetProgress method is thread-safe.
// IncrementDone and GetNextSubtask is thread-safe,
// but GetProgress is calculated in assumption that all subtasks, except lower-level subtasks with no subtasks,
// will be done consequently.
type Task struct {
	subtasks        []Task
	subtasksStarted int
	subtasksDone    int
	taskDone        bool
	mutex           sync.Mutex
}

// SetNumberOfSubtasks is used to set number of subtasks in the task
func (t *Task) SetNumberOfSubtasks(num int) {
	t.mutex.Lock()
	t.subtasks = make([]Task, num)
	t.mutex.Unlock()
}

// IncrementDone increments subtasksDone
func (t *Task) IncrementDone() {
	t.mutex.Lock()
	t.subtasksDone++
	t.mutex.Unlock()
}

// GetNextSubtask returns the pointer to next not done subtask
func (t *Task) GetNextSubtask() *Task {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.subtasksStarted >= len(t.subtasks) {
		return nil
	}
	taskIndex := t.subtasksStarted
	t.subtasksStarted++
	return &t.subtasks[taskIndex]
}

// SetAllDone sets the task done
func (t *Task) SetAllDone() {
	t.mutex.Lock()
	t.taskDone = true
	t.mutex.Unlock()
}

// floatToFraction returns fraction approximation to float value x.
// Returned values are numerator and denominator of that fraction.
func floatToFraction(x float64) (int32, int32) {
	denominator := 1000000
	numerator := x * float64(denominator)
	return int32(numerator), int32(denominator)
}

// GetProgress returns the fraction of work done.
// Assumes that all subtasks, except the subtasks, that don't have their own subtasks are done subseequently
// Returned values are numerator and denominator of that fraction.
// If not all work is done numerator is garantueed to be less than denominator.
func (t *Task) GetProgress() (int32, int32) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.taskDone {
		return 1, 1
	}

	fractionOfWorkDone := 0.0
	if len(t.subtasks) != 0 {
		oneSubtaskWeight := 1.0 / float64(len(t.subtasks))
		fractionOfWorkDone += float64(t.subtasksDone) * oneSubtaskWeight

		if t.subtasksDone < len(t.subtasks) {
			unfinishedSubtask := &t.subtasks[t.subtasksDone]
			done, all := unfinishedSubtask.GetProgress()
			fractionOfWorkDone += float64(done) / float64(all) * oneSubtaskWeight
		}
	}

	done, all := floatToFraction(fractionOfWorkDone)
	if done >= all {
		done = all - 1
	}
	return done, all
}
