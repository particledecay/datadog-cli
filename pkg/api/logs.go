package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type LogResponse struct {
	Status string `json:"status"`
	Logs   []Log  `json:"logs"`
}

type Log struct {
	ID      string `json:"id"`
	Content struct {
		Service    string        `json:"service"`
		Tags       []string      `json:"tags"`
		Timestamp  time.Time     `json:"timestamp"`
		Host       string        `json:"host"`
		Attributes LogAttributes `json:"attributes"`
		Message    string        `json:"message"`
	} `json:"content"`
}

type LogAttributes struct {
	Level     string  `json:"level"`
	Process   string  `json:"process"`
	Timestamp int64   `json:"timestamp"`
	Agent     string  `json:"agent"`
	Filename  string  `json:"filename"`
	Lineno    float32 `json:"lineno"`
}

type LogQueryOpts struct {
	Limit int
	Query string
	Sort  string
	From  string
	To    string
}

func (o *LogQueryOpts) defaults() {
	if o.Limit == 0 {
		o.Limit = 100
	}
	if o.Sort == "" {
		o.Sort = "desc"
	}
	now := time.Now()
	if o.From == "" {
		o.From = fmt.Sprint(now.UnixNano() / int64(time.Millisecond))
	}
	if o.To == "" {
		o.To = fmt.Sprint(now.Add(time.Minute*30).UnixNano() / int64(time.Millisecond))
	}
}

func Logs(opts *LogQueryOpts) ([]Log, error) {
	opts.defaults()
	body, err := MakePostRequest("logs-queries/list", opts)
	if err != nil {
		return nil, err
	}

	var response LogResponse
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Logs, nil
}
