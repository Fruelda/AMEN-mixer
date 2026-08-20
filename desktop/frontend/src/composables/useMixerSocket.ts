import {
    onUnmounted
} from "vue"


import {
    useRealtime
} from "./useRealtime"


import {
    mixerStore
} from "../stores/mixer"


import type {
    RealtimeMessage
} from "../types/mixer"




export function useMixerSocket() {


    const {
        connect,
        subscribe
    } = useRealtime()



    let unsubscribe:
        (() => void) | null =
        null





    // ============================================================
    // MESSAGE HANDLER
    // ============================================================

    function handleMessage(
        message: RealtimeMessage
    ) {


        console.log(
            "[REALTIME]",
            message
        )



        switch (
        message.type
        ) {



            case "STATE":


                mixerStore.setChannels(
                    message.channels
                )


                return





            case "CHANNEL_UPDATE":


                mixerStore.applyRemoteUpdate(
                    message
                )


                return





            case "COMMAND":


                mixerStore.handleCommand(
                    message.command
                )


                return





            // ====================================================
            // ESP32 HARDWARE COMMAND
            // ====================================================

            case "mixer.command":



                // ==========================
                // BUTTON
                // value 0 = press
                // ==========================

                if (
                    message.value === 0
                ) {


                    mixerStore.handleCommand({

                        type: "BTN",

                        channel:
                            message.channel,

                        value:
                            1

                    })


                    return

                }




                // ==========================
                // ENCODER
                // value +1 / -1
                // ==========================

                mixerStore.handleCommand({

                    type: "ENC",

                    channel:
                        message.channel,

                    value:
                        message.value

                })


                return





            case "DEVICES":


                mixerStore.setDevices(
                    message.devices
                )


                return





            case "DEVICE_STATUS":


                mixerStore.setConnected(
                    message.connected
                )


                return


        }


    }






    // ============================================================
    // START
    // ============================================================

    function start() {


        connect()



        if (
            unsubscribe
        ) {

            return

        }



        unsubscribe =
            subscribe(
                handleMessage
            )


    }






    // ============================================================
    // STOP
    // ============================================================

    function stop() {


        unsubscribe?.()


        unsubscribe =
            null


    }






    // ============================================================
    // CLEANUP
    // ============================================================

    onUnmounted(
        stop
    )






    return {

        start,

        stop

    }


}