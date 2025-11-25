import { createApp } from 'vue';
import './style.css';
import App from './App.vue';
import { client } from './client/client.gen';

export const BackendURL = import.meta.env.MODE === 'development' ? 'http://localhost:8156' : '';

client.setConfig({ baseUrl: BackendURL });

createApp(App).mount('#app');
