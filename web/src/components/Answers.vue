<script setup lang="ts">
import { useGlobalState } from '../store';
import type { QuestionAndAnswer } from '../client/types.gen';

defineProps<{ question: QuestionAndAnswer }>();
const { submitted, handleAnswerSelected } = useGlobalState();
</script>

<template>
  <div class="grid gap-4 px-5">
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
              ? 'bg-accent text-accent-content border-accent ring-accent'
              : // Disabled state (not selected, after submit)
                submitted
                ? 'bg-base-200/50 text-base-content border-base-300/50 opacity-60 cursor-not-allowed'
                : // Default
                  'text-base-content border-base-200/50 hover:bg-base-200 hover:border-base-300 active:bg-accent active:text-accent-content',
      ]"
      :disabled="submitted"
      @click="handleAnswerSelected(question.id, index + 1)"
    >
      {{ answer }}
    </button>
  </div>
</template>
