# Packet Flow

The sequence diagram below shows the start of the game packet and message flow

```mermaid
sequenceDiagram
    participant Server
    participant User

    create participant Ultima Client
    User ->> Ultima Client: Start Game
    Ultima Client ->> Server: Seed Info
    create participant Peer
    Server ->> Peer: Seed Message
    Note right of Server: If Seed Success, create Peer
    Ultima Client ->> Server: send Login Info (0x80)
    Note right of Server: If Auth Success, continues
    Server ->> Ultima Client: Send Game Server List (0xA8)
    Ultima Client ->> Server: Send Select Server (0xA0)
    Server ->> Ultima Client: Send Connect to Server (0x8C)
    create participant Ultima Client 2
    User ->> Ultima Client 2: Connect to Game Server
    Ultima Client 2 ->> Server: Send Game Server Login (0x91)
    Note right of Server: Associate Ultima Client 2 to Peer using Key
    Server ->> Peer: Send Client Features (0xB9)

```