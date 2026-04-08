# 🚀 WebSocket Server in Go (Gorilla WebSocket)

A scalable and production-ready WebSocket server built using **Go** and the **Gorilla WebSocket** package.
This project demonstrates a clean **Hub-based architecture** for managing real-time connections without relying on mutex locks.

---

## 📌 Architecture Overview

```
Client → HTTP → WebSocket Upgrade → Hub → Clients
```

### 🔄 Flow

1. Client sends an HTTP request
2. Server upgrades the connection to WebSocket
3. Client registers with the Hub
4. Hub handles:

   * Client registration
   * Client unregistration
   * Message broadcasting
5. All communication flows through the Hub (single goroutine)

---

## 🧠 Why Hub-Based Design?

Instead of managing shared state like:

* ❌ Maps with mutex locks
* ❌ Complex synchronization logic

This project uses:

* ✅ A **central Hub (event loop)**
* ✅ **Channels for communication**
* ✅ **Single goroutine for state management**

### Benefits

* No race conditions
* Cleaner code
* Easier to scale and debug
* Follows Go concurrency best practices

---

## 🔐 Role-Based Permissions

The server supports a flexible permission system for managing users in channels:

### 👑 Admin

* CanSend
* CanCreateChannel
* CanDeleteChannel
* CanInviteUser
* CanKick
* CanMute

### 🛡️ Moderator

* CanSend
* CanInviteUser
* CanKick
* CanMute

### 👤 Member

* CanSend

### 👀 Guest

* (Read-only / limited access — customizable)

---

## 🔑 JWT Authentication

This project integrates **JWT-based authentication** for secure communication.

### How it works:

* A shared `jwt_secret` is used between:

  * WebSocket server
  * Application layer (e.g., REST API / frontend)
* JWT is used to:

  * Authenticate users
  * Encode user identity & roles
  * Validate permissions before actions

### Example Use Cases:

* Authenticate WebSocket connections
* Attach user roles (Admin, Moderator, etc.)
* Secure channel access

---

## ⚙️ Key Features

* 🔄 Real-time communication via WebSockets
* 🧩 Hub-based architecture (no mutex needed)
* 🔐 JWT authentication & authorization
* 👥 Role-based access control (RBAC)
* 📡 Channel-based messaging system
* 🧵 Efficient concurrency using goroutines & channels

---
