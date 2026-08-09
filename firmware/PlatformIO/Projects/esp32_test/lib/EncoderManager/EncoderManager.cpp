#include "EncoderManager.h"

#include "AudioManager.h"
#include "Config.h"

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

    // simpan posisi awal

    lastCLK1 = digitalRead(EC1_CLK);
    lastCLK2 = digitalRead(EC2_CLK);
    lastCLK3 = digitalRead(EC3_CLK);
    lastCLK4 = digitalRead(EC4_CLK);
    lastCLK5 = digitalRead(EC5_CLK);
    lastCLK6 = digitalRead(EC6_CLK);

    Serial.println("Encoder Manager Ready");
}

void EncoderManager::update(AudioManager &audio)
{
    (void)audio;

    // Rotary

    readEncoder(
        EC1_CLK,
        EC1_DT,
        lastCLK1,
        1);

    readEncoder(
        EC2_CLK,
        EC2_DT,
        lastCLK2,
        2);

    readEncoder(
        EC3_CLK,
        EC3_DT,
        lastCLK3,
        3);

    readEncoder(
        EC4_CLK,
        EC4_DT,
        lastCLK4,
        4);

    readEncoder(
        EC5_CLK,
        EC5_DT,
        lastCLK5,
        5);

    readEncoder(
        EC6_CLK,
        EC6_DT,
        lastCLK6,
        6);

    // Button

    readButton(
        EC1_SW,
        lastButton1,
        1);

    // readButton(
    //     EC2_SW,
    //     lastButton2,
    //     2);

    readButton(
        EC3_SW,
        lastButton3,
        3);

    readButton(
        EC4_SW,
        lastButton4,
        4);

    readButton(
        EC5_SW,
        lastButton5,
        5);

    readButton(
        EC6_SW,
        lastButton6,
        6);
}

void EncoderManager::readEncoder(
    int clkPin,
    int dtPin,
    int &lastCLK,
    int encoderID)

{

    int clk = digitalRead(clkPin);

    if (clk != lastCLK)
    {

        if (digitalRead(dtPin) != clk)
        {

            Serial.printf(
                "ENC,%d,1\n",
                encoderID);
        }
        else
        {

            Serial.printf(
                "ENC,%d,-1\n",
                encoderID);
        }

        lastCLK = clk;
    }
}

void EncoderManager::readButton(
    int swPin,
    bool &lastButton,
    int encoderID)

{

    bool current = digitalRead(swPin) == LOW;

    if (current != lastButton)
    {

        if (millis() - lastDebounceTime > debounceDelay)
        {

            lastDebounceTime = millis();

            if (current)
            {

                Serial.printf(
                    "BTN,%d,0\n",
                    encoderID);
            }

            lastButton = current;
        }
    }
}