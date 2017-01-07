package generate

import (
	"github.com/Sirupsen/logrus"
	"io"
	"path"
	"text/template"
)

const tmplName = "haproxy.cfg.tmpl"

func Execute(cfg *Config, wr io.Writer) error {
	logrus.WithFields(logrus.Fields{
		"templatePath": cfg.TemplatePath,
		"certs":        cfg.Certs,
		"domains":      cfg.Domains,
		"userLists":    cfg.UserLists,
	}).Debug("Generating")

	tmpl, err := template.New(tmplName).ParseFiles(path.Join(cfg.TemplatePath, tmplName))
	if err != nil {
		logrus.WithError(err).WithField("templatePath", cfg.TemplatePath).
			Error("Unable to parse template file")

		return err
	}

	err = tmpl.Execute(wr, cfg)
	if err != nil {
		logrus.WithError(err).WithField("templatePath", cfg.TemplatePath).
			Error("Unable to execute template")

		return err
	}

	return nil
}
