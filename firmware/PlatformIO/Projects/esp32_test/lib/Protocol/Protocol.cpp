#include "Protocol.h"

#include <ArduinoJson.h>

#include "../Config/Config.h"

String createDeviceRegister()
{

    JsonDocument doc;

    doc["type"] =
        "device.register";

    doc["id"] =
        DEVICE_ID;

    doc["name"] =
        DEVICE_NAME;

    String output;

    serializeJson(
        doc,
        output);

    return output;
}

String createEncoderCommand(

    uint8_t channel,

    int value

)
{

    JsonDocument doc;

    doc["type"] =
        "mixer.command";

    doc["device"] =
        DEVICE_ID;

    doc["channel"] =
        channel;

    doc["value"] =
        value;

    String output;

    serializeJson(
        doc,
        output);

    return output;
}