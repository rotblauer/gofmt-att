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
	"github.com/rotblauer/gofmt-att/fmtatt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	// "github.com/pelletier/go-toml"
	"github.com/kr/pretty"
	"bufio"
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs shit.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var c = &fmtatt.DefaultFmtAttConfig
		// var c = &fmtatt.Config{}
		configPath := viper.ConfigFileUsed()
		if configPath != "" {
			log.Println("Reading config from file:", configPath)
			// err := viper.Unmarshal(c)
			bb, err := ioutil.ReadFile(configPath)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(bb, c)
			// // err = yaml.Unmarshal(bb, c)
			// err = toml.Unmarshal(bb, c)
			if err != nil {
				log.Fatalln("unable to decode into struct, %v", err)
			}
		} else {
			log.Fatalln("No config file found. Exiting.")
		}

		// make sure everything looks in order
		pretty.Logln(c)
		f := fmtatt.New(c)

		fmt.Println(" -> Look OK? (y/n)")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "y") {
				// ok = true
				break
			} else {
				fmt.Println("Done.")
				os.Exit(0)
			}
		}
		f.Go([3]fmtatt.DryRunT{fmtatt.WetRun, fmtatt.WetRun, fmtatt.WetRun})
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
