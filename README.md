# Proxy Rotator
A simple automatic proxy rotator that downloads free proxies from internet and then 
checks it's validity and apply to the system so that the system can use that proxy to get anonymity.

### How to use it?
```
func testProxy() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://ipinfo.io/json", nil)
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func main() {
	proxyHandler := ProxyHandler{}
	proxyHandler.applyNewProxy() # this will automatically apply the proxy to the system.
	testProxy()
}
```