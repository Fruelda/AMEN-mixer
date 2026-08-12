# AMEN Mixer Architecture

## Overview

AMEN Mixer uses a hub-based realtime architecture.

The Go backend running inside the Wails desktop application is the central communication point for all clients.

```text
                 ┌───────────────────────┐
                 │       Go Backend      │
                 │                       │
                 │  Realtime WebSocket   │
                 │       :8081/ws        │
                 │                       │
                 │     Audio Manager     │
                 │                       │
                 │    Serial Manager     │
                 └───────────┬───────────┘
                             │
             ┌───────────────┼───────────────┐
             │               │               │
             ▼               ▼               ▼
      Wails Desktop       Mobile UI        ESP32
          Vue          Browser/Android     Serial
```

Clients do not directly synchronize with other clients.

All authoritative state changes pass through the backend.

---

## Components

### Wails Desktop

Responsibilities:

* Starts the Go application
* Displays the native desktop UI
* Hosts the main mixer interface
* Provides direct Go bindings to Vue
* Runs the realtime backend
* Runs serial integration
* Starts mDNS discovery

Relevant files:

```text
desktop/app.go
desktop/app_audio.go
desktop/app_serial.go
desktop/main.go
```

---

## Realtime Server

The realtime server listens on:

```text
0.0.0.0:8081
```

HTTP health endpoint:

```text
/
```

WebSocket endpoint:

```text
/ws
```

Responsibilities:

* Accept WebSocket connections
* Register clients
* Register hardware devices
* Route messages
* Broadcast updates
* Track connected devices
* Send initial state to new clients

Relevant files:

```text
realtime.go
realtime_client.go
realtime_handler.go
realtime_broadcast.go
realtime_devices.go
realtime_server.go
state_sync.go
```

---

## Realtime File Responsibilities

### `realtime.go`

Contains:

* `WSClient`
* `RealtimeServer`
* WebSocket upgrader
* Server constructor
* Channel update callback registration

### `realtime_client.go`

Contains:

* WebSocket HTTP upgrade
* Client connection lifecycle
* Read loop
* Disconnect handling

### `realtime_handler.go`

Contains message routing and handlers.

Example incoming messages:

```text
client.register
device.register
CHANNEL_UPDATE
mixer.command
```

### `realtime_broadcast.go`

Contains:

* Raw broadcast
* JSON broadcast
* Channel broadcasts
* Device status broadcasts

### `realtime_devices.go`

Contains:

* Connected client list
* Device information
* Client removal
* Client count

### `realtime_server.go`

Starts:

```text
HTTP :8081
WebSocket /ws
Health /
```

### `state_sync.go`

Maintains the latest realtime channel snapshot and sends it to newly connected clients.

---

## Audio Layer

The Audio Manager is intended to be the authoritative mixer state.

Relevant files:

```text
backend/audio/manager.go
backend/audio/manager_state.go
backend/audio/manager_update.go
backend/audio/windows.go
backend/audio/events.go
```

### `manager.go`

Owns:

```text
channels
WindowsAudio
UpdateListener
```

### `manager_state.go`

Responsible for:

```text
GetChannels()
GetChannel()
volume normalization
```

### `manager_update.go`

Responsible for:

```text
ApplyChannelUpdate()
SetVolume()
SetMuted()
```

All client updates should eventually pass through:

```text
ApplyChannelUpdate()
```

This avoids multiple independent sources of mixer state.

---

## Mobile → Desktop Flow

When the user changes a slider from an iPhone:

```text
VolumeSlider
    ↓
Vue Store
    ↓
sendRealtime()
    ↓
WebSocket
    ↓
CHANNEL_UPDATE
    ↓
RealtimeServer.handleChannelUpdate()
    ↓
AudioManager.ApplyChannelUpdate()
    ↓
audioBridge.OnChannelUpdate()
    ↓
BroadcastChannelUpdate()
    ↓
Wails Desktop
    ↓
Vue Store
    ↓
Desktop Slider
```

The mobile client does not directly change the Wails UI.

It changes backend state, and the backend broadcasts the result.

---

## Desktop → Mobile Flow

A desktop change follows the same authoritative path:

