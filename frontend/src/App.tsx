import { useState } from 'react';
import { UrlForm } from './components/UrlForm';
import { ResultDisplay } from './components/ResultDisplay';
import { shortenUrl } from './api';
import { Zap } from 'lucide-react';

function App() {
  const [shortCode, setShortCode] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleShorten = async (longUrl: string) => {
    setIsLoading(true);
    setError(null);
    setShortCode(null);

    try {
      const code = await shortenUrl(longUrl);
      setShortCode(code);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong!');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="layout-container">
      <header className="mb-12 text-center transform -rotate-2">
        <div className="inline-block bg-[var(--color-maroon)] text-[var(--color-nude)] border-[4px] border-black p-4 shadow-[6px_6px_0px_0px_black]">
          <h1 className="text-6xl flex items-center gap-4">
            <Zap className="w-12 h-12 fill-[var(--color-nude)]" />
            ZIP IT!
          </h1>
        </div>
        <p className="mt-4 font-heading text-xl bg-white inline-block px-4 py-1 border-[3px] border-black shadow-[3px_3px_0px_0px_black] transform rotate-1">
          MAKE LONG URLS SHORT & SNAPPY!
        </p>
      </header>

      <main className="w-full flex flex-col items-center">
        <UrlForm onSubmit={handleShorten} isLoading={isLoading} />

        {error && (
          <div className="mt-6 bg-red-100 border-[3px] border-black p-4 text-red-800 font-bold shadow-[4px_4px_0px_0px_black] transform rotate-1 animate-pulse">
            ERR: {error.toUpperCase()}
          </div>
        )}

        {shortCode && <ResultDisplay shortCode={shortCode} />}
      </main>

      <footer className="mt-auto pt-12 text-center font-heading text-gray-500 opacity-60">
        <p>&copy; {new Date().getFullYear()} ZIP IT INC.</p>
      </footer>
    </div>
  );
}

export default App;
