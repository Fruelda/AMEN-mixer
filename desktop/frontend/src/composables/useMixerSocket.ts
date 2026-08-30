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
                //
                // IMPORTANT:
                //
                // On Wails desktop, the ESP32 serial command already reaches
                // App.vue through the Wails "serial-command" event.
                //
                // app.go also broadcasts the same command over WebSocket for
                // realtime visibility. Executing it again here would process
                // one physical action twice:
                //
                // BTN: false -> true -> false
                // ENC: volume can move twice per detent
                //
                // Therefore COMMAND messages received over WebSocket are not
                // executed here. The resulting CHANNEL_UPDATE is the message
                // used to synchronize all clients.
                // =========================================

                if (
                    message.type ===
                    "COMMAND" &&
                    message.command
                ) {

                    console.log(
                        "[REALTIME] ESP32 COMMAND mirrored over WebSocket. " +
                        "Desktop execution skipped because serial-command is authoritative.",
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
