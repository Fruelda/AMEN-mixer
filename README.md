# AMEN Mixer

AMEN Mixer is a custom desktop audio mixer application designed to
provide real-time control over Windows audio channels through a modern
UI interface.

The project combines a Wails desktop application, a Go backend, a Vue
frontend, and an ESP32 firmware layer for external hardware control.

## Features

-   Real-time audio channel management
-   Modern dark themed mixer interface
-   Per-application volume control
-   WebSocket based realtime communication
-   Serial communication support for external controllers
-   ESP32 firmware integration
-   Rotary encoder hardware support

## Project Structure

``` text
DEEJ/
├── desktop/
│   ├── backend/
│   │   ├── audio/       # Windows audio management
│   │   ├── serial/      # Serial communication layer
│   │   ├── protocol/    # Command definitions
│   │   └── models/      # Shared data models
│   ├── frontend/
│   │   ├── src/
│   │   │   ├── components/mixer/ # Mixer UI components
│   │   │   ├── stores/            # Application state
│   │   │   └── composables/       # Realtime hooks
│   │   └── package.json
│   ├── main.go
│   └── wails.json
│
└── firmware/
    └── PlatformIO/
        └── Projects/
            └── esp32_test/
                ├── lib/
                │   ├── AudioManager/
                │   ├── EncoderManager/
                │   └── Config/
                └── src/main.cpp
```

## Technology Stack

### Desktop Application

-   Go
-   Wails v2
-   Vue 3
-   Vite
-   Tailwind CSS

### Communication

-   WebSocket
-   Serial communication

### Hardware

-   ESP32
-   Rotary encoder input
-   PlatformIO firmware environment

## Development Setup

### Requirements

Install:

-   Go 1.25+
-   Node.js and npm
-   Wails CLI
-   PlatformIO (for firmware)

## Running Desktop Development

Navigate to:

``` bash
cd desktop
```

Install frontend dependencies:

``` bash
cd frontend
npm install
```

Run development mode:

``` bash
cd ..
wails dev
```

## Building Production Application

From the desktop directory:

``` bash
wails build
```

The generated application will be placed in the Wails build output
directory.

## Environment Configuration

Frontend environment variables are stored locally.

Example:

``` env
VITE_WS_URL=ws://localhost:8765/ws
```

Create your own `.env` file based on your environment.

## Architecture Overview

The system is divided into three main layers:

### Frontend Layer

Responsible for:

-   Mixer visualization
-   User interaction
-   Channel controls
-   Realtime state updates

### Backend Layer

Responsible for:

-   Audio session management
-   Windows audio control
-   Serial communication
-   WebSocket events

### Firmware Layer

Responsible for:

-   Reading physical controls
-   Sending encoder commands
-   Communicating with desktop application

## Future Development

Potential improvements:

-   Hardware enclosure design
-   More audio device support
-   User profiles
-   Persistent mixer presets
-   Cross-platform audio backend
-   Plugin system

## License

Add project license information here.
