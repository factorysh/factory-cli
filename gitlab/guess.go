package gitlab

import (
	"bufio"
	"io"
	"os"
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
	cmd = exec.Command("git", "remote", "show", line[:len(line)-1])
	cmd.Env = append(os.Environ(), "LANG=C")
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
		if strings.HasPrefix(strings.TrimSpace(line), "Fetch URL:") {
			slugs := strings.SplitN(line[12:], "@", 2)
			slugs = strings.SplitN(slugs[1], ":", 2)
			server = slugs[0]
			project = slugs[1][:len(slugs[1])-5]
			break
		}
	}
	return server, project, nil
}
