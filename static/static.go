package static

import (
	"embed"
)

//go:embed all:css
var CSS embed.FS

//go:embed all:icons
var Icons embed.FS

//go:embed all:js
var JS embed.FS

//go:embed index.html
var IndexTemplate string
