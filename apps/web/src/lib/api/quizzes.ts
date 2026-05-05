// Server-side API client for quiz management.

import type { Quiz, QuizDetail } from '@/types'
import { apiFetch } from './client'

export async function listQuizzes(): Promise<Quiz[]> {
  const res = await apiFetch('/quizzes')
  return res.json()
}

export async function createQuiz(name: string, description?: string): Promise<Quiz> {
  const res = await apiFetch('/quizzes', {
    method: 'POST',
    body: JSON.stringify({ name, description }),
  })
  return res.json()
}

export async function getQuiz(id: string): Promise<QuizDetail> {
  const res = await apiFetch(`/quizzes/${id}`)
  return res.json()
}

export async function updateQuiz(id: string, name: string, description?: string): Promise<Quiz> {
  const res = await apiFetch(`/quizzes/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ name, description }),
  })
  return res.json()
}

export async function deleteQuiz(id: string): Promise<void> {
  await apiFetch(`/quizzes/${id}`, { method: 'DELETE' })
}

export async function createRound(quizID: string, title?: string): Promise<void> {
  await apiFetch(`/quizzes/${quizID}/rounds`, {
    method: 'POST',
    body: JSON.stringify({ title }),
  })
}

export async function updateRound(quizID: string, roundID: string, title?: string): Promise<void> {
  await apiFetch(`/quizzes/${quizID}/rounds/${roundID}`, {
    method: 'PUT',
    body: JSON.stringify({ title }),
  })
}

export async function deleteRound(quizID: string, roundID: string): Promise<void> {
  await apiFetch(`/quizzes/${quizID}/rounds/${roundID}`, { method: 'DELETE' })
}

export async function setRoundQuestions(quizID: string, roundID: string, questionIDs: string[]): Promise<void> {
  await apiFetch(`/quizzes/${quizID}/rounds/${roundID}/questions`, {
    method: 'PUT',
    body: JSON.stringify({ question_ids: questionIDs }),
  })
}
