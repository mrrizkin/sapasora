const defaultCsrfHeaderName = 'X-CSRF-Token';
const defaultCsrfCookieName = 'fiber_csrf_token';

export const csrfHeaderName = import.meta.env.VITE_CSRF_KEY || defaultCsrfHeaderName;
export const csrfCookieName = import.meta.env.VITE_CSRF_COOKIE_NAME || defaultCsrfCookieName;

export function getCsrfToken(): string | undefined {
  if (typeof document === 'undefined') {
    return undefined;
  }

  const cookie = document.cookie.split('; ').find((entry) => entry.startsWith(`${csrfCookieName}=`));

  if (!cookie) {
    return undefined;
  }

  const value = cookie.slice(cookie.indexOf('=') + 1);
  try {
    return decodeURIComponent(value);
  } catch {
    return undefined;
  }
}

export function isSameOrigin(url?: string): boolean {
  if (typeof window === 'undefined') {
    return false;
  }

  return new URL(url || window.location.href, window.location.href).origin === window.location.origin;
}

export function csrfHeaders(url?: string): Record<string, string> {
  const token = isSameOrigin(url) ? getCsrfToken() : undefined;

  return token ? { [csrfHeaderName]: token } : {};
}
