package binance

import "fmt"

const (
	scheme    = "https"
	host      = "www.binance.com"
	pingPath  = "/api/v3/ping"
	depthPath = "/api/v3/depth?symbol=%s&limit=%d"
)

var (
	pingUrl  = fmt.Sprintf("%s://%s%s", scheme, host, pingPath)
	depthUrl = fmt.Sprintf("%s://%s%s", scheme, host, depthPath)
)
