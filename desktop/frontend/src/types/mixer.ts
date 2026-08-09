export interface MixerChannel {
  id: number
  name: string
  app: string
  volume: number
  muted: boolean
  connected: boolean
  selected: boolean
}

export interface MixerCommand {
  type: "ENC" | "BTN"
  channel: number
  value: number
}

export interface RealtimeChannel {
  id: number
  volume?: number
  muted?: boolean
  name?: string
  app?: string
  connected?: boolean
  selected?: boolean
}

export interface RealtimeMessage {
  type: string
  channel?: RealtimeChannel
}