<script setup lang="ts">
import type { QuestionAndAnswer } from '../client/types.gen';

defineProps<{ question: QuestionAndAnswer }>();
</script>

<template>
  <div class="grid gap-5">
    <h3>{{ question.question }}</h3>

    <div>
      <div class="flex items-center" v-for="answer in question.answers" :key="answer">
        <input type="radio" :id="answer" :value="answer" class="radio" />
        <label :for="answer" class="pl-4 w-full select-none">
          {{ answer }}
        </label>
      </div>
    </div>
  </div>
</template>
