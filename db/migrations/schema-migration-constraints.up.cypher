// Constraints. Run manually instead with migration.

CREATE CONSTRAINT FOR (a:SchemaMigration) REQUIRE a.version IS UNIQUE

CREATE CONSTRAINT user_email_constraint IF NOT EXISTS FOR (u:User) REQUIRE u.email IS UNIQUE
CREATE CONSTRAINT user_username_constraint IF NOT EXISTS FOR (u:User) REQUIRE u.username IS UNIQUE


// Show all constraints: SHOW CONSTRAINTS