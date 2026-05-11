package provider

import "testing"

func TestStoreGroupDataStartsEmpty(t *testing.T) {
	store := NewStore()

	groupID, err := store.CreateGroup()
	if err != nil {
		t.Fatalf("CreateGroup returned error: %v", err)
	}

	got, err := store.GroupData(groupID)
	if err != nil {
		t.Fatalf("GroupData returned error: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected empty group data, got %#v", got)
	}
}

func TestStorePutItemAddsEntryToGroup(t *testing.T) {
	store := NewStore()

	groupID, _ := store.CreateGroup()
	if err := store.PutItem(groupID, "item", "value"); err != nil {
		t.Fatalf("PutItem returned error: %v", err)
	}

	got, err := store.GroupData(groupID)
	if err != nil {
		t.Fatalf("GroupData returned error: %v", err)
	}

	if got["item"] != "value" {
		t.Fatalf("expected item=value, got %#v", got)
	}
}

func TestStorePutItemRejectsDuplicateKeys(t *testing.T) {
	store := NewStore()

	groupID, _ := store.CreateGroup()
	_ = store.PutItem(groupID, "item", "first")

	err := store.PutItem(groupID, "item", "second")
	if err == nil {
		t.Fatal("expected duplicate key error")
	}
}

func TestStoreDeleteItemRemovesEntry(t *testing.T) {
	store := NewStore()

	groupID, _ := store.CreateGroup()
	_ = store.PutItem(groupID, "item", "value")

	if err := store.DeleteItem(groupID, "item"); err != nil {
		t.Fatalf("DeleteItem returned error: %v", err)
	}

	got, err := store.GroupData(groupID)
	if err != nil {
		t.Fatalf("GroupData returned error: %v", err)
	}

	if _, ok := got["item"]; ok {
		t.Fatalf("expected item to be removed, got %#v", got)
	}
}

func TestStoreDeleteGroupRemovesGroup(t *testing.T) {
	store := NewStore()

	groupID, _ := store.CreateGroup()
	if err := store.DeleteGroup(groupID); err != nil {
		t.Fatalf("DeleteGroup returned error: %v", err)
	}

	if _, err := store.GroupData(groupID); err == nil {
		t.Fatal("expected missing group error")
	}
}
