#pragma once

#include <Arduino.h>

class NetworkManager
{

public:
    void begin();

    void loop();

    void sendEncoder(
        uint8_t channel,
        int value);

    bool isConnected();
};