package config

import (
	_ "embed"
)

//go:embed view.tmpl
var ViewTemplate string
