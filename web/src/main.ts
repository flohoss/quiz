import { createApp } from 'vue';
import './style.css';
import App from './App.vue';
import { client } from './client/client.gen';
import { getApp } from './client/sdk.gen';

export const BackendURL = import.meta.env.MODE === 'development' ? 'http://localhost:8156' : '';

client.setConfig({ baseUrl: BackendURL });

export const Setting = await getApp();
if (Setting.data?.Title) {
  document.title = Setting.data.Title;
}

createApp(App).mount('#app');
