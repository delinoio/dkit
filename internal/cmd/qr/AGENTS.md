# qr Command

## Purpose
Generate and decode QR codes directly in the terminal. Share URLs, WiFi credentials, contact information, and arbitrary data through QR codes for easy mobile device access.

## Command Signature
```bash
dkit qr [subcommand] [options]
```

## Subcommands

### generate - Generate QR Code

#### Purpose
Create a QR code from text, URLs, or structured data.

#### Command Signature
```bash
dkit qr generate <data> [options]
```

**Arguments:**
- `data` - Content to encode (text, URL, etc.)

**Options:**
- `--output <file>` - Save to image file (png, svg, jpg)
- `--size <small|medium|large>` - QR code size (default: medium)
- `--error-correction <L|M|Q|H>` - Error correction level (default: M)
- `--format <terminal|png|svg|ascii>` - Output format (default: terminal)
- `--inverted` - Invert colors (white on black)
- `--quiet-zone <n>` - Border size in modules (default: 4)

#### Error Correction Levels
- **L** (Low): ~7% error correction
- **M** (Medium): ~15% error correction (default)
- **Q** (Quartile): ~25% error correction
- **H** (High): ~30% error correction

#### Output Formats

**Terminal (default):**
```bash
dkit qr generate "https://example.com"

[dkit] QR Code:

████████████████████████████████
████████████████████████████████
████  ██████  ████  ██████  ████
████  ██  ██  ████  ██  ██  ████
████  ██  ██  ████  ██  ██  ████
████  ██████  ████  ██████  ████
████████████████████████████████
████      ██    ██      ████████
████  ████████████████  ████████
████████  ████  ████  ██████████
████  ██████████████████  ██████
████  ██  ██  ████  ██  ████████
████  ██████  ████  ██████  ████
████████████████████████████████

Scan with your phone camera
Data: https://example.com
Size: 29x29 modules
```

**ASCII (for copy-paste):**
```bash
dkit qr generate "Hello" --format ascii

██████████████  ████  ██████████████
██          ██  ██    ██          ██
██  ██████  ██    ██  ██  ██████  ██
██  ██████  ██  ██    ██  ██████  ██
██  ██████  ██  ████  ██  ██████  ██
██          ██      ████          ██
██████████████  ██  ████████████████
              ██  ██                
    ██  ██  ██      ██████      ██  
██    ██      ████        ██  ██████
  ██████████████    ██    ████  ██  
                ██  ████  ████      
██████████████      ████  ██    ████
██          ██  ██████  ████  ██    
██  ██████  ██    ██  ████    ██████
██  ██████  ██  ████████  ██████    
██  ██████  ██  ██      ████  ██  ██
██          ██        ██████  ██████
██████████████  ██    ██  ██    ██  
```

**PNG file:**
```bash
dkit qr generate "https://example.com" --output qr.png --size large

[dkit] ✓ QR code saved to qr.png
[dkit] Size: 500x500 pixels
[dkit] Format: PNG
```

**SVG file:**
```bash
dkit qr generate "Contact info" --output qr.svg --format svg

[dkit] ✓ QR code saved to qr.svg
[dkit] Format: SVG (scalable)
```

#### Examples

**Simple text:**
```bash
dkit qr generate "Hello, World!"
```

**URL:**
```bash
dkit qr generate "https://github.com/username/repo"
```

**From file:**
```bash
cat message.txt | dkit qr generate
```

**Large data with high error correction:**
```bash
dkit qr generate "$(cat data.json)" --error-correction H --size large
```

### wifi - Generate WiFi QR Code

#### Purpose
Create a QR code for WiFi network credentials. When scanned, automatically connects device to WiFi.

#### Command Signature
```bash
dkit qr wifi <ssid> [password] [options]
```

**Arguments:**
- `ssid` - WiFi network name
- `password` - WiFi password (optional for open networks)

**Options:**
- `--security <WPA|WEP|nopass>` - Security type (default: WPA)
- `--hidden` - Network is hidden
- `--output <file>` - Save to file
- `--size <small|medium|large>` - QR code size

#### WiFi QR Code Format
Uses standard WiFi QR code format: `WIFI:T:WPA;S:MyNetwork;P:MyPassword;;`

#### Output

```bash
dkit qr wifi "MyHomeWiFi" "secretpassword"

[dkit] WiFi QR Code:

████████████████████████████████████
████████████████████████████████████
████  ██████  ████    ██████  ██████
████  ██  ██  ██  ██  ██  ██  ██████
████  ██  ██  ████    ██  ██  ██████
████  ██████  ██████  ██████  ██████
████████████████████████████████████
████      ████  ██      ████████████
████  ██████  ████  ██████  ████████
████████████  ████  ████████████████
████  ████████  ██  ██  ██████  ████
████  ██  ██████  ██████  ██████████
████  ██████  ████  ████████  ██████
████████████████████████████████████

Scan to connect to WiFi
Network: MyHomeWiFi
Security: WPA/WPA2
```

**Open network:**
```bash
dkit qr wifi "GuestNetwork" --security nopass
```

