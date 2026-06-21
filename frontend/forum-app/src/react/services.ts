const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

export async function ping(): Promise<Response> {
  return fetch(BASE_URL);
}
