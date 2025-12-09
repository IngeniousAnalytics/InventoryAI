import Link from "next/link";

export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-indigo-500 to-purple-600 text-white">
      <main className="text-center space-y-8">
        <h1 className="text-6xl font-bold tracking-tight">
          Inventory<span className="text-black/30">AI</span>
        </h1>
        <p className="text-xl max-w-lg mx-auto text-indigo-100">
          The open-source, AI-powered inventory management platform for everyone.
        </p>

        <div className="flex gap-4 justify-center">
          <Link
            href="/auth/login"
            className="px-6 py-3 bg-white text-indigo-600 font-semibold rounded-lg shadow hover:bg-gray-100 transition"
          >
            Sign In
          </Link>
          <Link
            href="/auth/register"
            className="px-6 py-3 bg-indigo-800 text-white font-semibold rounded-lg shadow hover:bg-indigo-900 transition"
          >
            Get Started
          </Link>
        </div>
      </main>
    </div>
  );
}
