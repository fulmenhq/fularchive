package config

// SetEmbeddedAssets is called from main to inject build-time embedded assets.
// This avoids embed path issues — the //go:embed directives live in cmd/refbolt
// where they can reference repo-root paths.
func SetEmbeddedAssets(catalog, schema []byte) {
	embeddedCatalog = catalog
	embeddedSchema = schema
}

var (
	embeddedCatalog []byte
	embeddedSchema  []byte
)

// EmbeddedCatalog returns the built-in provider catalog YAML.
func EmbeddedCatalog() []byte {
	return embeddedCatalog
}

// EmbeddedSchema returns the built-in provider schema YAML.
func EmbeddedSchema() []byte {
	return embeddedSchema
}
