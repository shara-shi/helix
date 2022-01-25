package services

import (
	"fmt"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/shara/helix/database"
)

type Application struct {
	http.Server
	Name     string
	Version  string
	mux      *mux.Router
	listener net.Listener

	DB *database.DB
}

func NewApplication(name string, port int, version string, dns string) *Application {
	mx := mux.NewRouter()

	db, err := database.Open(dns)

	if err != nil {
		// 服务启动的时候，链接数据库错误！
		log.Fatal(err.Error())
	}

	app := &Application{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mx,
		},
		Name:     name,
		Version:  version,
		mux:      mx,
		listener: nil,

		DB: db,
	}
	app.Bootstrap()
	return app
}

//  starts up an Application
func (app *Application) Bootstrap() {
	log.Info("application : ", app.Name)
	app.mux.HandleFunc("/ping", pingHandler)
	app.mux.HandleFunc("/version", func(rw http.ResponseWriter, r *http.Request) { versionHandler(rw, r, app) })
	app.mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./endpoints/img_process/resources/"))))
}

// trivial indication that we are alive
func pingHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// report application git revision
func versionHandler(rw http.ResponseWriter, r *http.Request, app *Application) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.Write([]byte(app.Version))
}

// RunService brings up the service and begins listening for requests
func (app *Application) RunService() {
	// start listening
	l, err := net.Listen("tcp", app.Server.Addr)
	if err != nil {
		log.Fatal("failed to listen", "error", err)
	}
	app.listener = l

	// start serving requests
	log.Info("started listening for requests", "addr", app.Server.Addr)
	app.Serve(app.listener)

}

// 注册web servicesURI以"/vN"开始,例如"/v1/"
func (app *Application) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	// names that start with "/vN/..." (where N is 1,2) are REST endpoints and should be measured. Others
	if len(path) >= 2 && path[:2] == "/v" { // 对外提供服务
		app.mux.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
			handler(rw, r)
		})
	} else { // 对内提供服务
		app.mux.HandleFunc(path, handler)
	}
}
