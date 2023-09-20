package proxy_rotator

// proxy list download link : https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all
import (
	"fmt"
	"github.com/trananhtung/proxy-checker"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var proxy_list_file_name string = "proxy_list.txt"
var valid_proxy_list_file_name string = "valid_proxy_list.txt"
var proxy_download_link string = "https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all"

type ProxyHandler struct {
}

func (proxyHandler *ProxyHandler) ping(proxy string) {
	proxy = "http://" + strings.Trim(proxy, "\n")
	proxy = strings.Trim(proxy, " ")
	//prox, _ := url.Parse("http://47.98.206.38:80")
	prox, _ := url.Parse(proxy)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(prox)}}
	fmt.Println("value of proxy : ", proxy)
	//	os.Setenv("HTTP_PROXY",strings.Trim(proxy," "))

	//	os.Setenv("HTTP_PROXY",proxy)

	req, err := http.NewRequest("GET", "http://ipinfo.io/json", nil)
	url, _ := http.ProxyFromEnvironment(req)
	fmt.Println("Printing proxy from env : ", url)
	if err == nil {
		fmt.Println("Hyo")
	}
	//	ur,_ := http.ProxyFromEnvironment(req)
	//	fmt.Println("proxy url : ",ur)
	resp, _ := client.Do(req)
	if resp.StatusCode == 200 {
		fmt.Println("Status code : ", resp.StatusCode)

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println(string(body))
	} else {
		fmt.Println("The proxy is invalid")
	}

}

func (proxyHandler *ProxyHandler) download_proxy() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", proxy_download_link, nil)
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)

	file, _ := os.Create(proxy_list_file_name)
	file.Write(body)
}

// channel first element will be status (bool) and second element will be the result proxy ip
func (proxyHandler *ProxyHandler) checker_worker(proxy string, wg *sync.WaitGroup, workerChannel chan bool) {
	fmt.Println("Chekcing availability for proxy : ", proxy)
	defer wg.Done()
	defer close(workerChannel)

	targetURL := "http://ipinfo.io/json"
	timeout := uint(5)
	result := proxy_checker.ProxyTest(proxy, targetURL, timeout)

	if result {
		fmt.Println(proxy, " is valid")
		workerChannel <- true
	} else {
		fmt.Println(proxy, "is invalid")
		workerChannel <- false
	}
}

/*func check_proxy(proxy string) uint {
	workerChannel := make([]chan string,2)
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	checker_worker(proxy,&wg,&mutex,workerChannel)
	defer close(workerChannel)
	return 0
}*/

func (proxyHandler *ProxyHandler) writeValidProxy(validProxyMapper map[string]bool) {
	file, err := os.Create(valid_proxy_list_file_name)
	var buffer string

	for key, _ := range validProxyMapper {
		buffer += key + "\n"
	}

	if err != nil {
		fmt.Println(err)
	}

	file.Write([]byte(buffer))

	file.Close()

}

func (proxyHandler *ProxyHandler) getValidProxyFromValidProxyList() string {
	file, _ := os.Open(valid_proxy_list_file_name)
	data, _ := ioutil.ReadAll(file)

	proxy_list := strings.Split(string(data), "\n")

	// cleaning up the unnecessary control code
	for index, value := range proxy_list {
		tempValueData := strings.Trim(value, "\r")
		tempValueData = strings.Trim(value, " ")
		proxy_list[index] = tempValueData
	}

	for _, value := range proxy_list {
		urlTarget := "http://ipinfo.io/json"
		timeOut := uint(5)
		result := proxy_checker.ProxyTest(value, urlTarget, timeOut)

		if result {
			return value
		}
	}

	return ""
}

func isFileExists(fileName string) bool {
	_, err := os.Stat(fileName)

	if err == nil {
		return true
	}

	return false
}

func (proxyHandler *ProxyHandler) get_working_proxy() string {
	if isFileExists(valid_proxy_list_file_name) {
		validProxy := proxyHandler.getValidProxyFromValidProxyList()
		if len(validProxy) > 0 {
			return validProxy
		}
		err := os.Remove(valid_proxy_list_file_name)

		if err != nil {
			fmt.Println(err)
		}

		proxyHandler.get_working_proxy()
	}

	if isFileExists(proxy_list_file_name) {
		//var stringArray [2]string
		var channelProxyMapper = make(map[string]chan bool)
		var validProxyMapper = make(map[string]bool)
		proxy_file, _ := os.Open("proxy_list.txt")
		data, _ := ioutil.ReadAll(proxy_file)
		proxy_list := strings.Split(string(data), "\n")
		proxy_file.Close()

		wg := sync.WaitGroup{}
		//mutex := sync.Mutex{}
		proxyCount := len(proxy_list)

		wg.Add(proxyCount - 1)

		for _, value := range proxy_list {
			workerChannel := make(chan bool)
			if len(value) > 0 {
				dd := "http://" + value
				dd = strings.Trim(dd, "\r")

				go proxyHandler.checker_worker(dd, &wg, workerChannel)
				channelProxyMapper[dd] = workerChannel
			}
		}
		for key, value := range channelProxyMapper {
			status := <-value
			fmt.Println("For proxy ", key, " status is : ", <-value)
			if status {
				validProxyMapper[key] = status
			}
		}

		proxyHandler.writeValidProxy(validProxyMapper)
		fmt.Println("New valid proxy has been written...")

		wg.Wait()
	} else {
		// if the proxy file is not already available then download the file and restart the function again.
		fmt.Println("The proxy file is not available so downloading the proxy file.")
		proxyHandler.download_proxy()
		fmt.Println("The proxy file is downloaded successfully, so retrying the process again.")
		proxyHandler.get_working_proxy()
	}

	return ""
}

/*
*
Automatically rotate the proxy for each http request
@return newValidProxy string
*/
func (proxyHandler *ProxyHandler) rotateProxy() string {
	proxy := proxyHandler.get_working_proxy()
	if len(proxy) <= 0 {
		return proxyHandler.rotateProxy()
	}
	return proxy
}

func (proxyHandler *ProxyHandler) applyNewProxy() {
	newValidProxy := proxyHandler.rotateProxy()

	os.Setenv("HTTP_PROXY", newValidProxy)
}

func testProxy() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://ipinfo.io/json", nil)
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func main() {
	proxyHandler := ProxyHandler{}
	proxyHandler.applyNewProxy()
	testProxy()
	//download_proxy()
}
