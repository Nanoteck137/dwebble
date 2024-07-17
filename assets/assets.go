package assets

import (
	"embed"
)

//go:embed *.png
var DefaultImagesFS embed.FS
