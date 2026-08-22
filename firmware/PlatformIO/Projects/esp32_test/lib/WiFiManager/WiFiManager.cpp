#include "WiFiManager.h"

#include <WiFi.h>

#include "../Config/Config.h"

unsigned long lastAttempt = 0;

const unsigned long retryInterval = 5000;

void WiFiManager::begin()
{

    Serial.println(
        "Connecting WiFi...");

    WiFi.mode(
        WIFI_STA);

    WiFi.begin(
        WIFI_SSID,
        WIFI_PASSWORD);
}

void WiFiManager::loop()
{

    if (
        WiFi.status() ==
        WL_CONNECTED)
    {

        return;
    }

    if (
        millis() -
            lastAttempt <
        retryInterval)
    {

        return;
    }

    lastAttempt =
        millis();

    Serial.println(
        "WiFi reconnect...");

    WiFi.disconnect();

    WiFi.begin(
        WIFI_SSID,
        WIFI_PASSWORD);
}

bool WiFiManager::connected()
{

    return WiFi.status() ==
           WL_CONNECTED;
}