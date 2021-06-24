// Copyright © 2016 Geoff Bourne <itzgeoff@gmail.com>
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
	"github.com/Sirupsen/logrus"
	"github.com/itzg/haproxy-gen/generate"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the haproxy.cfg file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		config := loadConfigFromCommonFlags(cmd)

		disableCerts, err := cmd.Flags().GetBool(FlagDisableCerts)
		if err != nil {
			logrus.Fatal(err)
		}
		if disableCerts {
			logrus.Info("Disabling certs usage")
			config.Certs = ""
		}

		var writer io.Writer
		outFilename, err := cmd.Flags().GetString(FlagOutFile)
		if err != nil {
			logrus.Fatal(err)
		}
		if outFilename != "" {
			file, err := os.Create(outFilename)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.WithField("name", outFilename).Infoln("Writing to file")
			defer file.Close()
			writer = file
		} else {
			writer = os.Stdout
		}

		err = generate.Execute(config, writer)
		if err != nil {
			os.Exit(2)
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

	addCommonFlags(generateCmd)

	generateCmd.Flags().StringP(FlagOutFile, "o", "", "The name of a file where rendered results are written. If not provided, then results are rendered to stdout")
	generateCmd.MarkFlagFilename(FlagOutFile)

	generateCmd.Flags().BoolP(FlagDisableCerts, "", false, "Disables the certs configuration such as when bootstrapping SSL")
}
