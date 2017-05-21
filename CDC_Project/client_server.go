package main 

import (
	"fmt"
	"math/rand"
    "time"
	"net/http"
	"strconv"
	"strings"
	"html/template"
	"io/ioutil"
	"os"
	"bufio"	
)

type Page struct {
	Title string
}

var s_IP="127.0.0.1"
var c_IP="127.0.0.1"
var server_port= "8081"
var globe_port = ""
var user_preference = ""
var visited_URL_links = make(map[string]string)
var client_Port = ""

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	if strings.TrimRight(r.URL.Path[1:], "\n") != "" {

		fmt.Println("URL exists")
		fmt.Println(client_Port) 
		
		fmt.Println("-------------------------------------")


	} 	else {
		//fmt.Println("URL doesn't exixts")
	}


	//fmt.Fprintf(w, "Hi there, This is defult page %s!", r.URL.Path[1:])
}

func randomNumber( rangee int) int {

	rand.Seed(time.Now().UTC().UnixNano())
	num:=rand.Intn(rangee)
	num=num+2000
	return num

}

func handShake(port string, preference string ) string {
        
    url:="http://"+s_IP+":8081/saveConfig/"+c_IP+port+","
        
	if strings.TrimRight(preference, "\n") == "1" {
                url = url+ "Mid"
	} else if strings.TrimRight(preference, "\n") == "2" {
		 url = url+ "Exit"
	} else if strings.TrimRight(preference, "\n") == "3" { 
	         url = url+ "None"
	}

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

func browserClient_ConfigHandler(w http.ResponseWriter, r *http.Request) {
        //requestedUrl :=r.URL.Path[8:]
	//fmt.Println(requestedUrl)

	loadedPage := &Page{Title: "hell"}
	fmt.Println(r.RemoteAddr)

	editTemplate, _ := template.ParseFiles("config.html")
	editTemplate.Execute(w, loadedPage)       

}

func client_Request(w http.ResponseWriter, r *http.Request) {
           handShake(globe_port, user_preference) // default preference        
}

func makeRequest(ip,port,sender,requestedUrl string) (string,string,int){

        url:="http://"+ip+":"+port+"/"+sender+"/User Requested URL: "+requestedUrl
	req, _ := http.NewRequest("GET",url , nil)
        client := &http.Client{}
        resp, _ := client.Do(req)
        fmt.Println(resp.Status)
        
        defer resp.Body.Close()
        
        if resp.StatusCode == 200 { // OK
                bodyBytes, _ := ioutil.ReadAll(resp.Body)
                bodyString := string(bodyBytes)
                ip:=strings.Split(bodyString, ":")[0]
                port:=strings.Split(bodyString, ":")[1]
                return ip,port,resp.StatusCode
        } else{
                return "","",0
        }

}


func Entry_ClientToServer(w http.ResponseWriter, r *http.Request) {
    requestedUrl :=r.URL.Path[8:]
    fmt.Println(requestedUrl)
    // Irtiza CODE
    if len(parseURL_back(requestedUrl)) > 0 { // to handle www.google.com like sites
    	visited_URL_links[parseURL_back(requestedUrl)] = requestedUrl	
    }

    for _,val:= range(visited_URL_links){
     
           fmt.Println("*"+val)
    }

    url,_ := parseURL_front(requestedUrl)
    if val, ok := visited_URL_links[url]; ok {
		fmt.Println("asdasdasdasd  "+val)
		requestedUrl = combine_URL(val,requestedUrl)
	}

    //visited_URL_links[] 
    // Irtiza CODE END



	//fmt.Println(requestedUrl)// User requested page

	ip,port,_:=makeRequest(s_IP,server_port,"EntryToServer",requestedUrl)
    
        if strings.TrimRight(ip, "\n") != "" {
        
                for status := "" ; status !="200 OK" ; {
                
                        // time.Sleep(5 * time.Second)
                        fmt.Println("Assign Port=>"+port)
                        url:="http://"+ip+":"+port+"/middlelayer/"+requestedUrl
	                	req, _ := http.NewRequest("GET",url , nil)
                        client := &http.Client{}
                        resp, err := client.Do(req)
                        fmt.Println("-------------------------------------")
                        
                        if err == nil {
                                bodyBytes1, _ := ioutil.ReadAll(resp.Body)          // getting response from mid layer
                                bodyString1 := string(bodyBytes1)
                                fmt.Println("body from mid=>"+bodyString1)
                                fmt.Fprintf(w, string(bodyString1)) 
                                status ="200 OK"
                        } else {
                                status =""
                                 fmt.Println("Re new Requesting for Middle port ip...")
                                ip,port,_ =makeRequest(s_IP,server_port,"EntryToServer",requestedUrl)
                        }
                }      
            
        } else{
                fmt.Fprintf(w, string("No Middle Node in sever list")) 
        }
        
}


// middle relay functionality
func middle_layer_clientHandling(w http.ResponseWriter, r *http.Request) { 
        requestedUrl :=r.URL.Path[13:] // request URL from starting relay
	fmt.Println("middle laye conecting..."+requestedUrl)// User requested page
         
        ip,port,_:=makeRequest(s_IP,server_port,"middleToServer",requestedUrl) // Getting Exit lyaer port ip from server
        
        if strings.TrimRight(ip, "\n") != "" {
        
                for status := "" ; status !="200 OK" ; {
                
                        //fmt.Println("Waiting...")
                        //time.Sleep(5 * time.Second)
                        fmt.Println("Assign Port for exit=>"+port)
                        url:="http://"+ip+":"+port+"/Exitlayer/"+requestedUrl
	                req, _ := http.NewRequest("GET",url , nil)
                        client := &http.Client{}
                        resp, err := client.Do(req)
                        fmt.Println("-------------------------------------")
                        
                        if err == nil {
                
                                bodyBytes, _ := ioutil.ReadAll(resp.Body)          // getting response from exit layer
                                bodyString := string(bodyBytes)
                                fmt.Fprintf(w, string(bodyString))                // wrinting to entry node
                                status ="200 OK"
                                
                        } else {
                                status =""
                                fmt.Println("Re new Requesting for exit port ip...")
                                ip,port,_ =makeRequest(s_IP,server_port,"middleToServer",requestedUrl)
                        }
                }      
            
        } else{
                fmt.Fprintf(w, string("No Exit Node in sever list")) 
        }

}

func Exit_layer_clientHandling(w http.ResponseWriter, r *http.Request) {

                requestedUrl :=r.URL.Path[11:] 
                timeout := time.Duration(20 * time.Second)
                client := http.Client{
                Timeout: timeout,
                }
                req_url := ""        
                b := [6] string {"html", "htm", "aspx", "php" , "jsf" ,"jsp" }
                if contains_extension( b ,parseURL(requestedUrl)) {
                        fmt.Println("Extension exists")
                        req_url ="http://"+requestedUrl
                } else {
                        fmt.Println("Extension doesn't exists")
                        req_url ="http://"+requestedUrl+"/"
                }          

                //  response, err := client.Get("http://help.websiteos.com/websiteosexample_of_a_simple_html_page.htm")
                //"http://"+requestedUrl+"/")
                //response, err := client.Get("http://www.google.com/")
               
               
                fmt.Println("user req url=>"+req_url+"-"+string(len(req_url)))
                response, err := client.Get(req_url)
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
              //  fmt.Printf("%s\n", string(contents))
                fmt.Fprintf(w, "<html>"+string(contents)+"</html>")  
                }
}

