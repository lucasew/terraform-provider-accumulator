package provider

import (
	"context"
	"embed"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

//go:embed testdata/*
var testdata embed.FS

func TestAccCases(t *testing.T) {
	items, err := testdata.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read testdata directory: %s", err)
	}

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		item := item
		t.Run(item.Name(), func(t *testing.T) {
			data, err := testdata.ReadFile("testdata/" + item.Name() + "/main.tf")
			if err != nil {
				t.Fatalf("failed to read %s: %s", item.Name(), err)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config: string(data),
						ConfigStateChecks: []statecheck.StateCheck{
							expectOutputEqual("expected", "got"),
						},
					},
				},
			})
		})
	}
}

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"accumulator": providerserver.NewProtocol6WithError(New("test")()),
	}
}

func expectOutputEqual(expectedName, gotName string) statecheck.StateCheck {
	return outputEqualStateCheck{
		expectedName: expectedName,
		gotName:      gotName,
	}
}

type outputEqualStateCheck struct {
	expectedName string
	gotName      string
}

func (c outputEqualStateCheck) CheckState(_ context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	expectedOutput, err := lookupStateOutput(req.State, c.expectedName)
	if err != nil {
		resp.Error = err
		return
	}

	gotOutput, err := lookupStateOutput(req.State, c.gotName)
	if err != nil {
		resp.Error = err
		return
	}

	if !expectedOutput.Type.Equals(gotOutput.Type) {
		resp.Error = fmt.Errorf(
			"output %q type mismatch: expected %s, got %s",
			c.gotName,
			expectedOutput.Type.FriendlyName(),
			gotOutput.Type.FriendlyName(),
		)
		return
	}

	if diff := cmp.Diff(expectedOutput.Value, gotOutput.Value); diff != "" {
		resp.Error = fmt.Errorf("output %q mismatch (-expected +got):\n%s", c.gotName, diff)
	}
}

func lookupStateOutput(state *tfjson.State, name string) (*tfjson.StateOutput, error) {
	if state == nil || state.Values == nil {
		return nil, fmt.Errorf("terraform state is empty")
	}

	output, ok := state.Values.Outputs[name]
	if !ok {
		return nil, fmt.Errorf("output %q not found in terraform state", name)
	}

	return output, nil
}
