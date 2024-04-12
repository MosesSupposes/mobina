package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Input struct {
	Data map[string]interface{} `json:""`
}

type Output struct {
	Data map[string]interface{} `json:""`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path to your json file: ")
	scanner.Scan()
	filePath := scanner.Text()
	output, err := transformToJSON(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write output
	encoder := json.NewEncoder(os.Stdout)
	err = encoder.Encode(output.Data)
	if err != nil {
		fmt.Println(err)
	}
}
func transformToJSON(filepath string) (Output, error) {
	// Read input from file
	file, err := os.ReadFile(filepath)

	if err != nil {
		return Output{}, err
	}

	var input map[string]interface{}
	err = json.Unmarshal(file, &input)
	if err != nil {
		return Output{}, err
	}

	// Transform
	output := Output{Data: make(map[string]interface{})}
	for key, value := range input {
		key = strings.TrimSpace(key)

		// Omit fields with empty keys
		if key == "" {
			continue
		}

		if valueMap, ok := value.(map[string]interface{}); ok {
			for dataType, dataValue := range valueMap {
				// Use type assertion to access underlying data
				if dataValueStr, ok := dataValue.(string); ok {
					// Sanitize inputs before processing
					dataValueStr = strings.TrimSpace(dataValueStr)

					// Omit fields with empty values
					if dataValueStr == "" {
						continue
					}

					switch dataType {
					// String
					case "S":
						if transformedJsonData, hasValidStructure := handlePrimitiveValue(dataType, dataValueStr); hasValidStructure {
							output.Data[key] = transformedJsonData
						} else {
							continue
						}
					// Number
					case "N":
						if transformedJsonData, hasValidStructure := handlePrimitiveValue(dataType, dataValueStr); hasValidStructure {
							output.Data[key] = transformedJsonData
						} else {
							continue
						}
					case "BOOL":
						if transformedJsonData, hasValidStructure := handlePrimitiveValue(dataType, dataValueStr); hasValidStructure {
							output.Data[key] = transformedJsonData
						} else {
							continue
						}
					// Null
					case "NULL":
						if transformedJsonData, hasValidStructure := handlePrimitiveValue(dataType, dataValueStr); hasValidStructure {
							output.Data[key] = transformedJsonData
						} else {
							continue
						}
					// List
					case "L":
						if transformedJsonData, hasValidStructure, _ := handleCompositeValue(dataType, dataValueStr); hasValidStructure {
							output.Data[key] = transformedJsonData
						} else {
							continue
						}
					// Map
					case "M":
						if transformedJsonData, hasValidStructure, _ := handleCompositeValue(dataType, dataValueStr); hasValidStructure {
							output.Data[key] = transformedJsonData
						} else {
							continue
						}
					}
				}
			}
		}
	}

	return output, nil
}

func handlePrimitiveValue(dataType string, dataValue string) (interface{}, bool) {
	switch dataType {
	// String
	case "S":

		if t, err := time.Parse(time.RFC3339, dataValue); err == nil {
			// transform `RFC3339` formatted `Strings` to `Unix Epoch` in `Numeric` data type.
			return t.Unix(), true
		} else {
			return dataValue, true
		}
	// Number
	case "N":
		if n, err := strconv.ParseFloat(dataValue, 64); err == nil {
			return n, true
		}
	case "BOOL":
		if strings.ToLower(dataValue) == "1" || strings.ToLower(dataValue) == "t" || strings.ToLower(dataValue) == "T" || strings.ToLower(dataValue) == "true" || strings.ToLower(dataValue) == "TRUE" || strings.ToLower(dataValue) == "True" {
			return true, true
		} else if strings.ToLower(dataValue) == "0" || strings.ToLower(dataValue) == "f" || strings.ToLower(dataValue) == "F" || strings.ToLower(dataValue) == "false" || strings.ToLower(dataValue) == "FALSE" || strings.ToLower(dataValue) == "False" {
			return false, true
		}
	// Null
	case "NULL":
		if strings.ToLower(dataValue) == "1" || strings.ToLower(dataValue) == "t" || strings.ToLower(dataValue) == "T" || strings.ToLower(dataValue) == "true" || strings.ToLower(dataValue) == "TRUE" || strings.ToLower(dataValue) == "True" {
			return nil, true
		} else if strings.ToLower(dataValue) == "0" || strings.ToLower(dataValue) == "f" || strings.ToLower(dataValue) == "F" || strings.ToLower(dataValue) == "false" || strings.ToLower(dataValue) == "FALSE" || strings.ToLower(dataValue) == "False" {
			return nil, true
		}
	default:
		return nil, false
	}
	return nil, false
}

// Recursive function to handle complex data types
func handleCompositeValue(dataType string, dataValue interface{}) (interface{}, bool, error) {
	switch dataType {
	case "L":
		if list, ok := dataValue.([]interface{}); ok {
			newList := make([]interface{}, 0, len(list))
			for _, item := range list {
				if itemMap, ok := item.(map[string]interface{}); ok {
					for itemDataType, itemDataValue := range itemMap {
						if itemDataValueStr, ok := itemDataValue.(string); ok {
							if transformedJsonData, hasValidStructure := handlePrimitiveValue(itemDataType, itemDataValueStr); hasValidStructure {
								newList = append(newList, transformedJsonData)
							}
						}
					}
				}
			}
			return newList, true, nil
		}
	case "M":
		if mapValue, ok := dataValue.(map[string]interface{}); ok {
			newMap := make(map[string]interface{})
			for mapKey, mapItem := range mapValue {
				if mapItemMap, ok := mapItem.(map[string]interface{}); ok {
					for itemDataType, itemDataValue := range mapItemMap {
						if itemDataValueStr, ok := itemDataValue.(string); ok {
							if transformedJsonData, hasValidStructure := handlePrimitiveValue(itemDataType, itemDataValueStr); hasValidStructure {
								newMap[mapKey] = transformedJsonData
							}
						}
					}
				}
			}
			return newMap, true, nil
		}
	default:
		return nil, false, fmt.Errorf("unsupported data type: %s", dataType)
	}
	return nil, false, nil
}
