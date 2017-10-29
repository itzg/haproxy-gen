package generate_test

import (
	"bytes"
	"testing"

	"github.com/itzg/haproxy-gen/generate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
)

func TestGenerate_AcmePluginScenario(t *testing.T) {
	cfg := generate.NewConfig()
	cfg.TemplatePath = ".."
	cfg.Certs = "/certs"
	cfg.HttpHttpsRedirectCondition = "!{ or ssl_fc url_acme_http01 }"
	cfg.Extra.Global = "chroot /var/lib/haproxy"
	cfg.Extra.PreFrontend = "acl url_acme_http01 path_beg /.well-known/acme-challenge/"

	b := new(bytes.Buffer)
	err := generate.Execute(cfg, b)
	require.NoError(t, err)

	assert.Equal(t, "global # set default parameters to the modern configuration tune.ssl.default-dh-param 2048 ssl-default-bind-ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHAC ssl-default-bind-options no-sslv3 no-tlsv10 no-tlsv11 no-tls-tickets ssl-default-server-ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CH ssl-default-server-options no-sslv3 no-tlsv10 no-tlsv11 no-tls-tickets "+
		"chroot /var/lib/haproxy defaults mode http option forwardfor option http-server-close timeout client 30s timeout connect 4s timeout server 30s "+
		"frontend ft acl url_acme_http01 path_beg /.well-known/acme-challenge/ "+
		"bind :443 ssl crt /certs redirect scheme https code 301 if !{ or ssl_fc url_acme_http01 } "+
		"# HSTS (15768000 seconds = 6 months) http-response set-header Strict-Transport-Security max-age=15768000 bind :80 ",
		normalize(b.Bytes()))
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

func normalize(in []byte) string {
	re := regexp.MustCompile(`\s+`)

	return string(re.ReplaceAll(in, []byte(" ")))
}
