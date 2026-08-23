#pragma once

#include <Arduino.h>

#include "../Config/Config.h"

class AudioManager;

class EncoderManager
{

public:
    void begin();

    void update(
        AudioManager &audio);

private:
    void readEncoder(
        int clkPin,
        int dtPin,
        int &lastCLK,
        int encoderID);

    void readButton(
        int swPin,
        bool &lastButton,
        int encoderID);

private:
    int lastCLK1;
    int lastCLK2;
    int lastCLK3;
    int lastCLK4;
    int lastCLK5;
    int lastCLK6;

    bool lastButton1 = false;
    bool lastButton2 = false;
    bool lastButton3 = false;
    bool lastButton4 = false;
    bool lastButton5 = false;
    bool lastButton6 = false;

    unsigned long lastDebounceTime = 0;

    unsigned long debounceDelay = 50;
};