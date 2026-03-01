package pramana

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestPramanaObject_DefaultConstructor_HasNilGuid(t *testing.T) {
	obj := NewPramanaObject()
	if obj.PramanaGuid() != uuid.Nil {
		t.Errorf("expected uuid.Nil, got %s", obj.PramanaGuid())
	}
}

func TestPramanaObject_ConstructorWithID_SetsGuid(t *testing.T) {
	id := uuid.New()
	obj := NewPramanaObjectWithID(id)
	if obj.PramanaGuid() != id {
		t.Errorf("expected %s, got %s", id, obj.PramanaGuid())
	}
}

func TestPramanaObject_GenerateId_AssignsNonNilGuid(t *testing.T) {
	obj := NewPramanaObject()
	err := obj.GenerateId()
	if err != nil {
		t.Fatal(err)
	}
	if obj.PramanaGuid() == uuid.Nil {
		t.Error("expected non-nil GUID after GenerateId")
	}
}

func TestPramanaObject_GenerateId_ErrorsOnSecondCall(t *testing.T) {
	obj := NewPramanaObject()
	_ = obj.GenerateId()
	err := obj.GenerateId()
	if err == nil {
		t.Error("expected error on second GenerateId call")
	}
}

func TestPramanaObject_GenerateId_ErrorsWhenConstructedWithID(t *testing.T) {
	obj := NewPramanaObjectWithID(uuid.New())
	err := obj.GenerateId()
	if err == nil {
		t.Error("expected error when calling GenerateId on object with existing ID")
	}
}

func TestPramanaObject_PramanaId_IsEmpty_ForRegularObject(t *testing.T) {
	obj := NewPramanaObject()
	if obj.PramanaId() != "" {
		t.Errorf("expected empty PramanaId, got %q", obj.PramanaId())
	}
}

func TestPramanaObject_PramanaHashUrl_ContainsGuid(t *testing.T) {
	id := uuid.New()
	obj := NewPramanaObjectWithID(id)
	expected := fmt.Sprintf("https://pramana.dev/entity/%s", id)
	if obj.PramanaHashUrl() != expected {
		t.Errorf("PramanaHashUrl = %q, want %q", obj.PramanaHashUrl(), expected)
	}
}

func TestPramanaObject_PramanaUrl_EqualsHashUrl(t *testing.T) {
	obj := NewPramanaObjectWithID(uuid.New())
	if obj.PramanaUrl() != obj.PramanaHashUrl() {
		t.Error("PramanaUrl should equal PramanaHashUrl for regular objects")
	}
}

func TestPramanaObject_ClassId_EqualsRootId(t *testing.T) {
	expected := uuid.MustParse("10000000-0000-4000-8000-000000000001")
	if PramanaObjectRootID != expected {
		t.Errorf("RootID = %s, want %s", PramanaObjectRootID, expected)
	}
}

func TestPramanaObject_ClassUrl_UsesClassId(t *testing.T) {
	expected := fmt.Sprintf("https://pramana.dev/entity/%s", PramanaObjectRootID)
	if PramanaObjectClassURL() != expected {
		t.Errorf("ClassURL = %q, want %q", PramanaObjectClassURL(), expected)
	}
}

func TestPramanaObject_GetRoles_ReturnsEmpty(t *testing.T) {
	obj := NewPramanaObject()
	roles := obj.GetRoles()
	if len(roles) != 0 {
		t.Errorf("expected 0 roles, got %d", len(roles))
	}
}

func TestPramanaObject_ImplementsInterfaces(t *testing.T) {
	obj := NewPramanaObject()
	var _ PramanaLinkable = obj
	var _ PramanaRoleful = obj
}
