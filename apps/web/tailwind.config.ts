import type { Config } from 'tailwindcss'

const config: Config = {
  // Tell Tailwind which files to scan for class names.
  // It removes any classes not found here from the production build (tree-shaking).
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          blue: '#00338D',
          red: '#C60C30',
        },
      },
    },
  },
  plugins: [],
}

export default config
