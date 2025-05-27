package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"time"

	"hlf-controller/internal/config"
	types "hlf-controller/internal/types"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, cfg.NetworkScriptPath, "up")
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, cfg.NetworkScriptPath, "down")
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, cfg.NetworkScriptPath, "deployCC", "-ccn", req.Name, "-ccp", req.Path, "-ccl", req.Language, "-ccv", version, "-c", channel)
	output, err := cmd.CombinedOutput()
	log.Printf("[DEBUG] DeployChaincodeHandler: deployCC output: %s", string(output))
	if err != nil {
		log.Printf("[ERROR] DeployChaincodeHandler: %v", err)
	}
	successIndicators := []string{
		"Chaincode definition committed",
		"Chaincode is installed",
		"Committed chaincode definition",
		"Query chaincode definition successful",
		"Finished vendoring Go dependencies",
	}
	isSuccess := false
	for _, indicator := range successIndicators {
		if indicator != "" && string(output) != "" && (utils.ContainsIgnoreCase(string(output), indicator)) {
			isSuccess = true
			break
		}
	}
	resp := types.Response{Status: "success", Message: "Output: " + string(output) + "\nError: "}
	if err != nil {
		resp.Message += err.Error()
	}
	if err != nil && !isSuccess {
		resp.Status = "error"
		log.Printf("[ERROR] DeployChaincodeHandler: Deployment failed. Output: %s, Error: %v", string(output), err)
		utils.WriteJSON(w, http.StatusInternalServerError, resp)
		return
	}
	resp.Status = "success"
	log.Printf("[SUCCESS] DeployChaincodeHandler: Chaincode deployed successfully. Output: %s", string(output))
	utils.WriteJSON(w, http.StatusOK, resp)
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
	invokeSpec := map[string]interface{}{
		"function": req.Function,
		"Args":     req.Args,
	}
	invokeSpecBytes, _ := json.Marshal(invokeSpec)
	args := []string{
		"exec",
		"-e", "CORE_PEER_TLS_ENABLED=true",
		"-e", "CORE_PEER_LOCALMSPID=Org1MSP",
		"-e", "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt",
		"-e", "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
		"-e", "CORE_PEER_ADDRESS=peer0.org1.example.com:7051",
		"cli",
		"peer", "chaincode", "invoke",
		"-o", "orderer.example.com:7050",
		"--ordererTLSHostnameOverride", "orderer.example.com",
		"--tls",
		"--cafile", "/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem",
		"-C", req.Channel,
		"-n", req.Chaincode,
		"--peerAddresses", "peer0.org1.example.com:7051",
		"--tlsRootCertFiles", "/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt",
		"--peerAddresses", "peer0.org2.example.com:9051",
		"--tlsRootCertFiles", "/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt",
		"-c", string(invokeSpecBytes),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
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
	querySpec := map[string]interface{}{
		"function": req.Function,
		"Args":     req.Args,
	}
	querySpecBytes, _ := json.Marshal(querySpec)
	args := []string{
		"exec",
		"-e", "CORE_PEER_TLS_ENABLED=true",
		"-e", "CORE_PEER_LOCALMSPID=Org1MSP",
		"-e", "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt",
		"-e", "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
		"-e", "CORE_PEER_ADDRESS=peer0.org1.example.com:7051",
		"cli",
		"peer", "chaincode", "query",
		"-C", req.Channel,
		"-n", req.Chaincode,
		"-c", string(querySpecBytes),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, cfg.NetworkScriptPath, "createChannel", "-c", channel)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] CreateChannelHandler: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, types.Response{Status: "error", Message: string(output) + "\nError: " + err.Error()})
		return
	}
	log.Printf("[SUCCESS] CreateChannelHandler: Channel created successfully. Output: %s", string(output))
	utils.WriteJSON(w, http.StatusOK, types.Response{Status: "success", Message: string(output)})
}
