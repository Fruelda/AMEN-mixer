#pragma once

#include <Arduino.h>

String createDeviceRegister();

String createEncoderCommand(
    uint8_t channel,
    int value);