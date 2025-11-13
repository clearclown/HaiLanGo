export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">HaiLanGo</h1>
        <p className="text-xl text-gray-600 mb-8">
          AI-Powered Language Learning Platform
        </p>
        <div className="space-x-4">
          <a
            href="/books"
            className="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
          >
            マイ本
          </a>
          <a
            href="/learning"
            className="px-6 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600"
          >
            学習を始める
          </a>
        </div>
      </div>
    </main>
  );
}
