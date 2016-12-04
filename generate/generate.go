package generate

import (
	"github.com/Sirupsen/logrus"
	"os"
	"path"
	"strings"
	"text/template"
	"unicode"
)

type Config struct {
	TemplatePath string
	Domains      []string
	Backends     []string
	Stats        struct {
		Disabled bool
		BasePath string
		User     string
		Password string
	}
}

type Context struct {
	Config

	DomainRefs []string
}

const tmplName = "haproxy.cfg.tmpl"

func Execute(cfg *Config) {
	if len(cfg.Domains) != len(cfg.Backends) {
		logrus.Fatal("The number of domains and backends need to equal")
	}
	logrus.WithFields(logrus.Fields{
		"templatePath": cfg.TemplatePath,
		"domains":      cfg.Domains,
		"backends":     cfg.Backends,
	}).Info("Generating")

	var ctx Context = Context{Config: *cfg}
	ctx.DomainRefs = make([]string, len(cfg.Domains))
	for i, d := range cfg.Domains {
		ctx.DomainRefs[i] = convertDomainToRef(d)
	}

	tmpl, err := template.New(tmplName).ParseFiles(path.Join(cfg.TemplatePath, tmplName))
	if err != nil {
		logrus.WithError(err).WithField("templatePath", cfg.TemplatePath).
			Fatal("Unable to parse template file")
	}

	err = tmpl.Execute(os.Stdout, ctx)
	if err != nil {
		logrus.WithError(err).WithField("templatePath", cfg.TemplatePath).
			Fatal("Unable to execute template")
	}
}

func convertDomainToRef(domain string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			return r
		} else {
			return '_'
		}
	}, domain)
}
