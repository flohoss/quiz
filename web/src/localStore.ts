import { createGlobalState, usePreferredColorScheme, usePreferredLanguages, useStorage } from '@vueuse/core';

export const STORAGE_KEY = 'quiz';
export const LANGUAGES = ['de', 'en'];

const preferredLanguages = usePreferredLanguages();
const preferredLanguage = preferredLanguages.value[0] && preferredLanguages.value[0].split('-')[0];
const defaultLanguage = preferredLanguage && LANGUAGES.includes(preferredLanguage) ? preferredLanguage : 'de';

const preferredColorScheme = usePreferredColorScheme();
const darkMode = preferredColorScheme.value === 'dark';

const storage = useStorage<{ index: number; language: string | undefined; darkMode: boolean }>(
  STORAGE_KEY,
  {
    index: 0,
    language: defaultLanguage,
    darkMode: darkMode,
  },
  localStorage,
  { mergeDefaults: true }
);

if (typeof storage.value.index !== 'number') {
  storage.value.index = 0;
}

if (storage.value.language !== undefined && !LANGUAGES.includes(storage.value.language)) {
  storage.value.language = defaultLanguage;
}

if (typeof storage.value.darkMode !== 'boolean') {
  storage.value.darkMode = darkMode;
}

export const useLocalStore = createGlobalState(() => {
  return storage;
});
