package pramana

import (
	"fmt"

	"github.com/google/uuid"
)

// PramanaParticularClassID is the well-known class ID for PramanaParticular in the ontology.
var PramanaParticularClassID = uuid.MustParse("13000000-0000-4000-8000-000000000004")

// PramanaParticularClassURL returns the class-level URL in the Pramana graph.
func PramanaParticularClassURL() string {
	return fmt.Sprintf("https://pramana.dev/entity/%s", PramanaParticularClassID)
}

// PramanaParticular is a subclass of PramanaObject used for concrete
// entities in the Pramana OGM class hierarchy.
type PramanaParticular struct {
	PramanaObject
}

// NewPramanaParticular creates a new PramanaParticular with no ID.
func NewPramanaParticular() *PramanaParticular {
	return &PramanaParticular{PramanaObject: *NewPramanaObject()}
}

// NewPramanaParticularWithID creates a new PramanaParticular with the given UUID.
func NewPramanaParticularWithID(id uuid.UUID) *PramanaParticular {
	return &PramanaParticular{PramanaObject: *NewPramanaObjectWithID(id)}
}
