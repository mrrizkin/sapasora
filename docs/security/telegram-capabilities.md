# Telegram capability behavior

Telegram capabilities that are not in the current release scope return HTTP 501
with the stable error code prefix `telegram_capability_not_implemented:<capability>`.
They never panic, so an unsupported message type cannot crash an active request or
provider supervisor.

Currently implemented release paths include connection/status checks, user
lookup/checks, avatar lookup, and text messages. Audio, buttons, chat presence,
contacts, documents, images, lists, locations, stickers, video, and logout remain
explicitly unsupported until their TDLib behavior and contract tests are added.

A nil or still-connecting TDLib client is treated as disconnected for status and
safe disconnect operations. Usernames and sender identities tolerate missing or
empty `active_usernames` values.
