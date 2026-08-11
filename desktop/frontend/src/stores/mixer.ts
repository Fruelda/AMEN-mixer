import { reactive } from "vue"

import {
  GetChannels,
  SetVolume
} from "../../wailsjs/go/main/App"


import type {
  MixerChannel,
  MixerCommand,
  RealtimeMessage
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
    "__WAILS_RUNTIME__" in window
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


  channels:
    [] as MixerChannel[],



  loading: false,


  deviceConnected: false,



  // ============================================================
  // STATE
  // ============================================================


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

    this.loading = value

  },




  setConnected(
    value: boolean
  ) {

    this.deviceConnected = value

  },




  // ============================================================
  // LOAD
  // ============================================================


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





  // ============================================================
  // USER CHANGE
  // ============================================================


  async setVolume(
    id: number,
    volume: number,
    broadcast = true
  ) {


    const channel =
      this.channels.find(
        item => item.id === id
      )



    if (!channel)
      return



    const value =
      normalizeVolume(
        volume
      )




    // update UI

    channel.volume = value





    // WAILS AUDIO

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


    }
    else {


      console.log(
        `[BROWSER] ${channel.name}: ${value}%`
      )


    }




    // broadcast hanya dari user input

    if (
      broadcast
    ) {


      this.broadcastChannel(
        id
      )


    }



  },





  // ============================================================
  // NETWORK SEND
  // ============================================================


  broadcastChannel(
    id: number
  ) {


    const channel =
      this.channels.find(
        item => item.id === id
      )



    if (!channel)
      return




    sendRealtime({

      type: "CHANNEL_UPDATE",


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





  // ============================================================
  // RECEIVE NETWORK
  // ============================================================


  applyRemoteUpdate(
    message: RealtimeMessage
  ) {



    if (
      message.type !== "CHANNEL_UPDATE"
    )
      return



    if (
      !message.channel
    )
      return




    const channel =
      this.channels.find(
        item =>
          item.id === message.channel?.id
      )



    if (!channel)
      return





    // IMPORTANT
    // jangan panggil setVolume()
    // karena akan broadcast ulang



    if (
      typeof message.channel.volume === "number"
    ) {


      channel.volume =
        normalizeVolume(
          message.channel.volume
        )



    }




    if (
      typeof message.channel.muted === "boolean"
    ) {


      channel.muted =
        message.channel.muted


    }




  },





  // ============================================================
  // ESP32 COMMAND
  // ============================================================


  handleCommand(
    command: MixerCommand
  ) {


    const channel =
      this.channels.find(
        item =>
          item.id === command.channel
      )



    if (!channel)
      return





    // ENCODER

    if (
      command.type === "ENC"
    ) {


      channel.volume =
        normalizeVolume(
          channel.volume +
          command.value
        )



      // kirim hasil ESP32 ke client lain

      this.broadcastChannel(
        channel.id
      )



    }





    // BUTTON

    if (
      command.type === "BTN"
      &&
      command.value === 1
    ) {


      channel.muted =
        !channel.muted



      this.broadcastChannel(
        channel.id
      )


    }



  },





  setVolumeFromESP32(
    id: number,
    volume: number
  ) {


    this.setVolume(
      id,
      volume,
      false
    )


  }



})