import type { Quiz } from '@/types'
import { listQuizzes } from '@/lib/api/quizzes'
import QuizzesView from '@/components/QuizzesView'

export const metadata = { title: 'Quizzes — Quibble' }

export default async function QuizzesPage() {
  let quizzes: Quiz[] = []
  try {
    quizzes = await listQuizzes()
  } catch {
    // API might not be running during SSR; show empty state
  }

  return <QuizzesView quizzes={quizzes} />
}
