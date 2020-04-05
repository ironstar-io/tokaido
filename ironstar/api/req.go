package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type RawResponse struct {
	StatusCode int
	Body       []byte
}

type FailureBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
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

func Req(authToken, method, path string, payload map[string]string) (*RawResponse, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return ReqBytePayload(authToken, method, path, b)
}

// ReqBytePayload - Make a HTTP request to the Ironstar
func ReqBytePayload(authToken, method, path string, payload []byte) (*RawResponse, error) {
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

	var bodyBytes []byte
	if resp != nil && resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		bodyBytes = body
	}

	if err != nil {
		return nil, err
	}

	ir := &RawResponse{
		StatusCode: resp.StatusCode,
		Body:       bodyBytes,
	}

	defer resp.Body.Close()

	return ir, nil
}