**Hidden network:**
```bash
dkit qr wifi "HiddenNet" "password" --hidden
```

### vcard - Generate Contact QR Code

#### Purpose
Create a QR code with contact information (vCard format).

#### Command Signature
```bash
dkit qr vcard [options]
```

**Options:**
- `--name <name>` - Full name (required)
- `--email <email>` - Email address
- `--phone <phone>` - Phone number
- `--org <organization>` - Organization/Company
- `--title <title>` - Job title
- `--url <url>` - Website
- `--address <address>` - Physical address
- `--note <note>` - Additional notes
- `--output <file>` - Save to file

#### Output

```bash
dkit qr vcard --name "John Doe" --email "john@example.com" --phone "+1-555-0123"

[dkit] Contact QR Code:

████████████████████████████████
████████████████████████████████
████  ██████  ██  ██████  ██████
████  ██  ██  ████  ████  ██████
████  ██  ██  ████  ████  ██████
████  ██████  ████████████  ████
████████████████████████████████
████        ██  ████  ██████████
████  ██  ██████  ██  ████  ████
████████████  ██  ██████████████
████  ████    ████  ██  ████████
████  ██████████    ██  ████  ██
████  ██████  ████████████  ████
████████████████████████████████

Scan to add contact
Name: John Doe
Email: john@example.com
Phone: +1-555-0123
```

**Full contact card:**
```bash
dkit qr vcard \
  --name "Jane Smith" \
  --email "jane@company.com" \
  --phone "+1-555-9876" \
  --org "Acme Corp" \
  --title "Software Engineer" \
  --url "https://janesmith.dev" \
  --output contact.png
```

### url - Generate URL QR Code

#### Purpose
Create a QR code for a URL with optional URL shortening and analytics.

#### Command Signature
```bash
dkit qr url <url> [options]
```

**Arguments:**
- `url` - URL to encode

**Options:**
- `--shorten` - Shorten URL first (requires API key)
- `--title <title>` - Add title metadata
- `--output <file>` - Save to file

#### Output

```bash
dkit qr url "https://github.com/username/repo"

[dkit] URL QR Code:

████████████████████████████████
████████████████████████████████
████  ██████  ██  ████  ████████
████  ██  ██  ████  ████  ██████
████  ██  ██  ████  ██  ████████
████  ██████  ██████████████████
████████████████████████████████
████      ██  ████    ██████████
████  ████████  ██████  ████████
████  ██  ██  ██████  ██████████
████  ████  ████  ██████  ██████
████  ██████  ██  ████  ████  ██
████  ██████  ██████  ████  ████
████████████████████████████████

Scan to open URL
https://github.com/username/repo
```

### scan - Decode QR Code

#### Purpose
Read and decode QR codes from image files or camera.

#### Command Signature
```bash
dkit qr scan <image> [options]
```

**Arguments:**
- `image` - Path to image file containing QR code

**Options:**
- `--format <text|json>` - Output format
- `--camera` - Scan from camera (instead of file)
- `--copy` - Copy decoded data to clipboard

#### Supported Image Formats
- PNG
- JPEG/JPG
- GIF
- WebP
- BMP

#### Output Format

**Text (default):**
```bash
dkit qr scan qrcode.png

[dkit] QR Code decoded successfully

Type: URL
Data: https://example.com/page

Raw content:
https://example.com/page
```

**JSON:**
```bash
dkit qr scan qrcode.png --format json

{
  "success": true,
  "type": "url",
  "data": "https://example.com/page",
  "error_correction": "M",
  "version": 3
}
```

**WiFi QR code:**
```bash
dkit qr scan wifi-qr.png

[dkit] QR Code decoded successfully

Type: WiFi Network
SSID: MyHomeWiFi
Password: secretpassword
Security: WPA/WPA2
Hidden: false

To connect manually, use:
  Network: MyHomeWiFi
  Password: secretpassword
```

**vCard QR code:**
```bash
dkit qr scan contact.png

[dkit] QR Code decoded successfully

Type: Contact (vCard)
Name: John Doe
Email: john@example.com
Phone: +1-555-0123
Organization: Acme Corp

Full vCard:
BEGIN:VCARD
VERSION:3.0
FN:John Doe
EMAIL:john@example.com
TEL:+1-555-0123
ORG:Acme Corp
END:VCARD
```

**Camera scan (interactive):**
```bash
dkit qr scan --camera

[dkit] Opening camera...
[dkit] Position QR code in front of camera
[dkit] Press 'q' to quit

[dkit] ✓ QR code detected!
[dkit] Data: https://example.com
```

### batch - Generate Multiple QR Codes

#### Purpose
Generate multiple QR codes from a list of inputs.

#### Command Signature
```bash
dkit qr batch <input-file> [options]
```

**Arguments:**
- `input-file` - File containing data (one per line) or JSON array

**Options:**
- `--output-dir <dir>` - Directory to save QR codes (default: qr-codes/)
- `--format <png|svg>` - Output format
- `--prefix <prefix>` - Filename prefix (default: qr-)
- `--size <small|medium|large>` - QR code size

#### Input Formats

