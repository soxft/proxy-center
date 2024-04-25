package proxy

type DmResp struct {
	Code    int
	Success bool
	Msg     string
	Data    []struct {
		Ip      string
		Port    int
		Endtime string
		City    string
		Rosname string
	}
	Num int
}

type DyResp struct {
	Ret  int
	Msg  string
	Data []struct {
		Ip       string
		Port     string
		Prov     string
		City     string
		Isp      string
		Deadline string
		Outip    string
	}
}

type XqResp struct {
	Code int    `json:"ret"`
	Msg  string `json:"msg"`
	Data []struct {
		Ip        string `json:"IP"`
		Port      int    `json:"Port"`
		IpAddress string `json:"IpAddress"`
	} `json:"data"`
}

type qgResp struct {
	Code string `json:"code"`
	Data []struct {
		ProxyIp  string `json:"proxy_ip"`
		Server   string `json:"server"`
		Area     string `json:"area"`
		Isp      string `json:"isp"`
		Deadline string `json:"deadline"`
	}
}

type Proxy struct {
	Addr    string
	EndTime int64
	City    string
	Isp     string
}
