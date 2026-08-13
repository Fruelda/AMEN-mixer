import { isRealtimeWailsEnvironment } from "./environment"

interface ClientInfo {
    prefix: string
    name: string
    clientType:
    | "desktop"
    | "mobile"
    | "tablet"
    | "browser"
}

export interface ClientRegistration {
    id: string
    name: string
    clientType: ClientInfo["clientType"]
}

function getClientInfo(): ClientInfo {
    if (isRealtimeWailsEnvironment()) {
        return {
            prefix: "desktop",
            name: "Wails Desktop",
            clientType: "desktop",
        }
    }

    const userAgent = navigator.userAgent.toLowerCase()

    if (/iphone|ipod/.test(userAgent)) {
        return {
            prefix: "iphone",
            name: "iPhone",
            clientType: "mobile",
        }
    }

    const isIPad =
        /ipad/.test(userAgent) ||
        (
            navigator.platform === "MacIntel" &&
            navigator.maxTouchPoints > 1
        )

    if (isIPad) {
        return {
            prefix: "ipad",
            name: "iPad",
            clientType: "tablet",
        }
    }

    if (/android/.test(userAgent)) {
        const isTablet = !/mobile/.test(userAgent)

        return isTablet
            ? {
                prefix: "android-tablet",
                name: "Android Tablet",
                clientType: "tablet",
            }
            : {
                prefix: "android",
                name: "Android Phone",
                clientType: "mobile",
            }
    }

    return {
        prefix: "browser",
        name: "AMEN Browser",
        clientType: "browser",
    }
}

function createRandomID(prefix: string): string {
    try {
        if (
            typeof crypto !== "undefined" &&
            typeof crypto.randomUUID === "function"
        ) {
            return `${prefix}-${crypto.randomUUID()}`
        }
    } catch {
        // fallback ke timestamp
    }

    return (
        `${prefix}-` +
        `${Date.now()}-` +
        `${Math.random().toString(16).slice(2)}`
    )
}

function getClientID(prefix: string): string {
    const storageKey = "amen-mixer-client-id"

    try {
        const savedID = localStorage.getItem(storageKey)

        if (savedID) {
            return savedID
        }

        const id = createRandomID(prefix)

        localStorage.setItem(storageKey, id)

        return id
    } catch {
        return createRandomID(prefix)
    }
}

export function getClientRegistration(): ClientRegistration {
    const info = getClientInfo()

    return {
        id: getClientID(info.prefix),
        name: info.name,
        clientType: info.clientType,
    }
}