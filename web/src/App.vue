<script setup lang="ts">
import { ref } from 'vue';
import { getQuestions } from './client/sdk.gen';
import type { QuestionAndAnswer } from './client/types.gen';

const questions = ref<QuestionAndAnswer[] | null | undefined>(null);
getQuestions().then((resp) => {
  questions.value = resp.data;
});
</script>

<template>
  <div class="container">
    <div v-for="question in questions" :key="question.id">
      <h3>{{ question.question }}</h3>
      <p>{{ question.answers }}</p>
    </div>
  </div>
</template>
