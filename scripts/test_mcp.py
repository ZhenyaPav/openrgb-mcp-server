#!/usr/bin/env python3
import argparse
import json
import queue
import subprocess
import sys
import threading
import time

LATEST_PROTOCOL_VERSION = "2025-06-18"


def read_stdout(proc, out_q, verbose):
    for raw in proc.stdout:
        line = raw.strip()
        if not line:
            continue
        if verbose:
            print(f"<< {line}", file=sys.stderr)
        try:
            msg = json.loads(line)
        except json.JSONDecodeError:
            # Ignore non-JSON lines (server banners, logs, etc.)
            continue
        out_q.put(msg)


def send(proc, msg, verbose):
    payload = json.dumps(msg, separators=(",", ":"))
    if verbose:
        print(f">> {payload}", file=sys.stderr)
    proc.stdin.write(payload + "\n")
    proc.stdin.flush()


def wait_for_id(out_q, req_id, timeout_s):
    deadline = time.time() + timeout_s
    while time.time() < deadline:
        try:
            msg = out_q.get(timeout=0.1)
        except queue.Empty:
            continue
        if isinstance(msg, dict) and msg.get("id") == req_id:
            return msg
    raise TimeoutError(f"Timed out waiting for response id={req_id}")


def build_tool_arguments(tool, args):
    if tool in ("list_devices", "list_profiles"):
        return {}
    if tool == "set_device_color":
        if args.device_id is None:
            raise ValueError("--device-id is required for set_device_color")
        return {"device_id": args.device_id, "r": args.r, "g": args.g, "b": args.b}
    if tool == "set_all_color":
        return {"r": args.r, "g": args.g, "b": args.b}
    if tool == "set_profile":
        if not args.profile_name:
            raise ValueError("--profile-name is required for set_profile")
        return {"profile_name": args.profile_name}
    raise ValueError(f"Unknown tool: {tool}")


def main():
    parser = argparse.ArgumentParser(description="Minimal MCP stdio test client for OpenRGB MCP server")
    parser.add_argument("--server", default="./bin/openrgb-mcp-server", help="Path to MCP server binary")
    parser.add_argument("--server-arg", action="append", default=[], help="Extra arg for the server command")
    parser.add_argument("--call", choices=["list_devices", "set_device_color", "set_all_color", "list_profiles", "set_profile"], help="Tool to call after tools/list")
    parser.add_argument("--device-id", type=int, help="Device ID for set_device_color")
    parser.add_argument("--profile-name", help="Profile name for set_profile")
    parser.add_argument("--r", type=int, default=0, help="Red channel 0-255")
    parser.add_argument("--g", type=int, default=0, help="Green channel 0-255")
    parser.add_argument("--b", type=int, default=0, help="Blue channel 0-255")
    parser.add_argument("--timeout", type=float, default=5.0, help="Timeout in seconds for responses")
    parser.add_argument("--verbose", action="store_true", help="Print raw JSON-RPC messages")
    args = parser.parse_args()

    cmd = [args.server] + args.server_arg

    proc = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        text=True,
        bufsize=1,
    )

    out_q = queue.Queue()
    reader = threading.Thread(target=read_stdout, args=(proc, out_q, args.verbose), daemon=True)
    reader.start()

    try:
        send(
            proc,
            {
                "jsonrpc": "2.0",
                "id": 1,
                "method": "initialize",
                "params": {
                    "protocolVersion": LATEST_PROTOCOL_VERSION,
                    "capabilities": {},
                    "clientInfo": {"name": "openrgb-mcp-test", "version": "0.1.0"},
                },
            },
            args.verbose,
        )
        init_res = wait_for_id(out_q, 1, args.timeout)
        print("initialize:", json.dumps(init_res, indent=2))

        send(proc, {"jsonrpc": "2.0", "method": "notifications/initialized", "params": {}}, args.verbose)

        send(proc, {"jsonrpc": "2.0", "id": 2, "method": "tools/list"}, args.verbose)
        tools_res = wait_for_id(out_q, 2, args.timeout)
        print("tools/list:", json.dumps(tools_res, indent=2))

        if args.call:
            tool_args = build_tool_arguments(args.call, args)
            send(
                proc,
                {
                    "jsonrpc": "2.0",
                    "id": 3,
                    "method": "tools/call",
                    "params": {"name": args.call, "arguments": tool_args},
                },
                args.verbose,
            )
            call_res = wait_for_id(out_q, 3, args.timeout)
            print("tools/call:", json.dumps(call_res, indent=2))
    finally:
        if proc.stdin:
            proc.stdin.close()
        try:
            proc.wait(timeout=2)
        except subprocess.TimeoutExpired:
            proc.terminate()


if __name__ == "__main__":
    main()
