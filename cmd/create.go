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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mitchellh/go-homedir"
	"log"
	"path/filepath"
	"fmt"
)

var forceWriteConfig bool

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a config file from defaults",
	Long: `And writes it in .toml.

If provided, an argument will be uses as the file path.
Otherwise, it will write to the default location at $HOME/.gofmt-att.toml`,
	Run: func(cmd *cobra.Command, args []string) {

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err)
		}
		wantFile := filepath.Join(home, ".gofmt-att.toml")
		if len(args) > 0 {
			wantFile = args[0]
		}
		var configWriter func(string) error
		if forceWriteConfig {
			configWriter = viper.WriteConfigAs
		} else {
			configWriter = viper.SafeWriteConfigAs
		}
		if err := configWriter(wantFile); err != nil {
			log.Fatalln("could not write default config file:", err)
		} else {
			fmt.Println("Wrote config file:", wantFile)
		}
	},
}

func init() {
	createCmd.PersistentFlags().BoolVar(&forceWriteConfig, "force", false, "overwrite any existing config file")
	configCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
