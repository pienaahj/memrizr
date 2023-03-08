// import { fileURLToPath } from 'url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import dns from 'dns'

dns.setDefaultResultOrder('verbatim')

// https://vitejs.dev/config/
/*export default defineConfig({
  base: '/account/',
  plugins: [vue()],
})
*/
export default defineConfig({
  base: '/account/',
  plugins: [vue(),],
  resolve: {
    alias: {
      '@': path.resolve(__dirname,'./src'),
    },
  },
  server:{
    base: '/account/',
    host: "0.0.0.0",
      fs: {
        strict: false
      }
    },

})
