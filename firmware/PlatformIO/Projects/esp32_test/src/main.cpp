#include <Arduino.h>

#include <ArduinoOTA.h>

#include "NetworkManager.h"
#include "EncoderManager.h"
#include "AudioManager.h"

NetworkManager network;

EncoderManager encoder;

AudioManager audio;

// ============================================================
// SETUP
// ============================================================

void setup()
{

    Serial.begin(
        115200);

    delay(
        1000);

    Serial.println();

    Serial.println(
        "=== AMEN START ===");

    // ==========================
    // NETWORK
    // ==========================

    network.begin();

    // ==========================
    // OTA
    // ==========================

    ArduinoOTA.setHostname(
        "amen-mixer-01");

    // ArduinoOTA.setPassword(
    //     "amen123");

    ArduinoOTA.onStart(
        []()
        {
            Serial.println(
                "OTA START");
        });

    ArduinoOTA.onEnd(
        []()
        {
            Serial.println(
                "\nOTA END");
        });

    ArduinoOTA.onProgress(
        [](unsigned int progress,
           unsigned int total)
        {
            Serial.printf(
                "OTA Progress: %u%%\r",
                (progress * 100) / total);
        });

    ArduinoOTA.onError(
        [](ota_error_t error)
        {
            Serial.printf(
                "OTA Error[%u]\n",
                error);
        });

    ArduinoOTA.begin();

    Serial.println(
        "OTA READY");

    // ==========================
    // HARDWARE
    // ==========================

    encoder.begin();

    audio.begin();

    Serial.println(
        "BOOT");
}

// ============================================================
// LOOP
// ============================================================

void loop()
{

    // WebSocket
    network.loop();

    // OTA handler
    ArduinoOTA.handle();

    // Encoder
    encoder.update(
        audio);
}