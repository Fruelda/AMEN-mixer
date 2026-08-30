import { reactive } from "vue"

import {
  GetChannels,
  SetVolume
} from "../../wailsjs/go/main/App"

import type {
  MixerChannel,
  MixerCommand,
  RealtimeMessage,
  ConnectedDevice
} from "../types/mixer"

import {
  sendRealtime
} from "../composables/useRealtime"


// ============================================================
// MOCK DATA
// ============================================================

const mockChannels: MixerChannel[] = [

  {
    id: 1,
    name: "Master",
    app: "master",
    volume: 100,
    muted: false,
    connected: true,
    selected: false
  },

  {
    id: 2,
    name: "Browser",
    app: "browser",
    volume: 70,
    muted: false,
    connected: true,
    selected: false
  },

  {
    id: 3,
    name: "Spotify",
    app: "spotify",
    volume: 55,
    muted: false,
    connected: true,
    selected: false
  },

  {
    id: 4,
    name: "Discord",
    app: "discord",
    volume: 85,
    muted: false,
    connected: true,
    selected: false
  },

  {
    id: 5,
    name: "Valeton",
    app: "valeton",
    volume: 60,
    muted: false,
    connected: false,
    selected: false
  }

]


// ============================================================
// ENV
// ============================================================

function isWailsEnvironment() {

  if (
    typeof window === "undefined"
  ) {
    return false
  }

  return (
    typeof (window as any).runtime !== "undefined"
  )

}


// ============================================================
// HELPER
// ============================================================

function normalizeVolume(
  volume: number
) {

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


// ============================================================
// STORE
// ============================================================

export const mixerStore = reactive({

  // ========================================================
  // DATA
  // ========================================================

  channels:
    [] as MixerChannel[],

  devices:
    [] as ConnectedDevice[],

  loading:
    false,

  deviceConnected:
    false,


  // ========================================================
  // CHANNEL STATE
  // ========================================================

  setChannels(
    channels: MixerChannel[]
  ) {

    this.channels =
      channels.map(
        channel => ({

          ...channel,

          volume:
            normalizeVolume(
              channel.volume
            )

        })
      )

  },


  // ========================================================
  // DEVICE STATE
  // ========================================================

  setDevices(
    devices: ConnectedDevice[]
  ) {

    this.devices =
      devices

  },


  setLoading(
    value: boolean
  ) {

    this.loading =
      value

  },


  setConnected(
    value: boolean
  ) {

    this.deviceConnected =
      value

  },


  // ========================================================
  // LOAD CHANNEL
  // ========================================================

  async loadChannels() {

    this.setLoading(
      true
    )


    try {

      // =================================================
      // BROWSER MODE
      // =================================================

      if (
        !isWailsEnvironment()
      ) {

        this.setChannels(
          structuredClone(
            mockChannels
          )
        )

        return

      }


      // =================================================
      // WAILS MODE
      // =================================================

      const result =
        await GetChannels()


      if (
        Array.isArray(result) &&
        result.length > 0
      ) {

        this.setChannels(
          result
        )

        return

      }


      this.setChannels(
        mockChannels
      )

    }
    catch (error) {

      console.warn(
        "[MIXER] Load failed",
        error
      )


      this.setChannels(
        mockChannels
      )

    }
    finally {

      this.setLoading(
        false
      )

    }

  },


  // ========================================================
  // USER CHANGE
  // ========================================================

  async setVolume(
    id: number,
    volume: number,
    broadcast = true
  ) {

    const channel =
      this.channels.find(
        item =>
          item.id === id
      )


    if (
      !channel
    ) {
      return
    }


    const value =
      normalizeVolume(
        volume
      )


    // =====================================================
    // IGNORE SAME VALUE
    // =====================================================

    if (
      channel.volume === value
    ) {
      return
    }


    // =====================================================
    // UPDATE LOCAL UI
    // =====================================================

    channel.volume =
      value


    // =====================================================
    // WAILS AUDIO
    // =====================================================

    if (
      isWailsEnvironment()
    ) {

      try {

        await SetVolume(
          id,
          value
        )


        console.log(
          `[WAILS] ${channel.name}: ${value}%`
        )

      }
      catch (error) {

        console.warn(
          "[WAILS] SetVolume failed",
          error
        )

      }


      /*
      IMPORTANT

      Jangan broadcast lagi dari frontend Wails.

      SetVolume()
          ↓
      Go Audio Manager
          ↓
      Audio Event Bridge
          ↓
      BroadcastChannelUpdate()

      Broadcast sudah dilakukan backend.
      */

      return

    }


    // =====================================================
    // BROWSER / MOBILE
    // =====================================================

    console.log(
      `[BROWSER] ${channel.name}: ${value}%`
    )


    if (
      broadcast
    ) {

      this.broadcastChannel(
        id
      )

    }

  },


  // ========================================================
  // NETWORK SEND
  // ========================================================

  broadcastChannel(
    id: number
  ) {

    const channel =
      this.channels.find(
        item =>
          item.id === id
      )


    if (
      !channel
    ) {
      return
    }


    sendRealtime({

      type:
        "CHANNEL_UPDATE",

      channel: {

        id:
          channel.id,

        volume:
          channel.volume,

        muted:
          channel.muted

      }

    })

  },


  // ========================================================
  // RECEIVE CHANNEL UPDATE
  // ========================================================

  applyRemoteUpdate(
    message: RealtimeMessage
  ) {

    if (
      message.type !==
      "CHANNEL_UPDATE"
    ) {
      return
    }


    if (
      !message.channel
    ) {
      return
    }


    const channel =
      this.channels.find(
        item =>
          item.id ===
          message.channel?.id
      )


    if (
      !channel
    ) {
      return
    }


    /*
    IMPORTANT

    Jangan panggil:
    setVolume()
    broadcastChannel()

    Remote update hanya update local state.
    */


    if (
      typeof message.channel.volume ===
      "number"
    ) {

      channel.volume =
        normalizeVolume(
          message.channel.volume
        )

    }


    if (
      typeof message.channel.muted ===
      "boolean"
    ) {

      channel.muted =
        message.channel.muted

    }

  },


  // ========================================================
  // ESP32 COMMAND
  // ========================================================

  handleCommand(
    command: MixerCommand
  ) {

    /*
    COMMAND dari ESP32 hanya diproses oleh Wails.

    Browser / HP / iPad cukup menerima
    CHANNEL_UPDATE hasil akhirnya.

    Ini mencegah beberapa client
    mengirim update yang sama.
    */

    if (
      !isWailsEnvironment()
    ) {
      return
    }


    const channel =
      this.channels.find(
        item =>
          item.id ===
          command.channel
      )


    if (
      !channel
    ) {
      return
    }


    // =====================================================
    // ENCODER
    // =====================================================

    if (
      command.type === "ENC"
    ) {

      const newVolume =
        normalizeVolume(
          channel.volume +
          command.value
        )


      void this.setVolume(
        channel.id,
        newVolume,
        false
      )


      return

    }


    // =====================================================
    // BUTTON
    // =====================================================

    if (
      command.type === "BTN" &&
      command.value === 1
    ) {

      channel.muted =
        !channel.muted


      this.broadcastChannel(
        channel.id
      )

    }

  },


  // ========================================================
  // ESP32 DIRECT VOLUME
  // ========================================================

  setVolumeFromESP32(
    id: number,
    volume: number
  ) {

    void this.setVolume(
      id,
      volume,
      false
    )

  }

})