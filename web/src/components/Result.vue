<script setup lang="ts">
import { useGlobalState } from '../store';
import type { QuestionAndAnswer } from '../client/types.gen';

const { quiz } = useGlobalState();
function getAnswer(question: QuestionAndAnswer) {
  if (!question.answer) {
    return 'No answer selected';
  }
  return question.answers[question.answer - 1];
}
</script>

<template>
  <div class="grid gap-6">
    <div v-for="question in quiz.questions" :key="question.id" class="grid gap-2">
      <div class="text-2xl font-semibold">{{ question.question }}</div>
      <div></div>
      <span
        :class="{
          'text-success': question.correct,
          'text-error': !question.correct,
        }"
      >
        {{ getAnswer(question) }}
      </span>
    </div>
  </div>
</template>
