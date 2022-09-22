/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{html,ts}",
  ],
  theme: {
    extend: {
      colors: {
        'background': '#262626',
        'dark-font': '#272727',
        'light-font': '#f4f4f4',
      },
    },
  },
  plugins: [],
}
