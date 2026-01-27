# NodePing API Analysis

**Timestamp**: 2026-01-27T08:33:00Z  
**Meta-prompt version**: v3.8  
**Author**: Terraform Provider Engineering

---

## 1. API Overview

### Base URL
```
https://api.nodeping.com/api/1/
```

### URI Pattern
```
/api/:version/:resource/:id?querystring
```

- **:version** - Numeric API version (currently `1`)
- **:resource** - One of: `accounts`, `checks`, `contacts`, `contactgroups`, `schedules`, `results`
- **:id** - Resource-specific record ID

### HTTP Methods
| Method | Purpose |
|--------|---------|
| GET | Read record(s) |
| POST | Create new record |
| PUT | Update existing record |
| DELETE | Remove record |

Alternative: Use `action` query parameter (e.g., `?action=delete`) for environments that don't support all HTTP methods.

---

## 2. Authentication

### Token-based Authentication
A token is required for all API requests. Tokens are available only for primary accounts (not SubAccounts).

### Token Delivery Methods (in order of preference for security)
1. **HTTP Basic Auth** (username = token, password = ignored)
2. **JSON body**: `{"token": "abcdefghijklm"}`
3. **Query string**: `?token=abcdefghijklm` (least secure, avoid in production)

### Security Classification
- **Token**: RESTRICTED data - must never be logged, committed, or exposed

---

## 3. Error Handling

### HTTP Status Codes
| Code | Meaning |
|------|---------|
| 200 | Success |
| 400 | Invalid URL or parameter |
| 403 | Invalid token or insufficient permissions |
| 500/501 | Server-side bug |

### Error Response Format
```json
{"error": "Invalid time zone"}
```

---

## 4. Record ID Patterns

| Resource | ID Pattern | Example |
|----------|-----------|---------|
| Account | `{timestamp}{random}` | `201203232048C76FH` |
| Contact | `{account_id}-{random}` | `201203232048C76FH-BKPGH` |
| Contact Address | `{random}` (within contact) | `K5SP9CQP` |
| Contact Group | `{account_id}-G-{random}` | `201203232048C76FH-G-ISYB7` |
| Check | `{account_id}-{random}` | `201203232048C76FH-PXMVEFJN` |

**Important**: IDs are system-generated; cannot be set by client.

---

## 5. Contacts API

### Endpoints

| Operation | Method | Endpoint |
|-----------|--------|----------|
| List all | GET | `/api/1/contacts` |
| Get one | GET | `/api/1/contacts/{id}` |
| Create | POST | `/api/1/contacts` |
| Update | PUT | `/api/1/contacts/{id}` |
| Delete | DELETE | `/api/1/contacts/{id}` |
| Reset password | GET | `/api/1/contacts/{id}?action=RESETPASSWORD` |

### Contact Object Schema

```json
{
  "_id": "201205050153W2Q4C-BKPGH",
  "type": "contact",
  "customer_id": "201205050153W2Q4C",
  "name": "Foo Bar",
  "custrole": "owner",
  "addresses": {
    "K5SP9CQP": {
      "address": "foo@example.com",
      "type": "email",
      "status": "new"
    }
  }
}
```

### Contact Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `_id` | string | read-only | Contact ID |
| `customer_id` | string | read-only | Account ID |
| `name` | string | optional | Display label |
| `custrole` | string | optional | Permission level: `edit`, `view`, `notify` (default: `view`) |
| `addresses` | object | optional | Map of address ID to address object |

### Address Object Fields

| Field | Type | Description |
|-------|------|-------------|
| `address` | string | Email, phone, webhook URL, etc. |
| `type` | string | `email`, `sms`, `webhook`, `slack`, `hipchat`, `pushover`, `pagerduty`, `voice` |
| `suppressup` | bool | Suppress "up" notifications |
| `suppressdown` | bool | Suppress "down" notifications |
| `suppressfirst` | bool | Suppress "first result" notifications |
| `suppressdiag` | bool | Suppress diagnostic notifications |
| `suppressall` | bool | Suppress all notifications |
| `mute` | bool/int | `true` to mute, or millisecond timestamp for auto-unmute |

### Webhook Address Additional Fields

| Field | Type | Description |
|-------|------|-------------|
| `action` | string | HTTP method: `get`, `put`, `post`, `head`, `delete` |
| `data` | string | Request body (JSON, XML, etc.) |
| `headers` | object | HTTP headers |
| `querystrings` | object | Query string parameters |

### Pushover Address Additional Fields

| Field | Type | Description |
|-------|------|-------------|
| `priority` | int | -2 to 2 (default: 2) |

### Create/Update Parameters

