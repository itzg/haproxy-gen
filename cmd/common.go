package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/itzg/haproxy-gen/generate"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

const (
	FlagConfigFile = "in"
	FlagDomain     = "domain"
	FlagDomains    = "domains"
	FlagOutFile    = "out"
	DomainsDelim   = ","
)

var ReSimpleDomain = regexp.MustCompile(`(.*?)@(.*?:\d+)`)

func loadConfigFromCommonFlags(cmd *cobra.Command) *generate.Config {

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

	// from distinct args
	simpleDomains, err := cmd.Flags().GetStringSlice(FlagDomain)
	if err != nil {
		logrus.Fatal(err)
	}

	// from semicolon joined arg
	joinedSimpleDomains, err := cmd.Flags().GetString(FlagDomains)
	if err != nil {
		logrus.Fatal(err)
	}
	for _, part := range strings.Split(joinedSimpleDomains, DomainsDelim) {
		simpleDomains = append(simpleDomains, part)
	}

	for _, simpleDomain := range simpleDomains {
		parts := ReSimpleDomain.FindStringSubmatch(simpleDomain)
		if parts == nil {
			logrus.WithField("given", simpleDomain).Warn("Invalid simple domain format")
			continue
		}

		config.AddSimpleDomain(parts[1], parts[2])
	}

	return config
}

func addCommonFlags(cmd *cobra.Command) {

	cmd.Flags().StringP(FlagConfigFile, "i", "", "A YAML configuration file for haproxy-gen")
	cmd.MarkFlagFilename(FlagConfigFile, "yaml", "yml")

	cmd.Flags().StringSliceP(FlagDomain, "d", []string{}, "A domain definition formatted as FRONTEND_HOST@SERVER:PORT")
	cmd.Flags().String(FlagDomains, "", "Domain definitions separated by commas")
}
