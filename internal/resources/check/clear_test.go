package check

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nodeping/terraform-provider-nodeping/internal/client"
	"github.com/nodeping/terraform-provider-nodeping/internal/datasources/checkattr"
)

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    attr.Value
		want bool
	}{
		{"null string", types.StringNull(), true},
		{"empty string", types.StringValue(""), true},
		{"string", types.StringValue("x"), false},
		{"unknown is not empty", types.StringUnknown(), false},
		{"false", types.BoolValue(false), true},
		{"true", types.BoolValue(true), false},
		// 0 is a real value for some numbers, e.g. volumemin.
		{"zero number", types.Int64Value(0), false},
		{"null number", types.Int64Null(), true},
		{"empty list", types.ListValueMust(types.StringType, []attr.Value{}), true},
		{"list", types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a")}), false},
		{"empty map", types.MapValueMust(types.StringType, map[string]attr.Value{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isEmpty(tt.v); got != tt.want {
				t.Errorf("isEmpty(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func emptyModel() CheckResourceModel {
	return CheckResourceModel{
		ContentString:  types.StringNull(),
		Method:         types.StringNull(),
		PostData:       types.StringNull(),
		ServerName:     types.StringNull(),
		Description:    types.StringNull(),
		StatusCode:     types.Int64Null(),
		WarningDays:    types.Int64Null(),
		Regex:          types.BoolNull(),
		Invert:         types.BoolNull(),
		Follow:         types.BoolNull(),
		IPv6:           types.BoolNull(),
		Dep:            types.StringNull(),
		RunLocations:   types.ListNull(types.StringType),
		SendHeaders:    types.MapNull(types.StringType),
		ReceiveHeaders: types.MapNull(types.StringType),
		Public:         types.BoolNull(),
	}
}

// Every attribute the prior state had and the plan dropped gets its clearing
// value; nothing else is touched.
func TestClearRemovedSendsClearingValues(t *testing.T) {
	t.Parallel()

	state := emptyModel()
	state.ContentString = types.StringValue("ok")
	state.Method = types.StringValue("POST")
	state.PostData = types.StringValue("body")
	state.ServerName = types.StringValue("sni.example.com")
	state.Description = types.StringValue("was here")
	state.StatusCode = types.Int64Value(201)
	state.WarningDays = types.Int64Value(30)
	state.Regex = types.BoolValue(true)
	state.Invert = types.BoolValue(true)
	state.Follow = types.BoolValue(true)
	state.IPv6 = types.BoolValue(true)
	state.Dep = types.StringValue("201205050153W2Q4C-0J2HSIRF")
	state.RunLocations = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("nam")})
	plan := emptyModel()

	var req client.CheckCreateRequest
	clearRemoved(&req, &plan, &state)

	for name, got := range map[string]*string{
		"contentstring": req.ContentString, "method": req.Method,
		"postdata": req.PostData, "servername": req.ServerName,
	} {
		if got == nil || *got != "" {
			t.Errorf("%s: want a pointer to \"\", got %v", name, got)
		}
	}
	if req.Description != checkattr.ClearedDescription {
		t.Errorf("description = %q, want %q", req.Description, checkattr.ClearedDescription)
	}
	for name, got := range map[string]interface{}{"statuscode": req.StatusCode, "warningdays": req.WarningDays} {
		if got != "" {
			t.Errorf("%s = %#v, want \"\"", name, got)
		}
	}
	for name, got := range map[string]interface{}{
		"regex": req.Regex, "invert": req.Invert, "follow": req.Follow, "ipv6": req.IPv6, "dep": req.Dep,
	} {
		if got != false {
			t.Errorf("%s = %#v, want false", name, got)
		}
	}
	if rl, ok := req.RunLocations.([]string); !ok || len(rl) != 0 {
		t.Errorf("runlocations = %#v, want []string{}", req.RunLocations)
	}
}

// Something the check never had must never be written -- clearing it would
// be a pointless write, and for description one that changes the check.
func TestClearRemovedLeavesNeverSetValuesOut(t *testing.T) {
	t.Parallel()

	state, plan := emptyModel(), emptyModel()
	var req client.CheckCreateRequest
	clearRemoved(&req, &plan, &state)

	if req.ContentString != nil || req.Method != nil || req.PostData != nil || req.ServerName != nil {
		t.Error("a string the check never had was written")
	}
	if req.Description != "" {
		t.Errorf("description = %q, want it left out", req.Description)
	}
	if req.StatusCode != nil || req.Regex != nil || req.Dep != nil || req.RunLocations != nil {
		t.Error("a value the check never had was written")
	}
	if req.SendHeaders != nil || req.ReceiveHeaders != nil {
		t.Error("headers were written for a check that had none")
	}
}

// A description already cleared is held as a single space and must not be
// cleared again on every update.
func TestClearRemovedDoesNotReclearADescription(t *testing.T) {
	t.Parallel()

	state, plan := emptyModel(), emptyModel()
	state.Description = types.StringValue(checkattr.ClearedDescription)
	var req client.CheckCreateRequest
	clearRemoved(&req, &plan, &state)

	if req.Description != "" {
		t.Errorf("description = %q, want it left out", req.Description)
	}
}

// NodePing merges headers per name: only a header sent as null is deleted.
func TestDeleteRemovedKeys(t *testing.T) {
	t.Parallel()

	prior := types.MapValueMust(types.StringType, map[string]attr.Value{
		"keep": types.StringValue("1"), "drop": types.StringValue("2"),
	})
	planned := types.MapValueMust(types.StringType, map[string]attr.Value{
		"keep": types.StringValue("1"),
	})
	keep := "1"
	sent := map[string]*string{"keep": &keep}

	got := deleteRemovedKeys(sent, prior, planned)

	if v, ok := got["drop"]; !ok || v != nil {
		t.Errorf("drop = %v (present=%v), want it sent as null", v, ok)
	}
	if v := got["keep"]; v == nil || *v != "1" {
		t.Errorf("keep = %v, want \"1\"", v)
	}

	// Removing every header sends each one as null, not an empty map.
	all := deleteRemovedKeys(nil, prior, types.MapNull(types.StringType))
	if len(all) != 2 || all["keep"] != nil || all["drop"] != nil {
		t.Errorf("removing all headers = %#v, want both sent as null", all)
	}
}

func TestKeepClearedBooleans(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      types.Bool
		expected types.Bool
		want     types.Bool
	}{
		// NodePing keeps a cleared boolean as false; the plan says null.
		{"cleared reads back false, null stands", types.BoolValue(false), types.BoolNull(), types.BoolNull()},
		{"explicit false stays false", types.BoolValue(false), types.BoolValue(false), types.BoolValue(false)},
		// A true from NodePing is a real change and must show.
		{"true from NodePing is kept", types.BoolValue(true), types.BoolNull(), types.BoolValue(true)},
		{"absent stays absent", types.BoolNull(), types.BoolNull(), types.BoolNull()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			model := emptyModel()
			model.Regex = tt.got
			expected := emptyModel()
			expected.Regex = tt.expected

			keepClearedBooleans(&model, &expected)

			if !model.Regex.Equal(tt.want) {
				t.Errorf("regex = %v, want %v", model.Regex, tt.want)
			}
		})
	}
}

func TestPublicForUpdateSendsAString(t *testing.T) {
	t.Parallel()

	for _, v := range []bool{true, false} {
		plan := emptyModel()
		plan.Public = types.BoolValue(v)
		got := publicForUpdate(&plan)
		want := map[bool]string{true: "true", false: "false"}[v]
		if got != want {
			t.Errorf("public=%v: got %#v, want the string %q", v, got, want)
		}
	}
	plan := emptyModel()
	if got := publicForUpdate(&plan); got != nil {
		t.Errorf("unset public: got %#v, want nil", got)
	}
}

func TestHeadersToAPI(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics
	in := types.MapValueMust(types.StringType, map[string]attr.Value{"a": types.StringValue("1")})
	got := headersToAPI(context.Background(), in, &diags)
	if v := got["a"]; v == nil || *v != "1" {
		t.Errorf("headersToAPI = %#v", got)
	}
	if got := headersToAPI(context.Background(), types.MapNull(types.StringType), &diags); got != nil {
		t.Errorf("null map = %#v, want nil", got)
	}
}
