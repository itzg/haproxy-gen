package generate_test

import (
	"bytes"
	"github.com/alecthomas/assert"
	"github.com/itzg/haproxy-gen/generate"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestSanity_Encode(t *testing.T) {
	config := generate.NewConfig()

	config.Domains = []generate.Domain{
		{
			Domain:  "www.example.com",
			Servers: []string{"backend:8080"},
		},
	}

	expected := `templatepath: .
domains:
- domain: www.example.com
  servers:
  - backend:8080
`
	encoded, err := yaml.Marshal(config)
	require.NoError(t, err)
	assert.Equal(t, expected, string(encoded))
}

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

func TestExecute(t *testing.T) {
	config, err := generate.LoadFromYamlFile("testdata/typical.yml")
	require.NoError(t, err)

	b := new(bytes.Buffer)
	err = generate.Execute(config, b)

	assert.NoError(t, err)
	assert.Equal(t, "global\r\n"+
		"    # set default parameters to the modern configuration\r\n"+
		"    tune.ssl.default-dh-param 2048\r\n"+
		"    ssl-default-bind-ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHAC\r\n    ssl-default-bind-options no-sslv3 no-tlsv10 no-tlsv11 no-tls-tickets\r\n    ssl-default-server-ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CH\r\n    ssl-default-server-options no-sslv3 no-tlsv10 no-tlsv11 no-tls-tickets\r\n#GLOBAL\r\n\r\ndefaults\r\n    mode    http\r\n    option  forwardfor\r\n    option  http-server-close\r\n    timeout client 30s\r\n    timeout connect 4s\r\n    timeout server 30s\r\n#DEFAULTS\r\n\r\n\r\nuserlist u1\r\n  user admin password $1$B7LfUIdP$PQGZFB2JQ0Tq/BRQrCtG//\r\n\r\n\r\nfrontend ft\r\n    bind    :443 ssl crt /etc/certs\r\n    redirect scheme https code 301 if !{ ssl_fc }\r\n    # HSTS (15768000 seconds = 6 months)\r\n    http-response set-header Strict-Transport-Security max-age=15768000\r\n\r\n    bind    :80\r\n    use_backend one_example_com if { hdr(host) -i one.example.com }\r\n\r\n\r\nbackend one_example_com\r\n    \r\n    server 0 server1:8080\r\n    \r\n    \r\n    stats   enable\r\n    stats   uri /stats\r\n    stats   auth admin:admin\r\n\r\n    \r\n    acl AuthOkay_u1 http_auth(u1)\r\n    http-request auth realm u1 if !AuthOkay_u1\r\n    \r\n\r\n#END\r\n",
		b.String())
}

func TestDomain_Ref(t *testing.T) {
	d := generate.Domain{
		Domain: "www.example.com",
	}

	assert.Equal(t, "www_example_com", d.Ref())
}
