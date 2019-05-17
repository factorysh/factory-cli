package journaleux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/factorysh/go-longrun/longrun/sse"
	"gitlab.bearstech.com/factory/factory-cli/client"
)

type Journaleux struct {
	target       *url.URL
	privateToken string
	projects     map[string]*Project
}

type Project struct {
	journaleux *Journaleux
	session    *client.Session
}

func New(target, privateToken string) (*Journaleux, error) {
	t, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &Journaleux{
		target:       t,
		privateToken: privateToken,
		projects:     make(map[string]*Project),
	}, nil
}

func (j *Journaleux) Project(project string) *Project {
	p, ok := j.projects[project]
	if ok {
		return p
	}
	j.projects[project] = &Project{
		journaleux: j,
		session:    client.New(project, j.privateToken),
	}
	return j.projects[project]
}

func (p *Project) Hello() (string, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/hello", p.journaleux.target.String()), nil)
	if err != nil {
		return "", err
	}
	resp, err := p.session.Do(req)
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
	Project string `json:"project`
	Lines   int    `json:"lines`
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

func (p *Project) Logs(opts *LogsOpt, visitor func(*Event, error) error) error {
	buff, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/logs", p.journaleux.target.String()),
		bytes.NewReader(buff))
	if err != nil {
		return err
	}
	resp, err := p.session.Do(req)
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
