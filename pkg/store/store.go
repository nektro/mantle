package store

import (
	"github.com/nektro/mantle/pkg/store/local"
)

// global singleton
var (
	This *Store
)

// PreInit registers flags
func PreInit() {
}

// Init takes flag values and initializes datastore
	This = &Store{local.Get()}
}
