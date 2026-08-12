import {
  defineConfig,
} from "vite"

import vue from
  "@vitejs/plugin-vue"

import tailwindcss from
  "@tailwindcss/vite"


export default defineConfig({

  // ===================================================
  // PLUGINS
  // ===================================================

  plugins: [

    vue(),

    tailwindcss(),

  ],


  // ===================================================
  // DEV SERVER
  // ===================================================

  server: {

    // Bisa diakses device lain
    // dalam LAN.
    host:
      "0.0.0.0",

    // Port tetap supaya URL HP
    // tidak berubah.
    port:
      5173,

    // Jangan pindah ke 5174 kalau
    // port 5173 sedang terpakai.
    strictPort:
      true,

    // Izinkan hostname mDNS.
    allowedHosts: [

      "amen-mixer.local",

    ],

  },

})