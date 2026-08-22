#pragma once

#include <Arduino.h>

class WiFiManager
{

public:
    void begin();

    void loop();

    bool connected();
};