/*
 * @Author: reel
 * @Date: 2023-06-23 10:54:25
 * @LastEditors: reel
 * @LastEditTime: 2023-09-09 07:38:11
 * @Description: 请填写简介
 */
import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {

  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    plugins: [
      vue(),
      AutoImport({
        resolvers: [ElementPlusResolver()],
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      }),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
        'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.js',
      },
      extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue']
    },
    server: {
      host: '0.0.0.0',
      port: env.VITE_APP_PORT,
      proxy: {
          "api": {
            target: env.VITE_APP_API_BASEURL,
            changeOrigin: true,
            ws: true,
            rewrite: (path) => path.replace(/^\/api/, ""),
          }
      }
    },
    build: {
      assetsInlineLimit: 4096000, // 4000kb  超过会以base64字符串显示
      outDir: "mscui", // 输出名称
    },
    base: '/mscui',
  }
})
