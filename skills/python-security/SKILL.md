---
name: python-security
description: "Secure Python patterns for defensive coding. Trigger: writing security tools, implementing input validation, handling crypto, or building secure network apps."
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

# Python Security

Secure Python patterns, input validation, cryptographic best practices, and defensive coding for security tools and applications.

## When to Use This Skill

Use this skill when you need to:
- Write secure Python code for security tools and automation
- Implement input validation and sanitization
- Handle cryptographic operations safely
- Build secure network applications and APIs
- Prevent common Python security vulnerabilities (injection, path traversal, etc.)
- Use security-focused Python libraries correctly

## Core Security Libraries

### Cryptography

```python
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
import secrets

# Generate secure random token
def generate_token(length=32):
    """Generate cryptographically secure token"""
    return secrets.token_urlsafe(length)

# Generate secure random password
def generate_password(length=16):
    """Generate secure random password"""
    import string
    alphabet = string.ascii_letters + string.digits + string.punctuation
    return ''.join(secrets.choice(alphabet) for _ in range(length))

# Symmetric encryption with Fernet
def encrypt_data(data, key):
    """Encrypt data using Fernet symmetric encryption"""
    f = Fernet(key)
    return f.encrypt(data.encode())

def decrypt_data(encrypted_data, key):
    """Decrypt data"""
    f = Fernet(key)
    return f.decrypt(encrypted_data).decode()

# Generate encryption key
key = Fernet.generate_key()
```

### Secure Hashing

```python
import hashlib
import hmac

# Secure password hashing
def hash_password(password: str) -> str:
    """Hash password with salt using SHA-256"""
    salt = secrets.token_hex(16)
    hashed = hashlib.sha256((salt + password).encode()).hexdigest()
    return f"{salt}${hashed}"

def verify_password(stored: str, provided: str) -> bool:
    """Verify password against stored hash"""
    salt, hashed = stored.split('$')
    return hmac.compare_digest(
        hashlib.sha256((salt + provided).encode()).hexdigest(),
        hashed
    )

# HMAC for message authentication
def create_hmac(message: str, key: bytes) -> str:
    """Create HMAC signature"""
    return hmac.new(key, message.encode(), hashlib.sha256).hexdigest()
```

## Input Validation & Sanitization

### IP Address Validation

```python
import re
import ipaddress

def validate_ip(ip: str) -> bool:
    """Validate IP address using ipaddress module"""
    try:
        ipaddress.ip_address(ip)
        return True
    except ValueError:
        return False

def is_private_ip(ip: str) -> bool:
    """Check if IP is private/internal"""
    try:
        return ipaddress.ip_address(ip).is_private
    except ValueError:
        return False
```

### Input Sanitization

```python
import html
import re

def sanitize_html_input(user_input: str) -> str:
    """Escape HTML to prevent XSS"""
    return html.escape(user_input)

def sanitize_filename(filename: str) -> str:
    """Sanitize filename to prevent path traversal"""
    # Remove path separators and dangerous characters
    sanitized = re.sub(r'[\\/:\*\?"<>\|]', '', filename)
    # Remove leading/trailing whitespace and dots
    sanitized = sanitized.strip('. ')
    return sanitized or 'unnamed'

def validate_url(url: str, allowed_domains: list) -> bool:
    """Validate URL against allowlist to prevent SSRF"""
    from urllib.parse import urlparse
    parsed = urlparse(url)
    return parsed.hostname in allowed_domains
```

### SQL Injection Prevention

```python
import sqlite3

# SECURE: Parameterized queries
def get_user(db: sqlite3.Connection, user_id: int):
    """Secure parameterized query"""
    cursor = db.execute("SELECT * FROM users WHERE id = ?", (user_id,))
    return cursor.fetchone()

# NEVER DO THIS:
# cursor.execute(f"SELECT * FROM users WHERE id = {user_id}")
```

## Network Security

### Secure HTTP Requests

```python
import requests
from requests.adapters import HTTPAdapter
from urllib3.util.ssl_ import create_urllib3_context

def create_secure_session() -> requests.Session:
    """Create a session with TLS verification"""
    session = requests.Session()
    adapter = HTTPAdapter()
    session.mount('https://', adapter)
    return session

def secure_get(url: str, timeout: int = 10) -> requests.Response:
    """Make secure HTTP GET request"""
    response = requests.get(url, timeout=timeout, verify=True)
    response.raise_for_status()
    return response
```

