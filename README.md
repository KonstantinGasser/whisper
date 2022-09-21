# Ares: an in-memory message broker


## How it works

You can create a `Broker` which allows you to add topics and emit events/messages onto a specific topic


## Use-Cases

1. As a user I want to emit events/messages onto a specific topic which can be received and processed by `N` subscribers
2. As a module I want to able to subscribe to a topic without adding further dependencies to the module besides the `Broker` itself.
3. As a module I do not want to be blocked/in-wait if a push to the topic is blocking.
4. As a user I want to throttle (add a buffer to a specific topic).