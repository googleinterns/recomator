package automation

import ("regexp")

// TODO change to returning error
func extractFromURL (url, parameterName string) string {
	r, err := regexp.Compile("/" + parameterName + "/[a-zA-Z0-9]+/")
	
	if err != nil {
		return ""
	}

	partialResult := r.FindString(url)

	if len(partialResult) == 0 {
		return ""
	}
	
	stringResult := string(partialResult)
	result := stringResult[len(parameterName) + 2: len(stringResult) - 1]
	
	return result
}