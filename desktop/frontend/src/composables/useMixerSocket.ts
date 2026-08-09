import {
    onMounted,
    onUnmounted,
} from "vue"

import {
    mixerStore,
} from "../stores/mixer"

import type {
    MixerCommand,
    MixerSocketMessage,
} from "../types/mixer"


let socket: WebSocket | null = null

let reconnectTimer:
    ReturnType<typeof setTimeout> | null = null

let manuallyClosed = false


function getWebSocketUrl() {

    /*
    |--------------------------------------------------------------------------
    | VITE_WS_URL
    |--------------------------------------------------------------------------
    |
    | Kalau ada:
    |
    | VITE_WS_URL=ws://192.168.1.10:8765/ws
    |
    | pakai itu.
    |
    */

    const envUrl =
        import.meta.env.VITE_WS_URL

    if (envUrl) {
        return envUrl
    }


    /*
    |--------------------------------------------------------------------------
    | Browser
    |--------------------------------------------------------------------------
    */

    const host =
        window.location.hostname

    return `ws://${host}:8765/ws`

}


function scheduleReconnect() {

    if (manuallyClosed) {
        return
    }

    if (reconnectTimer) {
        return
    }

    reconnectTimer =
        setTimeout(() => {

            reconnectTimer = null

            connect()

        }, 2000)

}


function handleMessage(
    event: MessageEvent
) {

    try {

        const message =
            JSON.parse(
                event.data
            ) as MixerSocketMessage


        console.log(
            "WS ←",
            message
        )


        /*
        |--------------------------------------------------------------------------
        | FULL STATE
        |--------------------------------------------------------------------------
        */

        if (
            message.type === "STATE"
        ) {

            mixerStore.setChannels(
                message.channels
            )

            mixerStore.setLoading(
                false
            )

            return
        }


        /*
        |--------------------------------------------------------------------------
        | ESP32 COMMAND
        |--------------------------------------------------------------------------
        */

        if (
            message.type === "COMMAND"
        ) {

            mixerStore.handleCommand(
                message.command
            )

            return
        }

    } catch (error) {

        console.error(
            "Invalid WebSocket message:",
            error,
            event.data
        )

    }

}


function connect() {

    manuallyClosed = false

    const url =
        getWebSocketUrl()


    console.log(
        "WebSocket connecting:",
        url
    )


    try {

        socket =
            new WebSocket(url)


        socket.onopen = () => {

            console.log(
                "WebSocket connected"
            )

            mixerStore.setConnected(
                true
            )


            socket?.send(
                JSON.stringify({
                    type: "HELLO",
                    client:
                        "browser",
                })
            )

        }


        socket.onmessage =
            handleMessage


        socket.onerror = (
            error
        ) => {

            console.error(
                "WebSocket error:",
                error
            )

        }


        socket.onclose = () => {

            console.warn(
                "WebSocket disconnected"
            )

            // mixerStore.setConnected(
            //     false
            // )

            // mixerStore.setLoading(
            //     false
            // )

            socket = null

            scheduleReconnect()

        }

    } catch (error) {

        console.error(
            "WebSocket connection failed:",
            error
        )

        mixerStore.setConnected(
            false
        )

        scheduleReconnect()

    }

}


function disconnect() {

    manuallyClosed = true

    if (reconnectTimer) {

        clearTimeout(
            reconnectTimer
        )

        reconnectTimer = null

    }


    if (socket) {

        socket.close()

        socket = null

    }

}


/*
|--------------------------------------------------------------------------
| SEND
|--------------------------------------------------------------------------
*/

function send(
    message: object
) {

    if (
        !socket ||
        socket.readyState !==
        WebSocket.OPEN
    ) {

        console.warn(
            "WebSocket belum connected"
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


/*
|--------------------------------------------------------------------------
| TOUCH VOLUME
|--------------------------------------------------------------------------
*/

export function setVolume(
    channel: number,
    value: number
) {

    return send({

        type: "SET_VOLUME",

        channel,

        value: Math.max(
            0,
            Math.min(
                100,
                Math.round(value)
            )
        ),

    })

}


/*
|--------------------------------------------------------------------------
| MUTE
|--------------------------------------------------------------------------
*/

export function toggleMute(
    channel: number
) {

    return send({

        type: "TOGGLE_MUTE",

        channel,

    })

}


/*
|--------------------------------------------------------------------------
| GENERIC COMMAND
|--------------------------------------------------------------------------
*/

export function sendCommand(
    command: MixerCommand
) {

    return send({

        type: "COMMAND",

        command,

    })

}


/*
|--------------------------------------------------------------------------
| COMPOSABLE
|--------------------------------------------------------------------------
*/

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

        setVolume,

        toggleMute,

        sendCommand,

    }

}