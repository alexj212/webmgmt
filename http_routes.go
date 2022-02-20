package webmgmt

import (
	"fmt"
	"github.com/alexj212/gox/httpx"
	"net/http"
)

// handleServerVersion handles the server version rest handler
func (app *MgmtApp) handleServerVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("handleServerVersion\n")
	httpx.HttpNocacheJson(w)
	httpx.SendJson(w, r, AppBuildInfo)
}
