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
        (() => void) | null =
        null


    // ============================================================
    // MESSAGE
    // ============================================================

    function handleMessage(
        message: RealtimeMessage
    ) {

        console.log(
            "[REALTIME]",
            message
        )


        switch (
        message.type
        ) {

            case "STATE":

                mixerStore.setChannels(
                    message.channels
                )

                return


            case "CHANNEL_UPDATE":

                mixerStore.applyRemoteUpdate(
                    message
                )

                return


            case "COMMAND":

                mixerStore.handleCommand(
                    message.command
                )

                return


            case "DEVICES":

                mixerStore.setDevices(
                    message.devices
                )

                return


            case "DEVICE_STATUS":

                mixerStore.setConnected(
                    message.connected
                )

                return

        }

    }


    // ============================================================
    // START
    // ============================================================

    function start() {

        connect()


        if (
            unsubscribe
        ) {
            return
        }


        unsubscribe =
            subscribe(
                handleMessage
            )

    }


    // ============================================================
    // STOP
    // ============================================================

    function stop() {

        unsubscribe?.()

        unsubscribe =
            null

    }


    // ============================================================
    // CLEANUP
    // ============================================================

    onUnmounted(
        stop
    )


    return {
        start,
        stop
    }

}