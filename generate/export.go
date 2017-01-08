package generate

import (
	"fmt"
	"strings"
)

func (c *Config) ExportCertbotArgs() string {
	parts := make([]string, 0)

	for _, d := range c.Domains {
		parts = append(parts, fmt.Sprintf("-d %s", d.Domain))
	}

	return strings.Join(parts, " ")
}

func (c *Config) ExportPrimaryDomain() string {
	if len(c.Domains) == 0 {
		return ""
	}

	return c.Domains[0].Domain
}
