#pragma once

#include <Arduino.h>

class AudioManager;

class EncoderManager
{

public:
    void begin();

    void update(
        AudioManager &audio);

private:
    void readEncoder(
        int index);

    void readButton(
        int index);

private:
    int lastCLK[6];

    bool lastButton[6];

    unsigned long lastDebounceTime = 0;

    const unsigned long debounceDelay = 50;
};