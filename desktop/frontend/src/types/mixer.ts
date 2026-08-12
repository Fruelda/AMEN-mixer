/*
|--------------------------------------------------------------------------
| AMEN Mixer Data Contract
|--------------------------------------------------------------------------
|
| Semua komunikasi:
|
| Vue Frontend
|      ↓
| WebSocket
|      ↓
| Go Backend
|      ↓
| Serial
|      ↓
| ESP32
|
|--------------------------------------------------------------------------
*/


/*
|--------------------------------------------------------------------------
| CHANNEL
|--------------------------------------------------------------------------
*/

export interface MixerChannel {

  id:
  number

  name:
  string

  app:
  string

  volume:
  number

  muted:
  boolean

  connected:
  boolean

  selected:
  boolean

}


/*
|--------------------------------------------------------------------------
| AUDIO UPDATE
|--------------------------------------------------------------------------
*/

export interface ChannelUpdate {

  id:
  number

  volume?:
  number

  muted?:
  boolean

}


/*
|--------------------------------------------------------------------------
| ESP32 COMMAND
|--------------------------------------------------------------------------
*/

export type MixerCommandType =
  | "ENC"
  | "BTN"


export interface MixerCommand {

  type:
  MixerCommandType


  /*
  channel yang dikontrol

  contoh:
  1 = Master
  2 = Browser
  */

  channel:
  number


  /*
  ENC:
  +1 / -1

  BTN:
  1 = pressed
  */

  value:
  number

}


/*
|--------------------------------------------------------------------------
| CONNECTED DEVICE
|--------------------------------------------------------------------------
*/

export interface ConnectedDevice {

  id:
  string

  name:
  string

  clientType:
  string

  connected:
  boolean

}


/*
|--------------------------------------------------------------------------
| REALTIME MESSAGE
|--------------------------------------------------------------------------
*/

export type RealtimeMessageType =

  | "STATE"

  | "CHANNEL_UPDATE"

  | "COMMAND"

  | "DEVICE_STATUS"

  | "DEVICES"


export interface RealtimeMessage {

  type:
  RealtimeMessageType


  /*
  dipakai STATE
  */

  channels?:
  MixerChannel[]


  /*
  dipakai CHANNEL_UPDATE
  */

  channel?:
  ChannelUpdate


  /*
  dipakai COMMAND
  */

  command?:
  MixerCommand


  /*
  status ESP32
  */

  connected?:
  boolean


  /*
  daftar client/device
  yang sedang terhubung
  */

  devices?:
  ConnectedDevice[]

}


/*
|--------------------------------------------------------------------------
| DEVICE STATUS
|--------------------------------------------------------------------------
*/

export interface DeviceStatus {

  connected:
  boolean

  name?:
  string

  port?:
  string

}