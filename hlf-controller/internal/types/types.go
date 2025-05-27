package types

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DeployChaincodeRequest struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Language string `json:"language"`
	Version  string `json:"version,omitempty"`
	Channel  string `json:"channel,omitempty"`
}

type InvokeChaincodeRequest struct {
	Channel   string   `json:"channel"`
	Chaincode string   `json:"chaincode"`
	Function  string   `json:"function"`
	Args      []string `json:"args"`
}

type CreateChannelRequest struct {
	Channel string `json:"channel,omitempty"`
}
