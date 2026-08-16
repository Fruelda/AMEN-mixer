import { reactive } from "vue"

import { createState } from "./mixer/state"
import { createChannels } from "./mixer/channels"
import { createVolume } from "./mixer/volume"
import { createMute } from "./mixer/mute"
import { createDevices } from "./mixer/devices"
import { createRealtime } from "./mixer/realtime"
import { createCommands } from "./mixer/commands"


export const mixerStore = reactive({

  ...createState(),

  ...createChannels(),

  ...createVolume(),

  ...createMute(),

  ...createDevices(),

  ...createRealtime(),

  ...createCommands(),

})