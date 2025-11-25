<script setup lang="ts">
import { useGlobalState } from '../store';

const state = useGlobalState();
</script>

<template>{{ state.result }}</template>
