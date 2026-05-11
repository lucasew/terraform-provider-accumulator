# Accumulator MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a minimal `accumulator_group` + `accumulator_item` provider that materializes a computed map, rejects duplicate keys, and supports optional group typing, with tests runnable locally via `go test`.

**Architecture:** Keep Terraform-facing code thin and move aggregation behavior into a small in-memory store with explicit methods for creating groups, inserting items, deleting items, deleting groups, and rendering the final map. Resource implementations should use that store and expose stable Terraform schemas. The first milestone avoids acceptance tests and any dependency on external services or Terraform CLI downloads; verification should come from pure Go unit tests and focused resource-level tests.

**Tech Stack:** Go 1.25, Terraform Plugin Framework, `testing`, `go test`

---

### Task 1: Establish isolation boundaries

**Files:**
- Modify: `internal/provider/provider.go`
- Modify: `internal/provider/model.go`
- Create: `internal/provider/store.go`
- Test: `internal/provider/store_test.go`

- [ ] **Step 1: Write the failing store test for empty group creation**

```go
package provider

import "testing"

func TestStoreCreateGroupStartsEmpty(t *testing.T) {
	store := NewStore()

	groupID, err := store.CreateGroup("example", "")
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/provider -run TestStoreCreateGroupStartsEmpty -v`
Expected: FAIL with undefined `NewStore`, `CreateGroup`, or `GroupData`

- [ ] **Step 3: Write minimal store implementation**

```go
package provider

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Store struct {
	mu     sync.Mutex
	groups map[uuid.UUID]*GroupRecord
}

type GroupRecord struct {
	Name  string
	Type  string
	Items map[string]any
}

func NewStore() *Store {
	return &Store{
		groups: map[uuid.UUID]*GroupRecord{},
	}
}

func (s *Store) CreateGroup(name string, typeName string) (uuid.UUID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New()
	s.groups[id] = &GroupRecord{
		Name:  name,
		Type:  typeName,
		Items: map[string]any{},
	}

	return id, nil
}

func (s *Store) GroupData(id uuid.UUID) (map[string]any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	group, ok := s.groups[id]
	if !ok {
		return nil, fmt.Errorf("group %s not found", id)
	}

	out := make(map[string]any, len(group.Items))
	for k, v := range group.Items {
		out[k] = v
	}

	return out, nil
}
```

- [ ] **Step 4: Wire provider-scoped store into resources**

```go
type AccumulatorProvider struct {
	version string
	store   *Store
}

func (p *AccumulatorProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
	if p.store == nil {
		p.store = NewStore()
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AccumulatorProvider{
			version: version,
			store:   NewStore(),
		}
	}
}
```

```go
type GroupResource struct {
	store *Store
}

type ItemResource struct {
	store *Store
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./internal/provider -run TestStoreCreateGroupStartsEmpty -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/provider/provider.go internal/provider/model.go internal/provider/store.go internal/provider/store_test.go
git commit -m "refactor: add provider-scoped accumulator store"
```

### Task 2: Define stable resource models and schemas

**Files:**
- Modify: `internal/provider/model.go`
- Create: `internal/provider/types.go`
- Test: `internal/provider/model_test.go`

- [ ] **Step 1: Write the failing schema/model test**

```go
package provider

import "testing"

func TestGroupModelIncludesComputedDataAndOptionalType(t *testing.T) {
	var model groupModel

	if modelName := model.SchemaAttributeNames(); len(modelName) == 0 {
		t.Fatal("expected schema attribute names helper to exist")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/provider -run TestGroupModelIncludesComputedDataAndOptionalType -v`
Expected: FAIL with undefined `groupModel` or helper methods

- [ ] **Step 3: Replace ad hoc state handling with explicit models**

```go
type groupModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	Data types.Map    `tfsdk:"data"`
}

type itemModel struct {
	ID    types.String `tfsdk:"id"`
	Group types.String `tfsdk:"group"`
	Key   types.String `tfsdk:"key"`
	Value types.Dynamic `tfsdk:"value"`
}
```

```go
func (groupModel) SchemaAttributeNames() []string {
	return []string{"id", "name", "type", "data"}
}
```

- [ ] **Step 4: Update resource schemas to match the idea**

