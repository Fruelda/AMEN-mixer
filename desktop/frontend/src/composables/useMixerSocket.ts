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
        (() => void) | null = null








    // ============================================================
    // HANDLE REALTIME MESSAGE
    // ============================================================

    function handleMessage(
        message: RealtimeMessage | any
    ) {


        console.log(
            "[REALTIME RECEIVED]",
            message
        )



        if (
            !message ||
            !message.type
        ) {

            return

        }





        switch (
        message.type
        ) {





            // ====================================================
            // FULL STATE
            // ====================================================

            case "STATE":


                mixerStore.setChannels(
                    message.channels
                )


                return






            // ====================================================
            // UPDATE CHANNEL FROM BACKEND
            // ====================================================

            case "CHANNEL_UPDATE":



                mixerStore.applyRemoteUpdate({

                    id:
                        message.channel.id,


                    volume:
                        message.channel.volume,


                    muted:
                        message.channel.muted

                })


                return







            // ====================================================
            // OLD COMMAND FORMAT
            //
            // {
            //  type:"COMMAND",
            //  command:{
            //      type:"ENC",
            //      channel:1,
            //      value:1
            //  }
            // }
            //
            // ====================================================

            case "COMMAND":



                mixerStore.handleCommand(

                    message.command

                )


                return







            // ====================================================
            // ESP32 FORMAT
            //
            // Encoder:
            // value = 1 / -1
            //
            // Button:
            // value = 0
            //
            // ====================================================

            case "mixer.command":



                console.log(
                    "[MIXER COMMAND]",
                    message
                )




                // ==========================
                // BUTTON MUTE
                // ==========================

                if (
                    Number(message.value) === 0
                ) {



                    console.log(
                        "[MUTE PRESS]",
                        "CHANNEL",
                        message.channel
                    )



                    mixerStore.handleCommand({

                        type: "BTN",


                        channel:
                            message.channel,


                        value: 1

                    })



                    return

                }






                // ==========================
                // ENCODER VOLUME
                // ==========================

                mixerStore.handleCommand({

                    type: "ENC",


                    channel:
                        message.channel,


                    value:
                        message.value

                })



                return







            // ====================================================
            // DEVICES
            // ====================================================

            case "DEVICES":



                mixerStore.setDevices(

                    message.devices

                )


                return








            // ====================================================
            // DEVICE STATUS
            // ====================================================

            case "DEVICE_STATUS":



                mixerStore.setConnected(

                    message.connected

                )


                return







            default:



                console.warn(

                    "[UNKNOWN MESSAGE]",

                    message

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



        console.log(
            "[MIXER SOCKET] STARTED"
        )


    }









    // ============================================================
    // STOP
    // ============================================================

    function stop() {



        if (
            unsubscribe
        ) {

            unsubscribe()

        }



        unsubscribe = null



        console.log(
            "[MIXER SOCKET] STOPPED"
        )


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