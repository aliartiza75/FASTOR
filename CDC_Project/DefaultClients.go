package main 

import (
	"fmt"
	"math/rand"
        "time"
        "os"
	"net/http"
	"strconv"
	"io/ioutil"
	"strings"
	//"html/template"
        

)

type Page struct {
	Title string
}

var s_IP="127.0.0.1"
var c_IP="127.0.0.1"
var server_port= "8081"


func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, This is defult page%s!", r.URL.Path[1:])
}

func randomNumber( rangee int) int {

	rand.Seed(time.Now().UTC().UnixNano())
	num:=rand.Intn(rangee)
	num=num+2000
	return num

}

func handShake(port string, preference string ) string {

	url:="http://"+s_IP+":8081/saveConfig/"+c_IP+port+","+preference
	req, _ := http.NewRequest("GET",url , nil)
    client := &http.Client{}
    resp, _ := client.Do(req)
    fmt.Println(resp.Status)
    return resp.Status

}

func defaultPreference( preference int) string {
	if preference == 0 {
		return "Mid"
	} else if preference == 1 {
		return "Exit"
	}
	return "None"
}


func Entry_ClientToServer(w http.ResponseWriter, r *http.Request) {
        requestedUrl :=r.URL.Path[8:]
	fmt.Println(requestedUrl)// User requested page
	
	url:="http://"+s_IP+":"+server_port+"/EntryToServer/User Requested URL: "+requestedUrl
	req, _ := http.NewRequest("GET",url , nil)
        client := &http.Client{}
        resp, _ := client.Do(req)
        fmt.Println(resp.Status)
        
        defer resp.Body.Close()

        
        
        
        if resp.StatusCode == 200 { // OK
            bodyBytes, _ := ioutil.ReadAll(resp.Body)
            bodyString := string(bodyBytes)
            fmt.Println("response from server for middle layer : "+bodyString)
            
            // Code to send request to middle layer
                ip:=strings.Split(bodyString, ":")[0]
                port:=strings.Split(bodyString, ":")[1]
                
                
                
            	url:="http://"+ip+":"+port+"/middlelayer/"+requestedUrl
	        req, _ := http.NewRequest("GET",url , nil)
                client := &http.Client{}
                resp, _ := client.Do(req)
                fmt.Println(resp.Status)
                
                
                bodyBytes1, _ := ioutil.ReadAll(resp.Body)          // getting response from mid layer
                bodyString1 := string(bodyBytes1)
                fmt.Println("body from mid=>"+bodyString1)
                
                fmt.Fprintf(w, string(bodyString1)) 
                
                
            
        }

}


// middle relay functionality
func middle_layer_clientHandling(w http.ResponseWriter, r *http.Request) { 
        requestedUrl :=r.URL.Path[13:] // request URL from starting relay
	fmt.Println("middle laye conecting..."+requestedUrl)// User requested page
         
        
        url:="http://"+s_IP+":"+server_port+"/middleToServer/User Requested URL: "+requestedUrl
	req, _ := http.NewRequest("GET",url , nil)
        client := &http.Client{}
        resp, _ := client.Do(req)
        fmt.Println(resp.Status)
       
        
        
        
        if resp.StatusCode == 200 { // OK
         fmt.Println("In exit code two")
            bodyBytes, _ := ioutil.ReadAll(resp.Body)
            bodyString := string(bodyBytes)
            fmt.Println("response from server for exit layer: "+bodyString)
            
            // Code to send request to middle layer
                ip:=strings.Split(bodyString, ":")[0]
                port:=strings.Split(bodyString, ":")[1]
                
                
                
            	url:="http://"+ip+":"+port+"/Exitlayer/"+requestedUrl
	        req, _ := http.NewRequest("GET",url , nil)
                client := &http.Client{}
                resp, _ := client.Do(req)
                fmt.Println(resp.Status)
                
                bodyBytes1, _ := ioutil.ReadAll(resp.Body)          // getting response from exit layer
                bodyString1 := string(bodyBytes1)
                fmt.Println("body from exit=>"+bodyString1)
                
                fmt.Fprintf(w, string(bodyString1))                // wrinting to entry node
      
        }

}

func Exit_layer_clientHandling(w http.ResponseWriter, r *http.Request) {

                timeout := time.Duration(20 * time.Second)
                client := http.Client{
                Timeout: timeout,
                }

                

                response, err := client.Get("http://info.cern.ch/hypertext/WWW/TheProject.html")
                if err != nil {
                fmt.Printf("%s", err)
                os.Exit(1)
                } else {
                defer response.Body.Close()
                contents, err := ioutil.ReadAll(response.Body)
                if err != nil {
                fmt.Printf("%s", err)
                os.Exit(1)
                }
                fmt.Printf("%s\n", string(contents))
                fmt.Fprintf(w, "<html>"+string(contents)+"</html>")  
                }
}



func main() {

	num:=randomNumber(99)	
	port := ":"+strconv.Itoa(num)
	fmt.Println("System listening at=> "+c_IP+port)

	if handShake(port,defaultPreference(randomNumber(2)-2000))=="200 OK" {
		http.HandleFunc("/", defaultHandler)
		http.HandleFunc("/fastor/", Entry_ClientToServer)
		http.HandleFunc("/middlelayer/", middle_layer_clientHandling)
		http.HandleFunc("/Exitlayer/", Exit_layer_clientHandling)

		http.ListenAndServe(c_IP+port, nil)
		
	}	
}
