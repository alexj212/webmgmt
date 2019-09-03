package webmgmt

import (
    "bytes"
    "encoding/json"
    "net"
    "net/http"
    "strings"

    "github.com/gorilla/mux"
    "github.com/potakhov/loge"
)

func HttpNocacheContent(w http.ResponseWriter, content string) {
    w.Header().Set("Content-Type", content)
    w.Header().Set("Cache-Control", "no-cache, no-store")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
}

func HttpNocacheJson(w http.ResponseWriter) {
    HttpNocacheContent(w, "text/json")
}

// Put to log actual error, send 500 error code to the client with generic string
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
    w.Header().Set("Cache-Control", "no-cache, no-store")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
    loge.Error("Internal server error %s: %s", mux.CurrentRoute(r).GetName(), err.Error())
    http.Error(w, "Internal server error", 500)
}

func SendJson(w http.ResponseWriter, r *http.Request, val interface{}) {
    bytes, err := json.Marshal(val)
    if err != nil {
        loge.Error("Send json error: %v\n", err)
        InternalServerError(w, r, err)
        return
    }
    HttpNocacheContent(w, "text/json")
    _, err = w.Write(bytes)
    if err != nil {
        loge.Error("error calling w.Write() error: %v\n", err)
    }

}

func AddHeadersHandler(addHeaders map[string]string, h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        for key, value := range addHeaders {
            w.Header().Set(key, value)
        }

        h.ServeHTTP(w, r)
    })
}

// ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
    start net.IP
    end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
    // strcmp type byte comparison
    if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
        return true
    }
    return false
}

var privateRanges = []ipRange{
    ipRange{
        start: net.ParseIP("10.0.0.0"),
        end:   net.ParseIP("10.255.255.255"),
    },
    ipRange{
        start: net.ParseIP("100.64.0.0"),
        end:   net.ParseIP("100.127.255.255"),
    },
    ipRange{
        start: net.ParseIP("172.16.0.0"),
        end:   net.ParseIP("172.31.255.255"),
    },
    ipRange{
        start: net.ParseIP("192.0.0.0"),
        end:   net.ParseIP("192.0.0.255"),
    },
    ipRange{
        start: net.ParseIP("192.168.0.0"),
        end:   net.ParseIP("192.168.255.255"),
    },
    ipRange{
        start: net.ParseIP("198.18.0.0"),
        end:   net.ParseIP("198.19.255.255"),
    },
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func IsPrivateSubnet(ipAddress net.IP) bool {
    // my use case is only concerned with ipv4 atm
    if ipCheck := ipAddress.To4(); ipCheck != nil {
        // iterate over all our ranges
        for _, r := range privateRanges {
            // check if this ip is in a private range
            if inRange(r, ipAddress) {
                return true
            }
        }
    }
    return false
}

func GetIPAddress(r *http.Request) string {
    var ip = ""
    for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
        addresses := strings.Split(r.Header.Get(h), ",")
        // march from right to left until we get a public address
        // that will be the address right before our proxy.
        for i := len(addresses) - 1; i >= 0; i-- {
            ip = strings.TrimSpace(addresses[i])
            // header can contain spaces too, strip those out.
            realIP := net.ParseIP(ip)
            if !realIP.IsGlobalUnicast() || IsPrivateSubnet(realIP) {
                // bad address, go to next
                continue
            }
            return ip
        }
    }

    ip, _, _ = net.SplitHostPort(r.RemoteAddr)
    return ip
}
