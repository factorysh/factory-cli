package gitlab

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

// GitRemote returns server, project, error
func GitRemote() (string, string, error) {
	cmd := exec.Command("git", "remote", "show")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	err = cmd.Start()
	if err != nil {
		return "", "", err
	}
	buff := bufio.NewReader(stdout)
	line, err := buff.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	stdout.Close()
	err = cmd.Wait()
	if err != nil {
		return "", "", err
	}
	origin := strings.TrimRight(line, "\n")
	cmd = exec.Command("git", "remote", "-v")
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	err = cmd.Start()
	if err != nil {
		return "", "", err
	}
	buff = bufio.NewReader(stdout)
	//re := regexp.MustCompilePOSIX(`Fetch URL: \\w+@([^:]+):([a-zA-Z\-_./]+)\\.git`)
	var server, project string
	for {
		line, err = buff.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", "", err
		}
		if strings.HasPrefix(strings.TrimSpace(line), origin+"\t") {
			// match "origin\tgitlab@... (fetch)"
			slugs := strings.SplitN(line, "\t", 2)
			slugs = strings.SplitN(slugs[1], "@", 2)
			slugs = strings.SplitN(slugs[1], ":", 2)
			server = slugs[0]
			project = strings.SplitN(slugs[1], ".", 2)[0]
			break
		}
	}
	return server, project, nil
}
