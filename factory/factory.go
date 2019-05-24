package factory

import (
	"fmt"
	"net/url"

	"gitlab.bearstech.com/factory/factory-cli/client"
)

type Factory struct {
	gitlab_url   *url.URL
	privateToken string
	projects     map[string]*Project
}

func New(gitlab_url, privateToken string) (*Factory, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s", gitlab_url))
	if err != nil {
		return nil, err
	}
	return &Factory{
		gitlab_url:   u,
		privateToken: privateToken,
		projects:     make(map[string]*Project),
	}, nil
}

func (f *Factory) GitlabUrl() *url.URL {
	return f.gitlab_url
}

func (f *Factory) Project(project string) *Project {
	p, ok := f.projects[project]
	if ok {
		return p
	}
	f.projects[project] = &Project{
		factory: f,
		name:    project,
		session: client.New(project, f.privateToken),
	}
	return f.projects[project]
}

type Project struct {
	factory *Factory
	name    string
	session *client.Session
}

func (p *Project) Factory() *Factory {
	return p.factory
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) Session() *client.Session {
	return p.session
}

func (p *Project) Id() string {
	return url.PathEscape(p.name)
}
