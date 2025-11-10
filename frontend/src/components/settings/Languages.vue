<template>
  <select name="selectLanguage" @change="onChange" :value="locale">
    <option v-for="(language, value) in locales" :key="value" :value="value">
      {{ language }}
    </option>
  </select>
</template>

<script setup lang="ts">
defineProps<{
  locale: string;
}>();

const emit = defineEmits<{
  "update:locale": [value: string];
}>();

const locales = {
  he: "עברית",
  hr: "Hrvatski",
  hu: "Magyar",
  ar: "العربية",
  ca: "Català",
  cs: "Čeština",
  de: "Deutsch",
  el: "Ελληνικά",
  en: "English",
  es: "Español",
  fr: "Français",
  is: "Icelandic",
  it: "Italiano",
  ja: "日本語",
  ko: "한국어",
  "nl-be": "Dutch (Belgium)",
  no: "Norsk",
  pl: "Polski",
  "pt-br": "Português",
  pt: "Português (Brasil)",
  ro: "Romanian",
  ru: "Русский",
  sk: "Slovenčina",
  "sv-se": "Swedish (Sweden)",
  tr: "Türkçe",
  uk: "Українська",
  vi: "Tiếng Việt",
  "zh-cn": "中文 (简体)",
  "zh-tw": "中文 (繁體)",
};

function onChange(event: Event) {
  const target = event.target as HTMLSelectElement;
  emit("update:locale", target.value);
}
</script>
