package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Request struct {
	AuthToken        string
	Method           string
	Path             string
	MapStringPayload map[string]string
	BytePayload      []byte
}

type RawResponse struct {
	StatusCode int
	Body       []byte
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

func (r *Request) BuildBytePayload() error {
	if r.MapStringPayload != nil {
		b, err := json.Marshal(r.MapStringPayload)
		if err != nil {
			return err
		}

		r.BytePayload = b
	}

	return nil
}

// Send - Make a HTTP request to the Ironstar API
func (r *Request) Send() (*RawResponse, error) {
	err := r.BuildBytePayload()
	if err != nil {
		return nil, err
	}

	url := GetBaseURL() + r.Path
	req, err := http.NewRequest(r.Method, url, bytes.NewBuffer(r.BytePayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+r.AuthToken)

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
