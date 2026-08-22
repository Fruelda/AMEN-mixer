#include "MessageHandler.h"

#include <ArduinoJson.h>

void MessageHandler::handle(
    uint8_t *payload)
{

    Serial.print(
        "SERVER: ");

    Serial.println(
        (char *)payload);

    StaticJsonDocument<512> doc;

    DeserializationError error =
        deserializeJson(
            doc,
            payload);

    if (error)
    {

        Serial.println(
            "JSON ERROR");

        return;
    }

    const char *type =
        doc["type"];

    if (!type)
    {

        return;
    }

    Serial.print(
        "MESSAGE TYPE: ");

    Serial.println(
        type);
}