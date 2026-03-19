/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        bg:      "#0f0f0f",
        card:    "#1a1a1a",
        text:    "#e5e5e5",
        accent:  "#6366f1",
        danger:  "#ef4444",
        success: "#22c55e",
      },
      fontFamily: {
        sans: ["Noto Sans JP", "sans-serif"],
      },
    },
  },
  plugins: [],
}

