package config

import (
	_ "embed"
)

//go:embed templates/pub.key
var PubKey string

//go:embed templates/priv.key
var PrivKey string
