package utils

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	json "github.com/json-iterator/go"
)

// GetLocalAddress 获取外部接口ip
func GetLocalAddress(prober string) string {
	if prober == "" {
		prober = "8.8.8.8:80"
	}

	conn, err := net.Dial("udp", prober)
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[:idx]
}

// GetNonLocalIPsIfHostIsIPAny
func GetNonLocalIPsIfHostIsIPAny(host string, all bool) (bool, []string, error) {
	ip := net.ParseIP(host)
	// If this is not an IP, we are done
	if ip == nil {
		return false, nil, nil
	}
	// If this is not 0.0.0.0 or :: we have nothing to do.
	if !ip.IsUnspecified() {
		return false, nil, nil
	}

	var ips []string
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ipStr := ip.String()
			// Skip non global unicast addresses
			if !ip.IsGlobalUnicast() || ip.IsUnspecified() {
				ip = nil
				continue
			}
			ips = append(ips, ipStr)
			if !all {
				break
			}
		}
	}
	return true, ips, nil
}

// EncodeJson
func EncodeJson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// DecodeJson
func DecodeJson(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func Get32MD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Get16MD5Encode(data string) string {
	return Get32MD5Encode(data)[8:24]
}

// GetClientIP 获取客户端IP
func GetClientIP(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if r.Header.Get("X-Forwarded-For") != "" {
		remoteAddr = r.Header.Get("X-Forwarded-For")
	} else if r.Header.Get("Ali-Cdn-Real-Ip") != "" {
		remoteAddr = r.Header.Get("Ali-Cdn-Real-Ip")
	} else if r.Header.Get("Remote_addr") != "" {
		remoteAddr = r.Header.Get("Remote_addr")
	}

	remoteArr := strings.Split(remoteAddr, ",")
	remoteAddr = remoteArr[0]
	if strings.Contains(remoteAddr, ":") {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	return remoteAddr
}
