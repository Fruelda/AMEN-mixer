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

    ArduinoOTA.onStart([]()
                       {
                           Serial.println(
                               "OTA START");
                       });

    ArduinoOTA.onEnd([]()
                     {
                         Serial.println(
                             "OTA COMPLETE");
                     });

    ArduinoOTA.onProgress(
        [](unsigned int progress,
           unsigned int total)
        {
            Serial.printf(
                "OTA Progress: %u%%\n",
                (progress * 100) / total);
        });

    ArduinoOTA.onError(
        [](ota_error_t error)
        {
            Serial.printf(
                "OTA ERROR: %u\n",
                error);
        });

    ArduinoOTA.begin();

    Serial.println(
        "OTA READY");
}

void setup()
{

    Serial.begin(
        115200);

    delay(1000);

    Serial.println();

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