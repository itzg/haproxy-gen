package generate_test

import (
	"bytes"
	"github.com/itzg/haproxy-gen/generate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"regexp"
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

func normalize(in []byte) string {
	re := regexp.MustCompile(`\s+`)

	return string(re.ReplaceAll(in, []byte(" ")))
}

func TestExecute(t *testing.T) {
	config, err := generate.LoadFromYamlFile("testdata/typical.yml")
	require.NoError(t, err)

	b := new(bytes.Buffer)
	err = generate.Execute(config, b)

	assert.NoError(t, err)
	assert.Equal(t, "global # set default parameters to the modern configuration tune.ssl.default-dh-param 2048"+
		" ssl-default-bind-ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHAC"+
		" ssl-default-bind-options no-sslv3 no-tlsv10 no-tlsv11 no-tls-tickets ssl-default-server-ciphers"+
		" ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CH"+
		" ssl-default-server-options no-sslv3 no-tlsv10 no-tlsv11 no-tls-tickets #GLOBAL defaults mode http option"+
		" forwardfor option http-server-close timeout client 30s timeout connect 4s timeout server 30s #DEFAULTS"+
		" userlist u1 user admin password $1$B7LfUIdP$PQGZFB2JQ0Tq/BRQrCtG// frontend ft bind :443 ssl crt /etc/certs"+
		" redirect scheme https code 301 if !{ ssl_fc } # HSTS (15768000 seconds = 6 months) http-response set-header"+
		" Strict-Transport-Security max-age=15768000 bind :80 use_backend one_example_com if { hdr(host) -i one.example.com }"+
		" backend one_example_com server 0 server1:8080 stats enable stats uri /stats stats auth admin:admin"+
		" acl AuthOkay_u1 http_auth(u1) http-request auth realm u1 if !AuthOkay_u1 #END ",
		normalize(b.Bytes()))
}

func TestDomain_Ref(t *testing.T) {
	d := generate.Domain{
		Domain: "www.example.com",
	}

	assert.Equal(t, "www_example_com", d.Ref())
}
