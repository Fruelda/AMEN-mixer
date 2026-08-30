import {
  onMounted,
  onUnmounted,
} from "vue"

import {
  EventsOn,
  EventsOff,
} from "../../wailsjs/runtime/runtime"


export interface SerialCommand {
  type: "ENC" | "BTN"
  channel: number
  value: number
}


/*
|--------------------------------------------------------------------------
| DETECT WAILS
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
  | Wails runtime
  |--------------------------------------------------------------------------
  */

  return (
    typeof (window as any).runtime !== "undefined"
  )
}


/*
|--------------------------------------------------------------------------
| USE SERIAL
|--------------------------------------------------------------------------
*/

export function useSerial(
  callback: (
    command: SerialCommand
  ) => void
) {

  let wailsEnabled =
    false


  /*
  |--------------------------------------------------------------------------
  | MOUNT
  |--------------------------------------------------------------------------
  */

  onMounted(() => {

    /*
    |--------------------------------------------------------------------------
    | BROWSER MODE
    |--------------------------------------------------------------------------
    */

    if (
      !isWailsEnvironment()
    ) {

      console.log(
        "[SERIAL] Browser mode."
      )

      console.log(
        "[SERIAL] Wails serial disabled."
      )

      return
    }


    /*
    |--------------------------------------------------------------------------
    | WAILS MODE
    |--------------------------------------------------------------------------
    */

    wailsEnabled =
      true


    console.log(
      "[SERIAL] Listening for serial-command..."
    )


    EventsOn(
      "serial-command",
      (
        data: SerialCommand
      ) => {

        console.log(
          "[SERIAL] COMMAND:",
          data
        )


        callback(data)
      }
    )

  })


  /*
  |--------------------------------------------------------------------------
  | UNMOUNT
  |--------------------------------------------------------------------------
  */

  onUnmounted(() => {

    if (
      !wailsEnabled
    ) {
      return
    }


    console.log(
      "[SERIAL] Stop listening"
    )


    EventsOff(
      "serial-command"
    )

  })

}