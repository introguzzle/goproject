package http

import (
	"encoding/json"
)

func MapString(
	jsonStr string,
	mapping map[string]string,
) (string, error) {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	newData := mapKeys(data, mapping)

	modifiedJson, err := json.Marshal(newData)
	if err != nil {
		return "", err
	}

	return string(modifiedJson), nil
}

func mapKeys(
	data map[string]interface{},
	mapping map[string]string,
) map[string]interface{} {
	newData := make(map[string]interface{})
	for key, value := range data {
		newKey, exists := mapping[key]
		if !exists {
			newKey = key
		}

		if nestedMap, ok := value.(map[string]interface{}); ok {
			value = mapKeys(nestedMap, mapping)
		}
		newData[newKey] = value
	}

	return newData
}
