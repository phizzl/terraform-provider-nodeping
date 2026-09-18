resource "nodeping_contact" "oncall" {
  name = "On-call"

  address {
    type    = "email"
    address = "oncall@example.com"
  }
}

resource "nodeping_contact" "backup" {
  name = "Backup"

  address {
    type    = "email"
    address = "backup@example.com"
  }
}

# Members are contact *address* IDs, not contact IDs. A contact can hold
# several addresses, and a group references them individually.
resource "nodeping_contactgroup" "escalation" {
  name = "Escalation"

  members = [
    nodeping_contact.oncall.address[0].id,
    nodeping_contact.backup.address[0].id,
  ]
}

# A group is referenced from a check like any other contact.
resource "nodeping_check" "site" {
  type   = "HTTP"
  target = "https://example.com"
  label  = "Website"

  notifications {
    contact_id = nodeping_contactgroup.escalation.id
    delay      = 0
    schedule   = "All"
  }
}
