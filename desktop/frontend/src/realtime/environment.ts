// ============================================================
// WAILS ENVIRONMENT
// ============================================================

export function isWailsEnvironment(): boolean {
    return (
        typeof window !== "undefined" &&
        "runtime" in window
    )
}

export const isRealtimeWailsEnvironment =
    isWailsEnvironment


// ============================================================
// REALTIME URL
// ============================================================

export function getRealtimeURL(): string {
    const envURL =
        import.meta.env.VITE_WS_URL?.trim()

    if (envURL) {
        return envURL
    }

    if (isWailsEnvironment()) {
        return "ws://192.168.1.44:8081/ws"
    }

    return `ws://${window.location.hostname}:8081/ws`
}