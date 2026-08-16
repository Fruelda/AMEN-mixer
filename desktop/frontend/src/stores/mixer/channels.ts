import {
    GetChannels,
} from "../../../wailsjs/go/main/App"


import {
    isMixerWailsEnvironment,
    normalizeVolume,
} from "./helpers"


import {
    mockChannels,
} from "./mockChannels"



export function createChannels() {

    return {


        setChannels(
            channels: any[]
        ) {

            this.channels =
                channels.map(
                    (channel) => ({

                        ...channel,

                        volume:
                            normalizeVolume(
                                channel.volume
                            )

                    })
                )

        },



        async loadChannels() {

            this.loading = true


            try {


                if (
                    !isMixerWailsEnvironment()
                ) {

                    this.setChannels(
                        structuredClone(
                            mockChannels
                        )
                    )

                    return

                }



                const result =
                    await GetChannels()



                this.setChannels(

                    Array.isArray(result)
                        &&
                        result.length > 0

                        ?
                        result

                        :
                        mockChannels

                )



            } catch (error) {


                console.warn(
                    "[MIXER] Load failed",
                    error
                )


                this.setChannels(
                    mockChannels
                )


            }
            finally {

                this.loading = false

            }


        }


    }


}