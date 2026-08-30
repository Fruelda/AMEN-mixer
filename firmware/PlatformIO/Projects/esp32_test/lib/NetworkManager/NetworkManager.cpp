#include "NetworkManager.h"

#include <WiFi.h>
#include <WebSocketsClient.h>
#include <ArduinoJson.h>

#include "../Config/Config.h"

WebSocketsClient webSocket;

bool wsConnected = false;

namespace
{
constexpr unsigned long WIFI_CONNECT_TIMEOUT_MS = 6000;
}

// =================================================
// BEGIN
// =================================================

void NetworkManager::begin()
{
    if (!AMEN_ENABLE_WIFI_REALTIME)
    {
        Serial.println("[NETWORK] WiFi realtime disabled. USB serial mode active.");
        return;
    }

    Serial.println();
    Serial.println("=== AMEN NETWORK START ===");

    // Network is optional for the USB mixer path.
    // Never block encoder startup forever just because WiFi is unavailable.
    connectWiFi();

    if (WiFi.status() == WL_CONNECTED)
    {
        connectWebSocket();
    }
    else
    {
        Serial.println("[NETWORK] WiFi unavailable. USB serial mode stays active.");
    }
}

// =================================================
// WIFI
// =================================================

void NetworkManager::connectWiFi()
{
    if (WiFi.status() == WL_CONNECTED)
    {
        return;
    }

    Serial.println("Connecting WiFi...");

    WiFi.mode(WIFI_STA);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    const unsigned long startedAt = millis();

    while (WiFi.status() != WL_CONNECTED)
    {
        if (millis() - startedAt >= WIFI_CONNECT_TIMEOUT_MS)
        {
            Serial.println();
            Serial.println("[NETWORK] WiFi connect timeout.");
            return;
        }

        delay(250);
        Serial.print(".");
    }

    Serial.println();
    Serial.println("WiFi Connected");
    Serial.print("IP : ");
    Serial.println(WiFi.localIP());
}

// =================================================
// WEBSOCKET
// =================================================

void NetworkManager::connectWebSocket()
{
    if (WiFi.status() != WL_CONNECTED)
    {
        return;
    }

    Serial.println("Connecting AMEN Server...");
    Serial.print("Server : ");
    Serial.print(AMEN_HOST);
    Serial.print(":");
    Serial.println(AMEN_PORT);

    webSocket.onEvent(
        [](WStype_t type, uint8_t *payload, size_t length)
        {
            (void)length;

            switch (type)
            {
            case WStype_CONNECTED:
            {
                Serial.println("WebSocket Connected");
                wsConnected = true;

                StaticJsonDocument<256> doc;
                doc["type"] = "device.register";
                doc["id"] = DEVICE_ID;
                doc["name"] = DEVICE_NAME;

                String output;
                serializeJson(doc, output);
                webSocket.sendTXT(output);
                break;
            }

            case WStype_TEXT:
            {
                Serial.print("SERVER: ");
                Serial.println((char *)payload);
                break;
            }

            case WStype_DISCONNECTED:
            {
                Serial.println("WebSocket Disconnected");
                wsConnected = false;
                break;
            }

            default:
                break;
            }
        });

    webSocket.begin(AMEN_HOST, AMEN_PORT, AMEN_ENDPOINT);

    // Keep an intentionally enabled WiFi connection alive while idle.
    // Gorilla WebSocket uses the normal WebSocket ping/pong control frames,
    // so this is compatible with the desktop server.
    webSocket.enableHeartbeat(15000, 3000, 2);
    webSocket.setReconnectInterval(5000);
}

// =================================================
// LOOP
// =================================================

void NetworkManager::loop()
{
    if (!AMEN_ENABLE_WIFI_REALTIME)
    {
        return;
    }

    // USB serial control does not depend on this branch.
    if (WiFi.status() == WL_CONNECTED)
    {
        webSocket.loop();
    }
}

// =================================================
// SEND
// =================================================

void NetworkManager::sendEncoder(uint8_t channel, int value)
{
    if (!AMEN_ENABLE_WIFI_REALTIME)
    {
        return;
    }

    if (!wsConnected)
    {
        return;
    }

    StaticJsonDocument<256> doc;
    doc["type"] = "mixer.command";
    doc["device"] = DEVICE_ID;
    doc["channel"] = channel;
    doc["value"] = value;

    String output;
    serializeJson(doc, output);
    webSocket.sendTXT(output);
}

bool NetworkManager::isConnected()
{
    return wsConnected;
}
