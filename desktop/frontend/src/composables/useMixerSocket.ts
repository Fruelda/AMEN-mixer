import {
    onMounted,
    onUnmounted,
} from "vue"


import {
    mixerStore,
} from "../stores/mixer"


import type {
    MixerCommand,
    RealtimeMessage,
} from "../types/mixer"



let socket:
    WebSocket | null = null



let reconnectTimer:
    ReturnType<typeof setTimeout> | null = null



let manuallyClosed =
    false



// ============================================================
// WEBSOCKET URL
// ============================================================

function getWebSocketUrl() {

    const envUrl =
        import.meta.env.VITE_WS_URL



    if (envUrl) {

        return envUrl

    }



    const host =
        window.location.hostname



    return (
        `ws://${host}:8081/ws`
    )

}



// ============================================================
// RECONNECT
// ============================================================

function scheduleReconnect() {


    if (manuallyClosed) {

        return

    }



    if (reconnectTimer) {

        return

    }



    reconnectTimer =
        setTimeout(() => {


            reconnectTimer =
                null



            connect()


        }, 2000)

}



// ============================================================
// MESSAGE HANDLER
// ============================================================

function handleMessage(
    event: MessageEvent
) {


    try {


        const message =
            JSON.parse(
                event.data
            )




        console.log(
            "[WS RECEIVE]",
            message
        )



        // ======================================================
        // FULL STATE
        // ======================================================


        if (
            message.type === "STATE"
        ) {


            if (
                message.channels
            ) {

                mixerStore.setChannels(
                    message.channels
                )

            }



            mixerStore.setLoading(
                false
            )



            return

        }




        // ======================================================
        // CHANNEL UPDATE
        // ======================================================


        if (
            message.type ===
            "CHANNEL_UPDATE"
        ) {


            mixerStore.applyRemoteUpdate(
                message
            )



            return

        }




        // ======================================================
        // ESP32 COMMAND
        // ======================================================


        if (
            message.type ===
            "COMMAND"
        ) {


            if (
                message.command
            ) {

                mixerStore.handleCommand(
                    message.command
                )

            }



            return

        }




        // ======================================================
        // DEVICE STATUS
        // ======================================================


        if (
            message.type ===
            "DEVICE_STATUS"
        ) {


            if (
                typeof message.connected ===
                "boolean"
            ) {


                mixerStore.setConnected(
                    message.connected
                )

            }


            return

        }



    }
    catch (error) {


        console.error(
            "[WS] Invalid message",
            error,
            event.data
        )


    }


}



// ============================================================
// CONNECT
// ============================================================

function connect() {


    manuallyClosed =
        false



    const url =
        getWebSocketUrl()



    console.log(
        "[WS] Connecting:",
        url
    )



    try {


        socket =
            new WebSocket(
                url
            )




        socket.onopen =
            () => {


                console.log(
                    "[WS] Connected"
                )



                mixerStore.setConnected(
                    true
                )



                socket?.send(
                    JSON.stringify({

                        type:
                            "HELLO",


                        client:
                            "browser",

                    })
                )


            }




        socket.onmessage =
            handleMessage




        socket.onerror =
            (error) => {


                console.error(
                    "[WS] Error",
                    error
                )


            }




        socket.onclose =
            () => {


                console.warn(
                    "[WS] Closed"
                )



                mixerStore.setConnected(
                    false
                )



                mixerStore.setLoading(
                    false
                )



                socket =
                    null



                scheduleReconnect()


            }




    }
    catch (error) {


        console.error(
            "[WS] Failed",
            error
        )



        mixerStore.setConnected(
            false
        )



        scheduleReconnect()


    }

}



// ============================================================
// DISCONNECT
// ============================================================

function disconnect() {


    manuallyClosed =
        true



    if (reconnectTimer) {


        clearTimeout(
            reconnectTimer
        )



        reconnectTimer =
            null

    }




    if (socket) {


        socket.close()



        socket =
            null

    }


}



// ============================================================
// SEND
// ============================================================

function send(
    message: object
) {


    if (
        !socket ||
        socket.readyState !==
        WebSocket.OPEN
    ) {


        console.warn(
            "[WS] Not connected"
        )



        return false

    }



    socket.send(
        JSON.stringify(
            message
        )
    )



    return true

}



// ============================================================
// SEND COMMAND
// ============================================================

export function sendCommand(
    command: MixerCommand
) {


    return send({

        type:
            "COMMAND",


        command,

    })

}



// ============================================================
// COMPOSABLE
// ============================================================

export function useMixerSocket() {


    onMounted(() => {


        connect()


    })



    onUnmounted(() => {


        disconnect()


    })



    return {

        connect,

        disconnect,

        sendCommand,

    }


}