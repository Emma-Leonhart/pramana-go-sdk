package pramana

import "github.com/google/uuid"

// PramanaLinkable is the interface for objects that can be linked to
// entities in the Pramana knowledge graph. Provides identity and URL
// properties for graph integration.
type PramanaLinkable interface {
	// PramanaGuid returns the UUID (v4 or v5) identifying this entity in the Pramana graph.
	PramanaGuid() uuid.UUID

	// PramanaId returns the Pramana identifier string (e.g. "pra:num:3,1,2,1").
	// Returns empty string for objects that are not pseudo-class instances.
	PramanaId() string

	// PramanaHashUrl returns the Pramana entity URL using the hashed UUID,
	// e.g. "https://pramana.dev/entity/{PramanaGuid}".
	PramanaHashUrl() string

	// PramanaUrl returns the Pramana entity URL. For pseudo-class instances
	// this uses the PramanaId string; otherwise it falls back to PramanaHashUrl.
	PramanaUrl() string
}

// PramanaRoleful is the interface that all Pramana-mapped objects implement,
// providing access to the ontology roles the object participates in.
type PramanaRoleful interface {
	// GetRoles returns the PramanaRole instances that this object fulfils
	// within the Pramana ontology.
	GetRoles() []*PramanaRole
}