### Secure Socket Programming

```python
import socket
import ssl

def create_tls_socket(host: str, port: int) -> ssl.SSLSocket:
    """Create a TLS-secured socket"""
    context = ssl.create_default_context()
    context.check_hostname = True
    context.verify_mode = ssl.CERT_REQUIRED

    sock = socket.create_connection((host, port))
    secure_sock = context.wrap_socket(sock, server_hostname=host)
    return secure_sock
```

## Common Security Testing Patterns

### Port Scanning (Ethical Use Only)

```python
import socket

def scan_port(host: str, port: int, timeout: float = 1.0) -> bool:
    """Scan single port"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        result = sock.connect_ex((host, port))
        sock.close()
        return result == 0
    except socket.error:
        return False

def scan_ports(host: str, ports: list) -> list:
    """Scan multiple ports"""
    return [port for port in ports if scan_port(host, port)]
```

### Banner Grabbing

```python
import socket

def grab_banner(host: str, port: int, timeout: float = 2.0) -> str | None:
    """Grab service banner"""
    try:
        sock = socket.socket()
        sock.settimeout(timeout)
        sock.connect((host, port))
        banner = sock.recv(1024).decode().strip()
        sock.close()
        return banner
    except (socket.error, UnicodeDecodeError):
        return None
```

## Security Best Practices

### Never Hardcode Secrets

```python
# BAD: Hardcoded credentials
API_KEY = "sk-1234567890abcdef"
DB_PASSWORD = "supersecret"

# GOOD: Use environment variables
import os
API_KEY = os.environ.get("API_KEY")
DB_PASSWORD = os.environ.get("DB_PASSWORD")
```

### Use Secure Temp Files

```python
import tempfile

# GOOD: Secure temporary file
with tempfile.NamedTemporaryFile(delete=True, mode='w') as f:
    f.write("sensitive data")
    f.flush()
    # File is automatically deleted on close
```

### Avoid eval/exec

```python
# NEVER DO THIS:
# result = eval(user_input)
# exec(user_input)

# GOOD: Use safe alternatives
import ast
import json

# Parse JSON safely
data = json.loads(user_input)

# Parse Python literals safely
result = ast.literal_eval(user_input)
```

### Secure File Operations

```python
import os

def safe_file_read(filepath: str, allowed_dir: str) -> str:
    """Read file safely, preventing path traversal"""
    real_path = os.path.realpath(filepath)
    real_allowed = os.path.realpath(allowed_dir)

    if not real_path.startswith(real_allowed):
        raise ValueError(f"Access denied: {filepath} is outside allowed directory")

    with open(real_path, 'r') as f:
        return f.read()
```

## Ethical Security Guidelines

1. **Authorization**: ALWAYS get written permission before testing
2. **Scope**: Stay within defined boundaries
3. **Documentation**: Log all activities
4. **Responsible Disclosure**: Report vulnerabilities properly
5. **No Harm**: Never cause damage or disruption

## Installation

```bash
# Core security libraries
pip install cryptography requests

# Optional: for advanced security testing
pip install scapy impacket paramiko
```

## References

- [OWASP Python Security Cheatsheet](https://cheatsheetseries.owasp.org/cheatsheets/Python_Security_Cheatsheet.html)
- [Python Security Best Practices](https://python.readthedocs.io/en/stable/library/security_warnings.html)
- [Cryptography Documentation](https://cryptography.io/en/latest/)
- [Secrets Module](https://docs.python.org/3/library/secrets.html)

## Validation

To verify this skill works correctly:

1. **Load test**: Confirm the skill loads without frontmatter parsing errors.
2. **Code test**: Execute all Python code examples to verify they run without errors (cryptography, hashlib, ipaddress, html, requests, ssl, socket, tempfile, ast, json modules).
3. **Security test**: Verify the "NEVER DO THIS" examples are clearly marked as insecure and the "GOOD" alternatives are correct.
4. **Import test**: Confirm all import statements reference standard or correctly-named third-party libraries.
5. **Best practices test**: Validate that all code examples follow current Python security recommendations (no eval/exec, parameterized queries, TLS verification).
