import {
    isMixerWailsEnvironment,
    normalizeVolume
} from "./helpers"



export function createCommands() {


    return {



        handleCommand(
            command: any
        ) {



            console.log(
                "[HANDLE COMMAND]",
                command
            )



            if (
                !command
            ) {
                return
            }




            if (
                !isMixerWailsEnvironment()
            ) {

                return

            }



            console.log(
                "[AVAILABLE CHANNELS]",
                this.channels.map(
                    (x: any) => x.id
                )
            )


            console.log(
                "[COMMAND CHANNEL]",
                command.channel
            )

            const channel =
                this.channels.find(
                    (item: any) =>
                        item.id === command.channel
                )





            if (
                !channel
            ) {

                console.warn(
                    "[CHANNEL NOT FOUND]",
                    command.channel
                )

                return

            }







            if (

                command.type === "ENC"

            ) {



                const newVolume =

                    normalizeVolume(

                        channel.volume +

                        Number(command.value)

                    )





                console.log(

                    "[SET VOLUME]",

                    "CH",
                    channel.id,

                    channel.volume,

                    "=>",

                    newVolume

                )





                void this.setVolume(

                    channel.id,

                    newVolume,

                    false

                )



                return

            }







            if (

                command.type === "BTN"

                &&

                command.value === 1

            ) {


                this.toggleMute(

                    channel.id

                )


                return

            }



        },






        setVolumeFromESP32(
            id: number,
            volume: number
        ) {


            void this.setVolume(

                id,

                volume,

                false

            )


        }


    }


}