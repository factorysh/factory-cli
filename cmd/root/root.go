// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"fmt"
	"net/http"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_gitlab "gitlab.bearstech.com/factory/factory-cli/gitlab"
)

var (
	cfgFile   string
	GitlabUrl string
	Project   string
	Verbose   bool
	client    *http.Client
)

var RootCmd = &cobra.Command{
	Use:   "factory",
	Short: "Factory command line interface",
	Long: `
 _
| |             __            _
| | /| /| /|   / _| __ _  ___| |_ ___  _ __ _   _
| |/ |/ |/ |  | |_ / _' |/ __| __/ _ \| '__| | | |
|          |  |  _| (_| | (__| || (_) | |  | |_| |
+----------+  |_|  \__,_|\___|\__\___/|_|   \__, |
                                             |___/
`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.WithError(err).Error()
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig,
		func() {
			if Verbose {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}
		})
	filenameHook := filename.NewHook()
	log.AddHook(filenameHook)

	default_gitlab, default_project, _ := _gitlab.GitRemote()
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.factory-cli.yaml)")
	RootCmd.PersistentFlags().StringVarP(&GitlabUrl, "gitlab", "g", default_gitlab, "Gitlab server")
	RootCmd.PersistentFlags().StringVarP(&Project, "project", "P", default_project, "Project")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")

	client = &http.Client{}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".factory-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".factory-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
