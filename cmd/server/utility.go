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

package main

import "math/rand"

const characters = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Returns a random string of the given length
// generated using the given generator.
// This function will only be thread safe, if the given
// generator is thread safe.
func randomString(sequenceLen int, generator *rand.Rand) string {
	result := make([]byte, sequenceLen)
	for i := range result {
		result[i] = characters[generator.Intn(len(characters))]
	}

	return string(result)
}
