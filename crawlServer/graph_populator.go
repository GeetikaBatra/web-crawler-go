package crawlServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CreateGremlinQuery(base_url string, url string) ([]byte, error) {
	payload := make(map[string]string)
	//check if query is not empty
	payload["gremlin"] = ConstructGraphNodes(base_url, url)

	json_payload, _ := json.Marshal(payload)

	return json_payload, nil
}

func ConstructGraphNodes(base_url string, url string) string {
	if url != "" {
		base_str := fmt.Sprintf("base_url_node = g.V().has('url','name', '%s').tryNext().orElseGet{"+
			"g.addV('url').property('name', '%s').next()};",
			base_url, base_url)

		new_str := fmt.Sprintf("url_node = g.V().has('url','name', '%s').tryNext().orElseGet{"+
			"g.addV('url').property('name', '%s').next()};",
			url, url)

		edge_str := fmt.Sprintf("edge_c = g.V().has('url', '%s').has('url', '%s')"+
			".in('child_urls').tryNext()"+
			".orElseGet{base_url_node.addEdge('child_urls', url_node)};", base_url, url)
		fmt.Println(base_str + new_str + edge_str)
		return base_str + new_str + edge_str
	}

	return ""
}

func PostGraph(baseUrl string, links []string) {
	for _, link := range links {
		g_query, _ := CreateGremlinQuery(baseUrl, link)

		url := "http://localhost:8182"
		fmt.Println("URL:>", url)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(g_query))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))

	}
}
