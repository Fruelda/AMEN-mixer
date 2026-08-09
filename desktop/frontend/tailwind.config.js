/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        ui: ['Inter', 'Segoe UI', 'Arial', 'sans-serif'],
      },
      boxShadow: {
        panel: '0 0 0 1px rgba(88, 137, 166, .22), inset 0 1px 0 rgba(255,255,255,.025), 0 12px 36px rgba(0,0,0,.38)',
      },
    },
  },
  plugins: [],
}
