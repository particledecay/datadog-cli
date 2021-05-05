package api

import (
	"encoding/json"
)

type AuthCheck struct {
	Valid bool `json:"valid"`
}

func Validate() (bool, error) {
	body, err := MakeGetRequest("validate")
	if err != nil {
		return false, err
	}

	var auth AuthCheck
	err = json.Unmarshal(body, &auth)
	if err != nil {
		return false, err
	}

	return auth.Valid, nil
}
