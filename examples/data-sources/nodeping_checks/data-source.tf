data "nodeping_checks" "http" {
  type = "HTTP"
}

output "targets" {
  value = [for c in data.nodeping_checks.http.checks : c.target]
}

# Filtering on a check-type specific parameter.
output "following_redirects" {
  value = [
    for c in data.nodeping_checks.http.checks : c.label
    if c.follow == true
  ]
}

# Checks whose TLS warning window is shorter than two weeks.
data "nodeping_checks" "ssl" {
  type = "SSL"
}

output "short_warning_window" {
  value = [
    for c in data.nodeping_checks.ssl.checks : c.label
    if c.warningdays != null && c.warningdays < 14
  ]
}
