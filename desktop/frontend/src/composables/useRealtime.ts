import {
    onMounted,
    onUnmounted,
} from "vue"

import type {
    RealtimeMessage,
} from "../types/mixer"


let socket: WebSocket | null = null

let listeners:
    ((message: RealtimeMessage) => void)[] = []

let reconnectTimer:
    ReturnType<typeof setTimeout> | null = null


function getWebSocketURL() {

    const hostname =
        window.location.hostname

    //wails desktop app
    if (hostname === "wails.localhost" ||
        hostname === "localhost"
    ) {
        const url = "ws://127.0.0.1:8081/ws"
        console.log("[WS] Wails URL:", url)

        return url
    }
    // Browser / HP / LAN
    const url = `ws://${hostname}:8081/ws`

    console.log("[WS] Browser URL:", url)

    return url
    // return `ws://${hostname}:8081/ws`
}


function connect() {

    if (
        socket &&
        (
            socket.readyState ===
            WebSocket.OPEN ||

            socket.readyState ===
            WebSocket.CONNECTING
        )
    ) {
        return
    }


    const url =
        getWebSocketURL()


    console.log(
        "[WS] Connecting:",
        url
    )


    socket =
        new WebSocket(url)


    socket.onopen = () => {

        console.log(
            "[WS] CONNECTED:",
            url
        )

        window.dispatchEvent(
            new CustomEvent(
                "realtime-status",
                {
                    detail:
                        "CONNECTED",
                }
            )
        )
    }


    socket.onmessage = (
        event
    ) => {

        console.log(
            "[WS] MESSAGE:",
            event.data
        )


        try {

            const message =
                JSON.parse(
                    event.data
                ) as RealtimeMessage


            listeners.forEach(
                listener => {

                    listener(message)

                }
            )

        } catch (error) {

            console.error(
                "[WS] JSON ERROR:",
                error
            )
        }
    }


    socket.onerror = (
        error
    ) => {

        console.error(
            "[WS] ERROR:",
            error
        )
    }


    socket.onclose = () => {

        console.log(
            "[WS] CLOSED"
        )


        socket = null


        window.dispatchEvent(
            new CustomEvent(
                "realtime-status",
                {
                    detail:
                        "DISCONNECTED",
                }
            )
        )


        if (
            reconnectTimer
        ) {
            return
        }


        reconnectTimer =
            setTimeout(() => {

                reconnectTimer =
                    null

                if (
                    listeners.length >
                    0
                ) {

                    connect()

                }

            }, 2000)
    }
}


export function useRealtime(
    callback: (
        message: RealtimeMessage
    ) => void
) {

    onMounted(() => {

        listeners.push(
            callback
        )

        connect()
    })


    onUnmounted(() => {

        listeners =
            listeners.filter(
                listener =>
                    listener !==
                    callback
            )
    })
}


export function sendRealtime(
    message: RealtimeMessage
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


    console.log(
        "[WS] SEND:",
        message
    )


    socket.send(
        JSON.stringify(
            message
        )
    )


    return true
}