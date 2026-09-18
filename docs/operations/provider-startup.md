# Provider startup

Provider startup is explicitly opt-in per device.

- `m_devices.auto_connect` is a non-null boolean and defaults to `false`.
- `POST /api/v1/device` accepts `auto_connect`; it is `false` when omitted.
- `PUT /api/v1/device/{id}` accepts `auto_connect`. Omitting it preserves the existing value.
- WhatsApp and Telegram startup queries require `auto_connect = true`, active status, no soft delete, and an unexpired device.
- The `2026_09_19_100000_add_auto_connect_to_m_devices.go` migration deliberately defaults existing rows to `false`; upgrading cannot silently reconnect devices.

The startup worker isolates provider failures and uses bounded workers internally. Set `PROVIDER_STARTUP_CONCURRENCY` to a positive integer to configure the shared WhatsApp and Telegram startup worker bound; it defaults to `4`. Startup still waits for all eligible devices or cancellation, while one provider failure does not block the others.

The Telegram package requires TDLib development headers (`td/telegram/td_json_client.h`) to compile. Install TDLib development files in the build environment before running the full Go test/build suite.

Rollback follows the transaction and forward-fix guidance in [migration operations](./migrations.md). Dropping `auto_connect` is destructive to the opt-in setting, so verify a backup before rollback.
