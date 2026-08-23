#include <Arduino.h>

#include "NetworkManager.h"
#include "EncoderManager.h"
#include "AudioManager.h"
#include "OTAManager.h"

NetworkManager network;

EncoderManager encoder;

AudioManager audio;

OTAManager ota;

void setup()
{

    Serial.begin(
        115200);

    delay(1000);

    Serial.println();

    Serial.println(
        "=== AMEN START ===");

    network.begin();

    ota.begin();

    encoder.begin();

    audio.begin();

    Serial.println(
        "BOOT COMPLETE");
}

void loop()
{

    network.loop();

    ota.loop();

    encoder.update(
        audio);
}