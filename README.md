# Hadrian

THIS IS A WORK IN PROGRESS & NOT YET PRODUCTION READY

Hadrian is a CLI that replicates postgres database changes for application
consumption.

It does this by connecting to a target posgres server as a logical streaming
replication client, and outputs all database changes in a format applications
can understand, such as json.

Hadrian is intended to solve the dual write problem, described in Kleppmann's
Bottled Water blog post, which frequently occur when multiple systems need to
be notified of data changes in a given database.

Debezium is currently the industry standard solution for this problem, however
it requires both JVM and/or kafka expertise to use effectively.

Hadrian is intended to provide a subset of debezium's features, so that small
teams or single person projects without high availability or high throughput
requirements can replicate data from their postgres databases to downstream
systems without needing to change application code.

Credits:
Hadrian is a wrapper around @jackc's fantastic pglogrepl and was
inspired by both cainophile: and supabase's realtime.

References:
Postgres Streaming Replication:
https://www.postgresql.org/docs/current/protocol-replication.html.

Debezium:
https://debezium.io/

Bottled Water Blog Post:
https://www.confluent.io/blog/bottled-water-real-time-integration-of-postgresql-and-kafka/.

@jackc's pglogrepl:
https://github.com/jackc/pglogrepl

Cainophile:
https://github.com/cainophile/cainophile

Supabase Realtime:
https://github.com/supabase/realtime.
