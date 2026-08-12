# AMEN Mixer Development Guide

## Development Environment

Required:

```text
Go
Wails v2
Node.js
npm
Git
```

Optional:

```text
Android Studio
PlatformIO
ESP32 hardware
```

---

## First Setup

Enter the desktop project:

```bash
cd desktop
```

Install/update Go dependencies:

```bash
go mod tidy
```

Install frontend dependencies:

```bash
cd frontend
npm install
cd ..
```

Check Wails:

```bash
wails doctor
```

---

## mDNS Dependency

AMEN Mixer uses HashiCorp mDNS for:

```text
amen-mixer.local
```

If the dependency has not been installed:

```bash
go get github.com/hashicorp/mdns@v1.0.6
go mod tidy
```

---

## Start Development

From:

```text
desktop/
```

run:

```bash
wails dev
```

Expected services:

```text
Wails Desktop
Vite frontend :5173
Realtime server :8081
mDNS
Serial manager
```

---

## Expected Startup Logs

Example:

```text
===================================
AMEN MIXER STARTING...
===================================

[WS] =================================
[WS] AMEN REALTIME SERVER
[WS] Listening :8081

[MDNS] =================================
[MDNS] AMEN MIXER DISCOVERY
[MDNS] LAN IP : 192.168.x.x
[MDNS] HOST   : amen-mixer.local
[MDNS] UI     : http://amen-mixer.local:5173
[MDNS] WS     : ws://amen-mixer.local:8081/ws
[MDNS] =================================
```

Serial may report that hardware is unavailable.

That does not prevent browser/mobile development.

---

## Open Mobile UI

Preferred:

```text
http://amen-mixer.local:5173
```

The phone and development computer must be on the same local network.

---

## IP Fallback

If `.local` does not resolve on macOS:

```bash
ipconfig getifaddr en0
```

Example output:

```text
192.168.18.11
```

Open on the phone:

```text
http://192.168.18.11:5173
```

---

## Test Backend Connectivity

From the phone browser:

```text
http://amen-mixer.local:8081/
```

Expected:

```text
AMEN MIXER REALTIME SERVER OK
```

Using an IP:

```text
http://192.168.18.11:8081/
```

Expected response is the same.

If the health endpoint cannot be opened, troubleshoot the network before debugging Vue or WebSocket code.

---

## Test WebSocket Registration

When a phone connects, the Go terminal should show a new connection.

Expected flow:

```text
[WS] Incoming connection
[WS] Client connected
[CLIENT] mobile iPhone
```

The desktop should also register as a client.

---

## Test Mobile → Desktop

1. Start `wails dev`.
2. Open the mobile UI.
3. Move one mixer slider.
4. Watch the Go terminal.
5. Watch the Wails desktop slider.

Expected backend flow:

```text
[CHANNEL UPDATE]
[WINDOWS AUDIO]
[WS] Broadcasting
```

Expected UI behavior:

```text
Phone slider changes
        ↓
Desktop slider follows
```

---

## Test Desktop → Mobile

Move a slider from the Wails desktop.

Expected:

```text
Wails
 ↓
Audio Manager
 ↓
BroadcastChannelUpdate
 ↓
Mobile
```

The mobile slider should follow without refreshing the page.

---

## Browser Console

Useful messages:

```text
[WS] Connecting:
[WS] Connected:
[WS] Registered:
[WS RECEIVE]
```

Example mobile:

```text
[WS] Connecting: ws://amen-mixer.local:8081/ws
```

Example Wails:

```text
[WS] Connecting: ws://127.0.0.1:8081/ws
```

---

## Browser Serial Message

When opening AMEN Mixer from a normal browser, logs such as:

```text
[SERIAL] Browser mode.
[SERIAL] Wails serial disabled.
```

are expected.

The mobile browser does not have access to the Wails serial runtime.

Serial communication belongs to the desktop host.

---

## Favicon 404

This warning:

```text
/favicon.ico 404
```

does not affect mixer functionality.

To remove it, add:

```text
desktop/frontend/public/favicon.ico
```

---

## mDNS Troubleshooting

If:

```text
amen-mixer.local
```

does not resolve:

### Check IP

