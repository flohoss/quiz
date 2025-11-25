<script setup lang="ts">
import { ref } from 'vue';
import { getQuestions } from './client/sdk.gen';
import type { QuestionAndAnswer } from './client/types.gen';
import Question from './components/Question.vue';

const questions = ref<QuestionAndAnswer[] | null | undefined>(null);
getQuestions().then((resp) => {
  questions.value = resp.data;
});
</script>

<template>
  <div class="container">
    <div v-if="questions === null">Loading...</div>
    <div v-else-if="questions === undefined">Error loading questions.</div>
    <div v-else>
      <div v-for="question in questions" :key="question.id" class="mb-4">
        <Question :question="question" />
      </div>
    </div>
  </div>
</template>
