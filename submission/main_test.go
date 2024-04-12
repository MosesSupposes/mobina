package main

import (
	"reflect"
	"testing"
)

func TestTransformer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Output
	}{
		{
			name:     "transformer works on a simple JSON file",
			input:    "./testData/valid.json",
			expected: Output{Data: map[string]interface{}{"number_1": 1.5, "string_1": "784498", "string_2": int64(1405544146)}},
		},
		// TODO: Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the transformer function
			got, err := transformToJSON(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			// Compare the actual and expected outputs
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("transformToJSON() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHandleCompositeValue(t *testing.T) {
	tests := []struct {
		name      string
		dataType  string
		dataValue interface{}
		want      interface{}
		wantOk    bool
	}{
		{
			name:     "Test List",
			dataType: "L",
			dataValue: []interface{}{
				map[string]interface{}{"N": "1.23"},
				map[string]interface{}{"S": "test"},
			},
			want:   []interface{}{1.23, "test"},
			wantOk: true,
		},
		{
			name:     "Test Map",
			dataType: "M",
			dataValue: map[string]interface{}{
				"key1": map[string]interface{}{"N": "1.23"},
				"key2": map[string]interface{}{"S": "test"},
			},
			want:   map[string]interface{}{"key1": 1.23, "key2": "test"},
			wantOk: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk, err := handleCompositeValue(tt.dataType, tt.dataValue)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleCompositeValue() got = %v, want %v", got, tt.want)
			}
			if gotOk != tt.wantOk {
				t.Errorf("handleCompositeValue() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestHandlePrimitiveValue(t *testing.T) {
	tests := []struct {
		name      string
		dataType  string
		dataValue string
		want      interface{}
		wantOk    bool
	}{
		{
			name:      "Test String",
			dataType:  "S",
			dataValue: "test",
			want:      "test",
			wantOk:    true,
		},
		{
			name:      "Test Number",
			dataType:  "N",
			dataValue: "1.23",
			want:      1.23,
			wantOk:    true,
		},
		{
			name:      "Test Boolean",
			dataType:  "BOOL",
			dataValue: "true",
			want:      true,
			wantOk:    true,
		},
		{
			name:      "Test Null",
			dataType:  "NULL",
			dataValue: "true",
			want:      nil,
			wantOk:    true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk := handlePrimitiveValue(tt.dataType, tt.dataValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handlePrimitiveValue() got = %v, want %v", got, tt.want)
			}
			if gotOk != tt.wantOk {
				t.Errorf("handlePrimitiveValue() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
