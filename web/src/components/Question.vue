<script setup lang="ts">
import { ref } from 'vue';
import { useGlobalState } from '../store';

const { question, handleAnswerSelected } = useGlobalState();
const disabled = ref(false);
</script>

<template>
  <div class="grid gap-5">
    <div class="text-xl font-semibold">{{ question.question }}</div>

    <div class="grid gap-2">
      <div class="flex items-center" v-for="(answer, index) in question.answers" :key="index">
        <input
          type="radio"
          :id="`${question.id}-${index + 1}`"
          :name="`question-${question.id}`"
          :value="index + 1"
          class="radio"
          @change="handleAnswerSelected(question.id, index + 1)"
        />
        <label :for="`${question.id}-${index + 1}`" class="pl-4 w-full select-none" :class="{ 'opacity-60': disabled }">
          {{ answer }}
        </label>
      </div>
    </div>
  </div>
</template>
