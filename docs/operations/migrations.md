# Migration operations

## Transaction policy

`MigrationRunner.Run` executes each migration and its `migrations` record in one database transaction. A failed migration is rolled back and is not recorded as applied. Migrations are applied in name order, and a later migration is not attempted after an earlier failure.

PostgreSQL runners acquire the session-level advisory lock `hashtext('sapasora:migrations')` for the complete run, including creation of the migration table. Other supported databases rely on their transaction and deployment serialization mechanisms.

Migration DDL must therefore be compatible with the target database's transactional DDL behavior. External side effects (files, provider APIs, queues) must not be performed from a migration.

## Rollback

1. Stop application instances that can run migrations concurrently.
2. Take or verify a database backup before destructive rollback.
3. Run the migration rollback command for the last batch:

   ```sh
   go run ./cmd/toolbox migrate rollback
   ```

4. Verify the schema and migration records before restarting the application.
5. If a rollback is not safe, restore the backup and deploy a forward-fix migration instead.

Rollback operations are not automatically wrapped by the runner because some database DDL and external migration effects cannot be made atomic. A failed rollback stops the process and must be recovered manually or with a forward-fix migration.
