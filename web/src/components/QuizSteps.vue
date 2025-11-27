<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue';
import { useGlobalState } from '../store';

const { quiz, index, submitted, direction } = useGlobalState();
const stepsContainer = ref<HTMLElement | null>(null);

const steps = computed(() =>
  quiz.value.questions.map((q, i) => ({
    label: `#${i + 1}`,
    id: q.id,
    answered: typeof q.answer === 'number',
  }))
);

function goToStep(i: number) {
  if (!submitted.value) {
    return;
  }

  if (i + 1 > index.value) {
    direction.value = 'forward';
  } else if (i + 1 < index.value) {
    direction.value = 'backward';
  }
  index.value = i + 1;
}

function scrollToActive() {
  if (!stepsContainer.value) return;
  // Always scroll to the step representing the current index
  const stepElems = stepsContainer.value.querySelectorAll('.step');
  const currentStepElem = stepElems[index.value - 1] as HTMLElement | undefined;
  if (currentStepElem) {
    const containerRect = stepsContainer.value.getBoundingClientRect();
    const stepRect = currentStepElem.getBoundingClientRect();
    const offset = stepRect.left - containerRect.left - containerRect.width / 2 + stepRect.width / 2;
    stepsContainer.value.scrollBy({ left: offset, behavior: 'smooth' });
  }
}

onMounted(scrollToActive);
watch(index, scrollToActive);

function getStepClass(i: number) {
  const question = quiz.value.questions[i];
  const step = steps.value[i];

  if (!question || !step) {
    return '';
  }

  const answered = step.answered;
  const isCurrent = index.value === i + 1;

  if (!submitted.value) {
    return isCurrent ? 'step-primary' : '';
  }

  if (!answered) {
    return 'step-neutral opacity-60';
  }

  return question.answer === question.correct ? 'step-success' : 'step-error';
}
</script>

<template>
  <div ref="stepsContainer" class="steps w-full overflow-x-auto whitespace-nowrap scrollbar-hide" :class="{ 'pointer-events-none select-none': !submitted }">
    <div v-for="(step, i) in steps" :key="step.id" class="step px-2 cursor-pointer" :class="getStepClass(i)" @click="goToStep(i)"></div>
  </div>
</template>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
