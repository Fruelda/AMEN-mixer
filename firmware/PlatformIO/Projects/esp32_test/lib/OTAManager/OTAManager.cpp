#include "OTAManager.h"

#include <ArduinoOTA.h>
#include <WiFi.h>

#include "../Config/Config.h"

void OTAManager::begin()
{

    if (
        WiFi.status() != WL_CONNECTED)
    {

        Serial.println(
            "[OTA] WiFi not ready");

        return;
    }

    ArduinoOTA.setHostname(
        DEVICE_ID);

    ArduinoOTA.setPassword(
        OTA_PASSWORD);

    ArduinoOTA.onStart(
        []()
        {
            Serial.println(
                "[OTA] Start");
        });

    ArduinoOTA.onEnd(
        []()
        {
            Serial.println(
                "[OTA] Complete");
        });

    ArduinoOTA.onProgress(
        [](unsigned int progress,
           unsigned int total)
        {
            Serial.printf(
                "[OTA] %u%%\r",
                (progress * 100) / total);
        });

    ArduinoOTA.onError(
        [](ota_error_t error)
        {
            Serial.printf(
                "[OTA] Error %u\n",
                error);
        });

    ArduinoOTA.begin();

    Serial.println(
        "OTA READY");
}

void OTAManager::loop()
{

    ArduinoOTA.handle();
}