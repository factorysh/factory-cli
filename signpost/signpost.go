package signpost

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/factorysh/factory-cli/client"
	"github.com/factorysh/factory-cli/factory"
)

type SignPost struct {
	Project *factory.Project
}

func New(project *factory.Project) *SignPost {
	return &SignPost{
		Project: project,
	}
}

type Target struct {
	Target string `json:"target"`
}

func (s *SignPost) EnvURL(environment string) (*url.URL, error) {
	u := fmt.Sprintf("%s/api/factory/v1/projects/%s/environments/%s",
		s.Project.Factory().GitlabUrl().String(),
		s.Project.Id(),
		environment)
	l := log.WithField("url", u)
	url, err := url.Parse(u)
	if err != nil {
		l.WithError(err).Error()
		return nil, err
	}
	l.Debug()
	return url, nil
}

func (s *SignPost) Target(environment string) (*url.URL, error) {
	u, err := s.EnvURL(environment)
	l := log.WithField("url", u)
	if err != nil {
		l.WithError(err).Error()
		return nil, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		l.WithError(err).Error()
		return nil, err
	}
	resp, err := s.Project.Session().Do(req)
	if err != nil {
		l.WithError(err).Error()
		return nil, err
	}
	l = l.WithField("status", resp.Status)
	var t *Target
	err = client.ReadJson(resp, &t)
	if err != nil {
		l.WithError(err).Error()
		return nil, err
	}
	l = l.WithField("target", t.Target)
	r, err := url.Parse("https://" + t.Target)
	if err != nil {
		l.WithError(err).Error()
		return nil, err
	}
	l.Debug()
	return r, nil
}
