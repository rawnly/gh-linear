package utils

import (
	"encoding/json"
	"fmt"
)

func NewError(message string) error {
	return fmt.Errorf("❌ %s", message)
}

func PrettyJSON(obj interface{}) string {
	prettyJSON, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err.Error())
	}

	return string(prettyJSON)
}