```bash
ipconfig getifaddr en0
```

### Check backend

```text
http://<IP>:8081/
```

### Check Vite

```text
http://<IP>:5173
```

### Check network

Make sure both devices are on the same LAN.

mDNS may fail on networks using:

```text
Guest Isolation
Client Isolation
Multicast Blocking
```

Use direct IP access as the fallback.

---

## WebSocket Troubleshooting

### UI Opens but Mixer Does Not Sync

First check the browser console.

Expected:

```text
[WS] Connected
```

Then move a slider.

Expected:

```text
[WS RECEIVE]
```

On the backend:

```text
[CHANNEL UPDATE]
```

If the backend gets the message but another UI does not update, inspect the broadcast path.

---

## Realtime Debug Order

Always debug in this order:

```text
1. HTTP health
2. WebSocket connection
3. client.register
4. CHANNEL_UPDATE incoming
5. Audio Manager update
6. Broadcast
7. Remote Vue store
8. Remote component
```

This avoids debugging UI code when the actual problem is network-related.

---

## Development Ports

```text
5173  Vite frontend
8081  Go realtime server
```

Check port usage on macOS:

```bash
lsof -i :5173
```

and:

```bash
lsof -i :8081
```

---

## Format Go Code

After changing backend files:

```bash
gofmt -w *.go
```

For audio package:

```bash
gofmt -w backend/audio/*.go
```

For all split realtime files:

```bash
gofmt -w realtime*.go
```

---

## Go Validation

Run:

```bash
go test ./...
```

or at minimum:

```bash
go test ./backend/audio
```

Then:

```bash
wails dev
```

---

## Frontend Validation

From:

```text
desktop/frontend
```

run:

```bash
npm run build
```

This catches TypeScript/Vite build problems before Android synchronization.

---

## Android Development

Build frontend:

```bash
cd desktop/frontend
npm run build
```

Sync native project:

```bash
npx cap sync android
```

Open Android Studio:

```bash
npx cap open android
```

Debug APK:

```bash
cd android
./gradlew assembleDebug
```

---

## Android Network Model

The Android app is a remote mixer client.

It does not contain the Wails/Go backend.

Expected deployment:

```text
Windows/Desktop Host
        │
        │ Wi-Fi / LAN
        ▼
Android Device
```

The desktop application must be running for realtime mixer control.

---

## Firmware Development

ESP32 firmware is handled through PlatformIO.

The firmware project is separate from:

```text
desktop/
```

Do not attempt to compile PlatformIO firmware from Wails or Android Studio.

---

## Git Workflow

Check current branch:

```bash
git branch --show-current
```

Check changes:

```bash
git status
```

Commit:

```bash
git add .
git commit -m "docs: update realtime and mobile documentation"
```

First push for a new branch:

```bash
git push -u origin feature/android-app
```

Following pushes:

```bash
git push
```

---

## Files That Normally Should Not Be Committed

Check `.gitignore` for:

```text
node_modules/
.DS_Store
frontend/android/local.properties
```

Environment files should only be committed if they contain project-safe configuration and no machine-specific values or secrets.

With automatic hostname detection, a hardcoded LAN IP should not be required.

---

## Before Committing

Recommended checks:

```bash
gofmt -w *.go
gofmt -w backend/audio/*.go
go test ./...
```

Then:

```bash
cd frontend
npm run build
cd ..
```

Finally:

```bash
git status
```

---

## Current Development Limitation

During `wails dev`, mobile devices load the frontend from Vite:

```text
:5173
```

When a production Wails application is built, Vite development mode is not running.

The production mobile frontend therefore still needs a dedicated serving strategy.

Recommended future architecture:

```text
Vue build
   ↓
dist/
   ↓
Go HTTP server
   ↓
http://amen-mixer.local
```

That would allow mobile control without requiring:

```bash
wails dev
```

or a separate Vite server.

---

## Next Milestones

Recommended order:

```text
1. Confirm desktop ↔ mobile sync
2. Implement real Windows Core Audio
3. Route ESP32 commands through Audio Manager
4. Serve production mobile frontend from Go
5. Package Android client
6. Add reconnect/state recovery
7. Add presets/persistence
```
