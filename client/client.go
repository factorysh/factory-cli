package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"gitlab.bearstech.com/factory/factory-cli/version"
)

type Session struct {
	client       *http.Client
	project      string
	privateToken string
	headers      http.Header
}

func New(project string) *Session {
	s := &Session{
		client:  &http.Client{},
		project: project,
		headers: make(http.Header),
	}
	s.headers.Set("User-Agent", fmt.Sprintf("Factory-cli/%v", version.Version()))
	s.headers.Set("DNT", "1")
	s.headers.Set("Accept-Encoding", "gzip")
	s.headers.Set("Connection", "keep-alive")
	return s
}

func parseAuthenticate(txt string) (map[string]string, error) {
	// https://tools.ietf.org/html/rfc2617#page-8
	values := make(map[string]string)
	re := regexp.MustCompile(`(\w+)\s*=\s*"(.*)"`)
	for _, blob := range strings.Split(txt, ",") {
		m := re.FindStringSubmatch(blob)
		if len(m) != 3 {
			return nil, fmt.Errorf("Invalid header: %v", blob)
		}
		values[m[1]] = m[2]
	}
	return values, nil
}

func readJson(resp *http.Response, value interface{}) error {
	if resp.StatusCode%100 != 2 {
		return fmt.Errorf("Bad status: %v", resp.Status)
	}
	defer resp.Body.Close()
	enc := resp.Header.Get("Content-encoding")
	var reader io.Reader
	if enc != "" {
		if enc == "gzip" {
			var err error
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Wrong encoding: %s", enc)
		}
	} else {
		reader = resp.Body
	}
	raw, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, value)
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
	err = readJson(res, &jresp)
	if err != nil {
		return "", err
	}
	return jresp.Username, nil
}

func (s *Session) getToken(realm string) (string, error) {
	u, err := url.Parse(realm)
	if err != nil {
		return "", err
	}
	me, err := s.getMe(fmt.Sprintf("%s/%s", u.Scheme, u.Host))
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("client_id=gitlab_ci&service=container_registry&offline_token=true")
	fmt.Fprintf(buff, "&scope=repository:%v:%v", s.project, "pull") // FIXME why always pull?
	u.RawQuery = url.QueryEscape(buff.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	s.patchHeader(req)
	req.SetBasicAuth(me, s.privateToken)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	type t struct {
		Token string `json:"token"`
	}
	var token t
	err = readJson(resp, &token)
	if err != nil {
		return "", err
	}
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
