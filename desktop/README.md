# AMEN Mixer Desktop

Desktop host application for AMEN Mixer.

Built with:

* Wails
* Go
* Vue 3
* TypeScript
* Vite
* WebSocket
* mDNS
* Serial communication

## Responsibilities

The desktop application acts as the main AMEN Mixer host.

It is responsible for:

* Running the Go backend
* Running the realtime WebSocket server on port `8081`
* Managing mixer state
* Synchronizing desktop and mobile clients
* Handling ESP32 serial communication
* Advertising `amen-mixer.local` using mDNS
* Providing the Wails desktop interface

## Development

Run from this directory:

```bash
wails dev
```

During development:

```text
Wails Desktop
    ↓
Go Backend
    ↓
WebSocket :8081
```

Mobile devices can access the development frontend using:

```text
http://amen-mixer.local:5173
```

Fallback using the host LAN IP:

```text
http://<HOST-IP>:5173
```

## Realtime Server

Health endpoint:

```text
http://amen-mixer.local:8081/
```

WebSocket endpoint:

```text
ws://amen-mixer.local:8081/ws
```

The Wails desktop client connects locally using:

```text
ws://127.0.0.1:8081/ws
```

## Project Documentation

See the main documentation:

* [Project README](../README.md)
* [Architecture](../docs/architecture.md)
* [Development Guide](../docs/development.md)
