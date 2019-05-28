package root

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.bearstech.com/factory/factory-cli/signpost"
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
