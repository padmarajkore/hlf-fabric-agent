# Hyperledger Fabric Agent Suite

https://github.com/user-attachments/assets/75c3414a-3901-42ff-bf33-8f57b0ca96b4


Works for test-network.

A modular toolkit for managing Hyperledger Fabric test networks and chaincode lifecycle, featuring:
- **hlf-controller**: A Go-based REST API for network and chaincode operations.
- **hlf-mcp**: A Python-based MCP tool for automation, scripting, and agent integration.

---

## Features
- **Automated Fabric Setup**: The `hlf-controller` automatically downloads and configures Hyperledger Fabric binaries and samples on first run.
- **Dynamic Configuration**: Easily configure the entire Fabric network topology via a central `config.yaml` file.
- **Full Lifecycle Management**: Bring up/down the network, create channels, deploy, invoke, and query chaincode via a simple REST API.
- **Agent-Ready**: The `hlf-mcp` tool allows LLMs (like those in Cursor or Claude Desktop) to interact with your Fabric network seamlessly.

---

## Prerequisites
- **Go** (1.18+)
- **Python** (3.9+)
- **pip** (for Python dependencies)

---

## Setup Instructions

### 1. Clone the Repo
```sh
git clone <your-repo-url>
cd <repo-root>
```

### 2. Set Up and Run hlf-controller (Go REST API)
The controller handles its own prerequisites.
```sh
cd hlf-controller
go run main.go
```
- The first time you run this, it will check for `fabric-samples` and download them to your home directory if they are missing. The server will not start until this process is complete.
- The API will then be available at `http://localhost:8081`.

### 3. Set Up hlf-mcp (Python MCP Tool)
```sh
cd ../hlf-mcp
pip install httpx
# (Optional) Install any other agent/MCP dependencies
```

---

## Configuration

### hlf-controller
The Go controller is configured via the `hlf-controller/config.yaml` file. This is the primary way to define your network topology, including peer/orderer details, certificate paths, and timeouts.

### Environment Variables
- `HLF_CONFIG_PATH`: Path to a custom `config.yaml` file for the controller.
- `HLF_NETWORK_SCRIPT_PATH`: Overrides the path to your Fabric `network.sh`. If not set, the path from `config.yaml` is used, which in turn defaults to `~/fabric-samples/test-network/network.sh`.
- `HLF_API_BASE`: Base URL for the hlf-controller API (used by `hlf-mcp`, defaults to `http://localhost:8081`).

---

## How to Use

## Integration
- **Cursor:** Add the hlf-mcp tool to your `~/.cursor/mcp.json` (see hlf-mcp/README.md).
- **Claude Desktop:** Add to `claude_desktop_config.json` (see hlf-mcp/README.md).

---

## Example JSON Configurations for Integration

### Cursor (`~/.cursor/mcp.json`)
Add the following entry to your `~/.cursor/mcp.json` file to integrate the MCP tool with Cursor:

```json
"hlf-controller": {
    "command": "uv",
    "args": [
        "--directory",
        "/Users/padamarajkore/Desktop/hlf-mcp", //change it to your local path for hlf-mcp folder
        "run",
        "mcp_hlf_tool.py"
    ]
}
```

### Claude Desktop (`claude_desktop_config.json`)
Add the following entry to your `claude_desktop_config.json` file (usually found in `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
"hlf-controller": {
    "command": "/Users/padamarajkore/.local/bin/uv", //change it to your local path for uv binary
    "args": [
        "--directory",
        "/Users/padamarajkore/Desktop/hlf-mcp", //change it to your local path for hlf-mcp folder
        "run",
        "mcp_hlf_tool.py"
    ]
}
```

**Note:**
- Adjust the `--directory` and script name to match the actual location and filename of your MCP tool.
- Adjust the path to the `uv` command if it is installed elsewhere on your system.
- Place these entries inside the top-level JSON object, alongside your other tool definitions.

---

## Project Structure
```
.
├── hlf-controller/        # Go REST API server with its own config.yaml
├── hlf-mcp/               # Python MCP tool
├── README.md              # This file
```

---

## License
MIT
