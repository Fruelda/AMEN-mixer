#pragma once

#include <Arduino.h>

class WebSocketManager
{

public:
    void begin();

    void loop();

    void send(
        String message);

    bool connected();

    void setConnected(
        bool state);

private:
    bool wsConnected = false;
};