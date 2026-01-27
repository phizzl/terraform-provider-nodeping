# Risk Register

**Timestamp**: 2026-01-27T08:40:00Z  
**Meta-prompt version**: v3.8

---

## Risk Assessment Matrix

| Likelihood | Impact | Risk Level |
|------------|--------|------------|
| High | High | **Critical** |
| High | Medium | **High** |
| Medium | High | **High** |
| Medium | Medium | **Medium** |
| Low | High | **Medium** |
| Low | Medium | **Low** |
| Low | Low | **Low** |

---

## Security Risks

### SEC-001: API Token Exposure
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Critical |
| **Likelihood** | Medium |
| **Impact** | High |
| **Description** | API token could be exposed in logs, state, or error messages |
| **Mitigation** | - Mark `api_token` as Sensitive in schema<br>- Never log token value<br>- Use HTTP Basic Auth (not query string)<br>- Document state encryption best practices |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by design |

### SEC-002: Sensitive Data in State
| Attribute | Value |
|-----------|-------|
| **Risk Level** | High |
| **Likelihood** | High |
| **Impact** | Medium |
| **Description** | Terraform state contains contact addresses (email, phone) and check targets |
| **Mitigation** | - Mark sensitive fields in schema<br>- Document state encryption<br>- Warn in README about GDPR implications |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by documentation |

### SEC-003: Dependency Vulnerabilities
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Medium |
| **Likelihood** | Medium |
| **Impact** | Medium |
| **Description** | Third-party dependencies may contain security vulnerabilities |
| **Mitigation** | - Pin all dependency versions<br>- Run gosec in CI<br>- Generate SBOM<br>- Regular dependency updates |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by CI/CD |

### SEC-004: Credential Leakage via SubAccount
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Medium |
| **Likelihood** | Low |
| **Impact** | High |
| **Description** | Primary account token used for SubAccount operations could be misused |
| **Mitigation** | - Document least-privilege token usage<br>- Recommend separate tokens per SubAccount if available |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by documentation |

---

## Stability Risks

### STB-001: API Rate Limiting
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Medium |
| **Likelihood** | Medium |
| **Impact** | Medium |
| **Description** | NodePing may rate limit API requests, causing failures |
| **Mitigation** | - Implement client-side rate limiter<br>- Exponential backoff on 429 responses<br>- Configurable rate limit setting |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by design |

### STB-002: API Breaking Changes
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Low |
| **Likelihood** | Low |
| **Impact** | High |
| **Description** | NodePing API v1 could introduce breaking changes |
| **Mitigation** | - Pin to API v1<br>- Monitor NodePing changelog<br>- Acceptance tests catch regressions |
| **Owner** | Provider Engineering |
| **Status** | Accepted (low likelihood) |

### STB-003: Transient Network Failures
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Low |
| **Likelihood** | Medium |
| **Impact** | Low |
| **Description** | Network issues could cause intermittent failures |
| **Mitigation** | - Retry logic with exponential backoff<br>- Configurable retry settings<br>- Clear error messages |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by design |

### STB-004: State Drift Detection
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Medium |
| **Likelihood** | Medium |
| **Impact** | Medium |
| **Description** | External changes to NodePing resources may not be detected |
| **Mitigation** | - Implement thorough Read operations<br>- Compare all managed attributes<br>- Document refresh behavior |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by design |

---

## Operational Risks

### OPS-001: Complex Check Type Validation
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Medium |
| **Likelihood** | Medium |
| **Impact** | Medium |
| **Description** | 30+ check types with different required fields increases validation complexity |
| **Mitigation** | - Type-aware ValidateConfig<br>- Comprehensive test coverage per type<br>- Clear error messages |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by design |

### OPS-002: Address ID Management
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Medium |
| **Likelihood** | Medium |
| **Impact** | Medium |
| **Description** | Contact addresses have server-generated IDs that must be tracked |
| **Mitigation** | - Store address IDs in state<br>- Handle ID changes on update<br>- Test address add/remove scenarios |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by design |

### OPS-003: Import Complexity
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Low |
| **Likelihood** | Low |
| **Impact** | Medium |
| **Description** | Importing existing resources may miss configuration details |
| **Mitigation** | - Implement comprehensive ImportState<br>- Document import limitations<br>- Test import scenarios |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by design |

---

## Compliance Risks

### CMP-001: GDPR Data Subject Rights
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Medium |
| **Likelihood** | Low |
| **Impact** | High |
| **Description** | Contact data (email, phone) is personal data under GDPR |
| **Mitigation** | - Document data minimization<br>- `terraform destroy` removes managed data<br>- Document NodePing's data retention |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by documentation |

### CMP-002: Audit Trail
| Attribute | Value |
|-----------|-------|
| **Risk Level** | Low |
| **Likelihood** | Low |
| **Impact** | Medium |
| **Description** | Changes to contacts/checks may need audit trail for compliance |
| **Mitigation** | - Terraform state history provides audit trail<br>- Document logging configuration<br>- Recommend remote state with versioning |
| **Owner** | Provider Engineering |
| **Status** | Mitigated by documentation |

---

## Risk Summary

| Level | Count | IDs |
|-------|-------|-----|
| Critical | 1 | SEC-001 |
| High | 1 | SEC-002 |
| Medium | 7 | SEC-003, SEC-004, STB-001, STB-004, OPS-001, OPS-002, CMP-001 |
| Low | 4 | STB-002, STB-003, OPS-003, CMP-002 |

---

## Monitoring & Review

- Review risk register at each major release
- Update after security advisories
- Re-assess after NodePing API changes
- Track mitigation effectiveness through issue tracking

---

*Document generated as part of Phase 3: Planning*
