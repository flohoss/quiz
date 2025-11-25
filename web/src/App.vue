<script setup lang="ts">
import { ref } from 'vue';
import { getQuestions, validateAnswers } from './client/sdk.gen';
import type { QuestionAndAnswer, QuizAnswer, ValidationResult } from './client/types.gen';
import Question from './components/Question.vue';

const questions = ref<QuestionAndAnswer[] | null | undefined>(null);
const answers = ref<QuizAnswer[]>([]);
const validationResults = ref<ValidationResult[] | null>(null);
const isSubmitted = ref(false);
const isSubmitting = ref(false);

getQuestions().then((resp) => {
  questions.value = resp.data;
});

function handleAnswerSelected(question: number, answer: number) {
  if (isSubmitted.value) return; // Prevent changes after submission

  answers.value = answers.value.filter((a) => a.id !== question);
  answers.value.push({ id: question, answer: answer });
}

async function submitAnswers() {
  if (isSubmitting.value || !answers.value.length) return;

  isSubmitting.value = true;
  try {
    const response = await validateAnswers({
      body: answers.value,
    });
    validationResults.value = response.data || [];
    isSubmitted.value = true;
  } catch (error) {
    console.error('Error validating answers:', error);
    // You might want to show an error message to the user here
  } finally {
    isSubmitting.value = false;
  }
}

function getValidationResult(questionId: number): ValidationResult | null {
  return validationResults.value?.find((result) => result.id === questionId) || null;
}
</script>

<template>
  <div class="container py-10 grid gap-5">
    <div v-if="questions === null">Loading...</div>
    <div v-else-if="questions === undefined">Error loading questions.</div>
    <div v-else>
      <div class="grid gap-4">
        <div v-for="question in questions" :key="question.id">
          <Question
            :question="question"
            :validation-result="getValidationResult(question.id)"
            :disabled="isSubmitted"
            @answer-selected="handleAnswerSelected"
          />
        </div>

        <div class="flex justify-end">
          <button v-if="!isSubmitted" @click="submitAnswers" :disabled="isSubmitting || answers.length === 0" class="btn btn-primary">
            {{ isSubmitting ? 'Submitting...' : 'Submit Answers' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
