import {
    SetVolume,
} from "../../../wailsjs/go/main/App"


import {
    isMixerWailsEnvironment,
    normalizeVolume,
} from "./helpers"



export function createVolume() {

    return {


        async setVolume(
            id: number,
            volume: number,
            broadcast = true
        ) {


            const channel =
                this.channels.find(
                    (item: any) => item.id === id
                )


            if (!channel)
                return



            const value =
                normalizeVolume(volume)



            channel.volume = value




            if (
                isMixerWailsEnvironment()
            ) {


                try {


                    await SetVolume(
                        id,
                        value
                    )


                }
                catch (error) {


                    console.warn(
                        "[WAILS] SetVolume failed",
                        error
                    )


                }



                return


            }



            if (broadcast) {

                this.broadcastChannel(id)

            }



        }


    }


}