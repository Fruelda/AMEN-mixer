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



            // =========================================
            // INITIAL MIXER STATE
            // =========================================

            case "STATE":


                if (
                    message.channels
                ) {

                    mixerStore.setChannels(
                        message.channels
                    )

                }


                return






            // =========================================
            // DEVICE LIST
            // =========================================

            case "DEVICES":


                if (
                    message.devices &&
                    mixerStore.setDevices
                ) {

                    mixerStore.setDevices(
                        message.devices
                    )

                }


                return






            // =========================================
            // CHANNEL UPDATE
            // Android / Backend
            // =========================================

            case "CHANNEL_UPDATE":



                if (
                    !message.channel
                ) {

                    return

                }



                mixerStore.applyRemoteUpdate({

                    channel: {

                        id:
                            message.channel.id,


                        volume:
                            message.channel.volume,


                        muted:
                            message.channel.muted

                    }

                })



                return






            // =========================================
            // ESP32 COMMAND
            // =========================================

            case "mixer.command":


                console.log(
                    "[MIXER COMMAND]",
                    message
                )



                // BUTTON MUTE

                if (
                    Number(message.value) === 0
                ) {


                    mixerStore.handleCommand({

                        type:
                            "BTN",


                        channel:
                            message.channel,


                        value:
                            1

                    })


                    return

                }




                // ENCODER

                mixerStore.handleCommand({

                    type:
                        "ENC",


                    channel:
                        message.channel,


                    value:
                        message.value

                })


                return






            // =========================================
            // OLD COMMAND
            // =========================================

            case "COMMAND":


                mixerStore.handleCommand(
                    message.command
                )


                return






            default:


                console.log(
                    "[UNKNOWN MESSAGE]",
                    message
                )


                return

        }

    }







    function start() {


        if (
            unsubscribe
        ) {

            return

        }



        connect()



        unsubscribe =
            subscribe(
                handleMessage
            )



        console.log(
            "[MIXER SOCKET STARTED]"
        )

    }







    onUnmounted(() => {


        if (
            unsubscribe
        ) {

            unsubscribe()

            unsubscribe = null

        }


    })






    return {

        start

    }

}