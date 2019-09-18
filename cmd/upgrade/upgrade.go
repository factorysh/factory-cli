package upgrade

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/factorysh/factory-cli/cmd/root"
	"github.com/factorysh/factory-cli/version"
)

var (
	release_url string = "https://github.com/factorysh/factory-cli/releases/download/%v/factory-%v-%v-%v.gz"
)

func releaseUrl(tag string) string {
	return fmt.Sprintf(release_url, tag, version.Os(), version.Arch(), tag)
}

func init() {
	root.RootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to latest CLI version",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Get(
			"https://api.github.com/repos/factorysh/factory-cli/releases/latest")
		if err != nil {
			log.WithError(err).Error()
			return err
		}
		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close()
		latest := fmt.Sprintf("%v", data["name"].(string))
		if version.Version() == latest {
			fmt.Println("Nothing to upgrade")
			return nil
		}
		url := releaseUrl(latest)
		fmt.Printf("Downloading %v ...\n", url)
		resp, err = http.Get(url)
		defer resp.Body.Close()

		filename := os.Args[0] + "-" + latest
		tmpfile, err := os.Create(filename)
		if err != nil {
			log.WithError(err).Error()
			return err
		}
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.WithError(err).Error()
			return err
		}
		defer reader.Close()
		_, err = io.Copy(tmpfile, reader)
		if err != nil {
			log.WithError(err).Error()
			return err
		}
		err = os.Chmod(tmpfile.Name(), 0700)
		if err != nil {
			log.WithError(err).Error()
			return err
		}
		fmt.Println("New binary is available at " + filename)
		return nil
	},
}
