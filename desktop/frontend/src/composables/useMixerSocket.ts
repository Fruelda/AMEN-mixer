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



    function start() {


        connect()



        subscribe(
            (message) => {


                console.log(
                    "[REALTIME]",
                    message
                )



                if (
                    message.type === "CHANNEL_UPDATE"
                ) {

                    mixerStore.applyRemoteUpdate(
                        message
                    )

                }



                if (
                    message.type === "COMMAND"
                ) {

                    mixerStore.handleCommand(
                        message.command
                    )

                }



            }
        )


    }




    function updateChannel(
        id: number,
        volume: number
    ) {


        send({

            type: "CHANNEL_UPDATE",

            channel: {

                id,

                volume

            }

        })


    }



    return {

        start,

        updateChannel

    }


}