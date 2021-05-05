package api

import (
	"encoding/json"
	"fmt"
	"strings"
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
	Tags map[string]string
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
	Limit int    `json:"limit"`
	Query string `json:"query"`
	Sort  string `json:"sort"`
	Time  struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"time"`
}

func (o *LogQueryOpts) defaults() {
	if o.Limit == 0 {
		o.Limit = 100
	}
	if o.Sort == "" {
		o.Sort = "desc"
	}
	now := time.Now()
	if o.Time.From == "" {
		o.Time.From = fmt.Sprint(now.UnixNano() / int64(time.Millisecond))
	}
	if o.Time.To == "" {
		o.Time.To = fmt.Sprint(now.Add(time.Minute*30).UnixNano() / int64(time.Millisecond))
	}
}

func (l *Log) GetTag(name string) (string, error) {
	if l.Tags != nil {
		if tag, ok := l.Tags[strings.ToLower(name)]; !ok {
			return "", fmt.Errorf("tag '%s' not found", name)
		} else {
			return tag, nil
		}
	}

	l.loadAllTags()
	return l.GetTag(name)
}

func (l *Log) loadAllTags() {
	l.Tags = make(map[string]string)
	for _, tag := range l.Content.Tags {
		parts := strings.Split(tag, ":")
		key := strings.ToLower(parts[0])
		if key != "" {
			newVal := strings.Join(parts[1:], ":")
			oldVal, ok := l.Tags[key]   // duplicate tags?
			if ok && oldVal != newVal { // ... but no duplicate values
				newVal = strings.Join([]string{oldVal, newVal}, ", ")
			}
			l.Tags[key] = newVal
		}
	}
}

func Logs(opts *LogQueryOpts) ([]Log, error) {
	opts.defaults()
	body, err := MakePostRequest("logs-queries/list", opts)
	if err != nil {
		return nil, err
	}

	var response LogResponse
	err = json.Unmarshal(body, &response)

	return response.Logs, nil
}
