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


function isWailsEnvironment(): boolean {

  if (
    typeof window === "undefined"
  ) {
    return false
  }


  return (
    "__WAILS_RUNTIME__" in window
  )

}



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



export const mixerStore = reactive({

  channels:
    [] as MixerChannel[],


  loading:
    false,


  deviceConnected:
    false,



  /*
  |--------------------------------------------------------------------------
  | STATE CONTROL
  |--------------------------------------------------------------------------
  */


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



  /*
  |--------------------------------------------------------------------------
  | LOAD CHANNELS
  |--------------------------------------------------------------------------
  */


  async loadChannels() {

    this.setLoading(true)


    try {


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



      const result =
        await GetChannels()



      if (
        Array.isArray(result)
        &&
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


      this.setLoading(false)


    }

  },



  /*
  |--------------------------------------------------------------------------
  | VOLUME
  |--------------------------------------------------------------------------
  */


  setVolume(
    id: number,
    volume: number
  ) {


    const channel =
      this.channels.find(
        item =>
          item.id === id
      )


    if (!channel)
      return



    const value =
      normalizeVolume(
        volume
      )



    channel.volume =
      value



    sendRealtime({

      type:
        "CHANNEL_UPDATE",


      channel: {

        id:
          channel.id,


        volume:
          value,


        muted:
          channel.muted

      }

    })


  },



  /*
  |--------------------------------------------------------------------------
  | REMOTE UPDATE
  |--------------------------------------------------------------------------
  */


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



    if (!channel)
      return



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



  /*
  |--------------------------------------------------------------------------
  | ESP32 COMMAND
  |--------------------------------------------------------------------------
  */


  handleCommand(
    command: MixerCommand
  ) {


    const channel =
      this.channels.find(
        item =>
          item.id ===
          command.channel
      )


    if (!channel)
      return



    if (
      command.type ===
      "ENC"
    ) {

      this.setVolume(
        channel.id,
        channel.volume +
        command.value
      )

    }



    if (
      command.type ===
      "BTN"
      &&
      command.value === 1
    ) {

      channel.muted =
        !channel.muted


    }


  },



  setVolumeFromESP32(
    id: number,
    volume: number
  ) {

    this.setVolume(
      id,
      volume
    )

  },


})