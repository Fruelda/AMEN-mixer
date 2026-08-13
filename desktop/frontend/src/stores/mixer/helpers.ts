import {
    isWailsEnvironment
} from "../../realtime/environment"


// ============================================================
// ENVIRONMENT
// ============================================================

export function isMixerWailsEnvironment(): boolean {
    return isWailsEnvironment()
}


// ============================================================
// VOLUME
// ============================================================

export function normalizeVolume(
    volume: number
): number {

    if (
        !Number.isFinite(volume)
    ) {
        return 0
    }

    return Math.round(
        Math.max(
            0,
            Math.min(
                100,
                volume
            )
        )
    )
}