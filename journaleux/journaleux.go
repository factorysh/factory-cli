package journaleux

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"

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

	u := fmt.Sprintf("%s/api/v1/journal/logs", j.host.String())

	params := url.Values{}
	params.Add("project", opts.Project)
	params.Add("lines", fmt.Sprintf("%v", opts.Lines))
	if opts.Since != 0 {
		params.Add("since", strconv.FormatInt(opts.Since, 10))
	}
	if opts.Until != 0 {
		params.Add("until", strconv.FormatInt(opts.Until, 10))
	}
	params.Add("priority", strconv.FormatUint(uint64(opts.Priority), 10))
	if opts.Regexp != "" {
		params.Add("regexp", opts.Regexp)
	}
	for _, name := range opts.Fields {
		params.Add("fields", fmt.Sprintf("%v=%v", name, opts.Fields[name]))
	}

	u = u + "?" + params.Encode()

	log.Debug(u)
	req, err := http.NewRequest("GET", u, nil)
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
