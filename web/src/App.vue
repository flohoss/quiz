<script setup lang="ts">
import HeroCard from './components/HeroCard.vue';
import Navigation from './components/Navigation.vue';
import { useGlobalState } from './store';
import QuestionAndAnswers from './components/QuestionAndAnswers.vue';
import { BackendURL, AppSetting } from './main';
import QuizSteps from './components/QuizSteps.vue';

const { question, loading, direction } = useGlobalState();
</script>

<template>
  <HeroCard>
    <template v-slot:header>
      <div class="flex flex-col items-center w-full gap-5 py-3 sm:py-4">
        <div class="flex justify-center h-16 bg-auto bg-center bg-no-repeat w-full" :style="`background-image: url(${BackendURL + AppSetting.Logo})`"></div>
        <QuizSteps class="mt-2 w-full" />
      </div>
    </template>
    <div class="relative w-full h-full overflow-hidden p-3 sm:p-4">
      <Transition :name="direction === 'backward' ? 'swipe-right' : 'swipe-left'" mode="out-in" appear>
        <QuestionAndAnswers v-if="!loading && question" :question="question" :key="question.id" />
      </Transition>
    </div>
    <template v-slot:footer>
      <Navigation class="p-3 sm:p-4" />
    </template>
  </HeroCard>
</template>

<style scoped>
.swipe-left-enter-active,
.swipe-left-leave-active,
.swipe-right-enter-active,
.swipe-right-leave-active {
  position: absolute;
  left: 0;
  right: 0;
  width: 100%;
  transition: all 0.18s cubic-bezier(0.6, 0.05, 0.2, 1);
  pointer-events: none;
  will-change: transform, opacity;
}
.swipe-left-enter-from {
  opacity: 0;
  transform: translateX(40px);
}
.swipe-left-leave-to {
  opacity: 0;
  transform: translateX(-40px);
}
.swipe-left-enter-to,
.swipe-left-leave-from {
  opacity: 1;
  transform: translateX(0);
}
.swipe-right-enter-from {
  opacity: 0;
  transform: translateX(-40px);
}
.swipe-right-leave-to {
  opacity: 0;
  transform: translateX(40px);
}
.swipe-right-enter-to,
.swipe-right-leave-from {
  opacity: 1;
  transform: translateX(0);
}
</style>