**Text file (one per line):**
```
https://example.com/page1
https://example.com/page2
https://example.com/page3
```

**JSON:**
```json
[
  {"data": "https://example.com/1", "filename": "page1.png"},
  {"data": "https://example.com/2", "filename": "page2.png"}
]
```

#### Output

```bash
dkit qr batch urls.txt --output-dir ./codes

[dkit] Generating QR codes...

  [████████████████████] 100% (3/3)

Generated:
  ✓ codes/qr-1.png
  ✓ codes/qr-2.png
  ✓ codes/qr-3.png

Summary:
  3 QR codes generated
  Output directory: ./codes/
```

## Common Use Cases

### Development & Testing

**Share localhost URL to mobile:**
```bash
# Get local IP
IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1)

# Generate QR code
dkit qr generate "http://$IP:3000"
```

**Share API endpoint:**
```bash
dkit qr generate "https://api.example.com/v1/users" --output api-qr.png
```

### Deployment & Operations

**Share staging environment:**
```bash
dkit qr url "https://staging.example.com" --title "Staging Environment"
```

**Emergency contact info:**
```bash
dkit qr vcard \
  --name "On-Call Engineer" \
  --phone "+1-555-ONCALL" \
  --email "oncall@company.com" \
  --output oncall-contact.png
```

### Guest WiFi

**Print QR code for guests:**
```bash
dkit qr wifi "GuestNetwork" "guest2025" --output guest-wifi.png --size large
```

**Multiple locations:**
```bash
echo "Office-Main:password123
Office-Guest:guest456
Conference-Room:conf789" | while IFS=: read ssid pass; do
  dkit qr wifi "$ssid" "$pass" --output "wifi-$ssid.png"
done
```

### Event Management

**Event registration URLs:**
```bash
# Generate QR codes for different ticket types
dkit qr generate "https://event.com/register?type=vip" --output vip-ticket.png
dkit qr generate "https://event.com/register?type=regular" --output regular-ticket.png
```

**Table numbers with survey URLs:**
```bash
for i in {1..20}; do
  dkit qr generate "https://survey.com/table/$i" --output "table-$i-qr.png"
done
```

### Marketing

**Social media links:**
```bash
dkit qr generate "https://instagram.com/username" --output ig-qr.png
dkit qr generate "https://twitter.com/username" --output twitter-qr.png
```

**Product information:**
```bash
dkit qr generate "$(cat product-info.json)" --output product-qr.png
```

## Integration Examples

### With clipboard
```bash
# Generate and copy to clipboard
dkit qr generate "https://example.com" --format png --output - | dkit clipboard copy
```

### With web server
```bash
# Serve current directory and show QR code
IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1)
dkit qr generate "http://$IP:8080" &
dkit http serve -p 8080
```

### CI/CD Integration
```bash
# Generate QR code for deployment
DEPLOY_URL="https://app-$CI_COMMIT_SHA.example.com"
dkit qr generate "$DEPLOY_URL" --output deployment-qr.png
# Upload to artifacts
```

## Exit Codes
- `0` - Success
- `1` - Generation/decoding failed
- `2` - Invalid input data
- `3` - File I/O error
- `4` - No QR code found in image
- `127` - Invalid command usage

## Error Handling

### Data Too Large
```
[dkit] ERROR: Data too large for QR code
[dkit] Data size: 3,500 bytes
[dkit] Maximum for error correction level M: 2,953 bytes
[dkit] 
[dkit] Suggestions:
[dkit]   1. Use lower error correction (--error-correction L)
[dkit]   2. Reduce data size
[dkit]   3. Use URL shortener (--shorten)
```

### Invalid Image
```
[dkit] ERROR: Could not decode QR code from image
[dkit] File: blurry-qr.png
[dkit] Possible issues:
[dkit]   - Image quality too low
[dkit]   - QR code partially obscured
[dkit]   - Image contains no QR code
[dkit] 
[dkit] Try:
[dkit]   - Use higher resolution image
[dkit]   - Ensure QR code is fully visible
[dkit]   - Check image is not corrupted
```

### Camera Not Available
```
[dkit] ERROR: Camera not accessible
[dkit] Ensure camera permissions are granted
[dkit] On macOS: System Preferences → Security & Privacy → Camera
```

## Implementation Requirements

### Performance
- Fast QR code generation (< 100ms for typical data)
- Efficient image processing for scanning
- Minimal dependencies

### Correctness
- Comply with QR code specification (ISO/IEC 18004)
- Proper error correction implementation
- Accurate WiFi and vCard format encoding
- Reliable decoding from various image qualities

### User Experience
- Clear terminal rendering on all platforms
- Beautiful output suitable for printing
- Helpful error messages
- Progress indication for batch operations

### Cross-Platform
- Work on macOS, Linux, Windows
- Handle terminal color/Unicode support differences
- Camera access on different platforms

## Design Principles
- **Instant utility**: Generate QR codes in one command
- **Mobile-friendly**: QR codes designed for phone scanning
- **Flexible output**: Terminal, file, or clipboard
- **Smart defaults**: Sensible error correction and size
- **Standard formats**: WiFi, vCard, URL follow standards

