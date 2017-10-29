package generate_test

import (
	"github.com/itzg/haproxy-gen/generate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestSanity_Decode(t *testing.T) {
	config := generate.NewConfig()

	in := `templatepath: .
domains:
- domain: www.example.com
  servers:
  - backend:8080
`
	err := yaml.Unmarshal([]byte(in), config)
	require.NoError(t, err)
}

func TestConfig_Empty(t *testing.T) {

	config, err := generate.LoadFromYaml([]byte(""))
	require.NoError(t, err)

	assert.NotNil(t, config)
	assert.NotEmpty(t, config.TemplatePath)
}

func TestConfig_MissingUserList(t *testing.T) {
	in := `
domains:
- domain: www.example.com
  servers:
  - "backend:8080"
  userlistname: unknown
`

	config, err := generate.LoadFromYaml([]byte(in))
	require.NoError(t, err)

	assert.False(t, config.Validate())
}

func TestConfig_EmptyBackendServers(t *testing.T) {
	in := `
domains:
- domain: www.example.com
  servers:
`

	config, err := generate.LoadFromYaml([]byte(in))
	require.NoError(t, err)

	assert.False(t, config.Validate())
}

func TestConfig_MissingDomainName(t *testing.T) {
	in := `
domains:
- servers:
  - "backend:8080"
`

	config, err := generate.LoadFromYaml([]byte(in))
	require.NoError(t, err)

	assert.False(t, config.Validate())
}

func TestFromFile_NoDomains(t *testing.T) {
	config, err := generate.LoadFromYamlFile("testdata/almost_empty.yml")
	require.NoError(t, err)

	assert.NotNil(t, config)
	assert.Equal(t, "templates", config.TemplatePath)
	assert.Len(t, config.Domains, 0)
}

func TestFromFile_Typical(t *testing.T) {
	config, err := generate.LoadFromYamlFile("testdata/typical.yml")
	require.NoError(t, err)

	assert.NotNil(t, config)
	assert.NotEmpty(t, config.TemplatePath)
	assert.Equal(t, "..", config.TemplatePath)
	assert.Len(t, config.Domains, 1)

	assert.Equal(t, "one.example.com", config.Domains[0].Domain)

	assert.True(t, config.Domains[0].Stats.Enabled)
}

func TestFromFile_TypicalThenAdded(t *testing.T) {
	config, err := generate.LoadFromYamlFile("testdata/typical.yml")
	require.NoError(t, err)

	config.AddSimpleDomain("another.example.com", "another:80")

	assert.Len(t, config.Domains, 2)

	assert.Equal(t, "one.example.com", config.Domains[0].Domain)
	assert.Equal(t, "another.example.com", config.Domains[1].Domain)
}

func TestDomain_Ref(t *testing.T) {
	d := generate.Domain{
		Domain: "www.example.com",
	}

	assert.Equal(t, "www_example_com", d.Ref())
}
