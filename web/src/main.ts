import { createApp } from 'vue';
import './style.css';
import App from './App.vue';
import { client } from './client/client.gen';
import { getApp } from './client/sdk.gen';
import { useNavigatorLanguage, useDark } from '@vueuse/core';

export const IsDark = useDark({
  selector: 'html',
  attribute: 'data-theme',
  valueDark: 'business',
  valueLight: 'corporate',
  storageKey: 'quiz-theme',
});

export const BackendURL = import.meta.env.MODE === 'development' ? 'http://localhost:8156' : '';

client.setConfig({ baseUrl: BackendURL });

export const Setting = await getApp();
if (Setting.error || !Setting.data) {
  throw new Error('Failed to load app settings');
}

document.title = Setting.data.Title;

if (Setting.data.CSSVariables) {
  const root = document.documentElement;
  Object.entries(Setting.data.CSSVariables).forEach(([key, value]) => {
    root.style.setProperty(key, value);
  });
}

const { language } = useNavigatorLanguage();
let lang = '';
if (language.value) {
  lang = language.value.split('-')[0] ?? '';
}
if (Setting.data.Languages.includes(lang)) {
  document.documentElement.setAttribute('lang', lang);
} else if (Array.isArray(Setting.data.Languages) && Setting.data.Languages.length > 0 && typeof Setting.data.Languages[0] === 'string') {
  document.documentElement.setAttribute('lang', Setting.data.Languages[0]);
}

createApp(App).mount('#app');
