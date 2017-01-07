// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"regexp"
)

const (
	FlagConfigFile = "in"
	FlagDomain     = "domain"
	FlagOutFile    = "out"
)

var ReSimpleDomain = regexp.MustCompile(`(.*?)@(.*?:\d+)`)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the haproxy.cfg file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		configFile, err := cmd.Flags().GetString(FlagConfigFile)
		if err != nil {
			logrus.Fatal(err)
		}

		var config *generate.Config
		if configFile != "" {
			config, err = generate.LoadFromYamlFile(configFile)
			if err != nil {
				logrus.Fatal(err)
			}
		} else {
			config = generate.NewConfig()
		}

		simpleDomains, err := cmd.Flags().GetStringSlice(FlagDomain)
		if err != nil {
			logrus.Fatal(err)
		}

		for _, simpleDomain := range simpleDomains {
			parts := ReSimpleDomain.FindStringSubmatch(simpleDomain)
			if parts == nil {
				logrus.WithField("given", simpleDomain).Warn("Invalid simple domain format")
				continue
			}

			config.AddSimpleDomain(parts[1], parts[2])
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	generateCmd.Flags().StringP(FlagConfigFile, "i", "", "A YAML configuration file for haproxy-gen")
	generateCmd.MarkFlagFilename(FlagConfigFile, "yaml", "yml")

	generateCmd.Flags().StringSliceP(FlagDomain, "d", []string{}, "A domain definition formatted as FRONTEND_HOST@SERVER:PORT")

	generateCmd.Flags().StringP(FlagOutFile, "o", "", "The name of a file where rendered results are written. If not provided, then results are rendered to stdout")
	generateCmd.MarkFlagFilename(FlagOutFile)
}
