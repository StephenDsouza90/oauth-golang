package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/go-oauth2/oauth2/generates"
	"github.com/go-oauth2/oauth2/manage"
	"github.com/go-oauth2/oauth2/models"
	"github.com/go-oauth2/oauth2/server"
	"github.com/go-oauth2/oauth2/store"
	"github.com/go-session/session"
)

var (
	dumpvar      bool
	clientId     string
	clientSecret string
	clientDomain string
	portvar      int
)

func init() {
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
	flag.StringVar(&clientId, "i", "222222", "The client id being passed in")
	flag.StringVar(&clientSecret, "s", "22222222", "The client secret being passed in")
	flag.StringVar(&clientDomain, "r", "http://localhost:9094", "The domain of the redirect url")
	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

func tokenManager() *manage.Manager {
	// Manage storing and generating token
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewAccessGenerate())
	return manager
}

func getClientStore() *store.ClientStore {
	// A store for the client id, secret domain.
	clientStore := store.NewClientStore()
	clientStore.Set(clientId, &models.Client{
		ID:     clientId,
		Secret: clientSecret,
		Domain: clientDomain,
	})
	return clientStore
}

func getServer(manager *manage.Manager) *server.Server {
	server := server.NewServer(server.NewConfig(), manager)
	return server
}

func userAuthorizeHandler(writer http.ResponseWriter, response *http.Request) (userId string, err error) {
	store, err := session.Start(response.Context(), writer, response)
	if err != nil {
		return
	}

	loggedInUserID, ok := store.Get("LoggedInUserID")
	if !ok {
		if response.Form == nil {
			response.ParseForm()
		}

		store.Set("ReturnUri", response.Form)
		store.Save()

		writer.Header().Set("Location", "/login")
		writer.WriteHeader(http.StatusFound)
		return
	}

	userId = loggedInUserID.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	store, err := session.Start(request.Context(), writer, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if request.Method == "POST" {
		if request.Form == nil {
			if err := request.ParseForm(); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		store.Set("LoggedInUserID", request.Form.Get("username"))
		store.Save()

		writer.Header().Set("Location", "/auth")
		writer.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(writer, request, "server/static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "auth", r)
	}
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "server/static/auth.html")
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func main() {
	flag.Parse()

	manager := tokenManager()
	clientStore := getClientStore()
	manager.MapClientStorage(clientStore)

	server := getServer(manager)
	server.SetUserAuthorizationHandler(userAuthorizeHandler)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			dumpRequest(os.Stdout, "authorize", r)
		}

		store, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		r.Form = form

		store.Delete("ReturnUri")
		store.Save()

		err = server.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "token", r)
		}

		err := server.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Printf("Server is running at %d port.\n", portvar)
	log.Printf("Point your OAuth client Auth endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/token")

	http.ListenAndServe(fmt.Sprintf(":%d", portvar), nil)
}