func contains_extension(s [6]string, e string) bool { // function to check element exists in slice
    for _, a := range s {
        if strings.TrimRight(a, "\n") == e {
            return true 
        }
    }
    return false
}

func parseURL (URL string) string {
        list:=strings.Split(URL, ".")
        size:=len(list)
        return (list[size-1])
}


func Handler(w http.ResponseWriter, r *http.Request) {  
    fmt.Printf("Req: %s %s", r.URL.Host, r.URL.Path)
}


// Irtiza CODE START

func parseURL_front (URL string) (string, int) {
        list:=strings.Split(URL, "/")
        size:=len(list)
        size = size
        return (list[0]),len(list[0])
}
func parseURL_back (URL string) string {

        list:=strings.Split(URL, "/")
        size:=len(list)
        size = size
        return (list[size-1])

}

func compare_URLS ( old_URL string , new_URL string) bool{

	if strings.TrimRight(old_URL, "\n") == new_URL {
            return true 
    }
    return false
}

func combine_URL ( old_URL string, new_URL string) string {

	url,len :=parseURL_front (new_URL)
	if compare_URLS ( parseURL_back(old_URL) , url ) {
		return old_URL+new_URL[len:]
	}
	return old_URL

}



// Irtiza CODE END




func main() {

    // Getting user input
    fmt.Print("Enter you preference:\n 1. Middle \n 2. Exit \n 3. None\n => ")
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    user_preference = text   
                        

        
	num:=randomNumber(99)	
	port := ":"+strconv.Itoa(num)
	fmt.Println("System listening at=> "+ c_IP+port )
	client_Port = port
	globe_port=port
	handShake(globe_port,user_preference)
	
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/fastor/", Entry_ClientToServer)
	http.HandleFunc("/saveConfig/", client_Request)
	http.HandleFunc("/middlelayer/", middle_layer_clientHandling)
	http.HandleFunc("/Exitlayer/", Exit_layer_clientHandling)
	http.ListenAndServe(c_IP+port, nil)	
}
