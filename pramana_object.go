package pramana

import (
	"fmt"

	"github.com/google/uuid"
)

// PramanaObjectRootID is the well-known root ID for the PramanaObject class in the ontology.
var PramanaObjectRootID = uuid.MustParse("10000000-0000-4000-8000-000000000001")

// PramanaObjectClassURL returns the class-level URL in the Pramana graph.
func PramanaObjectClassURL() string {
	return fmt.Sprintf("https://pramana.dev/entity/%s", PramanaObjectRootID)
}

// PramanaObject is the base type for all objects mapped into the Pramana knowledge graph.
// It implements PramanaLinkable for graph identity and PramanaRoleful for ontology
// role participation.
//
// Friction by design: IDs are never auto-generated. A new PramanaObject starts
// with uuid.Nil and only receives a real UUID v4 when GenerateId is explicitly
// called. This prevents disposable or transient objects from polluting the graph
// with throw-away identifiers. Once assigned, the ID is immutable — calling
// GenerateId a second time returns a PramanaError.
type PramanaObject struct {
	pramanaGuid uuid.UUID
}

// NewPramanaObject creates a new PramanaObject with uuid.Nil (empty).
func NewPramanaObject() *PramanaObject {
	return &PramanaObject{pramanaGuid: uuid.Nil}
}

// NewPramanaObjectWithID creates a new PramanaObject with the given UUID.
func NewPramanaObjectWithID(id uuid.UUID) *PramanaObject {
	return &PramanaObject{pramanaGuid: id}
}

// PramanaGuid returns the UUID identifying this entity in the Pramana graph.
func (o *PramanaObject) PramanaGuid() uuid.UUID {
	return o.pramanaGuid
}

// PramanaId returns the Pramana identifier string. Regular objects do not
// belong to a pseudo-class, so this returns an empty string.
func (o *PramanaObject) PramanaId() string {
	return ""
}

// PramanaHashUrl returns the entity URL using the hashed UUID.
func (o *PramanaObject) PramanaHashUrl() string {
	return fmt.Sprintf("https://pramana.dev/entity/%s", o.pramanaGuid)
}

// PramanaUrl returns the Pramana entity URL. For regular objects this equals PramanaHashUrl.
func (o *PramanaObject) PramanaUrl() string {
	return o.PramanaHashUrl()
}

// GenerateId assigns a new UUID v4 to this object. Returns a PramanaError
// if the object already has a non-empty ID — IDs are write-once by design.
func (o *PramanaObject) GenerateId() error {
	if o.pramanaGuid == uuid.Nil {
		o.pramanaGuid = uuid.New()
		return nil
	}
	return NewPramanaError("Cannot reassign a PramanaObject ID once it has been set.")
}

// GetRoles returns the PramanaRole instances this object fulfils.
// For a base PramanaObject, this returns an empty slice.
func (o *PramanaObject) GetRoles() []*PramanaRole {
	return nil
}

// Compile-time interface checks.
var _ PramanaLinkable = (*PramanaObject)(nil)
var _ PramanaRoleful = (*PramanaObject)(nil)
