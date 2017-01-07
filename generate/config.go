package generate

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"unicode"
)

const DefaultStatsBasePath = "/_hastats"

type StatsConfig struct {
	Enabled bool
	// BasePath specifies the URL path where haproxy stats can be accessed. It defaults to DefaultStatsBasePath
	BasePath    string
	User        string
	Password    string
	Realm       string
	HideVersion bool
}

// ExtraConfig contains opaque text blocks that will be inserted into the respective sections.
type ExtraConfig struct {
	Global   string
	Defaults string
	// End content is inserted at the very end of the file
	End string
}

type Config struct {
	TemplatePath string
	// Certs is the path to a directory containing haproxy supported certificates. If not specified, only http on port 80
	// will be supported.
	Certs string

	FrontendStats StatsConfig
	Domains       []Domain
	UserLists     []UserList `yaml:",omitempty"`

	Extra ExtraConfig
}

type Domain struct {
	Domain string
	// Servers is the set of backend host:port servers that service this domain
	Servers []string
	// UserListName when non-empty indicates that HTTP basic auth should be enabled and restricted to the given
	// UserList under Config
	UserListName string `yaml:",omitempty"`
	Stats        StatsConfig
}

type User struct {
	Username string
	// EncPassword is a pre-encrypted password (with mkpasswd) for the user.
	EncPassword string
}

type UserList struct {
	// UserList is the name of the user list
	UserList string
	Users    []User
}

func NewConfig() *Config {
	return &Config{
		TemplatePath: ".",
	}
}

func refConversion(r rune) rune {
	if unicode.IsDigit(r) || unicode.IsLetter(r) {
		return r
	} else {
		return '_'
	}
}

func LoadFromYamlFile(filename string) (*Config, error) {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return LoadFromYaml(content)
}

func LoadFromYaml(content []byte) (*Config, error) {
	config := NewConfig()

	err := yaml.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) AddSimpleDomain(domain string, backendServer string) {
	c.Domains = append(c.Domains, Domain{
		Domain:  domain,
		Servers: []string{backendServer},
	})
}

func (d *Domain) Ref() string {
	return strings.Map(refConversion, d.Domain)
}
