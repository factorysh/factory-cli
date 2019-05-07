package journaleux

import (
	"fmt"
	"net/http"
	"net/url"

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

func (p *Project) Logs(lines int) {

}
