# Concurrent chat application with multiple chat rooms and private messages

We'll program this with Go. It's designed for concurrent programming, we've used it somewhat, and we've time constraints, so it's an obvious choice for this application.
The focus of this app is to utilize concurrency's most significant advantages which are the ability to handle multiple tasks concurrently and doing that with a minimal amount of waiting.
Challenge is how to make the most use of the aforementioned advantages. Thus, our design goals are preventing deadlocks/livelocks and minimizing congestion to critical sections.

## Identifying critical sections

### Handling connections

Things to keep in mind. Minimize read/write of connections. 
Different clients' Internet connection speeds vary a lot so very slow connections most not block the app.
Have to protect from simultaneous read/write.

#### Connections

Writing to data structure. Updating connection information to threads.

#### Disconnections

Inform threads when clients disconnects. Handle ungraceful disconnections so they won't block.

### Handling messages

Internet connection speeds vary and that cannot block the chat.

#### Sending messages

Sending most not block. Should not send to disconnected clients.

#### Receiving messages

Receiving most not block.

## Design