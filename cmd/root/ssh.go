package root

import (
	"fmt"
	"os"
	"strings"

	"github.com/factorysh/factory-cli/signpost"
	log "github.com/sirupsen/logrus"
)

// return current Signpost
func SignPost() (*signpost.SignPost, error) {
	log.Debug(GitlabUrl)
	log.Debug(Project)

	f, err := Factory()
	if err != nil {
		return nil, err
	}
	s := signpost.New(f.Project(Project))
	return s, nil
}

// SSHAddress return user@domain for an ssh connection
func SSHAddress() (string, error) {
	s, err := SignPost()
	if err != nil {
		return "", err
	}
	u, err := s.Target(Environment)
	if err != nil {
		return "", err
	}
	log.Debug(u)
	user := strings.Replace(Project, "/", "-", -1)
	return fmt.Sprintf("%s@%s", user, u.Hostname()), nil
}

// return an array with extra ssh options
// eg: -i pki/canary_key/id_rsa -o StrictHostKeyChecking=no
func SSHExtraArgs() []string {
	var args = []string{}
	value, ok := os.LookupEnv("SSH_IDENTITY_FILE")
	if ok {
		args = append(args, "-i")
		args = append(args, value)
	}
	value, ok = os.LookupEnv("SSH_STRICT_HOST_KEY_CHECKING")
	if ok && value == "no" {
		args = append(args, "-oStrictHostKeyChecking=no")
	}
	return args
}
