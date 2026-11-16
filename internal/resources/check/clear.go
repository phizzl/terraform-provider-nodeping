package check

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nodeping/terraform-provider-nodeping/internal/client"
	"github.com/nodeping/terraform-provider-nodeping/internal/datasources/checkattr"
)

// NodePing merges an update into the check and keeps whatever the update
// leaves out. An attribute removed from the configuration is simply null in
// the plan, so it used to be left out of the update: NodePing kept the value,
// the apply failed with "Provider produced inconsistent result after apply",
// and every later plan wanted the removal again.
//
// clearRemoved therefore looks at the prior state as well as the plan, and
// for each attribute the check had and the plan drops, sends the value
// NodePing stores as cleared. Something the check never had is never written,
// and a create still leaves unset values out.
//
// The clearing values are the ones the stefang-kahoot fork of this provider
// determined against the NodePing API on 2026-10-08:
//
//   - contentstring, method, postdata, servername: ""
//   - statuscode, warningdays: ""
//   - regex, invert, follow, ipv6: false. NodePing keeps the key, as false.
//   - dep: false, NodePing's documented way to remove a dependency.
//   - runlocations: []
//   - sendheaders, receiveheaders: each removed header as null. NodePing
//     merges headers per name, so a smaller map, {} or no map at all would
//     keep every header.
//   - description: a single space. NodePing ignores "", null, false and 0
//     for a description, so it can only be overwritten; the space reads back
//     as no description (checkattr.ClearedDescription).
//
// Attributes not listed (homeloc, port, the dns*, snmp* and database ones)
// were not probed and are left as before: removing one still keeps its value
// at NodePing.
func clearRemoved(req *client.CheckCreateRequest, plan, state *CheckResourceModel) {
	for _, s := range []struct {
		prior, planned attr.Value
		field          **string
	}{
		{state.ContentString, plan.ContentString, &req.ContentString},
		{state.Method, plan.Method, &req.Method},
		{state.PostData, plan.PostData, &req.PostData},
		{state.ServerName, plan.ServerName, &req.ServerName},
	} {
		if removed(s.prior, s.planned) {
			*s.field = stringPtr("")
		}
	}

	if removedDescription(state.Description, plan.Description) {
		req.Description = checkattr.ClearedDescription
	}

	for _, v := range []struct {
		prior, planned attr.Value
		field          *interface{}
		cleared        interface{}
	}{
		{state.StatusCode, plan.StatusCode, &req.StatusCode, ""},
		{state.WarningDays, plan.WarningDays, &req.WarningDays, ""},
		{state.Regex, plan.Regex, &req.Regex, false},
		{state.Invert, plan.Invert, &req.Invert, false},
		{state.Follow, plan.Follow, &req.Follow, false},
		{state.IPv6, plan.IPv6, &req.IPv6, false},
		{state.Dep, plan.Dep, &req.Dep, false},
		{state.RunLocations, plan.RunLocations, &req.RunLocations, []string{}},
	} {
		if removed(v.prior, v.planned) {
			*v.field = v.cleared
		}
	}

	req.SendHeaders = deleteRemovedKeys(req.SendHeaders, state.SendHeaders, plan.SendHeaders)
	req.ReceiveHeaders = deleteRemovedKeys(req.ReceiveHeaders, state.ReceiveHeaders, plan.ReceiveHeaders)
}

// publicForUpdate returns public the way an update has to send it. On an
// update NodePing ignores public sent as the boolean false (or 0); only the
// strings "false" and "0" switch public reports off, and "true" switches them
// on. A create keeps sending the boolean: a new check is not public unless
// told so.
func publicForUpdate(plan *CheckResourceModel) interface{} {
	if plan.Public.IsNull() || plan.Public.IsUnknown() {
		return nil
	}
	return strconv.FormatBool(plan.Public.ValueBool())
}

// deleteRemovedKeys adds to the headers an update sends each header the prior
// state has and the plan does not, as null, which deletes it.
func deleteRemovedKeys(sent map[string]*string, prior, planned types.Map) map[string]*string {
	if prior.IsNull() || prior.IsUnknown() || planned.IsUnknown() {
		return sent
	}
	kept := planned.Elements()
	for name := range prior.Elements() {
		if _, ok := kept[name]; ok {
			continue
		}
		if sent == nil {
			sent = make(map[string]*string)
		}
		sent[name] = nil
	}
	return sent
}

// keepClearedBooleans undoes what NodePing does to a cleared boolean. Once
// set, regex, invert, follow and ipv6 are kept as false rather than removed,
// so a check whose configuration dropped one reads back false where the plan
// or the prior state has null -- an inconsistent result after apply, and a
// permanent diff after it. Where the expected value is null and NodePing has
// false, the null stands. A true from NodePing is a real change and is kept.
//
// Not used on import, which has no expected value: there a stored false reads
// as false.
func keepClearedBooleans(model *CheckResourceModel, expected *CheckResourceModel) {
	for _, b := range []struct {
		got      *types.Bool
		expected types.Bool
	}{
		{&model.Regex, expected.Regex},
		{&model.Invert, expected.Invert},
		{&model.Follow, expected.Follow},
		{&model.IPv6, expected.IPv6},
	} {
		if b.expected.IsNull() && !b.got.IsNull() && !b.got.IsUnknown() && !b.got.ValueBool() {
			*b.got = types.BoolNull()
		}
	}
}

// headersToAPI converts configured headers to the request shape, or nil when
// there are none.
func headersToAPI(ctx context.Context, v types.Map, diags *diag.Diagnostics) map[string]*string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	headers := make(map[string]string)
	diags.Append(v.ElementsAs(ctx, &headers, false)...)
	out := make(map[string]*string, len(headers))
	for name, value := range headers {
		out[name] = stringPtr(value)
	}
	return out
}

func stringPtr(s string) *string {
	return &s
}

// removed reports whether the prior state holds a value the plan no longer
// does.
func removed(prior, planned attr.Value) bool {
	return !isEmpty(prior) && isEmpty(planned)
}

// removedDescription is removed for a description, where the cleared value
// itself also counts as empty: a description cleared earlier is held as
// ClearedDescription and must not be cleared again on every update.
func removedDescription(prior, planned types.String) bool {
	return !isEmptyDescription(prior) && isEmptyDescription(planned)
}

func isEmptyDescription(v types.String) bool {
	return isEmpty(v) || v.ValueString() == checkattr.ClearedDescription
}

// isEmpty reports whether a value is null or its type's empty value: "",
// false, or a list or map with nothing in it -- NodePing treats those alike.
// A number is empty only when null, since 0 is a real value for some, such as
// volumemin. An unknown value is not empty.
func isEmpty(v attr.Value) bool {
	if v.IsNull() {
		return true
	}
	if v.IsUnknown() {
		return false
	}
	switch v := v.(type) {
	case types.String:
		return v.ValueString() == ""
	case types.Bool:
		return !v.ValueBool()
	case types.List:
		return len(v.Elements()) == 0
	case types.Map:
		return len(v.Elements()) == 0
	}
	return false
}
