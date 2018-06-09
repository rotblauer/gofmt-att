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
	"github.com/mitchellh/go-homedir"
	"log"
	"path/filepath"
	"fmt"
	"io"
	"os"
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
		var err error
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err)
		}

		wantFile := filepath.Join(home, ".gofmt-att.json")
		if len(args) > 0 {
			wantFile = filepath.Clean(args[0])
		}

		var f io.Writer
		if forceWriteConfig {
			f, err = os.Create(wantFile)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			stat, err := os.Stat(wantFile)
			if err == nil {
				log.Println("file exists:", stat.Name())
			} else if !os.IsNotExist(err) {
				log.Fatalln("error checking file:", err)
			}
			f, err = os.Open(wantFile)
			if err != nil {
				log.Fatalln(err)
			}
		}
		err = writeDefaultConfig(f)
		if err != nil {
			log.Fatalln("could not write json file:", err)
		}
		fmt.Println("Wrote config file:", wantFile)
	},
}

func init() {
	createCmd.PersistentFlags().BoolVar(&forceWriteConfig, "force", false, "overwrite any existing config file")
	configCmd.AddCommand(createCmd)
}
