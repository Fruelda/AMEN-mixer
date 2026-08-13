export function isWailsEnvironment(): boolean {
    return (
        typeof window !== "undefined" &&
        "runtime" in window
    )
}


/*
|--------------------------------------------------------------------------
| BACKWARD COMPATIBILITY
|--------------------------------------------------------------------------
|
| useRealtime/client.ts yang sudah kita buat sebelumnya
| masih memakai nama ini.
|
*/

export const isRealtimeWailsEnvironment =
    isWailsEnvironment


/*
|--------------------------------------------------------------------------
| REALTIME URL
|--------------------------------------------------------------------------
*/

export function getRealtimeURL(): string {
    if (isWailsEnvironment()) {
        return "ws://127.0.0.1:8081/ws"
    }

    return (
        `ws://${window.location.hostname}:8081/ws`
    )
}