| Parameter | Create | Update | Description |
|-----------|--------|--------|-------------|
| `customerid` | optional | optional | SubAccount customer ID |
| `id` | ignored | required | Contact ID to update |
| `name` | optional | optional | Contact name |
| `custrole` | optional | optional | Permission role |
| `addresses` | N/A | optional | Update existing addresses (full list required) |
| `newaddresses` | optional | optional | Add new addresses (array format) |

---

## 6. Checks API

### Endpoints

| Operation | Method | Endpoint |
|-----------|--------|----------|
| List all | GET | `/api/1/checks` |
| Get one | GET | `/api/1/checks/{id}` |
| Get multiple | GET | `/api/1/checks/{id1},{id2}` |
| Create | POST | `/api/1/checks` |
| Update | PUT | `/api/1/checks/{id}` |
| Delete | DELETE | `/api/1/checks/{id}` |
| Disable all | PUT | `/api/1/checks?disableall=true` |

### Check Types

```
AGENT, AUDIO, CLUSTER, DOHDOT, DNS, FTP, HTTP, HTTPCONTENT, HTTPPARSE, 
HTTPADV, IMAP4, MONGODB, MTR, MYSQL, NTP, PGSQL, PING, POP3, PORT, 
PUSH, RBL, RDAP, RDP, REDIS, SIP, SMTP, SNMP, SPEC10DNS, SPEC10RDDS, 
SSH, SSL, WEBSOCKET, WHOIS
```

### Check Object Schema

```json
{
  "_id": "201205050153W2Q4C-0J2HSIRF",
  "_rev": "37-8776f919267df3973fdb33cba0a8dd09",
  "customer_id": "201205050153W2Q4C",
  "label": "Site 1",
  "interval": 1,
  "notifications": [],
  "type": "HTTP",
  "status": "assigned",
  "modified": 1336759793520,
  "enable": "active",
  "public": false,
  "parameters": {
    "target": "http://www.example.com/",
    "threshold": "5",
    "sens": "2"
  },
  "created": 1336185808566,
  "queue": "bINPckIRdv",
  "uuid": "4pybhg6m-4v1y-4enn-8tz5-tvywydu6h04k",
  "state": 0,
  "firstdown": 1336185868566
}
```

### Common Check Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `_id` | string | read-only | - | Check ID |
| `type` | string | **required** | - | Check type (see list above) |
| `target` | string | required* | - | URL or FQDN (*except AGENT, DNS, PUSH, SPEC10DNS, SPEC10RDDS) |
| `label` | string | optional | target | Display label |
| `interval` | number | optional | 15 | Check interval in minutes (0.25, 0.5, or 1+) |
| `enabled` | string | optional | `false` | `true`/`active` or `false` |
| `public` | string | optional | `false` | Enable public reports |
| `autodiag` | string | optional | `false` | Enable automated diagnostics |
| `threshold` | int | optional | 5 | Timeout in seconds |
| `sens` | int | optional | 2 | Rechecks before status change |
| `mute` | bool/int | optional | false | Mute notifications |
| `dep` | string | optional | - | Dependency check ID |
| `description` | string | optional | - | Arbitrary text (max 1000 chars) |
| `tags` | array | optional | - | Array of tag strings |

### Location Fields

| Field | Type | Description |
|-------|------|-------------|
| `runlocations` | string/array | Region (`nam`, `lam`, `eur`, `eao`, `wlw`) or probe array (`["ld","ca","au"]`) |
| `homeloc` | string/bool | Preferred probe location or `roam` or `false` |

### Notifications Format

```json
[
  {"contactkey1": {"delay": 0, "schedule": "schedule1"}},
  {"contactkey2": {"delay": 5, "schedule": "schedule2"}}
]
```

To remove: `{"contactkey1": "None"}`

### Type-Specific Fields

#### HTTP/HTTPCONTENT/HTTPADV
- `follow` (bool) - Follow redirects (up to 4)
- `ipv6` (bool) - Use IPv6
- `contentstring` (string) - String to match
- `regex` (bool) - Treat contentstring as regex
- `method` (string) - HTTP method (HTTPADV only)
- `statuscode` (int) - Expected status code (HTTPADV only)
- `sendheaders` (object) - Request headers
- `receiveheaders` (object) - Expected response headers
- `data` (object) - Request data
- `postdata` (string) - POST body
- `clientcert` (string) - Client certificate ID

#### DNS/DOHDOT
- `dnstoresolve` (string) - FQDN to query
- `dnstype` (string) - Query type: ANY, A, AAAA, CNAME, MX, NS, PTR, SOA, SRV, TXT
- `contentstring` (string) - Expected response
- `dnssection` (string) - Section to check: answer, authority, additional, edns_options
- `dnsrd` (bool) - Recursion Desired bit (default: true)
- `transport` (string) - Protocol: udp, tcp
- `verify` (bool) - Verify DNSSEC (DNS) or SSL (DOHDOT)
- `port` (int) - DNS server port

