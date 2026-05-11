package provider

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDynamicObjectRoundTrip(t *testing.T) {
	input := map[string]any{
		"string": "value",
		"bool":   true,
		"number": 42,
	}

	encoded := dynamicObjectFromGo(input)
	got, err := goMapFromDynamic(encoded)
	if err != nil {
		t.Fatalf("goMapFromDynamic returned error: %v", err)
	}

	if got["string"] != "value" {
		t.Fatalf("expected string round trip, got %#v", got)
	}

	if got["bool"] != true {
		t.Fatalf("expected bool round trip, got %#v", got)
	}

	number, ok := got["number"].(*big.Float)
	if !ok {
		t.Fatalf("expected number to round trip as *big.Float, got %T", got["number"])
	}

	if number.Cmp(big.NewFloat(42)) != 0 {
		t.Fatalf("expected number 42, got %s", number.String())
	}
}

func TestGoValueFromDynamicString(t *testing.T) {
	got, err := goValueFromDynamic(types.DynamicValue(types.StringValue("value")))
	if err != nil {
		t.Fatalf("goValueFromDynamic returned error: %v", err)
	}

	if got != "value" {
		t.Fatalf("expected value, got %#v", got)
	}
}
