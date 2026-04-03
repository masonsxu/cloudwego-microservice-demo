import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { TooltipProvider } from '@/components/ui/tooltip'

import App from './App.vue'
import router from './router'
import createI18n from './locales'
import './assets/styles/variables.css'
import './assets/styles/tailwind.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(createI18n())

app.component('TooltipProvider', TooltipProvider)

app.mount('#app')
