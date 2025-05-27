# hlf-controller

A modular, production-ready REST API server for managing Hyperledger Fabric test networks and chaincode lifecycle operations.

## Features
- Bring up/down a Fabric test network
- Create channels
- Deploy, invoke, and query chaincode
- Modular Go codebase with clear separation of concerns
- Configurable via environment variables
- Structured logging for all operations

## Project Structure
```
.
├── main.go                # Entrypoint, wires up the HTTP server
├── go.mod                 # Go module definition
├── internal/
│   ├── handlers/          # HTTP handler logic
│   ├── utils/             # Utility functions (JSON, logging, etc.)
│   ├── config/            # Configuration loading
│   └── types/             # Request/response types
```

## Setup
1. **Clone the repo**
2. **Install Go** (1.18+ recommended)
3. **Install Hyperledger Fabric test network** (see [fabric-samples](https://github.com/hyperledger/fabric-samples))
4. **Set environment variables** (optional, see below)
5. **Run the server:**
   ```sh
   go run main.go
   # or
   go build -o hlf-controller
   ./hlf-controller
   ```

## Environment Variables
- `HLF_NETWORK_SCRIPT_PATH` (optional): Path to your `network.sh` script. Default is `/Users/padamarajkore/fabric-samples/test-network/network.sh`.

## API Usage
All endpoints accept/return JSON. Example endpoints:

### Network
- `POST /network/up` — Bring the network up
- `POST /network/down` — Bring the network down

### Channel
- `POST /channel/create` — Create a channel
  - Body: `{ "channel": "mychannel" }`

### Chaincode
- `POST /chaincode/deploy` — Deploy chaincode
  - Body: `{ "name": "fabcar", "path": "../fabcar/go/", "language": "go", "version": "1.0", "channel": "mychannel" }`
- `POST /chaincode/invoke` — Invoke chaincode
  - Body: `{ "channel": "mychannel", "chaincode": "fabcar", "function": "initLedger", "args": [] }`
- `POST /chaincode/query` — Query chaincode
  - Body: `{ "channel": "mychannel", "chaincode": "fabcar", "function": "queryAllCars", "args": [] }`

## Contribution
- Fork the repo and create a feature branch
- Submit a pull request with clear description
- Ensure code is modular and well-documented

## License
MIT 