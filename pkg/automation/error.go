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

import "errors"

// Error returned when some type does not
// fit its asserted type
func typeError() error {
	return errors.New("Type error")
}

// Error returned when some parameter of the
// operation makes it impossible to perform
func operationUnsupportedError() error {
	return errors.New("The operation is not supported")
}
