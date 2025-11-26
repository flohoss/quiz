import { createApp } from 'vue';
import './style.css';
import App from './App.vue';
import { client } from './client/client.gen';
import { getUi } from './client/sdk.gen';

export const BackendURL = import.meta.env.MODE === 'development' ? 'http://localhost:8156' : '';

client.setConfig({ baseUrl: BackendURL });

export const UI = await getUi();
if (UI.data?.Title) {
  document.title = UI.data.Title;
}

createApp(App).mount('#app');
