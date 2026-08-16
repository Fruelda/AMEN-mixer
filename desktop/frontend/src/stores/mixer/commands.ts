import {
    isMixerWailsEnvironment,
    normalizeVolume
} from "./helpers"



export function createCommands() {

    return {


        handleCommand(
            command: any
        ) {


            if (
                !isMixerWailsEnvironment()
            )

                return




            const channel =
                this.channels.find(
                    (item: any) =>
                        item.id === command.channel
                )



            if (!channel)
                return




            if (
                command.type === "ENC"
            ) {


                void this.setVolume(

                    channel.id,

                    normalizeVolume(
                        channel.volume +
                        command.value
                    ),

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