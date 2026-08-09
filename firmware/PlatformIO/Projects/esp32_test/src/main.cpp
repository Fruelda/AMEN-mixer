#include <Arduino.h>

#include "EncoderManager.h"
#include "AudioManager.h"
#include "Utils.h"
#include "Config.h"

EncoderManager encoder;
AudioManager audio;

void setup()
{
  Serial.begin(BAUD_RATE);

  encoder.begin();
  audio.begin();

  Utils::printHeader();
}

void loop()
{
  encoder.update(audio);
}