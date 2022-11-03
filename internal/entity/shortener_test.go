package entity

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestShortenerData_UnmarshalJSON(t *testing.T) {
	sData := ShortenerData{}
	bytes := []byte(`null`)

	err := json.Unmarshal(bytes, &sData)
	if err != nil {
		t.Fatal("Unexpected json unmarshal error in test:", err)
	}

	bytes = []byte(`{"URL": 2}`)
	expectedErr := fmt.Errorf("json: cannot unmarshal number into Go struct field .URL of type string")

	err = json.Unmarshal(bytes, &sData)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}

	bytes = []byte(`{"URL": ""}`)
	expectedErr = fmt.Errorf("required field 'URL' not found")

	err = json.Unmarshal(bytes, &sData)
	if !reflect.DeepEqual(err, expectedErr) {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}

	bytes = []byte(`{"URL": "www.google.com"}`)
	expectedErr = fmt.Errorf("required field 'URL'(%s) not matching a valid URL", "www.google.com")

	err = json.Unmarshal(bytes, &sData)
	if !reflect.DeepEqual(err, expectedErr) {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}

	expectedURL := "https://www.google.com"
	bytes = []byte(`{"URL": "https://www.google.com"}`)

	err = json.Unmarshal(bytes, &sData)
	if err != nil {
		t.Fatal("Unexpected json unmarshal error in test:", err)
	}
	if sData.URL != expectedURL {
		t.Fatalf("Expected URL: %v, got: %v", expectedURL, sData.URL)
	}

}
