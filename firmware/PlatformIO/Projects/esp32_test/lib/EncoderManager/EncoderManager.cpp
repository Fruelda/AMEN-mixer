#include "EncoderManager.h"

#include "AudioManager.h"
#include "Config.h"
#include "NetworkManager.h"

extern NetworkManager network;

struct EncoderConfig
{
    uint8_t clk;
    uint8_t dt;
    uint8_t sw;
    uint8_t id;
    bool reverse;
};

EncoderConfig encoders[] =
    {
        {EC1_CLK, EC1_DT, EC1_SW, 1, EC1_REVERSE},
        {EC2_CLK, EC2_DT, EC2_SW, 2, EC2_REVERSE},
        {EC3_CLK, EC3_DT, EC3_SW, 3, EC3_REVERSE},
        {EC4_CLK, EC4_DT, EC4_SW, 4, EC4_REVERSE},
        {EC5_CLK, EC5_DT, EC5_SW, 5, EC5_REVERSE},
        {EC6_CLK, EC6_DT, EC6_SW, 6, EC6_REVERSE}};

#define ENCODER_COUNT 6

void EncoderManager::begin()
{

    for (int i = 0; i < ENCODER_COUNT; i++)
    {

        pinMode(encoders[i].clk, INPUT_PULLUP);
        pinMode(encoders[i].dt, INPUT_PULLUP);
        pinMode(encoders[i].sw, INPUT_PULLUP);

        lastCLK[i] =
            digitalRead(encoders[i].clk);

        lastButton[i] = false;
    }

    Serial.println(
        "Encoder Manager Ready");
}

void EncoderManager::update(AudioManager &audio)
{

    (void)audio;

    for (int i = 0; i < ENCODER_COUNT; i++)
    {
        readEncoder(i);
        readButton(i);
    }
}

void EncoderManager::readEncoder(int index)
{

    int clk =
        digitalRead(encoders[index].clk);

    if (clk == lastCLK[index])
        return;

    int direction =
        digitalRead(encoders[index].dt) != clk
            ? 1
            : -1;

    if (encoders[index].reverse)
        direction *= -1;

    Serial.printf(
        "ENC,%d,%d\n",
        encoders[index].id,
        direction);

    network.sendEncoder(
        encoders[index].id,
        direction);

    lastCLK[index] = clk;
}

void EncoderManager::readButton(int index)
{

    bool pressed =
        digitalRead(encoders[index].sw) == LOW;

    if (pressed == lastButton[index])
        return;

    if (
        millis() - lastDebounceTime <
        debounceDelay)
        return;

    lastDebounceTime = millis();

    if (pressed)
    {

        Serial.printf(
            "BTN,%d,0\n",
            encoders[index].id);

        network.sendEncoder(
            encoders[index].id,
            0);
    }

    lastButton[index] = pressed;
}