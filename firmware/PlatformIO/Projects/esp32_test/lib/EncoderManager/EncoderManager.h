#pragma once

#include <Arduino.h>

class AudioManager;

class EncoderManager
{

public:
    void begin();

    void update(AudioManager &audio);

private:
    // ==========================
    // Rotary Encoder State
    // ==========================

    int lastCLK1 = HIGH;
    int lastCLK2 = HIGH;
    int lastCLK3 = HIGH;
    int lastCLK4 = HIGH;
    int lastCLK5 = HIGH;
    int lastCLK6 = HIGH;

    // ==========================
    // Button State
    // ==========================

    bool lastButton1 = false;
    bool lastButton2 = false;
    bool lastButton3 = false;
    bool lastButton4 = false;
    bool lastButton5 = false;
    bool lastButton6 = false;

    // ==========================
    // Debounce
    // ==========================

    unsigned long lastDebounceTime = 0;

    const unsigned long debounceDelay = 50;

    // ==========================
    // Helper Function
    // ==========================

    void readEncoder(
        int clkPin,
        int dtPin,
        int &lastCLK,
        int encoderID);

    void readButton(
        int swPin,
        bool &lastButton,
        int encoderID);
};