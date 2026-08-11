import { ref } from "vue"


const connected = ref(false)


let socket: WebSocket | null = null


const listeners = new Set<
    (data: any) => void
>()



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
            "ws://192.168.1.44:8081/ws"
        )



    socket.onopen = () => {


        connected.value =
            true



        console.log(
            "[WS] Connected"
        )



        sendRealtime({

            type:
                "client.register",


            id:
                "browser-" + Date.now(),


            name:
                "AMEN Browser",


            clientType:
                "browser"

        })



    }




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

                    callback(data)

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





    socket.onerror = error => {


        console.error(
            "[WS ERROR]",
            error
        )


    }




    socket.onclose = () => {


        connected.value =
            false



        console.log(
            "[WS] Disconnected"
        )



        socket = null



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
        JSON.stringify(data)
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