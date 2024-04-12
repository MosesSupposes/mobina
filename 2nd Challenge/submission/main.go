package main

import (
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
	// Read input
	var input Input
	decoder := json.NewDecoder(os.Stdin)
	err := decoder.Decode(&input)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Transform
	output := Output{Data: make(map[string]interface{})}
	for key, value := range input.Data {
		key = strings.TrimSpace(key)

		// Omit fields with empty keys
		if key == "" {
			continue
		}

		if valueMap, ok := value.(map[string]interface{}); ok {
			for dataType, dataValue := range valueMap {
				// Use type assertion to access underlying data
				if dataValueStr, ok := dataValue.(string); ok {
					// Sanatize inputs before processing
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
						if transformedJsonData, hasValidStructure := handleCompositeValue(dataType, dataValueStr); hasValidStructure {
							output.Data[key] = transformedJsonData
						} else {
							continue
						}
					// Map
					case "M":
						if transformedJsonData, hasValidStructure := handleCompositeValue(dataType, dataValueStr); hasValidStructure {
							output.Data[key] = transformedJsonData
						} else {
							continue
						}
					}
				}
			}
			// Write output
			encoder := json.NewEncoder(os.Stdout)
			err = encoder.Encode(output.Data)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
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
func handleCompositeValue(dataType string, dataValue interface{}) (interface{}, bool) {
	switch dataType {
	case "L":
		list := dataValue.([]interface{})
		newList := make([]interface{}, 0, len(list))
		for _, item := range list {
			if itemMap, ok := item.(map[string]interface{}); ok {
				for k, v := range itemMap {
					if k != "NULL" && k != "L" && k != "M" {
						if dataValueStr, ok := v.(string); ok {
							compositePart, _ := handlePrimitiveValue(k, dataValueStr)
							newList = append(newList, compositePart)
						}
					}
				}
			}
		}
		if len(newList) == 0 {
			return nil, false
		} else {
			return newList, true
		}
	case "M":
		m := dataValue.(map[string]interface{})
		newMap := make(map[string]interface{})
		for k, v := range m {
			if vMap, ok := v.(map[string]interface{}); ok {
				for k2, v2 := range vMap {
					if dataValueStr, ok := v2.(string); ok {
						compositePart, _ := handlePrimitiveValue(k2, dataValueStr)
						newMap[k] = compositePart
					}
				}
			}
		}
		if len(newMap) == 0 {
			return nil, false
		} else {
			return newMap, true
		}
	default:
		return nil, false
	}
}
