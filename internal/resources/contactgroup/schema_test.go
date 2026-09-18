package contactgroup

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestContactGroupSchemaIsValid(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resp := &fwresource.SchemaResponse{}

	NewContactGroupResource().Schema(ctx, fwresource.SchemaRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation failed: %+v", diags)
	}
}

func TestContactGroupSchemaAttributes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resp := &fwresource.SchemaResponse{}
	NewContactGroupResource().Schema(ctx, fwresource.SchemaRequest{}, resp)

	attrs := resp.Schema.Attributes

	// id and customer_id are assigned by NodePing and must never be
	// user-settable.
	for _, name := range []string{"id", "customer_id"} {
		attr, ok := attrs[name]
		if !ok {
			t.Errorf("schema is missing the %q attribute", name)
			continue
		}
		if !attr.IsComputed() {
			t.Errorf("%q must be Computed", name)
		}
		if attr.IsRequired() {
			t.Errorf("%q must not be Required", name)
		}
	}

	// Both name and members are optional: the API accepts a group with
	// neither.
	for _, name := range []string{"name", "members"} {
		attr, ok := attrs[name]
		if !ok {
			t.Errorf("schema is missing the %q attribute", name)
			continue
		}
		if !attr.IsOptional() {
			t.Errorf("%q must be Optional", name)
		}
	}
}

func TestContactGroupSchemaAttributesAreUsable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resp := &fwresource.SchemaResponse{}
	NewContactGroupResource().Schema(ctx, fwresource.SchemaRequest{}, resp)

	for name, attr := range resp.Schema.Attributes {
		if !attr.IsRequired() && !attr.IsOptional() && !attr.IsComputed() {
			t.Errorf("attribute %q is neither Required, Optional nor Computed", name)
		}
		if attr.GetDescription() == "" && attr.GetMarkdownDescription() == "" {
			t.Errorf("attribute %q has no description", name)
		}
	}
}

func TestNewContactGroupResourceMetadata(t *testing.T) {
	t.Parallel()

	resp := &fwresource.MetadataResponse{}
	NewContactGroupResource().Metadata(context.Background(), fwresource.MetadataRequest{ProviderTypeName: "nodeping"}, resp)

	if want := "nodeping_contactgroup"; resp.TypeName != want {
		t.Errorf("TypeName = %q, want %q", resp.TypeName, want)
	}
}
