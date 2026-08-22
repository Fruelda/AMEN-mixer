#include "NetworkManager.h"

#include <WiFi.h>

#include "../Config/Config.h"

#include "../WebSocketManager/WebSocketManager.h"

#include "../Protocol/Protocol.h"

WebSocketManager socketManager;

void NetworkManager::begin()
{

    Serial.println(
        "=== AMEN NETWORK START ===");

    WiFi.mode(
        WIFI_STA);

    WiFi.begin(

        WIFI_SSID,

        WIFI_PASSWORD

    );

    Serial.println(
        "Connecting WiFi...");

    while (
        WiFi.status() !=
        WL_CONNECTED)
    {

        delay(500);

        Serial.print(".");
    }

    Serial.println();

    Serial.println(
        "WiFi Connected");

    Serial.print(
        "IP : ");

    Serial.println(
        WiFi.localIP());

    socketManager.begin();
}

void NetworkManager::loop()
{

    socketManager.loop();
}

void NetworkManager::sendEncoder(

    uint8_t channel,

    int value

)
{

    socketManager.send(

        createEncoderCommand(

            channel,

            value

            )

    );
}

bool NetworkManager::isConnected()
{

    return socketManager.connected();
}