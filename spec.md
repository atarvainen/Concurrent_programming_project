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

### Server

#### Handling connections

The server maintains connection information and starts a new goroutine for each connection (client). Multiple connection information data structures are required when multiple chat rooms are available. The server must be able to keep track of each room.

A single channel is provided for the goroutines (clients). The channel is used for passing information about the connections. The channel handles locking the connection information, so it's safe from parallel access. The server pushes the connection information to the channel only when it's altered.

When the client disconnects from a chat room it informs the server via the channel. If the client disconnects ungracefully an error should occur when trying to send a message to it. In this situation, the server should update the connection information.

#### Handling messages

Chat rooms may be implemented in several ways. A chat room could have a single channel which is used used by clients. Problem with this approach is that it's peer to peer pattern, which means slow connections will slow down the chat room.

Faster, provided that the server has enough resources, is that each client has own channel to the chat room mediated by the server. In this pattern, the clients are independent of each other and only communicate to the server via its own channel. When the server receives data from a channel it forwards it to every channel. This way slow connections won't block the chat room.

Buffered channels should be used to limit the memory burden on the server. Also when a channel blocks because of it's full, it can be a signal of a very bad connection which should be dropped. The exact size of the buffer is to be determined.

In short, the server's chat room connections information data structure should have a channel and addition to address and port.

### Chat rooms

#### How to create a chat room?