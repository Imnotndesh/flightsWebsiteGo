package Tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestUserEndpoint(t *testing.T) {
	endpoint := "localhost:9080/user/me"
	requestBody := map[string]interface{}{
		"username": "brian",
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Errorf("Error marshalling body: %v", err)
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(requestBodyBytes))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing request: %v", err)
	}
	defer resp.Body.Close()
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Error("Error unmarshalling response")
	}
	fmt.Println("Response: \n", response)
}
