import { reactive } from "vue"

import {
  GetChannels,
} from "../../wailsjs/go/main/App"

import type {
  MixerChannel,
  MixerCommand,
  RealtimeMessage,
} from "../types/mixer"

import {
  sendRealtime,
} from "../composables/useRealtime"


/*
|--------------------------------------------------------------------------
| MOCK CHANNELS
|--------------------------------------------------------------------------
|
| Dipakai ketika frontend dibuka melalui browser biasa.
|
| npm run dev
| http://localhost:5173
|
| Browser tidak mempunyai Wails Go binding.
|
|--------------------------------------------------------------------------
*/

const mockChannels: MixerChannel[] = [

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


/*
|--------------------------------------------------------------------------
| WAILS ENVIRONMENT CHECK
|--------------------------------------------------------------------------
*/

function isWailsEnvironment(): boolean {

  if (
    typeof window === "undefined"
  ) {
    return false
  }


  /*
  |--------------------------------------------------------------------------
  | Wails injects runtime into window.
  |--------------------------------------------------------------------------
  */

  return (
    "__WAILS_RUNTIME__" in window
  )
}


/*
|--------------------------------------------------------------------------
| NORMALIZE VOLUME
|--------------------------------------------------------------------------
*/

function normalizeVolume(
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


/*
|--------------------------------------------------------------------------
| STORE
|--------------------------------------------------------------------------
*/

export const mixerStore = reactive({

  /*
  |--------------------------------------------------------------------------
  | CHANNELS
  |--------------------------------------------------------------------------
  */

  channels:
    [] as MixerChannel[],


  /*
  |--------------------------------------------------------------------------
  | LOADING
  |--------------------------------------------------------------------------
  */

  loading:
    false,


  /*
  |--------------------------------------------------------------------------
  | DEVICE CONNECTION
  |--------------------------------------------------------------------------
  */

  deviceConnected:
    false,


  /*
  |--------------------------------------------------------------------------
  | LOAD CHANNELS
  |--------------------------------------------------------------------------
  */

  async loadChannels() {

    this.loading =
      true


    console.log(
      "[MIXER] Loading channels..."
    )


    try {

      /*
      |--------------------------------------------------------------------------
      | BROWSER MODE
      |--------------------------------------------------------------------------
      */

      if (
        !isWailsEnvironment()
      ) {

        console.log(
          "[MIXER] Browser mode detected."
        )

        console.log(
          "[MIXER] Using mock channels."
        )


        this.channels =
          structuredClone(
            mockChannels
          )


        return
      }


      /*
      |--------------------------------------------------------------------------
      | WAILS MODE
      |--------------------------------------------------------------------------
      */

      console.log(
        "[MIXER] Wails environment detected."
      )


      const result =
        await GetChannels()


      console.log(
        "[MIXER] Wails result:",
        result
      )


      /*
      |--------------------------------------------------------------------------
      | BACKEND RETURNED CHANNELS
      |--------------------------------------------------------------------------
      */

      if (
        Array.isArray(result) &&
        result.length > 0
      ) {

        this.channels =
          result.map(
            channel => ({

              id:
                channel.id,

              name:
                channel.name,

              app:
                channel.app,

              volume:
                normalizeVolume(
                  channel.volume
                ),

              muted:
                channel.muted,

              connected:
                channel.connected,

              selected:
                channel.selected,

            })
          )


        return
      }


      /*
      |--------------------------------------------------------------------------
      | EMPTY BACKEND RESULT
      |--------------------------------------------------------------------------
      */

      console.warn(
        "[MIXER] GetChannels returned empty."
      )

      console.log(
        "[MIXER] Using mock channels."
      )


      this.channels =
        structuredClone(
          mockChannels
        )

    } catch (error) {

      /*
      |--------------------------------------------------------------------------
      | BACKEND ERROR
      |--------------------------------------------------------------------------
      */

      console.warn(
        "[MIXER] Failed to load backend channels."
      )

      console.warn(
        error
      )


      /*
      |--------------------------------------------------------------------------
      | FALLBACK
      |--------------------------------------------------------------------------
      */

      this.channels =
        structuredClone(
          mockChannels
        )

    } finally {

      this.loading =
        false
    }

  },


  /*
  |--------------------------------------------------------------------------
  | SET VOLUME
  |--------------------------------------------------------------------------
  |
  | Dipakai untuk perubahan dari UI.
  |
  | Contoh:
  |
  | HP
  | Laptop
  | Wails
  |
  |--------------------------------------------------------------------------
  */

  setVolume(
    id: number,
    volume: number,
  ) {

    const channel =
      this.channels.find(
        item =>
          item.id === id
      )


    if (!channel) {

      console.warn(
        `[MIXER] Channel ${id} not found.`
      )

      return
    }


    const normalized =
      normalizeVolume(
        volume
      )


    /*
    |--------------------------------------------------------------------------
    | UPDATE LOCAL UI
    |--------------------------------------------------------------------------
    */

    channel.volume =
      normalized


    console.log(
      `[MIXER] LOCAL ${channel.name}: ${normalized}%`
    )


    /*
    |--------------------------------------------------------------------------
    | BROADCAST TO OTHER DEVICES
    |--------------------------------------------------------------------------
    */

    sendRealtime({

      type:
        "CHANNEL_UPDATE",

      channel: {

        id:
          channel.id,

        volume:
          normalized,

        muted:
          channel.muted,

      },

    })

  },


  /*
  |--------------------------------------------------------------------------
  | APPLY REMOTE UPDATE
  |--------------------------------------------------------------------------
  |
  | Dipanggil ketika:
  |
  | HP
  | Laptop
  | Wails lain
  |
  | mengirim perubahan melalui WebSocket.
  |
  | JANGAN panggil setVolume() di sini.
  |
  |--------------------------------------------------------------------------
  */

  applyRemoteUpdate(
    message: RealtimeMessage
  ) {

    /*
    |--------------------------------------------------------------------------
    | CHECK MESSAGE TYPE
    |--------------------------------------------------------------------------
    */

    if (
      message.type !==
      "CHANNEL_UPDATE"
    ) {
      return
    }


    /*
    |--------------------------------------------------------------------------
    | CHECK CHANNEL
    |--------------------------------------------------------------------------
    */

    if (
      !message.channel
    ) {
      return
    }


    const remote =
      message.channel


    /*
    |--------------------------------------------------------------------------
    | FIND LOCAL CHANNEL
    |--------------------------------------------------------------------------
    */

    const channel =
      this.channels.find(
        item =>
          item.id ===
          remote.id
      )


    if (!channel) {

      console.warn(
        `[MIXER] Remote channel ${remote.id} not found.`
      )

      return
    }


    /*
    |--------------------------------------------------------------------------
    | UPDATE VOLUME
    |--------------------------------------------------------------------------
    */

    if (
      typeof remote.volume ===
      "number"
    ) {

      channel.volume =
        normalizeVolume(
          remote.volume
        )
    }


    /*
    |--------------------------------------------------------------------------
    | UPDATE MUTE
    |--------------------------------------------------------------------------
    */

    if (
      typeof remote.muted ===
      "boolean"
    ) {

      channel.muted =
        remote.muted
    }


    console.log(
      `[MIXER] REMOTE ${channel.name}: ${channel.volume}%`
    )

  },


  /*
  |--------------------------------------------------------------------------
  | HANDLE ESP32 COMMAND
  |--------------------------------------------------------------------------
  */

  handleCommand(
    command: MixerCommand
  ) {

    console.log(
      "[MIXER] ESP32:",
      command
    )


    /*
    |--------------------------------------------------------------------------
    | FIND CHANNEL
    |--------------------------------------------------------------------------
    */

    const channel =
      this.channels.find(
        item =>
          item.id ===
          command.channel
      )


    if (!channel) {

      console.warn(
        `[MIXER] Channel ${command.channel} not found.`
      )

      return
    }


    /*
    |--------------------------------------------------------------------------
    | ROTARY ENCODER
    |--------------------------------------------------------------------------
    */

    if (
      command.type ===
      "ENC"
    ) {

      const newVolume =
        channel.volume +
        command.value


      this.setVolumeFromESP32(
        channel.id,
        newVolume
      )


      return
    }


    /*
    |--------------------------------------------------------------------------
    | BUTTON
    |--------------------------------------------------------------------------
    */

    if (
      command.type === "BTN" &&
      command.value === 1
    ) {

      channel.muted =
        !channel.muted


      console.log(
        `[MIXER] ESP32 ${channel.name} muted: ${channel.muted}`
      )


      /*
      |--------------------------------------------------------------------------
      | BROADCAST MUTE
      |--------------------------------------------------------------------------
      */

      sendRealtime({

        type:
          "CHANNEL_UPDATE",

        channel: {

          id:
            channel.id,

          volume:
            channel.volume,

          muted:
            channel.muted,

        },

      })
    }

  },


  /*
  |--------------------------------------------------------------------------
  | ESP32 SET VOLUME
  |--------------------------------------------------------------------------
  |
  | Update UI Wails.
  | Kemudian broadcast ke HP/laptop.
  |
  |--------------------------------------------------------------------------
  */

  setVolumeFromESP32(
    id: number,
    volume: number,
  ) {

    const channel =
      this.channels.find(
        item =>
          item.id === id
      )


    if (!channel) {

      console.warn(
        `[MIXER] ESP32 channel ${id} not found.`
      )

      return
    }


    const normalized =
      normalizeVolume(
        volume
      )


    /*
    |--------------------------------------------------------------------------
    | UPDATE WAILS UI
    |--------------------------------------------------------------------------
    */

    channel.volume =
      normalized


    console.log(
      `[MIXER] ESP32 ${channel.name}: ${normalized}%`
    )


    /*
    |--------------------------------------------------------------------------
    | BROADCAST TO BROWSER DEVICES
    |--------------------------------------------------------------------------
    */

    sendRealtime({

      type:
        "CHANNEL_UPDATE",

      channel: {

        id:
          channel.id,

        volume:
          normalized,

        muted:
          channel.muted,

      },

    })

  },

})