```text
Wails UI
    ↓
App.SetVolume()
    ↓
AudioManager.SetVolume()
    ↓
ApplyChannelUpdate()
    ↓
audioBridge
    ↓
BroadcastChannelUpdate()
    ↓
Mobile WebSocket
    ↓
Vue Store
    ↓
Mobile UI
```

This keeps desktop and mobile synchronized.

---

## Initial State Flow

When a client first connects:

```text
WebSocket connect
    ↓
client.register
    ↓
RealtimeServer
    ↓
SendState(client)
    ↓
STATE
    ↓
Client Store
```

This prevents a newly opened phone from starting with stale mixer values.

---

## Client Registration

Typical desktop registration:

```json
{
  "type": "client.register",
  "id": "desktop-...",
  "name": "Wails Desktop",
  "clientType": "desktop"
}
```

Typical mobile registration:

```json
{
  "type": "client.register",
  "id": "iphone-...",
  "name": "iPhone",
  "clientType": "mobile"
}
```

Client IDs are persisted in browser storage when available.

---

## Channel Update

Example:

```json
{
  "type": "CHANNEL_UPDATE",
  "channel": {
    "id": 2,
    "volume": 40,
    "muted": false
  }
}
```

`volume` and `muted` are treated as optional update fields.

This allows messages such as:

```json
{
  "type": "CHANNEL_UPDATE",
  "channel": {
    "id": 2,
    "volume": 50
  }
}
```

or:

```json
{
  "type": "CHANNEL_UPDATE",
  "channel": {
    "id": 2,
    "muted": true
  }
}
```

---

## Network Discovery

The desktop publishes:

```text
amen-mixer.local
```

using mDNS.

Development UI:

```text
http://amen-mixer.local:5173
```

Realtime server:

```text
ws://amen-mixer.local:8081/ws
```

### Desktop WebSocket

The Wails client does not need LAN routing.

It uses:

```text
ws://127.0.0.1:8081/ws
```

### Mobile WebSocket

The browser uses:

```text
window.location.hostname
```

Therefore opening:

```text
http://amen-mixer.local:5173
```

produces:

```text
ws://amen-mixer.local:8081/ws
```

Opening:

```text
http://192.168.18.11:5173
```

produces:

```text
ws://192.168.18.11:8081/ws
```

No LAN IP is hardcoded in the frontend.

---

## Serial Architecture

ESP32 communication is separate from WebSocket communication.

```text
ESP32
   │
   │ USB Serial
   ▼
Serial Manager
   │
   ├── Wails Event
   │
   └── Realtime Broadcast
```

This allows hardware input to eventually update the same mixer state used by desktop and mobile clients.

---

## Android Architecture

The Android application is based on the same Vue frontend using Capacitor.

```text
Vue
 ↓
Vite Build
 ↓
dist/
 ↓
Capacitor
 ↓
Android Project
 ↓
Android WebView
```

Android does not run the Go/Wails backend.

The desktop computer remains the mixer host.

```text
Android
   │
   │ LAN WebSocket
   ▼
Desktop Go Backend
```

---

## mDNS

`mdns.go` publishes the AMEN Mixer hostname on the local network.

At startup:

```text
startMDNS()
```

The backend discovers its active LAN IPv4 address and advertises:

```text
amen-mixer.local
```

If the network changes, restarting the application refreshes the advertised IP.

---

## Current Windows Audio Layer

Current architecture:

```text
AudioManager
    ↓
WindowsAudio
    ↓
SetVolume()
SetMuted()
```

The current `WindowsAudio` layer is the integration boundary.

The next backend milestone is implementing Windows Core Audio session control.

Target:

```text
WindowsAudio
    ↓
Windows Core Audio API
    ↓
Application Audio Session
    ↓
ISimpleAudioVolume
```

This should eventually support:

```text
Browser
Spotify
Discord
Game
Other application sessions
```

---

## Design Rule

The most important architectural rule is:

> The Go Audio Manager should be the authoritative source of mixer state.

Avoid this:

```text
Mobile state
Desktop state
ESP32 state
```

operating independently.

Prefer:

```text
Mobile ─┐
Desktop ├──► Audio Manager ───► Broadcast ───► All clients
ESP32  ─┘
```

This prevents realtime synchronization loops and state conflicts.
