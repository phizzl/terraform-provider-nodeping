data "nodeping_contactgroups" "all" {}

output "group_names" {
  value = [for g in data.nodeping_contactgroups.all.contactgroups : g.name]
}

# Look up a group by name instead of hard-coding its ID.
output "escalation_id" {
  value = one([
    for g in data.nodeping_contactgroups.all.contactgroups : g.id
    if g.name == "Escalation"
  ])
}
