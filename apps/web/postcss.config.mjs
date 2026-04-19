// PostCSS is the build tool that processes Tailwind's utility classes into real CSS.
const config = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {}, // adds vendor prefixes (-webkit-, -moz-, etc.) for browser compat
  },
}

export default config
