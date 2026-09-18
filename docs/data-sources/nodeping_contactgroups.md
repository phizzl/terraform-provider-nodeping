---
page_title: "nodeping_contactgroups Data Source - terraform-provider-nodeping"
subcategory: ""
description: |-
  Fetches all NodePing contact groups.
---

# nodeping_contactgroups (Data Source)

Fetches all NodePing contact groups for the account.

## Example Usage

```hcl
data "nodeping_contactgroups" "all" {}

output "group_names" {
  value = [for g in data.nodeping_contactgroups.all.contactgroups : g.name]
}
```

Looking a group up by name instead of hard-coding its ID:

```hcl
output "escalation_id" {
  value = one([
    for g in data.nodeping_contactgroups.all.contactgroups : g.id
    if g.name == "Escalation"
  ])
}
```

## Attribute Reference

- `contactgroups` - All contact groups, ordered by ID. Each entry contains:
  - `id` - The unique identifier of the contact group.
  - `customer_id` - The customer ID (account ID) that owns this contact group.
  - `name` - The display name of the contact group.
  - `members` - Contact **address** IDs belonging to this group.
