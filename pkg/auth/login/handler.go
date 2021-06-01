package login

import (
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

const exampleAppState = "todo state"

type Handler struct {
	oauth2         *oauth2.Config
	offlineAsScope bool
}

func New(oauth2 *oauth2.Config, offlineAsScope bool) *Handler {
	return &Handler{
		oauth2:         oauth2,
		offlineAsScope: offlineAsScope,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodHead, http.MethodGet:
		h.index(w, r)
	case http.MethodPost:
		h.login(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
}

func (h Handler) index(w http.ResponseWriter, r *http.Request) {
	renderIndex(w)
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	var scopes []string
	if extraScopes := r.FormValue("extra_scopes"); extraScopes != "" {
		scopes = strings.Split(extraScopes, " ")
	}
	var clients []string
	if crossClients := r.FormValue("cross_client"); crossClients != "" {
		clients = strings.Split(crossClients, " ")
	}
	for _, client := range clients {
		scopes = append(scopes, "audience:server:client_id:"+client)
	}
	connectorID := ""
	if id := r.FormValue("connector_id"); id != "" {
		connectorID = id
	}

	authCodeURL := ""
	scopes = append(scopes, "openid", "profile", "email")
	if r.FormValue("offline_access") != "yes" {
		authCodeURL = h.OAuth2(scopes).AuthCodeURL(exampleAppState)
	} else if h.offlineAsScope {
		scopes = append(scopes, "offline_access")
		authCodeURL = h.OAuth2(scopes).AuthCodeURL(exampleAppState)
	} else {
		authCodeURL = h.OAuth2(scopes).AuthCodeURL(exampleAppState, oauth2.AccessTypeOffline)
	}
	if connectorID != "" {
		authCodeURL = authCodeURL + "&connector_id=" + connectorID
	}

	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}

func (h Handler) OAuth2(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     h.oauth2.ClientID,
		ClientSecret: h.oauth2.ClientID,
		Endpoint:     h.oauth2.Endpoint,
		Scopes:       append(h.oauth2.Scopes, scopes...),
		RedirectURL:  h.oauth2.RedirectURL,
	}
}
