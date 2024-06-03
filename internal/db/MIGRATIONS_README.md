Short reminder
1.migrate create -ext sql -dir ./internal/db/migrations -seq name_of_your_migration_file
2.Edit the Migration Files
3.migrate -path ./internal/db/migrations -database sqlite3://./internal/db/database.db up



Adding a New Migration
If you need to add a new field to the posts table in your SQLite database, follow these steps:

1. Create a New Migration
Use the following command to generate a new migration:

migrate create -ext sql -dir ./internal/db/migrations -seq add_field_to_posts_table
This creates two files: 202202201600_add_field_to_posts_table.up.sql and 202202201600_add_field_to_posts_table.down.sql.

2. Edit the Migration Files
In 202202201600_add_field_to_posts_table.up.sql, add the SQL query to add the new field. For example:


ALTER TABLE posts
ADD COLUMN new_field TEXT;

Replace new_field with the actual name and data type of your new field.

In 202202201600_add_field_to_posts_table.down.sql, add the SQL query to revert the change. For example:


ALTER TABLE posts
DROP COLUMN new_field;я
3. Apply the Migration
Run the following command to apply the new migration:


migrate -path ./internal/db/migrations -database sqlite3://./internal/db/database.db up
Now, your posts table includes the newly added field. Ensure that you update your application's code to utilize this new field as needed.

