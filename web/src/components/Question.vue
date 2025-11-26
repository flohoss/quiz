<script setup lang="ts">
import { useGlobalState } from '../store';
import type { QuestionAndAnswer } from '../client/types.gen';

defineProps<{ question: QuestionAndAnswer | undefined }>();
const { submitted, handleAnswerSelected } = useGlobalState();
</script>

<template>
  <div v-if="question" class="grid gap-6 w-2xl">
    <div class="text-2xl font-semibold">{{ question.question }}</div>

    <div class="grid gap-4">
      <div class="flex items-center" v-for="(answer, index) in question.answers" :key="`${question.id}-${index + 1}`">
        <input
          type="radio"
          :id="`${question.id}-${index + 1}`"
          :name="`question-${question.id}`"
          :value="index + 1"
          class="radio"
          :class="{
            'opacity-60': submitted,
            'radio-success opacity-100': submitted && question.correct && question.answer === index + 1,
            'radio-error opacity-100': submitted && !question.correct && question.answer === index + 1,
          }"
          @change="handleAnswerSelected(question.id, index + 1)"
          :disabled="submitted"
          :checked="question.answer === index + 1"
        />
        <label
          :for="`${question.id}-${index + 1}`"
          class="w-full pl-4 select-none text-lg"
          :class="{
            'opacity-60': submitted,
            'text-success opacity-100': submitted && question.correct && question.answer === index + 1,
            'text-error opacity-100': submitted && !question.correct && question.answer === index + 1,
          }"
        >
          {{ answer }}
        </label>
      </div>
    </div>
  </div>
</template>
