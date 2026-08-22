#pragma once

#include <Arduino.h>

class AudioManager;

class NetworkManager
{

public:
    void begin();

    void loop();

    void sendEncoder(
        uint8_t channel,
        int value);

    bool isConnected();

private:
    void connectWiFi();

    void connectWebSocket();

    void checkWiFi();

    void checkWebSocket();

    unsigned long lastWifiAttempt = 0;

    unsigned long lastWsAttempt = 0;
};