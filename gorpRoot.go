package main

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/server"

	"os"
	"github.com/hashicorp/consul/api"
	"goji.io/pat"
	"net/http"
	"goji.io"
	"strings"
	"github.com/mitchellh/mapstructure"
	"log"
	"github.com/dghubble/sling"
	"fmt"
)

func main() {

	currDir, _ := os.Getwd()
	rpConf := conf.LoadConfig("", map[string]interface{}{"staticsPath": currDir})
	srv := server.New(rpConf)

	srv.AddRoute(func(router *goji.Mux) {
		router.Handle(pat.Get("/*"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nodeInfos, _ := srv.Sd.DoWithClient(func(client interface{}) (interface{}, error) {

				services, _, e := client.(*api.Client).Catalog().Services(&api.QueryOptions{})
				if nil != e {
					return nil, e
				}

				nodeInfos := make(map[string]*nodeInfo, len(services))

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
						ni.BaseUrl = fmt.Sprintf("http://%s:%d/", inst.Service.Address, inst.Service.Port)
						nodeInfos[k] = &ni
					}

				}

				return nodeInfos, nil
			})

			server.WriteJSON(w, 200, aggregateInfo(nodeInfos.(map[string]*nodeInfo)))

		}))
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

func aggregateInfo(nodeInfos map[string]*nodeInfo) map[string]interface{} {
	var aggregated = make(map[string]interface{}, len(nodeInfos))
	for node, info := range nodeInfos {
		var healthRS map[string]interface{}
		_, e := sling.New().Base(info.BaseUrl).Get(info.StatusPageUrlPath).ReceiveSuccess(&healthRS)
		if nil != e {
			log.Println(e)
			continue
		}
		aggregated[node] = healthRS
	}
	return aggregated

}

type nodeInfo struct {
	BaseUrl            string
	StatusPageUrlPath  string
	HealthCheckUrlPath string
}
