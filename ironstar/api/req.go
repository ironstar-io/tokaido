package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type IronstarAPIResponse struct {
	StatusCode int
	Body       string
}

// TODO - Change to real prod domain
const IronstarProductionAPIDomain = "https://localhost:8080"

func GetBaseURL() string {
	ipa := os.Getenv("IRONSTAR_PRODUCTION_API_DOMAIN")
	if ipa != "" {
		return ipa
	}

	return IronstarProductionAPIDomain
}

func Req(authToken, method, path string, payload map[string]string) (*IronstarAPIResponse, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return ReqBytePayload(authToken, method, path, b)
}

// ReqBytePayload - Make a HTTP request to the Ironstar
func ReqBytePayload(authToken, method, path string, payload []byte) (*IronstarAPIResponse, error) {
	url := GetBaseURL() + path
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+authToken)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)

	var bodyString string
	if resp != nil && resp.Body != nil {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		bodyString = string(bodyBytes)
	}

	if err != nil {
		return nil, err
	}

	ir := &IronstarAPIResponse{
		StatusCode: resp.StatusCode,
		Body:       bodyString,
	}

	defer resp.Body.Close()

	return ir, nil
}
