package static

import "embed"

//go:embed swagger-ui.html generated/swagger.swagger.json
var FS embed.FS
