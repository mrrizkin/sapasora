# CSRF protection

Web routes use a synchronizer-backed double-submit token:

- Fiber stores the CSRF token in the configured server session (`CSRF_SESSION=true`
  by default) and expires it according to `CSRF_EXPIRATION`. If session storage
  is disabled, Fiber still keeps the token in server-side middleware storage;
  the token is never trusted from the cookie alone.
- The CSRF cookie (`CSRF_COOKIE_NAME`) is only a readable same-origin token
  source. It is deliberately **not** `HttpOnly`, because browser JavaScript must
  copy its value into the configured header (`CSRF_KEY`).
- Unsafe web requests must provide both the cookie and the matching header. A
  cookie by itself is never accepted as the request token. The web pipeline does
  not use `KeyLookup: cookie:` or `CsrfFromCookie`.
- `CSRF_SAME_SITE`, `CSRF_SECURE`, cookie name, header name, expiration, and
  session mode are configured by the corresponding environment variables.
  `VITE_CSRF_KEY` and `VITE_CSRF_COOKIE_NAME` must match the server values.

Axios and Inertia add the CSRF header only for same-origin requests, using the
cookie value. They never add this token to cross-site requests. Pages must be
loaded from the application origin first so the initial safe response can set
the CSRF cookie.

The API and WebSocket channel pipelines do not include the web CSRF middleware.
Signed provider webhook endpoints must remain on a non-web pipeline and validate
the provider signature (including timestamp/replay protections) independently;
CSRF is not a replacement for webhook signature verification. Do not move a
webhook or signature-authenticated API route into the web pipeline merely to
reuse session middleware.
