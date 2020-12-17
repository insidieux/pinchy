package internal

import (
	// List of imports for registry extensions
	_ "github.com/insidieux/pinchy/internal/extension/registry/consul"
	// List of imports for source extensions
	_ "github.com/insidieux/pinchy/internal/extension/source/file"
)
