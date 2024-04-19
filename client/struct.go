package client

type resp struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    ProxyData `json:"data"`
}

type ProxyData struct {
	Addr    string `json:"addr"`
	City    string `json:"city"`
	EndTime int64  `json:"end_time"`
}
