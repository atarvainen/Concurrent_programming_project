# Concurrent chat application with multiple chat rooms and private messages

We'll program this with Go. It's designed for concurrent programming, we've used it somewhat, and we've time constraints, so it's an obvious choice for this application.
The focus of this app is to utilize concurrency's most significant advantages which are the ability to handle multiple tasks concurrently and doing that with a minimal amount of waiting.
Challenge is how to make the most use of the aforementioned advantages. Thus, our design goals are preventing deadlocks/livelocks and minimizing congestion to critical sections.

Works best on Linux.

## Identifying critical sections

### Handling connections

Things to keep in mind. Minimize read/write of connections.
Different clients' Internet connection speeds vary a lot so very slow connections must not block the app.
Have to protect from simultaneous read/write.

#### Connections

Writing to data structure. Updating connection information to threads.

#### Disconnections

Inform threads when clients disconnects. Handle ungraceful disconnections so they won't block.

### Handling messages

Internet connection speeds vary and that cannot block the chat.

#### Sending messages

Sending must not block. Should not send to disconnected clients.

#### Receiving messages

Receiving must not block.

## Design

### Server

#### Handling connections

The server maintains rooms and rooms maintains connection information about clients. The rooms starts a new goroutine for each connection (client). Multiple connection information data structures are required when multiple chat rooms are available (each room has it's own connections). This way connection info is in smaller and more accessible data structure. The server keeps track of each room.

A channel is provided for the goroutines (clients) to handle connection information. The channel is used for passing information about the connections. The channel handles locking the connection information, so it's safe from parallel access. The client pushes the connection information to the channel only when it's altered.

When the client disconnects from a chat room it informs the server via the channel. If the client disconnects ungracefully an error should occur when trying to send a message to it. In this situation, the room updates the connection information.

#### Handling messages

Chat rooms may be implemented in several ways. A chat room could have a single channel which is used used by clients. Problem with this approach is that it's peer to peer pattern, which means slow connections will slow down the chat room.

Faster, provided that the server has enough resources, is that each client has own channel to the chat room mediated by the server. In this pattern, the clients are independent of each other and only communicate to the server via its own channel. When the server receives data from a channel it forwards it to every channel. This way slow connections won't block the chat room.

Buffered channels should be used to limit the memory burden on the server. Also when a channel blocks because of it's full, it can be a signal of a very bad connection which should be dropped. The exact size of the buffer is to be determined.

In our current implementation the room receives messages from a message channel and then prints the message to the clients.

### Chat rooms

The focus of this work is in concurrent programming, not in creating the best possible chat server. Thus we won't spend much time in how to present messages from different chat rooms. We'll print messages from all the rooms the user is into the same screen.

We won't implement, for example, private rooms nor create one to one chats (direct messages).

#### Creating a chat room

Command ```/create <room name>``` creates a new channel.

#### Destroying a chat room

Not implemented.

#### Joining to a chat room

Command ```/join <room name>``` allows a user to receive messages from a given chat room.

#### Leaving a chat room

When a user joins to a new room, she disconnects from her current room. No separate leave-functionality is implemented.

## Implementation

As this is a concurrent programming course, we want to minimize and simplify other programming aspects as much as possible. Therefore, we'll use our teacher's example chat program as a starting point and refactor and expand from there.

## Client, connecting to a server

For example, just launch the server with ```go run chatserver.go```, open new terminals and use telnet to make client connections.

## Results

We are satisfied with the result, given our time constraints. We tested the implementation with GNU Terminator, which allowed us to launch multiple connections feasibly. What we would have liked to implement better was passing messages from a room to clients. Overall, this course was useful as was learning Go in the process.