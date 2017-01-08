package generate_test

import (
	"github.com/alecthomas/assert"
	"github.com/itzg/haproxy-gen/generate"
	"testing"
)

func TestConfig_ExportCertbotArgs(t *testing.T) {
	cfg := generate.NewConfig()
	cfg.Domains = []generate.Domain{
		{Domain: "one.example.com"},
		{Domain: "two.example.com"},
		{Domain: "three.example.com"},
	}

	result := cfg.ExportCertbotArgs()
	assert.Equal(t, "-d one.example.com -d two.example.com -d three.example.com", result)
}

func TestConfig_ExportPrimaryDomain(t *testing.T) {
	cfg := generate.NewConfig()
	cfg.Domains = []generate.Domain{
		{Domain: "one.example.com"},
		{Domain: "two.example.com"},
		{Domain: "three.example.com"},
	}

	result := cfg.ExportPrimaryDomain()
	assert.Equal(t, "one.example.com", result)
}
