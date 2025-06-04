"""
IMPORTANT: Always use double quotes (") for JSON keys and string values in payloads, not backticks (`).
Example of valid JSON:
{
    "name": "mycc",
    "path": "../asset-transfer-basic/chaincode-go/",
    "channel": "mychannel",
    "version": "1.0",
    "language": "go"
}
"""
from typing import Any
import httpx
from mcp.server.fastmcp import FastMCP
import os
import subprocess

# Initialize FastMCP server
mcp = FastMCP("hlf-fabric")

HLF_API_BASE = "http://localhost:8081"

async def post_json(endpoint: str, payload: dict[str, Any] = None) -> dict[str, Any]:
    async with httpx.AsyncClient() as client:
        try:
            resp = await client.post(f"{HLF_API_BASE}{endpoint}", json=payload, timeout=300.0)
            try:
                return resp.json()
            except Exception as e:
                return {"status": "error", "message": f"Invalid JSON: {e}, raw: {resp.text}"}
        except Exception as e:
            return {"status": "error", "message": str(e)}

@mcp.tool()
async def network_up() -> dict:
    """Bring the Hyperledger Fabric network up."""
    return await post_json("/network/up")

@mcp.tool()
async def network_down() -> dict:
    """Bring the Hyperledger Fabric network down."""
    return await post_json("/network/down")

@mcp.tool()
async def create_channel(channel: str = "mychannel") -> dict:
    """Create a channel. Default is 'mychannel'."""
    return await post_json("/channel/create", {"channel": channel})

@mcp.tool()
async def deploy_chaincode(
    name: str,
    path: str,
    language: str,
    version: str = "1.0",
    channel: str = "mychannel"
) -> dict:
    """Deploy chaincode to the network."""
    payload = {
        "name": name,
        "path": path,
        "language": language,
        "version": version,
        "channel": channel
    }
    return await post_json("/chaincode/deploy", payload)

@mcp.tool()
async def invoke_chaincode(
    channel: str,
    chaincode: str,
    function: str,
    args: list[str]
) -> dict:
    """Invoke a chaincode function."""
    payload = {
        "channel": channel,
        "chaincode": chaincode,
        "function": function,
        "args": args
    }
    return await post_json("/chaincode/invoke", payload)

@mcp.tool()
async def query_chaincode(
    channel: str,
    chaincode: str,
    function: str,
    args: list[str]
) -> dict:
    """Query a chaincode function."""
    payload = {
        "channel": channel,
        "chaincode": chaincode,
        "function": function,
        "args": args
    }
    return await post_json("/chaincode/query", payload)

@mcp.tool()
async def write_chaincode_file(name: str, code: str, path: str = "/Users/padamarajkore/Desktop/generated-chaincodes/") -> dict:
    """Write Go chaincode code to a new directory, add go.mod, and run go mod tidy."""
    try:
        dir_path = os.path.join(path, name)
        os.makedirs(dir_path, exist_ok=True)
        file_path = os.path.join(dir_path, "chaincode.go")
        with open(file_path, "w") as f:
            f.write(code)
        # Write go.mod
        go_mod_content = f"""module {name}

go 1.20

require github.com/hyperledger/fabric-contract-api-go v1.1.0
"""
        go_mod_path = os.path.join(dir_path, "go.mod")
        with open(go_mod_path, "w") as f:
            f.write(go_mod_content)
        # Run go mod tidy
        result = subprocess.run(["go", "mod", "tidy"], cwd=dir_path, capture_output=True, text=True)
        if result.returncode != 0:
            return {"status": "error", "message": f"go mod tidy failed: {result.stderr}"}
        return {"status": "success", "message": f"Chaincode and go.mod written to {dir_path} and go mod tidy succeeded."}
    except Exception as e:
        return {"status": "error", "message": str(e)}

if __name__ == "__main__":
    mcp.run(transport="stdio")
