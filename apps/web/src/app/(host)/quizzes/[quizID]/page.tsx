import { getQuiz } from '@/lib/api/quizzes'
import { listBanks } from '@/lib/api/banks'
import QuizBuilder from '@/components/QuizBuilder'

interface Props {
  params: Promise<{ quizID: string }>
}

export async function generateMetadata({ params }: Props) {
  const { quizID } = await params
  try {
    const quiz = await getQuiz(quizID)
    return { title: `${quiz.name} — Quibble` }
  } catch {
    return { title: 'Quiz Editor — Quibble' }
  }
}

export default async function QuizPage({ params }: Props) {
  const { quizID } = await params

  const [quiz, banks] = await Promise.all([
    getQuiz(quizID),
    listBanks().catch(() => []),
  ])

  return <QuizBuilder quiz={quiz} banks={banks} />
}
