# API documentation exposure

The Scalar API documentation route is registered only outside production. Production deployments must publish the OpenAPI document through an authenticated/internal documentation service or a separately protected static artifact; the public application does not expose `/docs/api`.

API operations use the `Authorization` header security scheme. The value is an `sk-dat-*` device token or `sk-dak-*` API key; credentials must not be placed in a URL. Every protected API operation is annotated with `@Security Authorization`, while `/api/v1/health` remains public.
