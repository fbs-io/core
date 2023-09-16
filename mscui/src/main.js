/*
 * @Author: reel
 * @Date: 2023-06-04 15:27:47
 * @LastEditors: reel
 * @LastEditTime: 2023-06-11 08:12:09
 * @Description: 请填写简介
 */
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import base from './base'
import i18n from './locales'
import store from './stores'

const app = createApp(App)

app.use(createPinia())
app.use(store)
app.use(router)
app.use(ElementPlus)
app.use(base)
app.use(i18n)
app.mount('#app')
