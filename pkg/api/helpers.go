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
	method = strings.ToUpper(method)
	url := fmt.Sprintf("https://api.datadoghq.com/api/v1/%s", endpoint)

	var req *http.Request
	if body != nil {
		bodyJSON, _ := json.Marshal(body)
		postBody := bytes.NewBuffer(bodyJSON)
		req, _ = http.NewRequest(method, url, postBody)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	// add DD headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("DD-API-KEY", config.APIKey)
	req.Header.Add("DD-APPLICATION-KEY", config.AppKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func MakeGetRequest(endpoint string) ([]byte, error) {
	return MakeRequest(endpoint, "get", nil)
}

func MakePostRequest(endpoint string, body Postable) ([]byte, error) {
	return MakeRequest(endpoint, "post", body)
}
