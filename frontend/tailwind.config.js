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
      gridTemplateRows: {
        // Complex site-specific row configuration
        'user-page': '20% 40% 40%',
        'admin-page': '20% 20% 60%',
      }
    },
  },
  plugins: [],
}
