# Windows Mixer UI — Vue + Vite + Tailwind

Frontend mixer yang meniru tampilan referensi: panel gelap, daftar aplikasi, segmented volume bars, dan persentase volume yang interaktif.

## Jalankan

```bash
npm install
npm run dev
```

## Build production

```bash
npm run build
npm run preview
```

## Struktur utama

```text
src/
├── assets/
│   ├── fonts/
│   ├── images/
│   └── main.css
├── components/
│   ├── common/
│   ├── layout/
│   └── mixer/
│       ├── AppIcon.vue
│       ├── ChannelRow.vue
│       ├── KnobRing.vue
│       ├── MixerBoard.vue
│       ├── MixerPanel.vue
│       ├── MuteButton.vue
│       ├── SimulatorPanel.vue
│       ├── VolumeMeter.vue
│       └── VolumeSlider.vue
├── composables/
│   └── useSerial.ts
├── stores/
│   ├── mixer.ts
│   ├── serial.ts
│   └── settings.ts
├── types/
│   └── mixer.ts
├── App.vue
└── main.js
```

`useSerial.ts` sudah disediakan sebagai titik awal untuk integrasi Web Serial ke hardware mixer.
