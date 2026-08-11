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

    if (
        WiFi.status() ==
        WL_CONNECTED)
    {
        return;
    }

    Serial.println(
        "Connecting WiFi...");

    WiFi.mode(
        WIFI_STA);

    WiFi.begin(

        WIFI_SSID,

        WIFI_PASSWORD

    );

    while (

        WiFi.status() !=
        WL_CONNECTED

    )
    {

        delay(500);

        Serial.print(
            ".");
    }

    Serial.println();

    Serial.println(
        "WiFi Connected");

    Serial.print(
        "IP : ");

    Serial.println(
        WiFi.localIP());
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

            {

                Serial.println(
                    "WebSocket Connected");

                wsConnected = true;

                StaticJsonDocument<256> doc;

                doc["type"] =
                    "device.register";

                doc["id"] =
                    DEVICE_ID;

                doc["name"] =
                    DEVICE_NAME;

                String output;

                serializeJson(

                    doc,

                    output

                );

                webSocket.sendTXT(

                    output

                );

                break;
            }

            case WStype_TEXT:

            {

                Serial.print(
                    "SERVER: ");

                Serial.println(

                    (char *)payload

                );

                break;
            }

            case WStype_DISCONNECTED:

            {

                Serial.println(
                    "WebSocket Disconnected");

                wsConnected = false;

                break;
            }

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

        5000

    );
}

// =================================================
// LOOP
// =================================================

void NetworkManager::loop()
{

    webSocket.loop();
}

// =================================================
// SEND
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

        output

    );

    webSocket.sendTXT(

        output

    );
}

bool NetworkManager::isConnected()
{

    return wsConnected;
}