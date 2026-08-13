import {
  onMounted,
  onUnmounted
} from "vue"

import {
  EventsOn
} from "../../wailsjs/runtime/runtime"

import {
  isWailsEnvironment
} from "../realtime/environment"

import type {
  MixerCommand
} from "../types/mixer"


export function useSerial(
  callback: (
    command: MixerCommand
  ) => void
) {

  let unsubscribe:
    (() => void) | null =
    null


  // ============================================================
  // START
  // ============================================================

  function start() {

    if (
      !isWailsEnvironment()
    ) {

      console.log(
        "[SERIAL] Browser mode - disabled"
      )

      return
    }


    if (
      unsubscribe
    ) {
      return
    }


    console.log(
      "[SERIAL] Listening for serial-command..."
    )


    unsubscribe =
      EventsOn(
        "serial-command",
        (
          data: MixerCommand
        ) => {

          console.log(
            "[SERIAL] COMMAND:",
            data
          )


          callback(
            data
          )

        }
      )

  }


  // ============================================================
  // STOP
  // ============================================================

  function stop() {

    if (
      !unsubscribe
    ) {
      return
    }


    console.log(
      "[SERIAL] Stop listening"
    )


    unsubscribe()

    unsubscribe =
      null

  }


  // ============================================================
  // LIFECYCLE
  // ============================================================

  onMounted(
    start
  )


  onUnmounted(
    stop
  )


  // ============================================================
  // RETURN
  // ============================================================

  return {
    start,
    stop
  }

}