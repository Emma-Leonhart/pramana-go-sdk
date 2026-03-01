package pramana

import (
	"testing"

	"github.com/google/uuid"
)

func TestPramanaRole_Constructor_SetsLabel(t *testing.T) {
	role := NewPramanaRole("Entity")
	if role.Label != "Entity" {
		t.Errorf("Label = %q, want %q", role.Label, "Entity")
	}
}

func TestPramanaRole_ConstructorWithID_SetsGuid(t *testing.T) {
	id := uuid.New()
	role := NewPramanaRoleWithID("Entity", id)
	if role.PramanaGuid() != id {
		t.Errorf("PramanaGuid = %s, want %s", role.PramanaGuid(), id)
	}
}

func TestPramanaRole_ConstructorWithoutID_HasNilGuid(t *testing.T) {
	role := NewPramanaRole("Entity")
	if role.PramanaGuid() != uuid.Nil {
		t.Errorf("expected uuid.Nil, got %s", role.PramanaGuid())
	}
}

func TestPramanaRole_IsPramanaObject(t *testing.T) {
	role := NewPramanaRole("Entity")
	// PramanaRole embeds PramanaObject, so it satisfies PramanaLinkable
	var _ PramanaLinkable = role
}

func TestPramanaRole_GetRoles_ReturnsSelf(t *testing.T) {
	role := NewPramanaRole("Entity")
	roles := role.GetRoles()
	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}
	if roles[0] != role {
		t.Error("GetRoles should return the role itself")
	}
}

func TestPramanaRole_ParentRoles_InitiallyEmpty(t *testing.T) {
	role := NewPramanaRole("Entity")
	if len(role.ParentRoles) != 0 {
		t.Errorf("expected 0 parent roles, got %d", len(role.ParentRoles))
	}
}

func TestPramanaRole_ChildRoles_InitiallyEmpty(t *testing.T) {
	role := NewPramanaRole("Entity")
	if len(role.ChildRoles) != 0 {
		t.Errorf("expected 0 child roles, got %d", len(role.ChildRoles))
	}
}

func TestPramanaRole_InstanceOf_DefaultsToNil(t *testing.T) {
	role := NewPramanaRole("Entity")
	if role.InstanceOf != nil {
		t.Error("InstanceOf should default to nil")
	}
}

func TestPramanaRole_SubclassOf_DefaultsToNil(t *testing.T) {
	role := NewPramanaRole("Entity")
	if role.SubclassOf != nil {
		t.Error("SubclassOf should default to nil")
	}
}

func TestPramanaRole_CanBuildHierarchy(t *testing.T) {
	parent := NewPramanaRole("Thing")
	child := NewPramanaRole("Person")
	parent.AddChild(child)

	if child.SubclassOf != parent {
		t.Error("child.SubclassOf should be parent")
	}
	if len(parent.ChildRoles) != 1 || parent.ChildRoles[0] != child {
		t.Error("parent.ChildRoles should contain child")
	}
	if len(child.ParentRoles) != 1 || child.ParentRoles[0] != parent {
		t.Error("child.ParentRoles should contain parent")
	}
}

func TestPramanaRole_InstanceOf_CanBeSet(t *testing.T) {
	classRole := NewPramanaRole("Class")
	instance := NewPramanaRole("MyClass")
	instance.InstanceOf = classRole
	if instance.InstanceOf != classRole {
		t.Error("InstanceOf should be classRole")
	}
}

func TestPramanaRole_GenerateId_Works(t *testing.T) {
	role := NewPramanaRole("Entity")
	err := role.GenerateId()
	if err != nil {
		t.Fatal(err)
	}
	if role.PramanaGuid() == uuid.Nil {
		t.Error("expected non-nil GUID after GenerateId")
	}
}
