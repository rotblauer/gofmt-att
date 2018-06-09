// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"github.com/mitchellh/go-homedir"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Echo default config to stderr.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// FIXME this is a shitty workaround because for some crazy reason viper
		// doesn't use an io.Writer for it's write config method.
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err)
		}
		tempPath := filepath.Join(home, fmt.Sprintf("%s", ".tmp.gofmt-att.toml"))
		err = ioutil.WriteFile(tempPath, []byte(""), os.ModePerm)
		if err != nil {
			log.Fatalln("can't create temp file:", err)
		}
		defer os.Remove(tempPath)
		if err := viper.WriteConfigAs(tempPath); err != nil {
			log.Fatalln("can't write config to temp file:", err)
		}
		b, err := ioutil.ReadFile(tempPath)
		if err != nil {
			log.Fatalln("cant read temp file:", err)
		}
		os.Stderr.Write(b)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
