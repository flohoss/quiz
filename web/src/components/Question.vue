<script setup lang="ts">
import type { QuestionAndAnswer, ValidationResult } from '../client/types.gen';

const props = defineProps<{
  question: QuestionAndAnswer;
  validationResult?: ValidationResult | null;
  disabled?: boolean;
}>();

const emit = defineEmits<{
  (e: 'answerSelected', question: number, answer: number): void;
}>();

function handleAnswerChange(answerIndex: number) {
  if (props.disabled) return;
  emit('answerSelected', props.question.id, answerIndex);
}
</script>

<template>
  <div class="grid gap-2">
    <h3>{{ question.question }}</h3>

    <div class="grid gap-1">
      <div class="flex items-center" v-for="(answer, index) in question.answers" :key="index">
        <input
          type="radio"
          :id="`${question.id}-${index + 1}`"
          :name="`question-${question.id}`"
          :value="index + 1"
          class="radio"
          :disabled="disabled"
          @change="handleAnswerChange(index + 1)"
        />
        <label :for="`${question.id}-${index + 1}`" class="pl-4 w-full select-none" :class="{ 'opacity-60': disabled }">
          {{ answer }}
        </label>
      </div>
    </div>
  </div>
</template>
