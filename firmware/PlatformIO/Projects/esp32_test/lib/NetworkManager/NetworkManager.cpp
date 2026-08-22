#include "NetworkManager.h"

#include "../WiFiManager/WiFiManager.h"

#include "../WebSocketManager/WebSocketManager.h"

#include "../Protocol/Protocol.h"

WiFiManager wifiManager;

WebSocketManager socketManager;

void NetworkManager::begin()
{

    Serial.println(
        "=== AMEN NETWORK START ===");

    wifiManager.begin();
}

void NetworkManager::loop()
{

    wifiManager.loop();

    if (
        wifiManager.connected())
    {

        socketManager.begin();

        socketManager.loop();
    }
}

void NetworkManager::sendEncoder(
    uint8_t channel,
    int value)
{

    socketManager.send(
        createEncoderCommand(
            channel,
            value));
}

bool NetworkManager::isConnected()
{

    return socketManager.connected();
}