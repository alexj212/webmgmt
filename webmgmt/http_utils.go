package webmgmt

import (
    "encoding/json"
    "net/http"

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
