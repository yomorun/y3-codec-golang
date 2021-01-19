package pprof

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"strings"
)

type pprofConf struct {
	Enabled    bool
	PathPrefix string
	Endpoint   string
}

const (
	pprofEnabled = "Y3_PPROF_ENABLED"
	pathPrefix   = "Y3_PPROF_PATH_PREFIX"
	endpoint     = "Y3_PPROF_ENDPOINT"
)

func newConf() pprofConf {
	conf := pprofConf{}
	conf.Enabled = getBool(pprofEnabled, true)
	conf.PathPrefix = getString(pathPrefix, "/debug/cpu/")
	conf.Endpoint = getString(endpoint, "0.0.0.0:6060")
	return conf
}

func Run() {
	conf := newConf()
	if conf.Enabled == false {
		return
	}

	mux := http.NewServeMux()
	pathPrefix := conf.PathPrefix
	mux.HandleFunc(pathPrefix,
		func(w http.ResponseWriter, r *http.Request) {
			name := strings.TrimPrefix(r.URL.Path, pathPrefix)
			if name != "" {
				pprof.Handler(name).ServeHTTP(w, r)
				return
			}
			pprof.Index(w, r)
		})
	mux.HandleFunc(pathPrefix+"cmdline", pprof.Cmdline)
	mux.HandleFunc(pathPrefix+"profile", pprof.Profile)
	mux.HandleFunc(pathPrefix+"symbol", pprof.Symbol)
	mux.HandleFunc(pathPrefix+"trace", pprof.Trace)

	server := http.Server{
		Addr:    conf.Endpoint,
		Handler: mux,
	}

	fmt.Printf("PProf server start... http://%s%s\n", conf.Endpoint, conf.PathPrefix)
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			fmt.Errorf("cpu server closed")
		} else {
			fmt.Errorf("cpu server error: %v", err)
		}
	}
}

func getBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if len(value) != 0 {
		flag, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return flag
	}
	return defaultValue
}

func getString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) != 0 {
		return value
	}
	return defaultValue
}
