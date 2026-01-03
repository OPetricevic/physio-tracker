import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Dev proxy so the frontend can call the Go API without CORS issues.
      '/api': {
        target: 'http://localhost:3600',
        changeOrigin: true,
      },
    },
  },
})
