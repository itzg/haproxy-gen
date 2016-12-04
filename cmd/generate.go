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
	"github.com/itzg/haproxy-gen/generate"
	"github.com/spf13/cobra"
)

var generateCfg = generate.Config{}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the haproxy.cfg file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		generate.Execute(&generateCfg)
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

	generateCmd.Flags().StringVar(&generateCfg.TemplatePath, "template", ".", "The directory that contains haproxy.cfg.tmpl")
	generateCmd.Flags().StringArrayVarP(&generateCfg.Domains, "domains", "d", nil, "The domains to proxy")
	generateCmd.Flags().StringArrayVarP(&generateCfg.Backends, "backends", "b", nil, "The host:port of the backends serving the domains")

	generateCmd.Flags().BoolVar(&generateCfg.Stats.Disabled, "stats-disabled", false, "Should stats endpoint be disabled")
	generateCmd.Flags().StringVar(&generateCfg.Stats.User, "stats-user", "admin", "Username of the stats endpoint")
	generateCmd.Flags().StringVar(&generateCfg.Stats.Password, "stats-password", "haproxy", "Password of the stats endpoint")
	generateCmd.Flags().StringVar(&generateCfg.Stats.BasePath, "stats-basepath", "/hastats", "Basepath of the stats endpoint")
}
