# Hyperledger Fabric Agent Suite

https://github.com/user-attachments/assets/75c3414a-3901-42ff-bf33-8f57b0ca96b4


Works for test-network.

A modular toolkit for managing Hyperledger Fabric test networks and chaincode lifecycle, featuring:
- **hlf-controller**: A Go-based REST API for network and chaincode operations
- **fabric-mcp**: A Python-based MCP tool for automation, scripting, and agent integration

---

## Features
- Bring up/down a Fabric test network
- Create channels
- Deploy, invoke, and query chaincode
- Write and scaffold Go chaincode with dependencies
- Integrate with agents, Cursor, or Claude Desktop
- **Seamlessly integrate the Python MCP tool with LLMs (such as via Cursor or Claude Desktop), enabling LLM agents to perform Fabric network and chainco


de operations automatically**

---

## Prerequisites
- **Go** (1.18+)
- **Python** (3.9+)
- **pip** (for Python dependencies)
- **Hyperledger Fabric test network** (see [fabric-samples](https://github.com/hyperledger/fabric-samples))

---

## Setup Instructions

### 1. Clone the Repo
```sh
git clone https://github.com/padmarajkore/hlf-fabric-agent.git
cd hlf-controller
```

### 2. Set Up hlf-controller (Go REST API)
```sh
cd hlf-controller
# (Optional) Set the path to your network.sh if not default
export HLF_NETWORK_SCRIPT_PATH=/path/to/network.sh

https://github.com/user-attachments/assets/a1090104-2a3f-4c2d-a2b1-8582068e1378


# Run the server
go run main.go
# or
go build -o hlf-controller
./hlf-controller
```
- The API will be available at `http://localhost:8081`

### 3. Set Up fabric-mcp (Python MCP Tool)
```sh
cd ../fabric-mcp
pip install httpx
# (Optional) Install any other agent/MCP dependencies
# Ensure Go is installed and in your PATH for chaincode writing
```

---

## Environment Variables
- `HLF_NETWORK_SCRIPT_PATH`: Path to your Fabric `network.sh` (for hlf-controller)
- `HLF_API_BASE`: Base URL for the hlf-controller API (for fabric-mcp, default is `http://localhost:8081`)

---

## How to Use

### Start the Network and Deploy Chaincode
1. **Bring up the network:**
   - POST `/network/up` (via hlf-controller API or fabric-mcp tool)
2. **Create a channel:**
   - POST `/channel/create` (default: `mychannel`)
3. **Deploy chaincode:**
   - POST `/chaincode/deploy` (specify name, path, language, etc.)
4. **Invoke/query chaincode:**
   - POST `/chaincode/invoke` and `/chaincode/query`

You can use the Python MCP tool to add into llm. so they can use the tool, to perform the operations.

---

## Integration
- **Cursor:** Add the fabric-mcp tool to your `~/.cursor/mcp.json` (see fabric-mcp/README.md)
- **Claude Desktop:** Add to `claude_desktop_config.json` (see fabric-mcp/README.md)

---

## Example JSON Configurations for Integration

### Cursor (`~/.cursor/mcp.json`)
Add the following entry to your `~/.cursor/mcp.json` file to integrate the MCP tool with Cursor:

```json
"hlf-controller": {
    "command": "uv",
    "args": [
        "--directory",
        "/Users/padamarajkore/Desktop/fabric-mcp",
        "run",
        "mcp_hlf_tool.py"
    ]
}
```

### Claude Desktop (`claude_desktop_config.json`)
Add the following entry to your `claude_desktop_config.json` file (usually found in `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
"hlf-controller": {
    "command": "/Users/padamarajkore/.local/bin/uv",
    "args": [
        "--directory",
        "/Users/padamarajkore/Desktop/fabric-mcp",
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
├── hlf-controller/        # Go REST API server
├── fabric-mcp/               # Python MCP tool
├── README.md              # This file
```

---

## License
MIT 
