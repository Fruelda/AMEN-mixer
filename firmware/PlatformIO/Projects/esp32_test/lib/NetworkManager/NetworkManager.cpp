#include "NetworkManager.h"

#include <WiFi.h>
#include <WebSocketsClient.h>
#include <ArduinoJson.h>

#include "../Config/Config.h"

WebSocketsClient webSocket;

bool wsConnected = false;

// =================================================
// BEGIN
// =================================================

void NetworkManager::begin()
{

    Serial.println();

    Serial.println(
        "=== AMEN NETWORK START ===");

    connectWiFi();

    connectWebSocket();
}

// =================================================
// WIFI
// =================================================

void NetworkManager::connectWiFi()
{

    Serial.println(
        "Connecting WiFi...");

    WiFi.mode(
        WIFI_STA);

    WiFi.begin(
        WIFI_SSID,
        WIFI_PASSWORD);

    while (
        WiFi.status() != WL_CONNECTED)
    {

        delay(500);

        Serial.print(
            ".");
    }

    Serial.println();

    Serial.println(
        "WiFi Connected");
}

// =================================================
// WEBSOCKET
// =================================================

void NetworkManager::connectWebSocket()
{

    Serial.println(
        "Connecting AMEN Server...");

    Serial.print(
        "Server : ");

    Serial.print(
        AMEN_HOST);

    Serial.print(
        ":");

    Serial.println(
        AMEN_PORT);

    webSocket.onEvent(

        [](WStype_t type,
           uint8_t *payload,
           size_t length)

        {
            switch (type)
            {

            case WStype_CONNECTED:

                Serial.println(
                    "WebSocket Connected");

                wsConnected = true;

                webSocket.sendTXT(

                    "{\"type\":\"device.register\",\"id\":\"amen-mixer-01\",\"name\":\"AMEN Hardware Mixer\"}"

                );

                break;

            case WStype_DISCONNECTED:

                Serial.println(
                    "WebSocket Disconnected");

                wsConnected = false;

                break;

            case WStype_TEXT:

                Serial.print(
                    "[SERVER] ");

                Serial.println(
                    (char *)payload);

                break;

            default:

                break;
            }
        }

    );

    webSocket.begin(

        AMEN_HOST,

        AMEN_PORT,

        AMEN_ENDPOINT

    );

    webSocket.setReconnectInterval(
        5000);
}

// =================================================
// LOOP
// =================================================

void NetworkManager::loop()
{

    webSocket.loop();
}

// =================================================
// SEND ENCODER
// =================================================

void NetworkManager::sendEncoder(

    uint8_t channel,

    int value

)
{

    if (
        !wsConnected)
    {
        return;
    }

    StaticJsonDocument<256> doc;

    doc["type"] =
        "mixer.command";

    doc["device"] =
        DEVICE_ID;

    doc["channel"] =
        channel;

    doc["value"] =
        value;

    String output;

    serializeJson(
        doc,
        output);

    webSocket.sendTXT(
        output);
}

// =================================================
// STATUS
// =================================================

bool NetworkManager::isConnected()
{

    return wsConnected;
}

bool NetworkManager::wifiConnected()
{

    return WiFi.status() == WL_CONNECTED;
}