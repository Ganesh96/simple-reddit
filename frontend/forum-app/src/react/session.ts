const STORE_KEY = 'simple_reddit_session';

export type Session = {
  username: string;
  authValue: string;
};

export function readSession(): Session | null {
  const raw = window.localStorage.getItem(STORE_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as Session;
  } catch {
    window.localStorage.removeItem(STORE_KEY);
    return null;
  }
}

export function saveSession(session: Session): void {
  window.localStorage.setItem(STORE_KEY, JSON.stringify(session));
}

export function clearSession(): void {
  window.localStorage.removeItem(STORE_KEY);
}
