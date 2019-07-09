# pubsub-dead-letter-log

`psdll` is the CLI tool to list/re-publish **p**ub**s**ub **d**ead-**l**etter **l**og.

## What's *pubsub dead-letter log*?

*pubsub dead-letter log* is the fail log for publishing the pubsub message.

name | type | description
-- | -- | --
message | `Message` | pubsub message
message.data | `string` (base64 encoded) | data of the pubsub message
message.attributes | `map[string]string` | attributes of the pubsub message
project | `string` | project ID for Google Cloud Platform
topic | `string` | topic name for Google Cloud Pub/Sub
publisher | `string` | application name publishing the message
pod_name | `string` | pod name which publisher application running
timestamp | `string` (ISO8601) | timestamp which the message published
error | `string` | error information

## How to use

### Installing

```bash
go get -u github.com/wantedly/pubsub-dead-letter-log/cmd/psdll
```

### List logs

```console
$ psdll list s3://wantedly-pubsub-dead-letter-log/dev-project
+-------------------------------+-------------+---------+----------------------------------------+
|           TIMESTAMP           |   PROJECT   |  TOPIC  |               ATTRIBUTES               |
+-------------------------------+-------------+---------+----------------------------------------+
| 2019-06-19 16:00:59 +0900 JST | dev-project | awesome | published_at=2019-06-19T16:00:56+09:00 |
+-------------------------------+-------------+---------+----------------------------------------+
```

### Re-publish logs

```console
$ psdll publish s3://wantedly-pubsub-dead-letter-log/dev-project --project=foo --topic=bar --attribute=retry_count=1
published: id=643425616190931, attributes=map[published_at:2019-06-19T16:00:56+09:00 retry_count:1]
```
