// Package assets embeds the provider catalog and schema into the binary.
//
// IMPORTANT: The files in this directory are DERIVED COPIES — do NOT edit them directly.
//
//	Source of truth:
//	  catalog.yaml  ← configs/providers.yaml
//	  schema.yaml   ← schemas/providers/v0/providers.schema.yaml
//
//	To update: run `make embed-assets` (also runs automatically as part of `make build`).
//	Edits to catalog.yaml or schema.yaml will be silently overwritten on the next build.
package assets

import _ "embed"

//go:embed catalog.yaml
var Catalog []byte

//go:embed schema.yaml
var Schema []byte
