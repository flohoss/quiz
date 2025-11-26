<script setup lang="ts">
import Question from './components/Question.vue';
import Navigation from './components/Navigation.vue';
import { useGlobalState } from './store';
import Result from './components/Result.vue';

const { submitted, question, page } = useGlobalState();
</script>

<template>
  <div class="h-screen flex flex-col">
    <div class="sticky top-0 flex w-full justify-center p-2 md:p-5 bg-base-200/50">
      <div class="badge badge-xl badge-secondary rounded-full">{{ page }}</div>
    </div>

    <div class="flex-1 flex flex-col min-h-0 overflow-y-auto overflow-x-hidden p-5 md:p-10 lg:p-20">
      <div class="w-full my-auto flex flex-col items-center">
        <Transition name="slide-forward" mode="out-in" appear>
          <Question v-if="!submitted" :question="question" />
          <Result v-else />
        </Transition>
      </div>
    </div>
    <div class="sticky bottom-0 p-2 md:p-5 bg-base-200/50">
      <div class="flex items-center w-full justify-between gap-5 mx-auto max-w-200">
        <Navigation />
      </div>
    </div>
  </div>
</template>

<style>
/* Forward animation (left to right) */
.slide-forward-enter-active {
  transition: all 0.6s ease-out;
}

.slide-forward-leave-active {
  transition: all 0.6s ease-in;
}

.slide-forward-enter-from {
  opacity: 0;
  transform: translateX(50px);
}

.slide-forward-enter-to {
  opacity: 1;
  transform: translateX(0);
}

.slide-forward-leave-from {
  opacity: 1;
  transform: translateX(0);
}

.slide-forward-leave-to {
  opacity: 0;
  transform: translateX(-50px);
}

/* Backward animation (right to left) */
.slide-backward-enter-active {
  transition: all 0.6s ease-out;
}

.slide-backward-leave-active {
  transition: all 0.6s ease-in;
}

.slide-backward-enter-from {
  opacity: 0;
  transform: translateX(-50px);
}

.slide-backward-enter-to {
  opacity: 1;
  transform: translateX(0);
}

.slide-backward-leave-from {
  opacity: 1;
  transform: translateX(0);
}

.slide-backward-leave-to {
  opacity: 0;
  transform: translateX(50px);
}
</style>
