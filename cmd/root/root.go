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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_gitlab "github.com/factorysh/factory-cli/gitlab"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	GitlabUrl   string
	GitlabToken string
	Project     string
	Verbose     bool
	Environment string
	client      *http.Client
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// this will run on all subcommands

		if Verbose {
			log.SetLevel(log.DebugLevel)
			log.Debug("Verbose mode on")
		} else {
			log.SetLevel(log.InfoLevel)
		}

		// its used to validate globale options
		if GitlabUrl == "" {
			fmt.Println("You must provide a valid gitlab url")
			os.Exit(1)
		}
		// get token from env if not already set via -t
		if GitlabToken == "" {
			GitlabToken = os.Getenv("PRIVATE_TOKEN")
			log.Debug(GitlabToken)
		}
		// get token from config if not already set
		if GitlabToken == "" {
			value := viper.Get("token")
			if value != nil {
				GitlabToken = value.(string)
			}
		}
		if GitlabToken == "" {
			fmt.Println("You must provide a valid gitlab token")
			os.Exit(1)
		} else {
			log.Debug(GitlabToken)
		}
	},
}

func loadPemFromEnv() {
	// check if we must add a CA from env
	pemPath := os.Getenv("CA_CERTIFICAT")
	if pemPath == "" {
		// also check config
		value := viper.Get("ca_certificat")
		if value != nil {
			pemPath = value.(string)
		}
	}
	if pemPath != "" {
		// read file
		pemData, err := ioutil.ReadFile(pemPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// append pem data to default transport's CAs
		certs := x509.NewCertPool()
		certs.AppendCertsFromPEM(pemData)
		tlsConfig := &tls.Config{}
		tlsConfig.RootCAs = certs
		http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, loadPemFromEnv)

	filenameHook := filename.NewHook()
	log.AddHook(filenameHook)

	default_gitlab, default_project, _ := _gitlab.GitRemote()
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.factory-cli.yaml)")
	RootCmd.PersistentFlags().StringVarP(&GitlabUrl, "gitlab", "g", default_gitlab, "Gitlab server url")
	RootCmd.PersistentFlags().StringVarP(&Project, "project", "p", default_project, "Gitlab project path")

	// show when token is set in env
	// do not set default token value (this appears in help)
	token_help := "Gitlab token"
	if os.Getenv("PRIVATE_TOKEN") != "" {
		token_help += " (default $PRIVATE_TOKEN)"
	}
	RootCmd.PersistentFlags().StringVarP(&GitlabToken, "token", "t", "", token_help)

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
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
