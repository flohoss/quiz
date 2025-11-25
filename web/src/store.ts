import { createGlobalState, useNavigatorLanguage } from '@vueuse/core';
import { computed, shallowRef } from 'vue';
import type { QuestionAndAnswer, QuizAnswer } from './client/types.gen';
import { getQuestions, validateAnswers } from './client/sdk.gen';

export const emptyQuestion: QuestionAndAnswer = { id: 0, question: '', answers: [] };
const { isSupported, language } = useNavigatorLanguage();

export const useGlobalState = createGlobalState(() => {
  const questions = shallowRef<QuestionAndAnswer[]>([emptyQuestion]);
  const index = shallowRef<number>(0);
  const amount = computed<number>(() => questions.value.length);
  const question = computed<QuestionAndAnswer>(() => questions.value[index.value] ?? emptyQuestion);
  const answers = shallowRef<QuizAnswer[]>([]);
  const selected = computed(() => answers.value.find((a: QuizAnswer) => a.id === question.value.id)?.answer);

  function loadQuestions() {
    let lang: 'de' | 'en' = 'de';
    if (isSupported.value && language.value) {
      lang = language.value.startsWith('en') ? 'en' : 'de';
    }
    getQuestions({ query: { lang: lang } }).then((resp) => {
      if (resp.error || !resp.data) {
        return;
      }
      questions.value = resp.data;
    });
  }
  loadQuestions();

  function nextIndex() {
    if (index.value < questions.value.length - 1) {
      index.value += 1;
    }
  }

  function previousIndex() {
    if (index.value > 0) {
      index.value -= 1;
    }
  }

  function handleAnswerSelected(question: number, answer: number) {
    answers.value = answers.value.filter((a) => a.id !== question);
    answers.value.push({ id: question, answer: answer });
  }

  async function submit() {
    const response = await validateAnswers({ body: answers.value });
    if (response.error) {
      console.error('Error validating answers:', response.error);
      return;
    }
    console.log('Validation result:', response.data);
  }

  return { questions, index, amount, question, answers, selected, loadQuestions, nextIndex, previousIndex, handleAnswerSelected, submit };
});
