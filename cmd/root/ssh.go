package root

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/factorysh/factory-cli/signpost"
)

// SSHAddress return user@domain for an ssh connection
func SSHAddress() (string, error) {
	log.Debug(GitlabUrl)
	log.Debug(Project)

	f, err := Factory()
	if err != nil {
		return "", err
	}
	s := signpost.New(f.Project(Project))
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
