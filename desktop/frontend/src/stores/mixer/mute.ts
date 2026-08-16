export function createMute() {

    return {


        toggleMute(
            id: number
        ) {


            const channel =
                this.channels.find(
                    (item: any) => item.id === id
                )


            if (!channel)
                return



            channel.muted =
                !channel.muted



            console.log(
                "[MUTE]",
                channel.name,
                channel.muted
            )



            this.broadcastChannel(
                id
            )



        }


    }


}