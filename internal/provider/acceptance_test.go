package provider

import (
	"embed"
	"testing"
)

//go:embed testdata/*
var testdata embed.FS

func TestAcceptance(t *testing.T) {
	items, err := testdata.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read testdata directory: %s", err)
	}

	for _, item := range items {
		t.Run(item.Name(), func(t *testing.T) {
			t.Logf("testing %v", item)
			data, err := testdata.ReadFile("testdata/" + item.Name() + "/main.tf")
			if err != nil {
				t.Fatalf("failed to read %s: %s", item.Name(), err)
			}
			t.Logf("data: %s", data)

		})
	}
}
