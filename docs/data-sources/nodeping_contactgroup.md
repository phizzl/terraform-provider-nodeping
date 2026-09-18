---
page_title: "nodeping_contactgroup Data Source - terraform-provider-nodeping"
subcategory: ""
description: |-
  Fetches a NodePing contact group by ID.
---

# nodeping_contactgroup (Data Source)

Fetches a NodePing contact group by ID.

## Example Usage

```hcl
data "nodeping_contactgroup" "escalation" {
  id = "201205050153W2Q4C-G-1ZIYU"
}

output "escalation_members" {
  value = data.nodeping_contactgroup.escalation.members
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the contact group.

## Attribute Reference

- `customer_id` - The customer ID (account ID) that owns this contact group.
- `name` - The display name of the contact group.
- `members` - Contact **address** IDs belonging to this group.
