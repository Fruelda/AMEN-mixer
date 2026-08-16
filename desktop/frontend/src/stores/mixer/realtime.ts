import {
    sendRealtime
} from "../../composables/useRealtime"


import {
    normalizeVolume
} from "./helpers"



export function createRealtime() {

    return {



        broadcastChannel(
            id: number
        ) {


            const channel =
                this.channels.find(
                    (item: any) => item.id === id
                )



            if (!channel)
                return




            sendRealtime({

                type: "CHANNEL_UPDATE",


                channel: {


                    id: channel.id,


                    volume: channel.volume,


                    muted: channel.muted


                }


            })


        },





        applyRemoteUpdate(
            message: any
        ) {


            if (
                message.type !== "CHANNEL_UPDATE"
                ||
                !message.channel
            )

                return



            const channel =
                this.channels.find(
                    (item: any) =>
                        item.id === message.channel.id
                )



            if (!channel)
                return




            if (
                typeof message.channel.volume === "number"
            ) {

                channel.volume =
                    normalizeVolume(
                        message.channel.volume
                    )

            }



            if (
                typeof message.channel.muted === "boolean"
            ) {

                channel.muted =
                    message.channel.muted

            }


        }



    }


}