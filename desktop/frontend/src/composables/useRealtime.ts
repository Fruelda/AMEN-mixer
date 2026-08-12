import {
    ref
} from "vue"


// =====================================================
// STATE
// =====================================================

const connected =
    ref(
        false
    )


let socket:
    WebSocket |
    null =
    null


let reconnectTimer:
    ReturnType<
        typeof setTimeout
    > |
    null =
    null


const listeners =
    new Set<
        (
            data: any
        ) => void
    >()


// =====================================================
// ENVIRONMENT
// =====================================================

function isWailsEnvironment() {

    if (
        typeof window ===
        "undefined"
    ) {

        return false

    }


    // Wails runtime tersedia
    // sebagai window.runtime.
    return (
        "runtime" in window
    )

}


// =====================================================
// REALTIME URL
// =====================================================

function getRealtimeURL() {

    // =================================================
    // WAILS DESKTOP
    // =================================================
    //
    // Backend Go berada di komputer
    // yang sama.
    //
    // Tidak perlu menggunakan LAN IP.
    //
    // =================================================

    if (
        isWailsEnvironment()
    ) {

        return (
            "ws://127.0.0.1:8081/ws"
        )

    }


    // =================================================
    // IPHONE / ANDROID / BROWSER
    // =================================================
    //
    // Misalnya halaman dibuka:
    //
    // http://amen-mixer.local:5173
    //
    // hostname:
    //
    // amen-mixer.local
    //
    // WebSocket otomatis:
    //
    // ws://amen-mixer.local:8081/ws
    //
    //
    // Kalau fallback pakai IP:
    //
    // http://192.168.18.11:5173
    //
    // WebSocket otomatis:
    //
    // ws://192.168.18.11:8081/ws
    //
    // =================================================

    const host =
        window.location.hostname


    return (
        `ws://${host}:8081/ws`
    )

}


// =====================================================
// CLIENT INFO
// =====================================================

function getClientInfo() {

    // =================================================
    // WAILS
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
                "desktop",

        }

    }


    // =================================================
    // BROWSER
    // =================================================

    const userAgent =
        navigator.userAgent
            .toLowerCase()


    // =================================================
    // IPHONE
    // =================================================

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
                "mobile",

        }

    }


    // =================================================
    // IPAD
    // =================================================

    if (
        /ipad/.test(
            userAgent
        ) ||
        (
            navigator.platform ===
            "MacIntel" &&

            navigator.maxTouchPoints >
            1
        )
    ) {

        return {

            prefix:
                "ipad",

            name:
                "iPad",

            clientType:
                "tablet",

        }

    }


    // =================================================
    // ANDROID
    // =================================================

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
                    "tablet",

            }

        }


        return {

            prefix:
                "android",

            name:
                "Android Phone",

            clientType:
                "mobile",

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
            "browser",

    }

}


// =====================================================
// RANDOM ID
// =====================================================

function createRandomID(
    prefix: string
) {

    try {

        if (
            typeof crypto !==
            "undefined" &&

            typeof crypto.randomUUID ===
            "function"
        ) {

            return (
                `${prefix}-${crypto.randomUUID()}`
            )

        }

    }
    catch {

        // fallback below

    }


    return (
        `${prefix}-${Date.now()}-${Math.random()
            .toString(16)
            .slice(2)}`
    )

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
            createRandomID(
                prefix
            )


        localStorage.setItem(
            storageKey,
            id
        )


        return id

    }
    catch {

        return (
            createRandomID(
                prefix
            )
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
            info.clientType,

    })


    console.log(
        "[WS] Registered:",
        info.name,
        id
    )

}


// =====================================================
// CLEAR RECONNECT
// =====================================================

function clearReconnectTimer() {

    if (
        reconnectTimer ===
        null
    ) {

        return

    }


    clearTimeout(
        reconnectTimer
    )


    reconnectTimer =
        null

}


// =====================================================
// RECONNECT
// =====================================================

function scheduleReconnect() {

    if (
        reconnectTimer !==
        null
    ) {

        return

    }


    console.log(
        "[WS] Reconnect in 3 seconds..."
    )


    reconnectTimer =
        setTimeout(
            () => {

                reconnectTimer =
                    null


                connectRealtime()

            },
            3000
        )

}


// =====================================================
// CONNECT
// =====================================================

export function connectRealtime() {

    // =================================================
    // ALREADY CONNECTED
    // =================================================

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


    clearReconnectTimer()


    // =================================================
    // URL
    // =================================================

    const url =
        getRealtimeURL()


    console.log(
        "[WS] Connecting:",
        url
    )


    // =================================================
    // CREATE SOCKET
    // =================================================

    const ws =
        new WebSocket(
            url
        )


    socket =
        ws


    // =================================================
    // OPEN
    // =================================================

    ws.onopen =
        () => {

            if (
                socket !==
                ws
            ) {

                return

            }


            connected.value =
                true


            console.log(
                "[WS] Connected:",
                url
            )


            registerClient()

        }


    // =================================================
    // RECEIVE
    // =================================================

    ws.onmessage =
        event => {

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

                        try {

                            callback(
                                data
                            )

                        }
                        catch (error) {

                            console.error(
                                "[WS LISTENER ERROR]",
                                error
                            )

                        }

                    }
                )

            }
            catch (error) {

                console.error(
                    "[WS PARSE ERROR]",
                    error
                )

            }

        }


    // =================================================
    // ERROR
    // =================================================

    ws.onerror =
        error => {

            console.error(
                "[WS ERROR]",
                error
            )

        }


    // =================================================
    // CLOSE
    // =================================================

    ws.onclose =
        event => {

            if (
                socket !==
                ws
            ) {

                return

            }


            connected.value =
                false


            socket =
                null


            console.log(
                "[WS] Disconnected:",
                event.code,
                event.reason
            )


            scheduleReconnect()

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
        socket.readyState !==
        WebSocket.OPEN
    ) {

        console.warn(
            "[WS] Not connected",
            data
        )

        return false

    }


    try {

        socket.send(
            JSON.stringify(
                data
            )
        )


        return true

    }
    catch (error) {

        console.error(
            "[WS SEND ERROR]",
            error
        )


        return false

    }

}


// =====================================================
// SUBSCRIBE
// =====================================================

export function onRealtimeMessage(
    callback:
        (
            data: any
        ) => void
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
            onRealtimeMessage,

    }

}