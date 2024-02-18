package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendMailFormarding(requestBody []byte) error {

	// Create a new POST request with the target URL and request body
	req, err := http.NewRequest("POST", os.Getenv("FORWARDING_URL"), bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Set the Content-Type header for JSON data
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request and capture the response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}

	// Print the response status code and body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(responseBody))
	return nil
}
