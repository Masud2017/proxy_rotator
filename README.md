# Proxy Rotator
A simple automatic proxy rotator that downloads free proxies from internet and then 
checks it's validity and apply to the system so that the system can use that proxy to get anonymity.

### How to use it?
The package can be downloaded by using the following command :
```
go get -u github.com/Masud2017/proxy_rotator
```

or it can be used as regular import inside the code

```
go mod tidy
```

```
import (
    "fmt"
    "github.com/Masud2017/proxy_rotator"
    "io/ioutil"
    "net/http"
)
```

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
	proxyHandler.ApplyNewProxy() # this will automatically apply the proxy to the system.
	testProxy()
	proxyHandler.ClearProxy() # after using the proxy it is important to remove it before it gets invalid.
}
```

#### If you face any issues then feel free to open an issue.