```go
resp.Schema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Stable group identifier.",
		},
		"name": schema.StringAttribute{
			Required:    true,
			Description: "Human-readable group name.",
		},
		"type": schema.StringAttribute{
			Optional:    true,
			Description: "Optional value type enforced for all items in the group.",
		},
		"data": schema.MapAttribute{
			Computed:    true,
			ElementType: types.DynamicType,
			Description: "Computed map assembled from all items in the group.",
		},
	},
}
```

```go
resp.Schema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"group": schema.StringAttribute{
			Required: true,
		},
		"key": schema.StringAttribute{
			Required: true,
		},
		"value": schema.DynamicAttribute{
			Required: true,
		},
	},
}
```

- [ ] **Step 5: Run the focused model test**

Run: `go test ./internal/provider -run TestGroupModelIncludesComputedDataAndOptionalType -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/provider/model.go internal/provider/types.go internal/provider/model_test.go
git commit -m "feat: define accumulator resource schemas"
```

### Task 3: Implement core store behaviors with tests first

**Files:**
- Modify: `internal/provider/store.go`
- Modify: `internal/provider/store_test.go`

- [ ] **Step 1: Add failing tests for insertion, duplicates, deletion, and missing groups**

```go
func TestStorePutItemAddsEntryToGroup(t *testing.T) {
	store := NewStore()
	groupID, _ := store.CreateGroup("example", "")

	if err := store.PutItem(groupID, "item", "value"); err != nil {
		t.Fatalf("PutItem returned error: %v", err)
	}

	got, _ := store.GroupData(groupID)
	if got["item"] != "value" {
		t.Fatalf("expected item=value, got %#v", got)
	}
}

func TestStorePutItemRejectsDuplicateKeys(t *testing.T) {
	store := NewStore()
	groupID, _ := store.CreateGroup("example", "")

	_ = store.PutItem(groupID, "item", "first")
	err := store.PutItem(groupID, "item", "second")
	if err == nil {
		t.Fatal("expected duplicate key error")
	}
}

func TestStoreDeleteItemRemovesEntry(t *testing.T) {
	store := NewStore()
	groupID, _ := store.CreateGroup("example", "")
	_ = store.PutItem(groupID, "item", "value")

	if err := store.DeleteItem(groupID, "item"); err != nil {
		t.Fatalf("DeleteItem returned error: %v", err)
	}

	got, _ := store.GroupData(groupID)
	if _, ok := got["item"]; ok {
		t.Fatalf("expected item to be removed, got %#v", got)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/provider -run 'TestStore(PutItemAddsEntryToGroup|PutItemRejectsDuplicateKeys|DeleteItemRemovesEntry)' -v`
Expected: FAIL with undefined methods or wrong behavior

- [ ] **Step 3: Implement the minimal store API**

```go
func (s *Store) PutItem(id uuid.UUID, key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	group, ok := s.groups[id]
	if !ok {
		return fmt.Errorf("group %s not found", id)
	}

	if _, exists := group.Items[key]; exists {
		return fmt.Errorf("duplicate key %q in group %s", key, id)
	}

	group.Items[key] = value
	return nil
}

func (s *Store) DeleteItem(id uuid.UUID, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	group, ok := s.groups[id]
	if !ok {
		return fmt.Errorf("group %s not found", id)
	}

	delete(group.Items, key)
	return nil
}

func (s *Store) DeleteGroup(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.groups, id)
}
```

- [ ] **Step 4: Run the store test suite**

Run: `go test ./internal/provider -run TestStore -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/provider/store.go internal/provider/store_test.go
git commit -m "feat: implement accumulator store behaviors"
```

### Task 4: Add optional type enforcement in the store

**Files:**
- Modify: `internal/provider/store.go`
- Modify: `internal/provider/store_test.go`

- [ ] **Step 1: Write failing tests for typed and untyped groups**

```go
func TestStorePutItemAllowsAnyTypeWhenGroupTypeIsEmpty(t *testing.T) {
	store := NewStore()
	groupID, _ := store.CreateGroup("example", "")

	if err := store.PutItem(groupID, "string", "value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := store.PutItem(groupID, "number", 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStorePutItemRejectsWrongTypeWhenGroupTypeIsString(t *testing.T) {
	store := NewStore()
	groupID, _ := store.CreateGroup("example", "string")

	err := store.PutItem(groupID, "number", 42)
	if err == nil {
		t.Fatal("expected type mismatch error")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/provider -run 'TestStorePutItem(AllowsAnyTypeWhenGroupTypeIsEmpty|RejectsWrongTypeWhenGroupTypeIsString)' -v`
