<script setup lang="ts">
import { useGlobalState } from '../store';
import type { QuestionAndAnswer } from '../client/types.gen';

const props = defineProps<{ question: QuestionAndAnswer }>();
const { submitted, handleAnswerSelected } = useGlobalState();

function getButtonClass(index: number) {
  const selected = props.question.answer === index + 1;
  const correct = props.question.correct;

  if (submitted.value && correct === index + 1) {
    return 'bg-success text-success-content ring-success font-bold border-2 border-success';
  }

  if (submitted.value && selected && correct !== index + 1) {
    return 'bg-error text-error-content ring-error border-2 border-error';
  }

  if (!submitted.value && selected) {
    return 'bg-secondary text-secondary-content ring-secondary';
  }

  if (submitted.value) {
    return 'bg-base-200/50 text-base-content opacity-60 cursor-not-allowed';
  }

  return 'bg-base-200/50 text-base-content hover:bg-base-200 active:bg-secondary active:text-secondary-content';
}
</script>

<template>
  <div class="grid gap-4 px-5">
    <button
      v-for="(answer, index) in question.answers"
      :key="`${question.id}-${index + 1}`"
      class="transition-all duration-200 rounded-lg px-5 py-3 w-full text-lg font-medium shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2"
      :class="getButtonClass(index)"
      :disabled="submitted"
      @click="handleAnswerSelected(question.id, index + 1)"
    >
      {{ answer }}
    </button>
  </div>
</template>
