package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fatih/color"
)

func HandleFailure(res *RawResponse) error {
	f := &FailureBody{}
	err := json.Unmarshal(res.Body, f)
	if err != nil {
		return err
	}

	fmt.Println()
	color.Red("Ironstar API authentication failed!")
	fmt.Println()
	fmt.Printf("Status Code: %+v\n", res.StatusCode)
	fmt.Println("Ironstar Code: " + f.Code)
	fmt.Println(f.Message)

	return errors.New("Ironstar API call was unsuccessful!")
}