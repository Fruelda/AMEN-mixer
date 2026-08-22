#include "WebSocketManager.h"

#include <WiFi.h>
#include <WebSocketsClient.h>

#include "../Config/Config.h"
#include "../Protocol/Protocol.h"
#include "../MessageHandler/MessageHandler.h"

WebSocketsClient webSocket;

WebSocketManager *instance = nullptr;

static void wsEvent(
    WStype_t type,
    uint8_t *payload,
    size_t length)
{

    if (!instance)
        return;

    switch (type)
    {

    case WStype_CONNECTED:

        Serial.println("WebSocket Connected");

        instance->setConnected(true);

        webSocket.sendTXT(
            createDeviceRegister().c_str());

        break;

    case WStype_DISCONNECTED:

        Serial.println("WebSocket Disconnected");

        instance->setConnected(false);

        break;

    case WStype_TEXT:

        MessageHandler::handle(payload);

        break;

    default:
        break;
    }
}

void WebSocketManager::begin()
{

    instance = this;

    webSocket.onEvent(wsEvent);

    webSocket.begin(
        AMEN_HOST,
        AMEN_PORT,
        AMEN_ENDPOINT);

    webSocket.setReconnectInterval(5000);
}

void WebSocketManager::loop()
{
    webSocket.loop();
}

void WebSocketManager::send(String message)
{

    if (!wsConnected)
        return;

    webSocket.sendTXT(
        message.c_str());
}

bool WebSocketManager::connected()
{
    return wsConnected;
}

void WebSocketManager::setConnected(bool state)
{
    wsConnected = state;
}