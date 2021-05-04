package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/particledecay/datadog-cli/config"
)

type Postable interface {
	defaults()
}

func MakeRequest(endpoint, method string, body Postable) ([]byte, error) {
	var postBody []byte
	var err error
	if body != nil {
		postBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(
		strings.ToUpper(method),
		fmt.Sprintf("https://api.datadoghq.com/api/v1/%s", endpoint),
		bytes.NewBuffer(postBody),
	)
	if err != nil {
		return nil, err
	}

	// add DD headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("DD-API-KEY", config.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

func MakeGetRequest(endpoint string) ([]byte, error) {
	return MakeRequest(endpoint, "get", nil)
}

func MakePostRequest(endpoint string, body Postable) ([]byte, error) {
	return MakeRequest(endpoint, "post", body)
}
