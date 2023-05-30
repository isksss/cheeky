package config

import (
	_ "embed"
	"text/template"
)

//go:embed templates/pub.key
var pubkey string

//go:embed templates/priv.key
var privkey string

// GetTemplates gets the embedded templates
func GetTemplates() (*template.Template, *template.Template, error) {
	pubTemplate, err := template.New("pubkey").Parse(pubkey)
	if err != nil {
		return nil, nil, err
	}
	privTemplate, err := template.New("privkey").Parse(privkey)
	if err != nil {
		return nil, nil, err
	}
	return pubTemplate, privTemplate, nil
}
