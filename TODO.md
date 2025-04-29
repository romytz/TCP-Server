# TODO List – TCP Server Project

## Near-Term Goals

- [x] Implement `/quit` command for client disconnection
- [x] Implement `/list` command to show connected clients
- [x] Refactor `handleMessage` into modular command handlers
- [ ] Track active client goroutines using `sync.WaitGroup`
- [ ] Add graceful shutdown handling (server shutdown, close all clients)
- [ ] Create Dockerfile and run server inside Docker container
- [ ] Polish README with full project description and instructions

## Stretch Goals (Optional)

- [ ] Implement `/name` command to let clients set a nickname
- [ ] Broadcast client messages to all other connected clients
- [ ] Improve logging with timestamps
