#pragma once

#include <Arduino.h>

#define BAUD_RATE 115200

// ==========================
// Encoder 1
// ==========================

#define EC1_CLK 4
#define EC1_DT 5
#define EC1_SW 6

// ==========================
// Encoder 2
// ==========================

#define EC2_CLK 7
#define EC2_DT 8
#define EC2_SW 9

// ==========================
// Encoder 3
// ==========================

#define EC3_CLK 10
#define EC3_DT 11
#define EC3_SW 12

// ==========================
// Encoder 4
// ==========================

#define EC4_CLK 13
#define EC4_DT 14
#define EC4_SW 15

// ==========================
// Encoder 5
// ==========================

#define EC5_CLK 16
#define EC5_DT 17
#define EC5_SW 18

// ==========================
// Encoder 6
// ==========================

#define EC6_CLK 21
#define EC6_DT 38
#define EC6_SW 39

// ===============================
// NETWORK MODE
// ===============================
//
// Desktop control is authoritative over USB serial.
// Keep WiFi realtime disabled by default so the ESP32 does not maintain a
// second connection to a host IP while it is already attached over USB.
// Set this to true only when intentionally running the hardware over WiFi.

#define AMEN_ENABLE_WIFI_REALTIME false

// ===============================
// WIFI CONFIG
// ===============================

#define WIFI_SSID "Yoga"
#define WIFI_PASSWORD "12345678"

// ===============================
// AMEN SERVER
// ===============================

#define AMEN_HOST "192.168.1.44"
#define AMEN_PORT 8081
#define AMEN_ENDPOINT "/ws"

// ===============================
// DEVICE
// ===============================

#define DEVICE_ID "amen-mixer-01"
#define DEVICE_NAME "AMEN Hardware Mixer"

// ===============================
// SYSTEM
// ===============================

#define DEBUG_SERIAL true
