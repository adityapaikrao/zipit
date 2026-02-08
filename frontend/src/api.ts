
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

interface ShortenResponse {
    short_code: string;
}

export const shortenUrl = async (longUrl: string): Promise<string> => {
    const response = await fetch(`${API_BASE_URL}/shorten`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ long_url: longUrl }),
    });

    if (!response.ok) {
        if (response.status === 400) {
            throw new Error('Invalid URL provided');
        }
        throw new Error('Failed to shorten URL');
    }

    const data: ShortenResponse = await response.json();
    return data.short_code;
};
