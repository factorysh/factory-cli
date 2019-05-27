package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.bearstech.com/factory/factory-cli/version"
)

type Session struct {
	client       *http.Client
	project      string
	privateToken string
	headers      http.Header
}

func New(project string, privateToken string) *Session {
	s := &Session{
		client:       &http.Client{},
		project:      project,
		privateToken: privateToken,
		headers:      make(http.Header),
	}
	s.headers.Set("User-Agent", fmt.Sprintf("Factory-cli/%v", version.Version()))
	s.headers.Set("DNT", "1")
	s.headers.Set("Accept-Encoding", "gzip")
	s.headers.Set("Connection", "keep-alive")
	return s
}

func (s *Session) getMe(_url string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v4/user", _url), nil)
	if err != nil {
		return "", err
	}
	s.patchHeader(req)
	req.Header.Set("PRIVATE-TOKEN", s.privateToken)
	res, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	type u struct {
		Username string `json:"username"`
	}
	var jresp u
	err = ReadJson(res, &jresp)
	if err != nil {
		return "", err
	}
	return jresp.Username, nil
}

func (s *Session) getToken(realm string) (string, error) {
	if !strings.HasPrefix(realm, "http") {
		realm = "https://" + realm
	}
	l := log.WithField("realm", realm)
	u, err := url.Parse(realm)
	if err != nil {
		l.WithError(err).Error()
		return "", err
	}
	me, err := s.getMe(fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	if err != nil {
		l.WithError(err).Error()
		return "", err
	}
	l = l.WithField("me", me)
	buff := bytes.NewBufferString("client_id=gitlab_ci&service=container_registry&offline_token=true")
	fmt.Fprintf(buff, "&scope=repository:%v:%v", url.QueryEscape(s.project), "pull") // FIXME why always pull?
	u.RawQuery = buff.String()
	l = l.WithField("url", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		l.WithError(err).Error()
		return "", err
	}
	s.patchHeader(req)
	req.SetBasicAuth(me, s.privateToken)
	resp, err := s.client.Do(req)
	if err != nil {
		l.WithError(err).Error()
		return "", err
	}
	type t struct {
		Token string `json:"token"`
	}
	var token t
	err = ReadJson(resp, &token)
	if err != nil {
		l.WithError(err).Error()
		return "", err
	}
	l = l.WithField("token", token)
	l.Debug("getToken")
	return token.Token, nil
}

func (s *Session) patchHeader(req *http.Request) {
	for k, vs := range s.headers {
		for n, v := range vs {
			if n == 0 {
				req.Header.Set(k, v)
			} else {
				req.Header.Add(k, v)
			}
		}
	}
}

func (s *Session) Do(req *http.Request) (*http.Response, error) {
	s.patchHeader(req)
	r, err := s.client.Do(req)
	if err != nil {
		return r, err
	}
	if r.StatusCode == 401 { // Need auth
		authenticate := r.Header.Get("www-authenticate") // OAuth
		if authenticate == "" || !strings.HasPrefix(authenticate, "Bearer ") {
			return r, err
		}
		m, err := parseAuthenticate(authenticate[7:])
		if err != nil {
			return nil, err
		}
		realm, ok := m["realm"]
		if !ok {
			return nil, fmt.Errorf("Realm is mandatory: %v", authenticate)
		}
		token, err := s.getToken(realm)
		if err != nil {
			return nil, err
		}
		s.headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		// ok, I'm lazily authenticated, lets try once again
		s.patchHeader(req)
		return s.client.Do(req)
	}
	return r, err
}
