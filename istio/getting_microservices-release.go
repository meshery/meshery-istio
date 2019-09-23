package main

import (
	"net/http"
	"io/ioutil"

)




func main() {
	resp, err := http.Get("https://raw.githubusercontent.com/GoogleCloudPlatform/microservices-demo/master/release/istio-manifests.yaml")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := ioutil.WriteFile("merged.yaml", body, 0644); err != nil {
		panic(err)
	}
	
		resp, e := http.Get("https://raw.githubusercontent.com/GoogleCloudPlatform/microservices-demo/master/release/kubernetes-manifests.yaml")
	if e != nil {
		panic(e)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	
	d, err := ioutil.ReadFile("merged.yaml")
	if err := ioutil.WriteFile("merged.yaml", append(d, b...), 0644); err != nil {
		panic(err)
	}

	
}
