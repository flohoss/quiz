<script setup lang="ts">
import HeroCard from './components/HeroCard.vue';
import Navigation from './components/Navigation.vue';
import { useGlobalState } from './store';
import QuestionAndAnswers from './components/QuestionAndAnswers.vue';
import { Setting, BackendURL } from './main';

const { quiz, question, loading, submitted, direction } = useGlobalState();
</script>

<template>
  <HeroCard>
    <template v-slot:header>
      <div class="flex justify-center h-16 bg-auto bg-center bg-no-repeat" :style="`background-image: url(${BackendURL + Setting.data?.Logo})`"></div>
    </template>
    <div class="relative w-full h-full min-h-[200px] overflow-x-hidden">
      <div v-if="submitted" class="grid gap-8 lg:gap-12">
        <QuestionAndAnswers v-for="question in quiz.questions" :question="question" :key="question.id" />
      </div>
      <Transition v-else :name="direction === 'backward' ? 'swipe-right' : 'swipe-left'" mode="out-in" appear>
        <QuestionAndAnswers v-if="!loading && question" :question="question" :key="question.id" />
      </Transition>
    </div>
    <template v-slot:footer>
      <Navigation />
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
