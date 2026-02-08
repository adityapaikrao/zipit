import React, { useState } from 'react';
import { Scissors } from 'lucide-react';

interface UrlFormProps {
    onSubmit: (url: string) => Promise<void>;
    isLoading: boolean;
}

export const UrlForm: React.FC<UrlFormProps> = ({ onSubmit, isLoading }) => {
    const [url, setUrl] = useState('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (url.trim()) {
            let submitUrl = url.trim();
            // Auto-prepend https:// if protocol is missing
            if (!/^https?:\/\//i.test(submitUrl)) {
                submitUrl = `https://${submitUrl}`;
            }
            onSubmit(submitUrl);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="w-full max-w-lg flex flex-col gap-4">
            <div className="comic-box p-4 bg-white transform -rotate-1">
                <label htmlFor="url-input" className="block font-heading text-2xl mb-2 text-maroon">
                    PASTE YOUR LONG URL HERE:
                </label>
                <input
                    id="url-input"
                    type="text"
                    required
                    placeholder="https://super-long-boring-url.com/..."
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    className="comic-input"
                    disabled={isLoading}
                />
            </div>

            <button
                type="submit"
                disabled={isLoading || !url.trim()}
                className="comic-button self-center flex items-center gap-2 transform rotate-1 hover:rotate-0"
            >
                {isLoading ? (
                    'ZIPPING...'
                ) : (
                    <>
                        <Scissors className="w-6 h-6" />
                        ZIP IT!
                    </>
                )}
            </button>
        </form>
    );
};
