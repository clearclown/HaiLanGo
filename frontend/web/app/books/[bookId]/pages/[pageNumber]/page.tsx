export default function BookPageDetail({
  params,
}: {
  params: { bookId: string; pageNumber: string }
}) {
  return (
    <div className="min-h-screen bg-background-secondary p-8">
      <h1 className="text-2xl font-bold">学習ページ</h1>
      <p className="mt-4 text-text-secondary">
        Book ID: {params.bookId}, Page: {params.pageNumber}
      </p>
      <p className="mt-2 text-text-secondary">Coming soon...</p>
    </div>
  )
}
