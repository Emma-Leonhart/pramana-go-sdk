package pramana

import "github.com/google/uuid"

// PramanaRole represents a role (interface) in the Pramana ontology.
// Roles form a hierarchy via SubclassOf and InstanceOf, and track their
// position in the role graph through ParentRoles and ChildRoles.
type PramanaRole struct {
	PramanaObject

	// Label is the human-readable name for this role.
	Label string

	// InstanceOf is the role that this role is an instance of.
	InstanceOf *PramanaRole

	// SubclassOf is the role that this role is a subclass of.
	SubclassOf *PramanaRole

	// ParentRoles are the parent roles of this role in the hierarchy.
	ParentRoles []*PramanaRole

	// ChildRoles are the child roles of this role in the hierarchy.
	ChildRoles []*PramanaRole
}

// NewPramanaRole creates a new PramanaRole with the given label and no ID.
func NewPramanaRole(label string) *PramanaRole {
	return &PramanaRole{
		PramanaObject: *NewPramanaObject(),
		Label:         label,
	}
}

// NewPramanaRoleWithID creates a new PramanaRole with the given label and UUID.
func NewPramanaRoleWithID(label string, id uuid.UUID) *PramanaRole {
	return &PramanaRole{
		PramanaObject: *NewPramanaObjectWithID(id),
		Label:         label,
	}
}

// GetRoles returns this role itself as the list of roles it fulfils.
func (r *PramanaRole) GetRoles() []*PramanaRole {
	return []*PramanaRole{r}
}

// AddChild adds a child role to this role's ChildRoles and sets up
// the bidirectional relationship (child.SubclassOf and child.ParentRoles).
func (r *PramanaRole) AddChild(child *PramanaRole) {
	child.SubclassOf = r
	child.ParentRoles = append(child.ParentRoles, r)
	r.ChildRoles = append(r.ChildRoles, child)
}
