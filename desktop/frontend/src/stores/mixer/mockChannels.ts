import type {
    MixerChannel,
} from "../../types/mixer"

export const mockChannels: MixerChannel[] = [
    {
        id: 1,
        name: "Master",
        app: "master",
        volume: 100,
        muted: false,
        connected: true,
        selected: false,
    },

    {
        id: 2,
        name: "Browser",
        app: "browser",
        volume: 70,
        muted: false,
        connected: true,
        selected: false,
    },

    {
        id: 3,
        name: "Spotify",
        app: "spotify",
        volume: 55,
        muted: false,
        connected: true,
        selected: false,
    },

    {
        id: 4,
        name: "Discord",
        app: "discord",
        volume: 85,
        muted: false,
        connected: true,
        selected: false,
    },

    {
        id: 5,
        name: "Valeton",
        app: "valeton",
        volume: 60,
        muted: false,
        connected: false,
        selected: false,
    },
]