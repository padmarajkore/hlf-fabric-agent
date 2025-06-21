package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"

	"hlf-controller/internal/config"
	"hlf-controller/internal/types"
	"hlf-controller/internal/utils"
)

func UpNetworkHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] UpNetworkHandler called. Method: %s, Path: %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] UpNetworkHandler: Method not allowed: %s", r.Method)
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	log.Println("[INFO] Bringing network UP...")
	cfg := config.LoadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Network)
	defer cancel()
	cmd := exec.CommandContext(ctx, cfg.Network.ScriptPath, "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] UpNetworkHandler: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, types.Response{Status: "error", Message: string(output) + "\nError: " + err.Error()})
		return
	}
	log.Printf("[SUCCESS] Network UP. Output: %s", string(output))
	utils.WriteJSON(w, http.StatusOK, types.Response{Status: "success", Message: string(output)})
}

func DownNetworkHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] DownNetworkHandler called. Method: %s, Path: %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] DownNetworkHandler: Method not allowed: %s", r.Method)
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	log.Println("[INFO] Bringing network DOWN...")
	cfg := config.LoadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Network)
	defer cancel()
	cmd := exec.CommandContext(ctx, cfg.Network.ScriptPath, "down")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] DownNetworkHandler: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, types.Response{Status: "error", Message: string(output) + "\nError: " + err.Error()})
		return
	}
	log.Printf("[SUCCESS] Network DOWN. Output: %s", string(output))
	utils.WriteJSON(w, http.StatusOK, types.Response{Status: "success", Message: string(output)})
}

func DeployChaincodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] DeployChaincodeHandler called. Method: %s, Path: %s", r.Method, r.URL.Path)
	var req types.DeployChaincodeRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		log.Printf("[ERROR] DeployChaincodeHandler: Invalid request body: %v", err)
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("[INFO] DeployChaincodeHandler: Params: Name=%s, Path=%s, Language=%s, Version=%s, Channel=%s", req.Name, req.Path, req.Language, req.Version, req.Channel)
	if req.Name == "" || req.Path == "" || req.Language == "" {
		log.Printf("[ERROR] DeployChaincodeHandler: Missing required fields")
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields: name, path, language")
		return
	}
	version := req.Version
	if version == "" {
		version = "1.0"
	}
	channel := req.Channel
	if channel == "" {
		channel = "mychannel"
	}
	cfg := config.LoadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Deploy)
	defer cancel()
	cmd := exec.CommandContext(ctx, cfg.Network.ScriptPath, "deployCC", "-ccn", req.Name, "-ccp", req.Path, "-ccl", req.Language, "-ccv", version, "-c", channel)
	output, err := cmd.CombinedOutput()
	log.Printf("[DEBUG] DeployChaincodeHandler: deployCC output: %s", string(output))
	if err != nil {
		log.Printf("[ERROR] DeployChaincodeHandler: Deployment failed. Output: %s, Error: %v", string(output), err)
		utils.WriteJSON(w, http.StatusInternalServerError, types.Response{Status: "error", Message: string(output) + "\nError: " + err.Error()})
		return
	}
	log.Printf("[SUCCESS] DeployChaincodeHandler: Chaincode deployed successfully. Output: %s", string(output))
	utils.WriteJSON(w, http.StatusOK, types.Response{Status: "success", Message: string(output)})
}

func InvokeChaincodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] InvokeChaincodeHandler called. Method: %s, Path: %s", r.Method, r.URL.Path)
	var req types.InvokeChaincodeRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		log.Printf("[ERROR] InvokeChaincodeHandler: Invalid request body: %v", err)
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("[INFO] InvokeChaincodeHandler: Params: Channel=%s, Chaincode=%s, Function=%s, Args=%v", req.Channel, req.Chaincode, req.Function, req.Args)
	if req.Channel == "" || req.Chaincode == "" || req.Function == "" {
		log.Printf("[ERROR] InvokeChaincodeHandler: Missing required fields")
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields: channel, chaincode, function")
		return
	}

	cfg := config.LoadConfig()
	invokeSpec := map[string]interface{}{
		"function": req.Function,
		"Args":     req.Args,
	}
	invokeSpecBytes, err := json.Marshal(invokeSpec)
	if err != nil {
		log.Printf("[ERROR] InvokeChaincodeHandler: Failed to marshal invoke spec: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create invoke spec")
		return
	}

	args := []string{
		"exec",
		"-e", "CORE_PEER_TLS_ENABLED=true",
		"-e", "CORE_PEER_LOCALMSPID=" + cfg.Network.CLI.MSP_ID,
		"-e", "CORE_PEER_TLS_ROOTCERT_FILE=" + cfg.Network.CLI.TLSRootCertFile,
		"-e", "CORE_PEER_MSPCONFIGPATH=" + cfg.Network.CLI.MSPConfigPath,
		"-e", "CORE_PEER_ADDRESS=" + cfg.Network.CLI.PeerAddress,
		"cli",
		"peer", "chaincode", "invoke",
		"-o", cfg.Network.Orderer.Address,
		"--ordererTLSHostnameOverride", cfg.Network.Orderer.HostnameOverride,
		"--tls",
		"--cafile", cfg.Network.Orderer.TLSCaCert,
		"-C", req.Channel,
		"-n", req.Chaincode,
	}

	for _, peer := range cfg.Network.Peers {
		args = append(args, "--peerAddresses", peer.Address, "--tlsRootCertFiles", peer.TLSRootCertFile)
	}

	args = append(args, "-c", string(invokeSpecBytes))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Invoke)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] InvokeChaincodeHandler: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, types.Response{Status: "error", Message: string(output) + "\nError: " + err.Error()})
		return
	}
	log.Printf("[SUCCESS] InvokeChaincodeHandler: Chaincode invoke successful. Output: %s", string(output))
	utils.WriteJSON(w, http.StatusOK, types.Response{Status: "success", Message: string(output)})
}

func QueryChaincodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] QueryChaincodeHandler called. Method: %s, Path: %s", r.Method, r.URL.Path)
	var req types.InvokeChaincodeRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		log.Printf("[ERROR] QueryChaincodeHandler: Invalid request body: %v", err)
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("[INFO] QueryChaincodeHandler: Params: Channel=%s, Chaincode=%s, Function=%s, Args=%v", req.Channel, req.Chaincode, req.Function, req.Args)
	if req.Channel == "" || req.Chaincode == "" || req.Function == "" {
		log.Printf("[ERROR] QueryChaincodeHandler: Missing required fields")
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields: channel, chaincode, function")
		return
	}

	cfg := config.LoadConfig()
	querySpec := map[string]interface{}{
		"function": req.Function,
		"Args":     req.Args,
	}
	querySpecBytes, err := json.Marshal(querySpec)
	if err != nil {
		log.Printf("[ERROR] QueryChaincodeHandler: Failed to marshal query spec: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create query spec")
		return
	}

	args := []string{
		"exec",
		"-e", "CORE_PEER_TLS_ENABLED=true",
		"-e", "CORE_PEER_LOCALMSPID=" + cfg.Network.CLI.MSP_ID,
		"-e", "CORE_PEER_TLS_ROOTCERT_FILE=" + cfg.Network.CLI.TLSRootCertFile,
		"-e", "CORE_PEER_MSPCONFIGPATH=" + cfg.Network.CLI.MSPConfigPath,
		"-e", "CORE_PEER_ADDRESS=" + cfg.Network.CLI.PeerAddress,
		"cli",
		"peer", "chaincode", "query",
		"-C", req.Channel,
		"-n", req.Chaincode,
		"-c", string(querySpecBytes),
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Query)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] QueryChaincodeHandler: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, types.Response{Status: "error", Message: string(output) + "\nError: " + err.Error()})
		return
	}
	log.Printf("[SUCCESS] QueryChaincodeHandler: Chaincode query successful. Output: %s", string(output))
	utils.WriteJSON(w, http.StatusOK, types.Response{Status: "success", Message: string(output)})
}

func CreateChannelHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] CreateChannelHandler called. Method: %s, Path: %s", r.Method, r.URL.Path)
	var req types.CreateChannelRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		log.Printf("[ERROR] CreateChannelHandler: Invalid request body: %v", err)
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("[INFO] CreateChannelHandler: Params: Channel=%s", req.Channel)
	channel := req.Channel
	if channel == "" {
		channel = "mychannel"
	}
	cfg := config.LoadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Channel)
	defer cancel()
	cmd := exec.CommandContext(ctx, cfg.Network.ScriptPath, "createChannel", "-c", channel)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] CreateChannelHandler: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, types.Response{Status: "error", Message: string(output) + "\nError: " + err.Error()})
		return
	}
	log.Printf("[SUCCESS] CreateChannelHandler: Channel created successfully. Output: %s", string(output))
	utils.WriteJSON(w, http.StatusOK, types.Response{Status: "success", Message: string(output)})
}
