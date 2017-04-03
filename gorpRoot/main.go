package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/avarabyeu/goRP/commons"
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/registry"
	"github.com/avarabyeu/goRP/server"
	"github.com/dghubble/sling"
	"github.com/gorilla/handlers"
	"github.com/hashicorp/consul/api"
	"goji.io"
	"goji.io/pat"
)

func main() {

	rpConf := conf.LoadConfig("", map[string]interface{}{})
	rpConf.AppName = "gorproot"

	srv := server.New(rpConf)

	srv.AddRoute(func(router *goji.Mux) {
		router.Use(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		})

		router.HandleFunc(pat.Get("/composite/info"), func(w http.ResponseWriter, r *http.Request) {
			commons.WriteJSON(http.StatusOK, aggregateInfo(getNodesInfo(srv.Sd, true)), w)
		})
		router.HandleFunc(pat.Get("/composite/health"), func(w http.ResponseWriter, r *http.Request) {
			commons.WriteJSON(http.StatusOK, aggregateHealth(getNodesInfo(srv.Sd, false)), w)
		})
		router.HandleFunc(pat.New("/"), func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/ui/", http.StatusMovedPermanently)
		})

		u, e := url.Parse("http://" + rpConf.Consul.Address)
		if e != nil {
			log.Fatal("Cannot parse consul URL")
		}
		proxy := httputil.NewSingleHostReverseProxy(u)
		router.Handle(pat.Get("/consul/*"), http.StripPrefix("/consul/", proxy))
		router.Handle(pat.Get("/v1/*"), proxy)
	})
	srv.StartServer()
}

func parseKVTag(tags []string, tagsMap map[string]string) {
	for _, tag := range tags {
		kv := strings.Split(tag, "=")
		if 2 == len(kv) {
			tagsMap[kv[0]] = kv[1]
		}
	}
}

func aggregateHealth(nodesInfo map[string]*nodeInfo) map[string]interface{} {
	var aggregated = make(map[string]interface{}, len(nodesInfo))
	for node, info := range nodesInfo {
		var rs map[string]interface{}

		if "" != info.getHealthCheckURL() {
			_, e := sling.New().Base(info.BaseURL).Get(info.getHealthCheckURL()).Receive(&rs, &rs)
			if nil != e {
				rs = make(map[string]interface{}, 1)
				rs["status"] = "DOWN"
			}
		} else {
			rs = make(map[string]interface{}, 1)
			rs["status"] = "UNKNOWN"
		}

		aggregated[node] = rs
	}
	return aggregated
}

func aggregateInfo(nodesInfo map[string]*nodeInfo) map[string]interface{} {
	var aggregated = make(map[string]interface{}, len(nodesInfo))
	for node, info := range nodesInfo {
		var rs map[string]interface{}
		_, e := sling.New().Base(info.BaseURL).Get(info.getStatusPageURL()).ReceiveSuccess(&rs)
		if nil != e {
			log.Println(e)
			continue
		}
		if nil != rs {
			aggregated[node] = rs
		}

	}
	return aggregated
}

func getNodesInfo(discovery registry.ServiceDiscovery, passing bool) map[string]*nodeInfo {
	nodesInfo, _ := discovery.DoWithClient(func(client interface{}) (interface{}, error) {
		services, _, e := client.(*api.Client).Catalog().Services(&api.QueryOptions{})
		if nil != e {
			return nil, e
		}
		nodesInfo := make(map[string]*nodeInfo, len(services))
		for k := range services {
			instances, _, e := client.(*api.Client).Health().Service(k, "", passing, &api.QueryOptions{})
			if nil != e {
				return nil, e
			}
			for _, inst := range instances {
				tagsMap := map[string]string{}
				parseKVTag(inst.Service.Tags, tagsMap)

				var ni nodeInfo
				ni.BaseURL = fmt.Sprintf("http://%s:%d/", inst.Service.Address, inst.Service.Port)
				ni.Tags = tagsMap
				nodesInfo[strings.ToUpper(k)] = &ni
			}

		}

		return nodesInfo, nil
	})
	return nodesInfo.(map[string]*nodeInfo)
}

type nodeInfo struct {
	BaseURL string
	Tags    map[string]string
}

func (ni *nodeInfo) getStatusPageURL() string {
	return ni.Tags["statusPageUrlPath"]
}
func (ni *nodeInfo) getHealthCheckURL() string {
	return ni.Tags["healthCheckUrlPath"]
}
