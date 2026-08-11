import {
    createApp,
} from "vue"

import App from "./App.vue"

import "./assets/main.css"

import {
    connectRealtime
} from './composables/useRealtime.ts'

createApp(App)
    .mount("#app")