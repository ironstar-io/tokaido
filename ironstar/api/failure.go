package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fatih/color"
)

type APIError struct {
	StatusCode   int
	IronstarCode string
	Message      string
}

func (err *APIError) Error() {
	fmt.Println()

	switch err.StatusCode {
	case 400:
		color.Red("Ironstar API call failed! (Bad Request)")
	case 401:
		color.Red("Ironstar API call failed! (Unauthorized)")
	case 403:
		color.Red("Ironstar API call failed! (Forbidden)")
	case 404:
		color.Red("Ironstar API call failed! (Not Found)")
	case 500:
		color.Red("Ironstar API call failed! (Server Error)")
		fmt.Println()
		color.Yellow("Please contact Ironstar Support - support@ironstar.io")
	default:
		color.Red("Ironstar API call failed!")
	}

	fmt.Println()
	fmt.Printf("Status Code: %+v\n", err.StatusCode)
	fmt.Println("Ironstar Code: " + err.IronstarCode)
	fmt.Println(err.Message)
}

var ErrIronstarAPICall = errors.New("Ironstar API call was unsuccessful!")

func (res *RawResponse) HandleFailure() error {
	f := &FailureBody{}
	err := json.Unmarshal(res.Body, f)
	if err != nil {
		return err
	}

	apiErr := &APIError{
		StatusCode:   res.StatusCode,
		IronstarCode: f.Code,
		Message:      f.Message,
	}

	apiErr.Error()

	return ErrIronstarAPICall
}

type FailureBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
