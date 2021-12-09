package jarvis

import _ "embed"

//go:embed tmpl/Dockerfile.tmpl
var dockerfileTmpl string
