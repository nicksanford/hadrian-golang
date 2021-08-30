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

## Install:
```
go install github.com/nicksanford/hadrian

hadrian --help
```

## Local Development & Testing:

The following will get you set up with a postgres container with the correct
wal_level, set up a postgres publication and replecation slot, and start
replecation. It will then run the postgres built in benchmarking utility
`pgbench` which is a convenient way of creating some db activity for testing
purposes.

#### shell 1: database
*THIS IS ONLY OK FOR TEST PURPOSES* Change your passwords if you are using this for anything outside of testing.
```bash
docker run -e POSTGRES_PASSWORD=password -p 5432:5432 postgres:13.4 -c wal_level=logical
```


####  shell 2: replication
```bash
go run main.go create publication hadrian  'postgres://postgres:password@localhost/postgres?replication=database'
go run main.go create slot hadrian  'postgres://postgres:password@localhost/postgres?replication=database'
go run main.go replicate 'postgres://postgres:password@localhost/postgres?replication=database' -s hadrian -p hadrian > output
```
#### shell 3: load test
```bash
docker exec -it postgres_hadrian_test pgbench -h localhost -p 5432 -i -U postgres
```

#### Example current output (from running pgbench):
```jsonlines
{"FinalLSN":59833168,"CommitTime":"2021-08-29T20:42:51.64758-04:00","Xid":516}
{"RelationID":16435,"Namespace":"public","RelationName":"pgbench_accounts","ReplicaIdentity":100,"ColumnNum":4,"Columns":[{"Flags":0,"Name":"aid","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"bid","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"abalance","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"filler","DataType":1042,"TypeModifier":88}]}
{"RelationID":16438,"Namespace":"public","RelationName":"pgbench_branches","ReplicaIdentity":100,"ColumnNum":3,"Columns":[{"Flags":0,"Name":"bid","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"bbalance","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"filler","DataType":1042,"TypeModifier":92}]}
{"RelationID":16429,"Namespace":"public","RelationName":"pgbench_history","ReplicaIdentity":100,"ColumnNum":6,"Columns":[{"Flags":0,"Name":"tid","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"bid","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"aid","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"delta","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"mtime","DataType":1114,"TypeModifier":4294967295},{"Flags":0,"Name":"filler","DataType":1042,"TypeModifier":26}]}
{"RelationID":16432,"Namespace":"public","RelationName":"pgbench_tellers","ReplicaIdentity":100,"ColumnNum":4,"Columns":[{"Flags":0,"Name":"tid","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"bid","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"tbalance","DataType":23,"TypeModifier":4294967295},{"Flags":0,"Name":"filler","DataType":1042,"TypeModifier":88}]}
{"RelationNum":4,"Option":0,"RelationIDs":[16435,16438,16429,16432]}
{"RelationID":16438,"Tuple":{"ColumnNum":3,"Columns":[{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"Mg=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"Mw=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"NA=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"NQ=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"Ng=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"Nw=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"OA=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"OQ=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16432,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":2,"Data":"MTA="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":110,"Length":0,"Data":null}]}}
{"RelationID":16435,"Tuple":{"ColumnNum":4,"Columns":[{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MQ=="},{"DataType":116,"Length":1,"Data":"MA=="},{"DataType":116,"Length":84,"Data":"ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg"}]}}
{"Flags":0,"CommitLSN":59833168,"TransactionEndLSN":59833472,"CommitTime":"2021-08-29T20:42:51.64758-04:00"}
```

## TODO (MVP in roughly priority order):

- [ ] Currently streaming replication messages are written as they come in, however I believe that means that it is possible for uncommitted data to be written to  the output. To handle this I think I'm going to need to keep track of what data has & has not yet been comitted & only write & ack messages which have been committed.
- [ ] There are no tests currently
- [ ] Currently we are serializing the `pglogrepl` structs to json which means that the content is always base64 encoded (which is not necessary for many datatypes) and also the datatypes are described using numbers, where as describing them using strings would be much more descriptive
- [ ] Dropping replecation slots is unimplemented
- [ ] Creating temporary replecation slots is currently not supported
- [ ] Replecating using temporary replication slots is currently not supported
- [ ] Allow the client to be configured via env var (for secure values such as the postgres url) and also through yaml file for the rest of the values

## TODO Post MVP:

- [ ] All business logic is in the `cmd/` directory, I need to split business logic out of the cli modules & make it so that cmd only calls into that business logc
- [ ] Provide more output adapters than just stdout, consider also SNS to start.
- [ ] Allow more serialization formats than just json, maybe also MessagePack?

## Credits:
Hadrian is a wrapper around @jackc's fantastic pglogrepl and was
inspired by both cainophile: and supabase's realtime.

## References:
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

