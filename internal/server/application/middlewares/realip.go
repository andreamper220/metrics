package middlewares

import (
	"net"
	"net/http"
)

func WithIpCheck(hf http.HandlerFunc, trustedSubnetStr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, trustedSubnet, err := net.ParseCIDR(trustedSubnetStr)
		if err != nil {
			http.Error(w, "Server error with trusted subnet", http.StatusInternalServerError)
			return
		}

		ipStr := r.Header.Get("X-Real-IP")
		ip := net.ParseIP(ipStr)
		if !trustedSubnet.Contains(ip) {
			http.Error(w, "Access denied for not trusted subnet", http.StatusForbidden)
			return
		}

		hf(w, r)
	}
}
