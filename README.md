# WebSocket Server in Go (Gorilla WebSocket)

A simple and scalable WebSocket server built using Go and the Gorilla WebSocket package. This project demonstrates how to manage real-time client connections using a **Hub-based architecture**, avoiding manual locking with mutexes.

---

## Architecture Overview

```
Client → HTTP → WebSocket Upgrade → Hub → Clients
```

### Flow:

1. Client sends HTTP request
2. Server upgrades connection to WebSocket
3. Client is registered in the Hub
4. Hub manages:

   * client registration
   * client unregistration
   * message broadcasting
5. Messages are routed through the Hub (single goroutine)

---

## Why Hub-Based Design?

Instead of:

* managing shared maps with mutex
