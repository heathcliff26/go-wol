package static

import (
	"embed"
)

//go:embed all:css
//go:embed all:icons
//go:embed all:js
var Assets embed.FS

//go:embed index.html
var IndexTemplate string
