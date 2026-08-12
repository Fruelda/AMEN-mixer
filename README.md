# AMEN Mixer

AMEN Mixer is a realtime audio mixer system built with Wails, Go, Vue, WebSocket, mobile clients, and optional ESP32 hardware controls.

The desktop application acts as the main host. Mobile devices connect to the desktop over the local network and stay synchronized through the Go realtime backend.

## Features

Current implementation includes:

* Wails desktop application
* Vue 3 + TypeScript frontend
* Go realtime backend
* WebSocket communication
* Desktop ↔ mobile synchronization
* Realtime volume and mute updates
* Initial mixer state synchronization
* Connected device/client tracking
* mDNS local discovery using `amen-mixer.local`
* ESP32 serial integration layer
* PlatformIO firmware project
* Capacitor Android project

## Architecture

```text
                     AMEN Mixer
                         │
                         ▼
                ┌─────────────────┐
                │   Go Backend    │
                │                 │
                │ WebSocket :8081 │
                │ Audio Manager   │
                │ Serial Manager  │
                └────────┬────────┘
                         │
          ┌──────────────┼──────────────┐
          │              │              │
          ▼              ▼              ▼
     Wails Desktop    Mobile UI       ESP32
        Vue         Browser/Android   Hardware
```

The Go backend is the central source of mixer state.

Clients do not communicate directly with each other.

## Project Structure

```text
DEEJ/
├── README.md
├── docs/
│   ├── architecture.md
│   └── development.md
│
├── desktop/
│   ├── backend/
│   ├── frontend/
│   ├── app.go
│   ├── app_audio.go
│   ├── app_serial.go
│   ├── mdns.go
│   ├── realtime.go
│   ├── realtime_client.go
│   ├── realtime_handler.go
│   ├── realtime_broadcast.go
│   ├── realtime_devices.go
│   ├── realtime_server.go
│   ├── state_sync.go
│   └── README.md
│
└── firmware/
```

## Technology Stack

Desktop:

* Go
* Wails v2
* Vue 3
* TypeScript
* Vite

Realtime:

* WebSocket
* JSON messaging
* mDNS / Bonjour

Mobile:

* Browser client
* Capacitor
* Android Studio

Hardware:

* ESP32
* PlatformIO
* Serial communication

## Development

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

Run development mode:

```bash
wails dev
```

## Mobile Access

Preferred local URL:

```text
http://amen-mixer.local:5173
```

If mDNS is unavailable, get the host IP:

```bash
ipconfig getifaddr en0
```

Then open:

```text
http://<HOST-IP>:5173
```

Example:

```text
http://192.168.18.11:5173
```

## Realtime Server

Health endpoint:

```text
http://amen-mixer.local:8081/
```

Expected response:

```text
AMEN MIXER REALTIME SERVER OK
```

WebSocket endpoint:

```text
ws://amen-mixer.local:8081/ws
```

The Wails desktop application connects locally using:

```text
ws://127.0.0.1:8081/ws
```

## Realtime Flow

Mobile update:

```text
Mobile Slider
    ↓
CHANNEL_UPDATE
    ↓
Go Realtime Server
    ↓
Audio Manager
    ↓
BroadcastChannelUpdate
    ↓
Desktop + Mobile Clients
```

Desktop update:

```text
Wails UI
    ↓
Audio Manager
    ↓
BroadcastChannelUpdate
    ↓
Mobile Clients
```

## Android

The Android project is located at:

```text
desktop/frontend/android
```

Build the frontend:

```bash
cd desktop/frontend
npm run build
```

Sync Capacitor:

```bash
npx cap sync android
```

Open Android Studio:

```bash
npx cap open android
```

## ESP32

Firmware is maintained separately under:

```text
firmware/
```

PlatformIO handles ESP32 compilation and upload.

The ESP32 firmware is not bundled into the desktop or Android application.

## Current Limitation

The realtime state architecture is implemented, but actual Windows per-application audio control still needs the Windows Core Audio implementation inside:

```text
desktop/backend/audio/windows.go
```

Target architecture:

```text
Audio Manager
    ↓
Windows Audio Layer
    ↓
Windows Core Audio
    ↓
Application Audio Session
```

## Documentation

Detailed documentation:

* [Architecture](docs/architecture.md)
* [Development Guide](docs/development.md)
* [Desktop README](desktop/README.md)

## Git

Example feature branch:

```bash
git checkout -b feature/android-app
```

First push:

```bash
git push -u origin feature/android-app
```

Next pushes:

```bash
git push
```

## Roadmap

* Windows Core Audio integration
* Full desktop ↔ mobile mixer synchronization
* ESP32 mixer control
* Production mobile frontend serving
* Android packaging
* Presets and persistence
* Automatic audio session discovery
