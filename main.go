package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type PolicyDocument struct {
	PolicyName     string          `json:"PolicyName"`
	PolicyDocument json.RawMessage `json:"PolicyDocument"`
}

func verifyJSON(filePath string) (bool, error) {

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return false, errors.New("input data is not valid JSON")
	}

	var policy PolicyDocument
	if err := json.Unmarshal(jsonData, &policy); err != nil {
		return false, errors.New("input JSON does not have the required fields")
	}

	var document map[string]interface{}
	if err := json.Unmarshal(policy.PolicyDocument, &document); err != nil {
		return false, errors.New("unable to unmarshal PolicyDocument field")
	}

	statements, ok := document["Statement"].([]interface{})
	if !ok {
		return false, errors.New("policydocument does not contain 'statement' field")
	}

	for _, statement := range statements {
		statementMap, ok := statement.(map[string]interface{})
		if !ok {
			return false, errors.New("statement is not a valid JSON object")
		}

		if resource, ok := statementMap["Resource"].(string); ok {
			if resource == "*" {
				return false, nil
			}
		}
	}

	return true, nil
}

func main() {
	filePath := "json_with_asterisk.json"
	isValid, err := verifyJSON(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(isValid)
}