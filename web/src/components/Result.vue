<script setup lang="ts">
import { useGlobalState } from '../store';
import Question from './Answers.vue';

const { quiz } = useGlobalState();
</script>

<template>
  <div class="grid gap-10">
    <div v-for="question in quiz.questions" :key="question.id" class="grid gap-2">
      <Question :question="question" />
    </div>
  </div>
</template>
