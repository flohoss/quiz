<script setup lang="ts">
import { useGlobalState } from '../store';
import ThemeChanger from './ThemeChanger.vue';
import { AppSetting } from '../main';

const { start, end, question, submitted, colorless, previousIndex, nextIndex, submit, reset } = useGlobalState();
</script>

<template>
  <div class="flex justify-between">
    <button @click="previousIndex" :disabled="start || submitted" class="nav-btn btn-secondary">
      <div class="size-8" v-html="AppSetting.icons['previous']"></div>
    </button>
    <div class="flex items-center gap-5">
      <ThemeChanger />
    </div>
    <button v-if="end && !submitted" :disabled="!question?.answer" @click="submit" class="nav-btn btn-primary">
      <div class="size-8" v-html="AppSetting.icons['submit']"></div>
    </button>
    <button v-else-if="!submitted" :disabled="!question?.answer" @click="nextIndex" class="nav-btn btn-primary">
      <div class="size-8" v-html="AppSetting.icons['next']"></div>
    </button>
    <button v-else @click="reset" :disabled="colorless" class="nav-btn btn-secondary">
      <div class="size-8" v-html="AppSetting.icons['restart']"></div>
    </button>
  </div>
</template>
