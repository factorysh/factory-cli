package signpost

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"gitlab.bearstech.com/factory/factory-cli/client"
	"gitlab.bearstech.com/factory/factory-cli/factory"
)

type SignPost struct {
	project *factory.Project
}

func New(project *factory.Project) *SignPost {
	return &SignPost{
		project: project,
	}
}

type Target struct {
	Target string `json:"target"`
}

func (s *SignPost) Target(environment string) (*url.URL, error) {
	u := fmt.Sprintf("%s/api/factory/v1/projects/%s/environments/%s",
		s.project.Factory().GitlabUrl().String(),
		s.project.Id(),
		environment)
	l := log.WithField("url", u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		l.WithError(err).Error()
		return nil, err
	}
	resp, err := s.project.Session().Do(req)
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
