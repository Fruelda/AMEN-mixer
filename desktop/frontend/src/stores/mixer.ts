import { reactive } from "vue"

import {
  GetChannels,
  SetVolume,
} from "../../wailsjs/go/main/App"

import {
  sendRealtime,
} from "../composables/useRealtime"

import type {
  ConnectedDevice,
  MixerChannel,
  MixerCommand,
  RealtimeMessage,
} from "../types/mixer"

import {
  isMixerWailsEnvironment,
  normalizeVolume,
} from "./mixer/helpers"

import {
  mockChannels,
} from "./mixer/mockChannels"

export const mixerStore = reactive({
  channels: [] as MixerChannel[],

  devices: [] as ConnectedDevice[],

  loading: false,

  deviceConnected: false,

  setChannels(channels: MixerChannel[]) {
    this.channels = channels.map((channel) => ({
      ...channel,

      volume: normalizeVolume(
        channel.volume,
      ),
    }))
  },

  setDevices(devices: ConnectedDevice[]) {
    this.devices = devices
  },

  setLoading(value: boolean) {
    this.loading = value
  },

  setConnected(value: boolean) {
    this.deviceConnected = value
  },

  async loadChannels() {
    this.setLoading(true)

    try {
      if (!isMixerWailsEnvironment()) {
        this.setChannels(
          structuredClone(mockChannels),
        )

        return
      }

      const result = await GetChannels()

      this.setChannels(
        Array.isArray(result) &&
          result.length > 0
          ? result
          : mockChannels,
      )
    } catch (error) {
      console.warn(
        "[MIXER] Load failed",
        error,
      )

      this.setChannels(mockChannels)
    } finally {
      this.setLoading(false)
    }
  },

  async setVolume(
    id: number,
    volume: number,
    broadcast = true,
  ) {
    const channel = this.channels.find(
      (item) => item.id === id,
    )

    if (!channel) {
      return
    }

    const value = normalizeVolume(volume)

    if (channel.volume === value) {
      return
    }

    channel.volume = value

    if (isMixerWailsEnvironment()) {
      try {
        await SetVolume(
          id,
          value,
        )

        console.log(
          `[WAILS] ${channel.name}: ${value}%`,
        )
      } catch (error) {
        console.warn(
          "[WAILS] SetVolume failed",
          error,
        )
      }

      /*
       * Backend Wails sudah melakukan
       * broadcast setelah SetVolume().
       */
      return
    }

    console.log(
      `[BROWSER] ${channel.name}: ${value}%`,
    )

    if (broadcast) {
      this.broadcastChannel(id)
    }
  },

  broadcastChannel(id: number) {
    const channel = this.channels.find(
      (item) => item.id === id,
    )

    if (!channel) {
      return
    }

    sendRealtime({
      type: "CHANNEL_UPDATE",

      channel: {
        id: channel.id,
        volume: channel.volume,
        muted: channel.muted,
      },
    })
  },

  applyRemoteUpdate(
    message: RealtimeMessage,
  ) {
    if (
      message.type !== "CHANNEL_UPDATE" ||
      !message.channel
    ) {
      return
    }

    const channel = this.channels.find(
      (item) =>
        item.id === message.channel?.id,
    )

    if (!channel) {
      return
    }

    /*
     * Remote update hanya mengubah
     * local state agar tidak terjadi
     * broadcast loop.
     */
    if (
      typeof message.channel.volume ===
      "number"
    ) {
      channel.volume = normalizeVolume(
        message.channel.volume,
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

  handleCommand(
    command: MixerCommand,
  ) {
    /*
     * COMMAND dari ESP32 hanya
     * diproses oleh Wails.
     *
     * Client lain menerima hasil
     * melalui CHANNEL_UPDATE.
     */
    if (!isMixerWailsEnvironment()) {
      return
    }

    const channel = this.channels.find(
      (item) =>
        item.id === command.channel,
    )

    if (!channel) {
      return
    }

    if (command.type === "ENC") {
      const newVolume = normalizeVolume(
        channel.volume +
        command.value,
      )

      void this.setVolume(
        channel.id,
        newVolume,
        false,
      )

      return
    }

    if (
      command.type === "BTN" &&
      command.value === 1
    ) {
      channel.muted =
        !channel.muted

      this.broadcastChannel(
        channel.id,
      )
    }
  },

  setVolumeFromESP32(
    id: number,
    volume: number,
  ) {
    void this.setVolume(
      id,
      volume,
      false,
    )
  },
})