data "nodeping_contactgroup" "escalation" {
  id = "201205050153W2Q4C-G-1ZIYU"
}

output "escalation_name" {
  value = data.nodeping_contactgroup.escalation.name
}

output "escalation_members" {
  value = data.nodeping_contactgroup.escalation.members
}