#### SSL
- `warningdays` (int) - Days before expiry to fail
- `servername` (string) - SNI server name

#### PING
- `ipv6` (bool) - Use IPv6

#### PORT
- `port` (int) - **required** - Port number
- `invert` (bool) - Invert match

#### SSH
- `port` (int) - SSH port
- `username` (string) - SSH username
- `password` (string) - SSH password
- `sshkey` (string) - SSH private key ID
- `contentstring` (string) - Expected output
- `invert` (bool) - Invert match

#### SMTP/IMAP4/POP3
- `port` (int) - Server port
- `username` (string) - Auth username
- `password` (string) - Auth password
- `secure` (string) - `false`, `ssl`, or `starttls` (SMTP only)
- `verify` (bool) - Verify SSL certificate
- `warningdays` (int) - Days before cert expiry to fail
- `email` (string) - Email address (SMTP only)

#### Database Checks (MYSQL/PGSQL/MONGODB)
- `database` (string) - Database name
- `query` (string) - Query to execute
- `fields` (object) - Fields to parse from response
- `namespace` (string) - Collection (MONGODB only)

#### CLUSTER
- `data` (object) - Map of check IDs to include (`"checkid": "1"`)
- `threshold` (int) - Number of checks that must pass

#### PUSH/AGENT
- `checktoken` (string) - Read-only, set to `reset` to regenerate
- `oldresultfail` (bool) - Fail if results are old
- `fields` (object) - Fields to parse (PUSH only)

---

## 7. Contact Groups API

### Endpoints

| Operation | Method | Endpoint |
|-----------|--------|----------|
| List all | GET | `/api/1/contactgroups` |
| Get one | GET | `/api/1/contactgroups/{id}` |
| Create | POST | `/api/1/contactgroups` |
| Update | PUT | `/api/1/contactgroups/{id}` |
| Delete | DELETE | `/api/1/contactgroups/{id}` |

### Contact Group Object Schema

```json
{
  "_id": "201205050153W2Q4C-G-3QJWG",
  "type": "group",
  "customer_id": "201205050153W2Q4C",
  "name": "Example Group",
  "members": ["SLS78SDG", "9ZODE0VF"]
}
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `_id` | string | read-only | Group ID |
| `customer_id` | string | read-only | Account ID |
| `name` | string | optional | Group label |
| `members` | array | optional | Array of address IDs |

---

## 8. Pagination

**No explicit pagination documented.** The API returns all records for list operations. For the target scale (~50 contacts + ~50 checks), this is acceptable.

---

## 9. Rate Limits

**No explicit rate limits documented.** However, best practices apply:
- Implement exponential backoff for 5xx errors
- Avoid rapid successive requests
- Cache responses where appropriate

**Recommendation**: Implement configurable rate limiting with sensible defaults (e.g., 10 req/sec).

---

## 10. SubAccount Support

All resources support `customerid` parameter for SubAccount operations:
- Omit for primary account operations
- Include SubAccount's customer ID for SubAccount operations

This enables multi-account management through provider aliases.

---

## 11. Data Classification Summary

| Data Element | Classification |
|--------------|----------------|
| API Token | RESTRICTED |
| Contact addresses (email, phone) | CONFIDENTIAL |
| Check targets (URLs, IPs) | CONFIDENTIAL |
| Contact/Check IDs | CONFIDENTIAL |
| Check types, intervals, thresholds | INTERNAL |
| Documentation, schemas | PUBLIC |

---

## 12. Terraform Provider Implications

### Resources to Implement

1. **`nodeping_contact`** - Full CRUD for contacts with nested addresses
2. **`nodeping_check`** - Full CRUD for checks (all types via single resource with type-specific attributes)

### Data Sources to Implement

1. **`nodeping_contact`** - Read single contact by ID
2. **`nodeping_contacts`** - List all contacts
3. **`nodeping_check`** - Read single check by ID
4. **`nodeping_checks`** - List all checks

### Import Strategy

Both contacts and checks use predictable ID formats that can be used for `terraform import`:
- Contact: `{contact_id}` or `{customer_id}:{contact_id}` for SubAccounts
- Check: `{check_id}` or `{customer_id}:{check_id}` for SubAccounts

### State Considerations

- Terraform state will contain CONFIDENTIAL data (contact addresses, check targets)
- Documentation must warn users about state security
- Sensitive fields (passwords, tokens) must be marked as `Sensitive` in schema

---

## 13. Open Questions / Assumptions

1. **Rate limits**: Assumed no hard limits; implementing conservative defaults
2. **Pagination**: Assumed not needed for target scale
3. **Webhook secrets**: Not documented; assuming no HMAC validation available
4. **API versioning**: Currently v1; no deprecation timeline documented

---

*Document generated as part of Phase 1: Analysis*
