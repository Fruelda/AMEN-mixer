import type {
    MixerChannel,
    ConnectedDevice,
} from "../../types/mixer"


export function createState() {

    return {

        channels:
            [] as MixerChannel[],


        devices:
            [] as ConnectedDevice[],


        loading: false,


        deviceConnected: false,


    }

}