#include <Arduino.h>

#include <ArduinoOTA.h>

#include "NetworkManager.h"

#include "EncoderManager.h"

#include "AudioManager.h"

NetworkManager network;

EncoderManager encoder;

AudioManager audio;

void setupOTA()
{

    ArduinoOTA.setHostname(
        "amen-mixer");

    ArduinoOTA.setPassword(
        "amen123");

    ArduinoOTA.begin();

    Serial.println(
        "OTA READY");
}

void setup()
{

    Serial.begin(
        115200);

    delay(1000);

    Serial.println(
        "=== AMEN START ===");

    network.begin();

    encoder.begin();

    audio.begin();

    setupOTA();

    Serial.println(
        "BOOT");
}

void loop()
{

    ArduinoOTA.handle();

    network.loop();

    encoder.update(
        audio);
}