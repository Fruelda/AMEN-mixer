import { ref } from "vue"

const connected = ref(false)

let socket: WebSocket | null = null

const listeners = new Set<
    (data: any) => void
>()


// =====================================================
// CONFIG
// =====================================================

const WS_URL =
    "ws://192.168.1.44:8081/ws"


// =====================================================
// ENVIRONMENT
// =====================================================

function isWailsEnvironment() {

    if (
        typeof window === "undefined"
    ) {
        return false
    }

    return (
        "__WAILS_RUNTIME__" in window
    )
}


// =====================================================
// CLIENT INFO, baca device yang konek
// =====================================================

function getClientInfo() {

    // =================================================
    // WAILS DESKTOP
    // =================================================

    if (
        isWailsEnvironment()
    ) {

        return {

            prefix:
                "desktop",

            name:
                "Wails Desktop",

            clientType:
                "desktop"

        }

    }


    // =================================================
    // BROWSER DEVICE
    // =================================================

    const userAgent =
        navigator.userAgent.toLowerCase()


    // iPhone

    if (
        /iphone|ipod/.test(
            userAgent
        )
    ) {

        return {

            prefix:
                "iphone",

            name:
                "iPhone",

            clientType:
                "mobile"

        }

    }


    // iPad

    if (
        /ipad/.test(
            userAgent
        ) ||
        (
            navigator.platform === "MacIntel" &&
            navigator.maxTouchPoints > 1
        )
    ) {

        return {

            prefix:
                "ipad",

            name:
                "iPad",

            clientType:
                "tablet"

        }

    }


    // Android

    if (
        /android/.test(
            userAgent
        )
    ) {

        const isTablet =
            !/mobile/.test(
                userAgent
            )

        if (
            isTablet
        ) {

            return {

                prefix:
                    "android-tablet",

                name:
                    "Android Tablet",

                clientType:
                    "tablet"

            }

        }

        return {

            prefix:
                "android",

            name:
                "Android Phone",

            clientType:
                "mobile"

        }

    }


    // =================================================
    // GENERIC BROWSER
    // =================================================

    return {

        prefix:
            "browser",

        name:
            "AMEN Browser",

        clientType:
            "browser"

    }
}


// =====================================================
// PERSISTENT CLIENT ID
// =====================================================

function getClientID(
    prefix: string
) {

    const storageKey =
        "amen-mixer-client-id"


    try {

        const savedID =
            localStorage.getItem(
                storageKey
            )


        if (
            savedID
        ) {

            return savedID

        }


        const id =
            `${prefix}-${crypto.randomUUID()}`


        localStorage.setItem(
            storageKey,
            id
        )


        return id

    }
    catch {


        // fallback jika localStorage
        // tidak tersedia

        return (
            `${prefix}-${Date.now()}`
        )

    }
}


// =====================================================
// REGISTER CLIENT
// =====================================================

function registerClient() {

    const info =
        getClientInfo()


    const id =
        getClientID(
            info.prefix
        )


    sendRealtime({

        type:
            "client.register",

        id:
            id,

        name:
            info.name,

        clientType:
            info.clientType

    })


    console.log(
        "[WS] Registered:",
        info.name,
        id
    )
}


// =====================================================
// CONNECT
// =====================================================

export function connectRealtime() {

    if (
        socket &&
        (
            socket.readyState === WebSocket.OPEN ||
            socket.readyState === WebSocket.CONNECTING
        )
    ) {

        console.log(
            "[WS] Already connecting/connected"
        )

        return

    }


    socket =
        new WebSocket(
            WS_URL
        )


    // =================================================
    // OPEN
    // =================================================

    socket.onopen = () => {

        connected.value =
            true


        console.log(
            "[WS] Connected"
        )


        registerClient()

    }


    // =================================================
    // RECEIVE
    // =================================================

    socket.onmessage = event => {

        try {

            const data =
                JSON.parse(
                    event.data
                )


            console.log(
                "[WS RECEIVE]",
                data
            )


            listeners.forEach(
                callback => {

                    callback(
                        data
                    )

                }
            )

        }
        catch (error) {

            console.error(
                "[WS ERROR]",
                error
            )

        }

    }


    // =================================================
    // ERROR
    // =================================================

    socket.onerror = error => {

        console.error(
            "[WS ERROR]",
            error
        )

    }


    // =================================================
    // CLOSE
    // =================================================

    socket.onclose = () => {

        connected.value =
            false


        console.log(
            "[WS] Disconnected"
        )


        socket =
            null


        setTimeout(
            () => {

                connectRealtime()

            },
            3000
        )

    }

}


// =====================================================
// SEND
// =====================================================

export function sendRealtime(
    data: any
) {

    if (
        !socket ||
        socket.readyState !== WebSocket.OPEN
    ) {

        console.warn(
            "[WS] not connected"
        )

        return

    }


    socket.send(
        JSON.stringify(
            data
        )
    )

}


// =====================================================
// SUBSCRIBE
// =====================================================

export function onRealtimeMessage(
    callback: (data: any) => void
) {

    listeners.add(
        callback
    )


    return () => {

        listeners.delete(
            callback
        )

    }

}


// =====================================================
// COMPOSABLE
// =====================================================

export function useRealtime() {

    return {

        connected,

        connect:
            connectRealtime,

        send:
            sendRealtime,

        subscribe:
            onRealtimeMessage

    }

}