package contactgroup

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func stringList(t *testing.T, vals ...string) types.List {
	t.Helper()
	// A variadic call with no arguments yields a nil slice, and ListValueFrom
	// turns nil into a *null* list -- which would silently make the
	// "empty list" cases below test the null path instead.
	if vals == nil {
		vals = []string{}
	}
	list, diags := types.ListValueFrom(context.Background(), types.StringType, vals)
	if diags.HasError() {
		t.Fatalf("failed to build list: %+v", diags)
	}
	return list
}

func TestMembersFromPlan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   types.List
		want []string
	}{
		{
			name: "values are read through",
			in:   stringList(t, "AD1", "AD2"),
			want: []string{"AD1", "AD2"},
		},
		{
			// An omitted attribute must stay nil. Turning it into an empty
			// slice would send "remove every member" on create.
			name: "null stays nil",
			in:   types.ListNull(types.StringType),
			want: nil,
		},
		{
			name: "unknown stays nil",
			in:   types.ListUnknown(types.StringType),
			want: nil,
		},
		{
			name: "explicitly empty list is empty, not nil",
			in:   stringList(t),
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics
			got := membersFromPlan(context.Background(), tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %+v", diags)
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("expected nil, got %#v", got)
				}
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("index %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// The API omits `members` for an empty group. Whether that maps back to null or
// to an empty list depends on what the configuration asked for -- getting this
// wrong produces a permanent diff.
func TestMembersToModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		api       []string
		planned   types.List
		wantNull  bool
		wantElems []string
	}{
		{
			name:      "members are mapped through",
			api:       []string{"AD1", "AD2"},
			planned:   stringList(t, "AD1", "AD2"),
			wantElems: []string{"AD1", "AD2"},
		},
		{
			name:     "no members and none configured stays null",
			api:      nil,
			planned:  types.ListNull(types.StringType),
			wantNull: true,
		},
		{
			// Config said [], API returned nothing: that is agreement, so it
			// must come back as [] rather than null.
			name:      "no members but an empty list was configured",
			api:       nil,
			planned:   stringList(t),
			wantElems: []string{},
		},
		{
			name:      "API members win over a configured empty list",
			api:       []string{"AD9"},
			planned:   stringList(t),
			wantElems: []string{"AD9"},
		},
		{
			name:     "unknown planned value with no API members stays null",
			api:      nil,
			planned:  types.ListUnknown(types.StringType),
			wantNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics
			got := membersToModel(context.Background(), tt.api, tt.planned, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %+v", diags)
			}
			if tt.wantNull {
				if !got.IsNull() {
					t.Errorf("expected null list, got %v", got)
				}
				return
			}
			if got.IsNull() {
				t.Fatalf("expected a list, got null")
			}

			var elems []string
			if d := got.ElementsAs(context.Background(), &elems, false); d.HasError() {
				t.Fatalf("failed to read list: %+v", d)
			}
			if len(elems) != len(tt.wantElems) {
				t.Fatalf("got %#v, want %#v", elems, tt.wantElems)
			}
			for i := range tt.wantElems {
				if elems[i] != tt.wantElems[i] {
					t.Errorf("index %d: got %q, want %q", i, elems[i], tt.wantElems[i])
				}
			}
		})
	}
}
