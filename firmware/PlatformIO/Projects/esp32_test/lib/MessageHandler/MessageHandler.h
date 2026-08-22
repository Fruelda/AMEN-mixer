#pragma once

#include <Arduino.h>

class MessageHandler
{

public:
    static void handle(
        uint8_t *payload);
};