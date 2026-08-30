#include "EncoderManager.h"

#include "AudioManager.h"
#include "Config.h"
#include "NetworkManager.h"

extern NetworkManager network;

void EncoderManager::begin()
{
    // ==========================
    // Encoder 1
    // ==========================
    pinMode(EC1_CLK, INPUT_PULLUP);
    pinMode(EC1_DT, INPUT_PULLUP);
    pinMode(EC1_SW, INPUT_PULLUP);

    // ==========================
    // Encoder 2
    // ==========================
    pinMode(EC2_CLK, INPUT_PULLUP);
    pinMode(EC2_DT, INPUT_PULLUP);
    pinMode(EC2_SW, INPUT_PULLUP);

    // ==========================
    // Encoder 3
    // ==========================
    pinMode(EC3_CLK, INPUT_PULLUP);
    pinMode(EC3_DT, INPUT_PULLUP);
    pinMode(EC3_SW, INPUT_PULLUP);

    // ==========================
    // Encoder 4
    // ==========================
    pinMode(EC4_CLK, INPUT_PULLUP);
    pinMode(EC4_DT, INPUT_PULLUP);
    pinMode(EC4_SW, INPUT_PULLUP);

    // ==========================
    // Encoder 5
    // ==========================
    pinMode(EC5_CLK, INPUT_PULLUP);
    pinMode(EC5_DT, INPUT_PULLUP);
    pinMode(EC5_SW, INPUT_PULLUP);

    // ==========================
    // Encoder 6
    // ==========================
    pinMode(EC6_CLK, INPUT_PULLUP);
    pinMode(EC6_DT, INPUT_PULLUP);
    pinMode(EC6_SW, INPUT_PULLUP);

    // Store initial rotary state.
    lastCLK1 = digitalRead(EC1_CLK);
    lastCLK2 = digitalRead(EC2_CLK);
    lastCLK3 = digitalRead(EC3_CLK);
    lastCLK4 = digitalRead(EC4_CLK);
    lastCLK5 = digitalRead(EC5_CLK);
    lastCLK6 = digitalRead(EC6_CLK);

    // Store initial button state so boot does not create a false press.
    lastButton1 = digitalRead(EC1_SW) == LOW;
    lastButton2 = digitalRead(EC2_SW) == LOW;
    lastButton3 = digitalRead(EC3_SW) == LOW;
    lastButton4 = digitalRead(EC4_SW) == LOW;
    lastButton5 = digitalRead(EC5_SW) == LOW;
    lastButton6 = digitalRead(EC6_SW) == LOW;

    Serial.println("Encoder Manager Ready");
}

void EncoderManager::update(AudioManager &audio)
{
    (void)audio;

    // ==========================
    // Rotary Encoder
    // ==========================
    readEncoder(EC1_CLK, EC1_DT, lastCLK1, 1);
    readEncoder(EC2_CLK, EC2_DT, lastCLK2, 2);
    readEncoder(EC3_CLK, EC3_DT, lastCLK3, 3);
    readEncoder(EC4_CLK, EC4_DT, lastCLK4, 4);
    readEncoder(EC5_CLK, EC5_DT, lastCLK5, 5);
    readEncoder(EC6_CLK, EC6_DT, lastCLK6, 6);

    // ==========================
    // Button
    // ==========================
    readButton(EC1_SW, lastButton1, 1);
    readButton(EC2_SW, lastButton2, 2);
    readButton(EC3_SW, lastButton3, 3);
    readButton(EC4_SW, lastButton4, 4);
    readButton(EC5_SW, lastButton5, 5);
    readButton(EC6_SW, lastButton6, 6);
}

void EncoderManager::readEncoder(
    int clkPin,
    int dtPin,
    int &lastCLK,
    int encoderID)
{
    const int clk = digitalRead(clkPin);

    if (clk == lastCLK)
    {
        return;
    }

    int direction;

    if (digitalRead(dtPin) != clk)
    {
        direction = 1;
    }
    else
    {
        direction = -1;
    }

    // This is the USB serial protocol consumed by the Wails backend.
    Serial.printf("ENC,%d,%d\n", encoderID, direction);

    // Network is optional. If no WebSocket exists this simply returns.
    network.sendEncoder(encoderID, direction);

    lastCLK = clk;
}

void EncoderManager::readButton(
    int swPin,
    bool &lastButton,
    int encoderID)
{
    const bool current = digitalRead(swPin) == LOW;

    if (current == lastButton)
    {
        return;
    }

    if (millis() - lastDebounceTime <= debounceDelay)
    {
        return;
    }

    lastDebounceTime = millis();
    lastButton = current;

    if (!current)
    {
        return;
    }

    // Desktop protocol defines 1 as a button press.
    Serial.printf("BTN,%d,1\n", encoderID);

    // Keep existing network behavior for compatibility.
    network.sendEncoder(encoderID, 0);
}
