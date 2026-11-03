# Demo Message Broker

An in-memory publish-subscribe message broker

## Run
Requires Go 1.25 or later.

``sh
go run .
``

## Design Decisions
- Broker keeps only one append only log, Publishing appends to it and reading never removes from it. Each subscriber stores the offset of the next message it will read.

