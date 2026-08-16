import type {
    ConnectedDevice
} from "../../types/mixer"



export function createDevices() {

    return {


        setDevices(
            devices: ConnectedDevice[]
        ) {

            this.devices =
                devices

        },



        setConnected(
            value: boolean
        ) {

            this.deviceConnected =
                value

        }


    }


}