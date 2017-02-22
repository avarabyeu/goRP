package main

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/server"
	"fmt"
	"github.com/dghubble/sling"
	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
	"goji.io"
	"goji.io/pat"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/avarabyeu/goRP/registry"
)

func main() {

	currDir, _ := os.Getwd()
	rpConf := conf.LoadConfig("", map[string]interface{}{"staticsPath": currDir})
	srv := server.New(rpConf)

	srv.AddRoute(func(router *goji.Mux) {
		router.HandleFunc(pat.Get("/composite/info"), func(w http.ResponseWriter, r *http.Request) {
			server.WriteJSON(w, 200, aggregateInfo(getNodesInfo(srv.Sd)))
		})
		router.HandleFunc(pat.Get("/composite/health"), func(w http.ResponseWriter, r *http.Request) {
			server.WriteJSON(w, 200, aggregateHealth(getNodesInfo(srv.Sd)))
		})
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

func aggregateHealth(nodeInfos map[string]*nodeInfo) map[string]interface{} {
	var aggregated = make(map[string]interface{}, len(nodeInfos))
	for node, info := range nodeInfos {
		var rs map[string]interface{}

		if "" != info.HealthCheckURLPath {
			_, e := sling.New().Base(info.BaseURL).Get(info.HealthCheckURLPath).Receive(&rs, &rs)
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

func aggregateInfo(nodeInfos map[string]*nodeInfo) map[string]interface{} {
	var aggregated = make(map[string]interface{}, len(nodeInfos))
	for node, info := range nodeInfos {
		var rs map[string]interface{}
		_, e := sling.New().Base(info.BaseURL).Get(info.StatusPageURLPath).ReceiveSuccess(&rs)
		if nil != e {
			log.Println(e)
			continue
		}
		aggregated[node] = rs
	}
	return aggregated
}

func getNodesInfo(discovery registry.ServiceDiscovery) map[string]*nodeInfo {
	nodesInfo, _ := discovery.DoWithClient(func(client interface{}) (interface{}, error) {
		services, _, e := client.(*api.Client).Catalog().Services(&api.QueryOptions{})
		if nil != e {
			return nil, e
		}
		nodesInfo := make(map[string]*nodeInfo, len(services))
		for k := range services {
			instances, _, e := client.(*api.Client).Health().Service(k, "", true, &api.QueryOptions{})
			if nil != e {
				return nil, e
			}
			for _, inst := range instances {
				tagsMap := map[string]string{}
				parseKVTag(inst.Service.Tags, tagsMap)

				var ni nodeInfo
				mapstructure.Decode(tagsMap, &ni)
				ni.BaseURL = fmt.Sprintf("http://%s:%d/", inst.Service.Address, inst.Service.Port)
				nodesInfo[k] = &ni
			}

		}

		return nodesInfo, nil
	})
	return nodesInfo.(map[string]*nodeInfo)
}

type nodeInfo struct {
	BaseURL            string
	StatusPageURLPath  string
	HealthCheckURLPath string
}