Expected: FAIL because type validation does not exist

- [ ] **Step 3: Implement a deliberately small type system**

```go
func validateValue(typeName string, value any) error {
	switch typeName {
	case "":
		return nil
	case "string":
		if _, ok := value.(string); ok {
			return nil
		}
	case "number":
		switch value.(type) {
		case int, int64, float64:
			return nil
		}
	case "bool":
		if _, ok := value.(bool); ok {
			return nil
		}
	}

	return fmt.Errorf("value does not conform to type %q", typeName)
}
```

```go
if err := validateValue(group.Type, value); err != nil {
	return err
}
```

- [ ] **Step 4: Run the store test suite**

Run: `go test ./internal/provider -run TestStore -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/provider/store.go internal/provider/store_test.go
git commit -m "feat: enforce optional group typing"
```

### Task 5: Implement `accumulator_group` resource lifecycle

**Files:**
- Modify: `internal/provider/model.go`
- Create: `internal/provider/group_resource_test.go`

- [ ] **Step 1: Write a failing resource test for group create/read/delete behavior**

```go
func TestGroupResourceCreateStoresComputedData(t *testing.T) {
	t.Skip("replace with framework-backed create/read lifecycle test")
}
```

- [ ] **Step 2: Run test to verify it fails meaningfully**

Run: `go test ./internal/provider -run TestGroupResourceCreateStoresComputedData -v`
Expected: FAIL or SKIP indicating missing resource lifecycle coverage

- [ ] **Step 3: Implement group resource methods against the store**

```go
func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := r.store.CreateGroup(plan.Name.ValueString(), plan.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("creating group failed", err.Error())
		return
	}

	plan.ID = types.StringValue(groupID.String())
	plan.Data = types.MapNull(types.DynamicType)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
```

```go
func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	data, err := r.store.GroupData(groupID)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Data = dynamicMapFromGo(data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

```go
func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	r.store.DeleteGroup(groupID)
}
```

- [ ] **Step 4: Run package tests**

Run: `go test ./internal/provider -v`
Expected: PASS except any intentionally skipped lifecycle tests

- [ ] **Step 5: Commit**

```bash
git add internal/provider/model.go internal/provider/group_resource_test.go
git commit -m "feat: implement accumulator group resource"
```

### Task 6: Implement `accumulator_item` resource lifecycle

**Files:**
- Modify: `internal/provider/model.go`
- Create: `internal/provider/item_resource_test.go`

- [ ] **Step 1: Write failing tests for item insertion, duplicate rejection, and deletion**

```go
func TestItemResourceCreateAddsValueToGroup(t *testing.T) {
	t.Skip("replace with focused resource method tests around store interactions")
}

