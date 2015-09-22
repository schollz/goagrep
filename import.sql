.echo ON

PRAGMA cache_size = 800000;
PRAGMA synchronous = OFF;
PRAGMA journal_mode = OFF;
PRAGMA locking_mode = EXCLUSIVE;
PRAGMA count_changes = OFF;
PRAGMA temp_store = MEMORY;
PRAGMA auto_vacuum = NONE;


BEGIN;
.read words.sql
COMMIT;

.exit
