// Copyright © 2018 Patrick Nuckolls <nuckollsp+poorbox at gmail.com>
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

package cmd

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"errors"

	"github.com/spf13/cobra"
	pdb "github.com/yourfin/poorbox/poorboxdb"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "poorbox",
		Short: "The poorbox backend cli",
		Long: `This command handles everything poorbox related after install

Built with love by [nuckolls] and [guptayas]

Poorbox is designed to be run with postgress running under docker
(typically exposed on port 5432), and will break without it.
An internet connection is also required to grab information from themoviedb.org.

See README.md for more details. Happy streaming :)
`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}

	pgIdFilePath, pgEndpoint string
	tmdbApiKey string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&pgEndpoint,
		"pg-endpoint",
		"e",
		"localhost:5432",
		"The network location of postgres")

	// Why this is handled with a file:
	// http://web.archive.org/web/20180421040440/https://www.netmeister.org/blog/passing-passwords.html
	rootCmd.PersistentFlags().StringVarP(
		&pgIdFilePath,
		"pg-identity-file",
		"s",
		"./pg-secret",
		`The file containing the postgress username and password
seperated by a newline`)
}

func parseSecretsFile(filepath string, lines int) (results []string, err error) {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	scanned := true
	for ii := 0; scanned && ii < lines; ii++ {
		scanner.Scan()
		results = append(results, strings.TrimSpace(scanner.Text()))
	}
	if !scanned {
		errorNotice := "Didn't parse enough lines from secret file: %v"
		return make([]string, 0), errors.New(fmt.Sprintf(errorNotice, filepath))
	}
	return
}

func parsePgSecrets() (password string, username string, err error) {
	parsed, err := parseSecretsFile(pgIdFilePath, 2)
	if err != nil {
		return
	}
	username = parsed[0]
	password = parsed[1]
	return
}

func parseTmdbApiSecret() (apiKey string, err error) {
	parsed, err := parseSecretsFile(tmdbApiKey, 1)
	if err != nil {
		return
	}
	apiKey = parsed[0]
	return
}

func parseConnect() error {
	username, password, err := parsePgSecrets()
	if err != nil {
		return err
	}
	err = pdb.Connect(username, password, pgEndpoint, pdb.SSL_DISABLE)
	if err != nil {
		return err
	}
	return nil
}