func TestItemResourceCreateRejectsDuplicateKey(t *testing.T) {
	t.Skip("replace with focused resource method tests around diagnostics")
}
```

- [ ] **Step 2: Run tests to verify they capture the gap**

Run: `go test ./internal/provider -run 'TestItemResource(CreateAddsValueToGroup|CreateRejectsDuplicateKey)' -v`
Expected: FAIL or SKIP indicating missing item lifecycle coverage

- [ ] **Step 3: Implement item resource methods**

```go
func (r *ItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan itemModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(plan.Group.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	value, err := goValueFromDynamic(plan.Value)
	if err != nil {
		resp.Diagnostics.AddError("invalid item value", err.Error())
		return
	}

	if err := r.store.PutItem(groupID, plan.Key.ValueString(), value); err != nil {
		resp.Diagnostics.AddError("creating item failed", err.Error())
		return
	}

	plan.ID = types.StringValue(groupID.String() + ":" + plan.Key.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
```

```go
func (r *ItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state itemModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(state.Group.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	if err := r.store.DeleteItem(groupID, state.Key.ValueString()); err != nil {
		resp.Diagnostics.AddError("deleting item failed", err.Error())
	}
}
```

- [ ] **Step 4: Run package tests**

Run: `go test ./internal/provider -v`
Expected: PASS except any intentionally skipped lifecycle tests

- [ ] **Step 5: Commit**

```bash
git add internal/provider/model.go internal/provider/item_resource_test.go
git commit -m "feat: implement accumulator item resource"
```

### Task 7: Add value conversion helpers and end-to-end local smoke coverage

**Files:**
- Create: `internal/provider/dynamic.go`
- Create: `internal/provider/dynamic_test.go`
- Modify: `examples/teste.tf`

- [ ] **Step 1: Write failing conversion tests**

```go
func TestDynamicRoundTripString(t *testing.T) {
	input := map[string]any{"item": "value"}

	encoded := dynamicMapFromGo(input)
	got, err := goMapFromDynamic(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["item"] != "value" {
		t.Fatalf("expected round trip value, got %#v", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/provider -run TestDynamicRoundTripString -v`
Expected: FAIL with undefined helpers

- [ ] **Step 3: Implement minimal supported conversions**

```go
func dynamicMapFromGo(input map[string]any) types.Map {
	if len(input) == 0 {
		return types.MapValueMust(types.DynamicType, map[string]attr.Value{})
	}

	values := make(map[string]attr.Value, len(input))
	for k, v := range input {
		values[k] = dynamicValueFromGo(v)
	}

	return types.MapValueMust(types.DynamicType, values)
}
```

```go
func dynamicValueFromGo(v any) attr.Value {
	switch x := v.(type) {
	case string:
		return types.DynamicValue(types.StringValue(x))
	case bool:
		return types.DynamicValue(types.BoolValue(x))
	case int:
		return types.DynamicValue(types.NumberValue(big.NewFloat(float64(x))))
	case float64:
		return types.DynamicValue(types.NumberValue(big.NewFloat(x)))
	default:
		return types.DynamicNull()
	}
}
```

- [ ] **Step 4: Update the example to reflect the supported MVP contract**

```terraform
resource "accumulator_group" "example" {
  name = "example-accumulator"
  type = "string"
}

resource "accumulator_item" "item" {
  group = accumulator_group.example.id
  key   = "item"
  value = "example value"
}
```

- [ ] **Step 5: Run the full local test suite**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/provider/dynamic.go internal/provider/dynamic_test.go examples/teste.tf
git commit -m "test: add local dynamic conversion coverage"
```

### Task 8: Clean up provider identity and docs after behavior is real

**Files:**
- Modify: `main.go`
- Modify: `README.md`
- Modify: `docs/resources/example.md`
- Modify: `docs/index.md`

- [ ] **Step 1: Write a small failing metadata test**

```go
func TestProviderMetadataTypeNameIsAccumulator(t *testing.T) {
	provider := New("test")()
	if provider == nil {
		t.Fatal("expected provider instance")
	}
}
```

- [ ] **Step 2: Run test to verify the package still compiles around provider identity**

Run: `go test ./internal/provider -run TestProviderMetadataTypeNameIsAccumulator -v`
Expected: PASS or minimal edit needed for provider metadata coverage

- [ ] **Step 3: Update identity and human-facing docs**

```go
opts := providerserver.ServeOpts{
	Address: "registry.terraform.io/lucasew/accumulator",
	Debug:   debug,
}
```

```md
## Supported MVP behavior

- `accumulator_group` exposes computed `data`
- `accumulator_item` contributes one key/value pair
- duplicate keys fail
- supported typed groups: `string`, `number`, `bool`
```

- [ ] **Step 4: Run the full local test suite again**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add main.go README.md docs/resources/example.md docs/index.md
git commit -m "docs: align provider identity and mvp behavior"
```

### Task 9: Define the safe verification envelope for future work

**Files:**
- Modify: `README.md`
- Modify: `GNUmakefile`

- [ ] **Step 1: Add a documented local-only verification target**

```make
.PHONY: test-unit
test-unit:
	go test ./...
```

- [ ] **Step 2: Add a separate explicitly unsafe acceptance target**

```make
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v
```

- [ ] **Step 3: Document the rule in README**

```md
## Testing policy

- `make test-unit` is the default verification path and must remain local-only.
- `make testacc` is reserved for explicit acceptance coverage and should not be used for routine iteration.
- The provider should be designed so most behavior is covered by deterministic unit tests.
```

- [ ] **Step 4: Run the local verification target**

Run: `make test-unit`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add README.md GNUmakefile
git commit -m "build: define safe local verification targets"
```
