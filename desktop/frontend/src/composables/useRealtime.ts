import { ref } from "vue"

import { getClientRegistration } from "../realtime/client"
import { getRealtimeURL } from "../realtime/environment"

import type { ClientRegistration } from "../realtime/client"
import type {
    ChannelUpdateMessage,
    RealtimeMessage,
} from "../types/mixer"


type RealtimeListener = (
    message: RealtimeMessage
) => void

type ClientRegisterMessage =
    ClientRegistration & {
        type: "client.register"
    }

type RealtimeOutgoingMessage =
    | ClientRegisterMessage
    | ChannelUpdateMessage


const connected = ref(false)

const listeners =
    new Set<RealtimeListener>()

let socket: WebSocket | null = null

let reconnectTimer:
    ReturnType<typeof setTimeout> | null =
    null


function registerClient() {
    const client =
        getClientRegistration()

    sendRealtime({
        type: "client.register",
        ...client,
    })
}


function handleMessage(
    rawData: string
) {
    try {
        const message =
            JSON.parse(rawData) as RealtimeMessage

        console.log(
            "[WS RECEIVE]",
            message
        )

        listeners.forEach(
            listener => listener(message)
        )
    } catch (error) {
        console.error(
            "[WS PARSE ERROR]",
            error
        )
    }
}


function scheduleReconnect() {
    if (reconnectTimer) return

    reconnectTimer =
        setTimeout(
            () => {
                reconnectTimer = null
                connectRealtime()
            },
            3000
        )
}


export function connectRealtime() {
    if (
        socket?.readyState ===
        WebSocket.OPEN ||
        socket?.readyState ===
        WebSocket.CONNECTING
    ) {
        return
    }

    if (reconnectTimer) {
        clearTimeout(reconnectTimer)
        reconnectTimer = null
    }

    const url =
        getRealtimeURL()

    const ws =
        new WebSocket(url)

    socket = ws


    ws.onopen = () => {
        if (socket !== ws) return

        connected.value = true

        console.log(
            "[WS] Connected:",
            url
        )

        registerClient()
    }


    ws.onmessage = event => {
        if (
            typeof event.data === "string"
        ) {
            handleMessage(event.data)
        }
    }


    ws.onerror = error => {
        console.error(
            "[WS ERROR]",
            error
        )
    }


    ws.onclose = () => {
        if (socket !== ws) return

        connected.value = false
        socket = null

        scheduleReconnect()
    }
}


export function sendRealtime(
    message: RealtimeOutgoingMessage
): boolean {
    if (
        socket?.readyState !==
        WebSocket.OPEN
    ) {
        return false
    }

    try {
        socket.send(
            JSON.stringify(message)
        )

        return true
    } catch (error) {
        console.error(
            "[WS SEND ERROR]",
            error
        )

        return false
    }
}


export function onRealtimeMessage(
    listener: RealtimeListener
) {
    listeners.add(listener)

    return () =>
        listeners.delete(listener)
}


export function useRealtime() {
    return {
        connected,
        connect: connectRealtime,
        send: sendRealtime,
        subscribe: onRealtimeMessage,
    }
}