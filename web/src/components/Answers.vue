<script setup lang="ts">
import { useGlobalState } from '../store';
import type { QuestionAndAnswer } from '../client/types.gen';

defineProps<{ question: QuestionAndAnswer }>();
const { submitted, handleAnswerSelected } = useGlobalState();
</script>

<template>
  <div class="grid gap-4">
    <button
      v-for="(answer, index) in question.answers"
      :key="`${question.id}-${index + 1}`"
      class="transition-all duration-200 rounded-lg px-5 py-3 w-full text-lg font-medium border shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2"
      :class="[
        // Success state
        question.answer === index + 1 && submitted && question.correct
          ? 'bg-success text-success-content border-success ring-success'
          : // Error state
            question.answer === index + 1 && submitted && !question.correct
            ? 'bg-error text-error-content border-error ring-error'
            : // Active state (selected, not yet submitted)
              question.answer === index + 1 && !submitted
              ? 'bg-primary text-primary-content border-primary ring-primary'
              : // Disabled state (not selected, after submit)
                submitted
                ? 'bg-base-300 text-base-content border-base-300 opacity-60 cursor-not-allowed'
                : // Default
                  'bg-base-200 text-base-content border-base-300 hover:bg-base-300 hover:border-primary active:bg-primary active:text-primary-content',
      ]"
      :disabled="submitted"
      @click="handleAnswerSelected(question.id, index + 1)"
    >
      {{ answer }}
    </button>
  </div>
</template>
