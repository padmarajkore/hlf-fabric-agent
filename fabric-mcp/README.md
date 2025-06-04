# hlf-mcp

A Python-based MCP (Multi-Chain Platform) tool for automating and orchestrating Hyperledger Fabric network operations and chaincode lifecycle management via the hlf-controller REST API.

## Features
- Bring up/down a Fabric test network
- Create channels
- Deploy, invoke, and query chaincode
- Write and scaffold Go chaincode with dependencies
- Async Python API for integration with agents or automation

## Project Structure
```
.
├── mcp_hlf_tool.py        # Main MCP tool with all Fabric operations
```

## Setup
1. **Clone the repo**
2. **Install Python 3.9+**
3. **Install dependencies:**
   ```sh
   pip install -r requirements.txt
   ```
4. **Ensure Go is installed** (for chaincode writing):
   ```sh
   go version
   # Should print your Go version
   ```
5. **Ensure hlf-controller is running** (see its README)

## Environment Variables
- `HLF_API_BASE` (optional): Base URL for the hlf-controller API. Default is `http://localhost:8081` (set in code).

## Usage
You can use the tool functions in `mcp_hlf_tool.py` as part of an agent, or run the MCP server directly:

### Example Tool Functions
- `network_up()` — Bring the network up
- `network_down()` — Bring the network down
- `create_channel(channel="mychannel")` — Create a channel
- `deploy_chaincode(name, path, language, version="1.0", channel="mychannel")` — Deploy chaincode
- `invoke_chaincode(channel, chaincode, function, args)` — Invoke a chaincode function
- `query_chaincode(channel, chaincode, function, args)` — Query a chaincode function
- `write_chaincode_file(name, code, path=...)` — Write Go chaincode, add go.mod, run go mod tidy

### Example: Deploy and Invoke Chaincode
```python
from mcp_hlf_tool import deploy_chaincode, invoke_chaincode

# Deploy
await deploy_chaincode(
    name="fabcar",
    path="../fabcar/go/",
    language="go",
    version="1.0",
    channel="mychannel"
)

# Invoke
await invoke_chaincode(
    channel="mychannel",
    chaincode="fabcar",
    function="initLedger",
    args=[]
)
```

## Cursor Integration
To use this tool with Cursor, add the following entry to your `~/.cursor/mcp.json` file:

```json
"hlf-controller": {
    "command": "uv",
    "args": [
        "--directory",
        "/Users/padamarajkore/Desktop/weather",
        "run",
        "mcp_hlf_tool.py"
    ]
}
```

**Important:**
- The `--directory` argument should point to the folder containing your MCP tool script.
- The last argument (e.g., `hlf_controller.py`) should be the filename of your MCP tool's entrypoint script.
- For example, if your tool is named `mcp_hlf_tool.py` and is located in `/Users/padamarajkore/Desktop/hlf-mcp`, use:

```json
"hlf-controller": {
    "command": "uv",
    "args": [
        "--directory",
        "/Users/padamarajkore/Desktop/hlf-mcp",
        "run",
        "mcp_hlf_tool.py"
    ]
}
```

Place this inside the top-level JSON object, alongside your other tool definitions. Adjust the directory and script name as needed for your setup.

## Claude Desktop Integration

To use this tool with Claude Desktop, add the following entry to your `claude_desktop_config.json` file (usually found in `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
"hlf-controller": {
    "command": "/Users/padamarajkore/.local/bin/uv",
    "args": [
        "--directory",
        "/Users/padamarajkore/Desktop/weather",
        "run",
        "mcp_hlf_tool.py"
    ]
}
```

**Important:**
- The `--directory` argument should point to the folder containing your MCP tool script.
- The last argument (e.g., `hlf_controller.py`) should be the filename of your MCP tool's entrypoint script.
- Adjust the path to the `uv` command if it is installed elsewhere on your system.
- Place this inside the top-level JSON object, alongside your other tool definitions.

If your tool is named `mcp_hlf_tool.py` and is located in `/Users/padamarajkore/Desktop/hlf-mcp`, use:

```json
"hlf-controller": {
    "command": "/Users/padamarajkore/.local/bin/uv",
    "args": [
        "--directory",
        "/Users/padamarajkore/Desktop/hlf-mcp",
        "run",
        "mcp_hlf_tool.py"
    ]
}
```

## Contribution
- Fork the repo and create a feature branch
- Submit a pull request with clear description
- Ensure code is modular and well-documented

## License
MIT
