import { createGlobalState, useNavigatorLanguage } from '@vueuse/core';
import { computed, shallowRef } from 'vue';
import type { Quiz } from './client/types.gen';
import { getQuestions, validateAnswers } from './client/sdk.gen';

export const emptyQuiz: Quiz = {
  questions: [],
  total: 0,
};
const { isSupported, language } = useNavigatorLanguage();

export const useGlobalState = createGlobalState(() => {
  let lang = 'en';
  if (isSupported.value && language.value) {
    lang = language.value.split('-')[0] ?? 'en';
  }

  const quiz = shallowRef<Quiz>(emptyQuiz);
  const index = shallowRef<number>(1);
  const question = computed(() => quiz.value.questions[index.value - 1]);
  const start = computed(() => index.value === 1);
  const end = computed(() => index.value === quiz.value.total);
  const submitted = shallowRef(false);

  function loadQuiz() {
    getQuestions({ query: { language: lang } }).then((resp) => {
      if (resp.error || !resp.data) {
        return;
      }
      quiz.value = resp.data;
    });
  }
  loadQuiz();

  function nextIndex() {
    if (!end.value) {
      index.value += 1;
    }
  }

  function previousIndex() {
    if (index.value > 1) {
      index.value -= 1;
    }
  }

  function handleAnswerSelected(id: number, answer: number) {
    quiz.value = {
      ...quiz.value,
      questions: quiz.value.questions.map((q) => (q.id === id && q.answer !== answer ? { ...q, answer } : q)),
    };
  }

  function submit() {
    const answers = quiz.value.questions.filter((q) => typeof q.answer === 'number').map((q) => ({ id: q.id, answer: q.answer as number }));

    validateAnswers({ query: { language: lang }, body: answers }).then((resp) => {
      if (resp.error || !resp.data) {
        return;
      }
      quiz.value = resp.data;
      submitted.value = true;
    });
  }

  return { quiz, index, question, start, end, submitted, nextIndex, previousIndex, handleAnswerSelected, submit };
});
