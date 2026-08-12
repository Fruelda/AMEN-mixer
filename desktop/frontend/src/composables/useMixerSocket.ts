import {
    useRealtime
} from "./useRealtime"

import {
    mixerStore
} from "../stores/mixer"


export function useMixerSocket() {

    const {
        connect,
        subscribe,
        send
    } = useRealtime()


    // =====================================================
    // START
    // =====================================================

    function start() {

        connect()


        subscribe(
            (message) => {

                console.log(
                    "[REALTIME]",
                    message
                )


                // =========================================
                // CHANNEL UPDATE
                // =========================================

                if (
                    message.type ===
                    "CHANNEL_UPDATE"
                ) {

                    mixerStore.applyRemoteUpdate(
                        message
                    )

                    return

                }


                // =========================================
                // ESP32 COMMAND
                // =========================================

                if (
                    message.type ===
                    "COMMAND" &&
                    message.command
                ) {

                    mixerStore.handleCommand(
                        message.command
                    )

                    return

                }


                // =========================================
                // CONNECTED DEVICES
                // =========================================

                if (
                    message.type ===
                    "DEVICES" &&
                    Array.isArray(
                        message.devices
                    )
                ) {

                    mixerStore.setDevices(
                        message.devices
                    )


                    console.log(
                        "[DEVICES]",
                        message.devices
                    )


                    return

                }


                // =========================================
                // DEVICE STATUS
                // =========================================

                if (
                    message.type ===
                    "DEVICE_STATUS" &&
                    typeof message.connected ===
                    "boolean"
                ) {

                    mixerStore.setConnected(
                        message.connected
                    )

                    return

                }

            }
        )

    }


    // =====================================================
    // UPDATE CHANNEL
    // =====================================================

    function updateChannel(
        id: number,
        volume: number
    ) {

        send({

            type:
                "CHANNEL_UPDATE",

            channel: {

                id,

                volume

            }

        })

    }


    // =====================================================
    // RETURN
    // =====================================================

    return {

        start,

        updateChannel

    }

}