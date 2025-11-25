<script setup lang="ts">
import Question from './components/Question.vue';
import Navigation from './components/Navigation.vue';
import { useGlobalState } from './store';
import Result from './components/Result.vue';

const state = useGlobalState();
</script>

<template>
  <div class="container">
    <div v-if="state.result.value === null" class="hero min-h-screen">
      <div class="hero-content">
        <div class="card w-xl h-screen max-h-96 shadow-sm">
          <div class="card-body">
            <Question />
          </div>
          <div class="card-actions justify-end">
            <Navigation />
          </div>
        </div>
      </div>
    </div>
    <Result v-else />
  </div>
</template>
