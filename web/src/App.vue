<script setup lang="ts">
import HeroCard from './components/HeroCard.vue';
import Navigation from './components/Navigation.vue';
import { useGlobalState } from './store';
import QuestionAndAnswers from './components/QuestionAndAnswers.vue';
import { UI } from './main';
import LoadingDots from './components/LoadingDots.vue';

const { quiz, question, loading, submitted } = useGlobalState();
</script>

<template>
  <HeroCard>
    <template v-slot:header>
      <div class="flex justify-center h-16" v-html="UI.data?.Logo"></div>
    </template>
    <div class="relative w-full h-full min-h-[200px]">
      <div v-if="submitted" class="grid gap-8 lg:gap-12">
        <QuestionAndAnswers v-for="question in quiz.questions" :question="question" :key="question.id" />
      </div>
      <Transition v-else name="crossfade" mode="out-in" appear>
        <QuestionAndAnswers v-if="!loading && question" :question="question" key="question" />
        <LoadingDots v-else />
      </Transition>
    </div>
    <template v-slot:footer>
      <Navigation />
    </template>
  </HeroCard>
</template>

<style scoped>
.crossfade-enter-active,
.crossfade-leave-active {
  position: absolute;
  left: 0;
  right: 0;
  width: 100%;
  transition: opacity 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  pointer-events: none;
}
.crossfade-enter-from,
.crossfade-leave-to {
  opacity: 0;
}
.crossfade-enter-to,
.crossfade-leave-from {
  opacity: 1;
}
</style>
