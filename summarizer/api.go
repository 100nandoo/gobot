package summarizer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gobot/config"
	"net/http"
	"net/url"
	"os"
)

const (
	baseURL     = "https://api.smmry.com/"
	paramAPIKey = "SM_API_KEY"
	paramURL    = "SM_URL"
	paramInput  = "sm_api_input"
)

func SummarizeURL(link string) (*SmmryResponse, *ErrorResponse) {
	params := url.Values{}
	params.Set(paramAPIKey, os.Getenv(config.SmmryKey))
	params.Set(paramURL, link)

	reqURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, &ErrorResponse{SmAPIMessage: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &ErrorResponse{SmAPIMessage: fmt.Sprintf("HTTP Status: %s", resp.Status)}
	}

	var smmryResponse SmmryResponse
	if err := json.NewDecoder(resp.Body).Decode(&smmryResponse); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, &ErrorResponse{SmAPIMessage: "Error decoding JSON response"}
	}

	return &smmryResponse, nil
}

func SummarizeText(text string) (*SmmryResponse, *ErrorResponse) {
	params := url.Values{}
	params.Set(paramAPIKey, os.Getenv(config.SmmryKey))
	reqURL := baseURL + "?" + params.Encode()

	textBody := paramInput + "=" + text
	body := []byte(textBody)
	r, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(body))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Expect", "")
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	smmryResponse := &SmmryResponse{}

	derr := json.NewDecoder(res.Body).Decode(smmryResponse)
	if derr != nil {
		panic(derr)
	}

	return smmryResponse, nil
}
