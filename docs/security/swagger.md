# API documentation exposure

The Scalar API documentation route is registered only outside production. Production deployments must publish the OpenAPI document through an authenticated/internal documentation service or a separately protected static artifact; the public application does not expose `/docs/api`.
