#include <Arduino.h>

#include "NetworkManager.h"
#include "EncoderManager.h"
#include "AudioManager.h"

NetworkManager network;

EncoderManager encoder;

AudioManager audio;

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

  Serial.println(
      "BOOT");
}

void loop()
{

  network.loop();

  encoder.update(
      audio);
}