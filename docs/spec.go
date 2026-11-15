package docs

import _ "embed"

// OpenAPISpec menyimpan seluruh isi file openapi.json yang di-embed ke dalam binary saat build.
//
//go:embed openapi.json
var OpenAPISpec []byte
