# NSQ Tester

### This repo is for testing cases for NSQ
This will spin up all the components to run NSQ with a producer and consumer.

By default the producer will push out 1k messages to NSQ (1msg/millisecond).
Producer will log out how many messages were sent on success and failure.

Consumer by default will read the entire NSQ payload including the message and number of attempts. You can limit this to just the message body if you like. Every 10 seconds the Consumer will notify via log output how many messages it has consumed.

NSQ is using TLS, so please run the generate-certs.sh prior to starting services.
NSQ is configured to use a data directory to persist data between restarts.

## Use
- docker-compose up   `# start all services`
- docker-compose down `# stop all services`
- docker-compose build `# build producer and consumer containers` 

## Services
 - nsqd
 - nsqadmin
 - nsqlookupd
 - producer
 - consumer
