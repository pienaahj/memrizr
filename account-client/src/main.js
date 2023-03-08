import { createApp } from 'vue';
import App from './App.vue';
import { createAuthStore } from '../src/store/auth';
import router from './routes.js';
import './style.css';
import './validators';

const authStore = createAuthStore({
    onAuthRoute: '/',
    requireAuthRoute: '/authenticate',
  });
// const app = createApp(App);
// app.use(router);
// app.mount('#app')
createApp(App).use(authStore).use(router).mount('#app');
