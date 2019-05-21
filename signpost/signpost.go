package signpost

import (
	"fmt"
	"net/http"
	"net/url"

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
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/api/factory/v1/projects/%s/environments/%s",
			s.project.Factory().Target().String(),
			s.project.Id(),
			environment), nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.project.Session().Do(req)
	if err != nil {
		return nil, err
	}
	var t *Target
	err = client.ReadJson(resp, &t)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(t.Target)
	return u, err
}
