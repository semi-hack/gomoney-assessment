# Deno Message Broker

An in-memory publish-subscribe message broker written in Go.

## Requirements

Go 1.25 or later 

## Run

```sh
go run .
```

This runs the example sequence from the brief:

1. Create a broker and a publisher.
2. Publish messages 1–5.
3. Create Subscriber 1 and subscribe it.
4. Publish messages 6–10.
5. Create Subscriber 2 and subscribe it.
6. Read and print every message for both subscribers.

Expected output (20 lines):

```
Subscriber 1 received: Message 1
...
Subscriber 1 received: Message 10
Subscriber 2 received: Message 1
...
Subscriber 2 received: Message 10
```                   


## Design decisions

**Append-only.**
Broker keeps only one append only log, Publishing appends to it and reading never removes from it. Each subscriber stores the offset of the next message it will read. so offsets increase in the order messages were published.

**Subscribers start from the beginning of the log.**
A new subscriber starts at offset 0 and can read every message published before it subscribed. 

**Generic message payloads.**
`Broker[T]`, `Publisher[T]` and `Subscriber[T]` are generic, so a broker can carry strings, structs or any other type.

**The broker assigns offsets.**
Each `Message` records its `Offset`, set by the broker when the message is appended. `Publish` returns the stored message, so the caller knows its position in the log. Multiple publishers can share a broker, and every message still gets one total order.



