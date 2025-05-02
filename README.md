# TCP Server in Go

A simple TCP server implemented in Go, designed to accept multiple client connections, read incoming messages, and respond dynamically.

This project was built to practice client-server networking concepts and to demonstrate Go’s concurrency model using Goroutines, channels, and synchronization primitives like WaitGroups.

## 📦 Project Structure

- `main.go` — Entry point; starts the server and handles shutdown  
- `server.go` — Initializes the listener and manages lifecycle  
- `client_manager.go` — Tracks and manages connected clients  
- `handlers.go` — Delegates client message handling  
- `command/` — Modular command implementations:
  - `interface.go` — Command interface definition
  - `list.go`, `quit.go`, `default.go` — Specific command logic
- `message/message.go` — Defines the `Message` struct and format
- `Dockerfile` — Containerization instructions
- `README.md` — Project documentation

## 🚀 How to Run

### Run locally (no Docker)

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/tcp-server.git
   cd tcp-server
2. Run the server:
   ```bash
   go run .
3. Connect to the server using telnet:
   ```bash
   telnet localhost 3000
4. Type a message and press Enter — the server will respond.
5. Type /quit to disconnect, or /list to see connected clients.
6. Press ENTER in the server terminal to gracefully shut down.

## 🐳 Docker Support

This project can also run in a Docker container for easy portability and isolation.

1. Build the Docker image:
    ```bash
    docker build -t go-tcp-server .
    ```
2. Run the container:
    ```bash
    docker run -it -p 3000:3000 go-tcp-server
    ```
3. Connect to the server using Telnet (in a separate terminal):
    ```bash
    telnet localhost 3000
    ```
4. Type a message and press Enter — the server will respond.
5. Type `/quit` to disconnect, or `/list` to see connected clients.
6. Press ENTER in the Docker container terminal to gracefully shut down the server.

## 🛠 Features
- Accepts multiple simultaneous client connections
- Reads full messages until newline (\n)
- Responds to custom commands:
  - /quit — Disconnects the client
  - /list — Lists all connected clients
- Graceful shutdown with sync.WaitGroup (no panics)
- Uses channels to safely communicate between Goroutines
- Clean modular structure with Go best practices
- Can be run inside a Docker container

## 💬 Example Interaction
```ruby
$ telnet localhost 3000
[Server]: Welcome [::1]:55590!
Hello
[Server]: Thank you for your message, [::1]:55590!
/list
[Server]: 2 clients connected:
[::1]:55595
[::1]:55590
/quit
[Server]: Goodbye, [::1]:55590!

```

## 📚 Technologies Used
- Go (Golang)
- TCP networking (net package)
- Goroutines
- Channels
- sync.WaitGroup
- Docker (for containerization)

## 🎯 Learning Goals
- Understand TCP server-client architecture
- Practice concurrency and safe communication in Go
- Structure Go projects into modular, maintainable files
- Learn how to containerize and run Go applications using Docker

## Contributors 👥

- **Romy Tzafrir** - [GitHub Profile](https://github.com/romytz)

## License 📜

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
