# hlf-controller

A modular, production-ready REST API server for managing Hyperledger Fabric test networks and chaincode lifecycle operations.

## Features
- **Dynamic Configuration**: Configure all network, peer, and orderer details via a central `config.yaml` file.
- **Automated Setup**: Automatically downloads and sets up Hyperledger Fabric binaries and `fabric-samples` if not found on startup.
- **Lifecycle Management**: Bring up/down a Fabric test network, create channels, and manage the full chaincode lifecycle.
- **RESTful API**: Simple, clean REST endpoints for all operations.
- **Robust & Modular**: Built with a clear separation of concerns, making it easy to extend and maintain.
- **Structured Logging**: Detailed logging for all operations.

## Project Structure
```
.
├── main.go                # Entrypoint, wires up the HTTP server
├── go.mod                 # Go module definition
├── config.yaml            # Central configuration file for the network
├── internal/
│   ├── handlers/          # HTTP handler logic
│   ├── utils/             # HTTP helper functions
│   ├── config/            # Configuration loading logic
│   └── types/             # Request/response types
```

## Setup
1.  **Install Go** (1.18+ recommended).
2.  **Run the server:**
    ```sh
    go run main.go
    ```
3.  **That's it!** The first time you run the server, it will automatically check for Hyperledger Fabric prerequisites (binaries and `fabric-samples`) and download them to your home directory if they are missing. The server will not start until this setup is complete.

## Configuration
The controller is configured via the `config.yaml` file in the project root. This file allows you to define your entire network topology, including peer/orderer addresses, certificate paths, and timeouts.

### Environment Variables
While most configuration is done in `config.yaml`, you can override file paths with these environment variables:

-   `HLF_CONFIG_PATH`: Path to a custom `config.yaml` file. Defaults to `config.yaml` in the working directory.
-   `HLF_NETWORK_SCRIPT_PATH`: Overrides the path to your `network.sh` script. If not set, the path from `config.yaml` is used. If that is also empty, it defaults to `~/fabric-samples/test-network/network.sh`.

## API Usage
All endpoints accept/return JSON.

### Network
- `POST /network/up` — Bring the network up.
- `POST /network/down` — Bring the network down.

### Channel
- `POST /channel/create` — Create a channel.
  - Body: `{ "channel": "mychannel" }`

### Chaincode
- `POST /chaincode/deploy` — Deploy chaincode.
  - Body: `{ "name": "fabcar", "path": "../chaincode/fabcar/go", "language": "go", "version": "1.0", "channel": "mychannel" }`
- `POST /chaincode/invoke` — Invoke a chaincode function.
  - Body: `{ "channel": "mychannel", "chaincode": "fabcar", "function": "InitLedger", "args": [] }`
- `POST /chaincode/query` — Query a chaincode function.
  - Body: `{ "channel": "mychannel", "chaincode": "fabcar", "function": "QueryAllCars", "args": [] }`

## Contribution
- Fork the repo and create a feature branch.
- Submit a pull request with a clear description of your changes.
- Ensure code is modular and well-documented.

## License
MIT 