<script setup lang="ts">
import type { QuestionAndAnswer } from '../client/types.gen';
import Answers from './Answers.vue';

defineProps<{ question: QuestionAndAnswer }>();
</script>

<template>
  <div class="grid gap-4 lg:gap-6 text-center">
    <div class="text-xl font-semibold">{{ question.question }}</div>

    <Answers :question="question" />
  </div>
</template>
