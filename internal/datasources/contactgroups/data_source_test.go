package contactgroups

import (
	"context"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestDataSourceMetadata(t *testing.T) {
	t.Parallel()

	resp := &fwdatasource.MetadataResponse{}
	NewContactGroupsDataSource().Metadata(context.Background(), fwdatasource.MetadataRequest{ProviderTypeName: "nodeping"}, resp)

	if want := "nodeping_contactgroups"; resp.TypeName != want {
		t.Errorf("TypeName = %q, want %q", resp.TypeName, want)
	}
}

func TestDataSourceSchemaIsValid(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resp := &fwdatasource.SchemaResponse{}

	NewContactGroupsDataSource().Schema(ctx, fwdatasource.SchemaRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation failed: %+v", diags)
	}
}

func TestDataSourceSchemaAttributesAreUsable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resp := &fwdatasource.SchemaResponse{}
	NewContactGroupsDataSource().Schema(ctx, fwdatasource.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) == 0 && len(resp.Schema.Blocks) == 0 {
		t.Fatal("schema exposes neither attributes nor blocks")
	}

	for name, attr := range resp.Schema.Attributes {
		if !attr.IsRequired() && !attr.IsOptional() && !attr.IsComputed() {
			t.Errorf("attribute %q is neither Required, Optional nor Computed", name)
		}
	}
}
