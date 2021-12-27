package gitlab

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

func gitUrl(url string) (string, string) {
	slugs := strings.SplitN(url, "@", 2)
	slugs = strings.SplitN(slugs[1], ":", 2)
	server := slugs[0]
	project := strings.SplitN(slugs[1], ".", 2)[0]
	return server, project
}

func httpUrl(url string) (string, string) {
	slugs := strings.SplitN(url, "/", 3)
	if (strings.Contains(url, "@")) {
		slugs = strings.SplitN(url, "@", 2)
		slugs = strings.SplitN(slugs[1], "/", 2)
	} else {
		slugs = strings.SplitN(slugs[2], "/", 2)
	}
	server := slugs[0]
	project := strings.SplitN(slugs[1], ".", 2)[0]
	return server, project
}

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
			slugs := strings.SplitN(line, "\t", 2)
			slugs = strings.SplitN(slugs[1], " ", 2)
			url := slugs[0]
			if strings.HasPrefix(strings.TrimSpace(url), "http") {
				// match "origin\thttp...gitlab@... (fetch)"
				server, project = httpUrl(url)
			} else {
				// match "origin\tgitlab@... (fetch)"
				server, project = gitUrl(url)
			}
			break
		}
	}
	return server, project, nil
}
