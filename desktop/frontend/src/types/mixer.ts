// ============================================================
// CHANNEL
// ============================================================

export interface MixerChannel {
  id: number
  name: string
  app: string
  volume: number
  muted: boolean
  connected: boolean
  selected: boolean
}


// ============================================================
// CHANNEL UPDATE
// ============================================================

export interface ChannelUpdate {
  id: number
  volume?: number
  muted?: boolean
}


// ============================================================
// MIXER COMMAND
// ============================================================

export type MixerCommandType =
  | "ENC"
  | "BTN"

export interface MixerCommand {
  type: MixerCommandType
  channel: number
  value: number
}


// ============================================================
// CONNECTED DEVICE
// ============================================================

export type ClientType =
  | "desktop"
  | "mobile"
  | "tablet"
  | "browser"
  | "hardware"

export interface ConnectedDevice {
  id: string
  name: string
  clientType: ClientType
  connected: boolean
}


// ============================================================
// STATE MESSAGE
// ============================================================

export interface StateMessage {
  type: "STATE"
  channels: MixerChannel[]
}


// ============================================================
// CHANNEL UPDATE MESSAGE
// ============================================================

export interface ChannelUpdateMessage {
  type: "CHANNEL_UPDATE"
  channel: ChannelUpdate
}


// ============================================================
// COMMAND MESSAGE
// ============================================================

export interface CommandMessage {
  type: "COMMAND"
  command: MixerCommand
}


// ============================================================
// DEVICES MESSAGE
// ============================================================

export interface DevicesMessage {
  type: "DEVICES"
  devices: ConnectedDevice[]
}


// ============================================================
// DEVICE STATUS MESSAGE
// ============================================================

export interface DeviceStatusMessage {
  type: "DEVICE_STATUS"
  connected: boolean
}

// ============================================================
// ESP32 MIXER COMMAND MESSAGE
// ============================================================

export interface ESP32MixerCommandMessage {

  type: "mixer.command"

  device: string

  channel: number

  value: number

  seource?: string

}

// ============================================================
// REALTIME MESSAGE
// ============================================================

export type RealtimeMessage =

  | StateMessage

  | ChannelUpdateMessage

  | CommandMessage

  | ESP32MixerCommandMessage

  | DevicesMessage

  | DeviceStatusMessage