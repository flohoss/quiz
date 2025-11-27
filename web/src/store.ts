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
  const loading = shallowRef(true);
  const direction = shallowRef<'forward' | 'backward'>('forward');
  const colorless = shallowRef(false);

  async function loadQuiz() {
    loading.value = true;
    const resp = await getQuestions({ query: { language: lang } });
    if (resp.error || !resp.data) {
      loading.value = false;
      return;
    }
    quiz.value = resp.data;
    loading.value = false;
  }
  loadQuiz();

  function nextIndex() {
    if (!end.value) {
      direction.value = 'forward';
      index.value += 1;
    }
  }

  function previousIndex() {
    if (index.value > 1) {
      direction.value = 'backward';
      index.value -= 1;
    }
  }

  function handleAnswerSelected(id: number, answer: number) {
    quiz.value = {
      ...quiz.value,
      questions: quiz.value.questions.map((q) => (q.id === id && q.answer !== answer ? { ...q, answer } : q)),
    };
  }

  function setColorless(val: boolean) {
    colorless.value = val;
  }

  function submit() {
    setColorless(true);
    submitted.value = true;
    const answers = quiz.value.questions.filter((q) => typeof q.answer === 'number').map((q) => ({ id: q.id, answer: q.answer as number }));

    const minTime = new Promise((resolve) => setTimeout(resolve, 500));
    const validation = validateAnswers({ query: { language: lang }, body: answers });

    Promise.all([minTime, validation]).then(([, resp]) => {
      if (resp.error || !resp.data) {
        setColorless(false);
        return;
      }
      quiz.value = resp.data;
      setColorless(false);
    });
  }

  function reset() {
    index.value = 1;
    loadQuiz();
    submitted.value = false;
  }

  return {
    quiz,
    index,
    question,
    start,
    end,
    submitted,
    loading,
    direction,
    nextIndex,
    previousIndex,
    handleAnswerSelected,
    submit,
    reset,
    colorless,
    setColorless,
  };
});
