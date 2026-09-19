data "nodeping_check" "example" {
  id = "201205050153W2Q4C-0J2HSIRF"
}

output "check_target" {
  value = data.nodeping_check.example.target
}

# Check-type specific parameters are available too.
output "expected_content" {
  value = data.nodeping_check.example.contentstring
}

output "expected_status" {
  value = data.nodeping_check.example.statuscode
}
