<script setup lang="ts">
import type { QuestionAndAnswer } from '../client/types.gen';

defineProps<{ question: QuestionAndAnswer }>();
const emit = defineEmits<{
  (e: 'answerSelected', question: number, answer: number): void;
}>();
</script>

<template>
  <div class="grid gap-5">
    <h3>{{ question.question }}</h3>

    <div>
      <div class="flex items-center" v-for="(answer, id) in question.answers" :key="answer">
        <input type="radio" :id="answer" :name="question.question" :value="answer" class="radio" @change="emit('answerSelected', question.id, id + 1)" />
        <label :for="answer" class="pl-4 w-full select-none">
          {{ answer }}
        </label>
      </div>
    </div>
  </div>
</template>
