package journaleux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/factorysh/factory-cli/client"
	"github.com/factorysh/factory-cli/factory"
	"github.com/factorysh/go-longrun/longrun/sse"
)

type Journaleux struct {
	project *factory.Project
	host    *url.URL
}

func New(project *factory.Project, host *url.URL) *Journaleux {
	return &Journaleux{
		project: project,
		host:    host,
	}
}

func (j *Journaleux) Hello() (string, error) {
	url := fmt.Sprintf("%s/api/v1/journal/hello", j.host.String())
	log.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := j.project.Session().Do(req)
	if err != nil {
		return "", err
	}
	type r struct {
		Msg string `json:"msg"`
	}
	var rr r
	err = client.ReadJson(resp, &rr)
	if err != nil {
		return "", err
	}
	return rr.Msg, nil
}

type LogsOpt struct {
	Project  string            `json:"project"`
	Lines    int               `json:"lines"`
	Follow   bool              `json:"follow"`
	Since    int64             `json:"since"`
	Until    int64             `json:"until"`
	Priority uint              `json:"priority"`
	Regexp   string            `json:"regexp"`
	Fields   map[string]string `json:"fields"`
}

type Event struct {
	Monotonic uint64            `json:"monotonic"`
	Realtime  uint64            `json:"realtime"`
	Message   string            `json:"message"`
	Priority  uint              `json:"priority"`
	Fields    map[string]string `json:"fields"`
}

type EventOrError struct {
	Event *Event `json:"event"`
	Error error  `json:"error"`
}

func (j *Journaleux) Logs(opts *LogsOpt, visitor func(*Event, error) error) error {
	buff, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/journal/logs", j.host.String())
	log.Debug(url)
	req, err := http.NewRequest("POST", url, bytes.NewReader(buff))
	if err != nil {
		return err
	}
	resp, err := j.project.Session().Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad status: %v", resp.Status)
	}
	defer resp.Body.Close()
	err = sse.Reader(resp.Body, func(evt *sse.Event) error {
		var event EventOrError
		err := json.Unmarshal([]byte(evt.Data), &event)
		if err != nil {
			return err
		}
		return visitor(event.Event, event.Error)
	})
	if err != nil {
		return err
	}
	return nil
}
