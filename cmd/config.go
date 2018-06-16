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
	"io"
	"os"

	"log"
	"encoding/json"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Echo default config to stderr.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		writeDefaultConfig(os.Stderr)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func writeDefaultConfig(w io.Writer) error {
	c := fmtatt.DefaultFmtAttConfig
	// b, err := yaml.Marshal(c)
	// b, err := toml.Marshal(c)
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
		log.Fatalln("could not marshal json:", err)
	}
	_, err = w.Write(b)
	return err

	// enc := toml.NewEncoder(w)
	// return enc.Encode(c)
}
