package main 

import (
	"fmt"
	"net/http"
	"math/rand"
        "time"
	"strings"
)

var mid_client_list []Client
var exit_client_list []Client

type Page struct {
	Title string
}

type Client struct{
	_ip string
	_port string
	relayType string
}

func randomNumber( rangee int) int {

	rand.Seed(time.Now().UTC().UnixNano())
	num:=rand.Intn(rangee)
	return num
}

func saveClients(w http.ResponseWriter, req *http.Request) {            // Save clients in respective list

	req_parametre := req.URL.Path[len("/saveConfig/"):]
	
	ip:=strings.Split(req_parametre, ":")[0]
	port:=strings.Split(req_parametre, ":")[1]
	port=strings.Split(port, ",")[0]
	pref:=strings.Split(req_parametre, ",")[1]

	fmt.Println(ip+" "+port+ " "+pref)

	client_obj:=Client{_ip:ip,_port:port,relayType:pref}
	
	if pref=="Mid"{
		mid_client_list = append(mid_client_list, client_obj)
	} else if pref=="Exit"{
		exit_client_list = append(exit_client_list, client_obj)
	}
	
	
}



// for middle node list
func request_response_middle (w http.ResponseWriter, r *http.Request){ //this code handles clients requested URL and return response
       user_Requested_URL :=r.URL.Path[13:]
       fmt.Println("Mid url=>"+user_Requested_URL)
        
	size:=len(mid_client_list)
	if size !=0{
	         middle_client_index := randomNumber(size)                     
	         fmt.Fprintf(w,mid_client_list[middle_client_index]._ip+":"+mid_client_list[middle_client_index]._port) // edit this line
	}else{
	        fmt.Fprintf(w,":")
	}
	
        
        
}


// for exit node list
func request_response_exit (w http.ResponseWriter, r *http.Request){ //this code handles clients requested URL and return response
       
    user_Requested_URL :=r.URL.Path[14:]
	fmt.Println("exit url=>"+user_Requested_URL)
	size:=len(exit_client_list)
	if size !=0{
	         exit_client_index := randomNumber(size)                     
	         fmt.Fprintf(w,exit_client_list[exit_client_index]._ip+":"+exit_client_list[exit_client_index]._port) // edit this line 
	}else{
	        fmt.Fprintf(w,":")
	}
}




/************************************************** Refreshing Clients List ***************************************/


func UpdateClientLists() {
  for {
    time.Sleep(1 * time.Second)
    go UpdatingClientLists()
  }
}


func UpdatingClientLists() {

        var mid_temp_client_list []Client                       // Middle Temp list for updating

        for _,val:= range(mid_client_list){
     
           flag := makeRequest(val._ip,val._port)
           if flag == true {
                mid_temp_client_list = append(mid_temp_client_list, val)
           } 
        }
        mid_client_list=mid_temp_client_list
        
        var exit_temp_client_list []Client                      // Exit Temp list for updating

        for _,val:= range(exit_client_list){
     
           flag := makeRequest(val._ip,val._port)
           if flag == true {
                exit_temp_client_list = append(exit_temp_client_list, val)
           } 
        }
        exit_client_list=exit_temp_client_list
        
        for _,val:= range(mid_client_list){
     
           fmt.Println(val._port)
        }
         for _,val:= range(exit_client_list){
     
           fmt.Println(val._port)
        }
}
func makeRequest(ip string ,port string) bool {
    fmt.Println("New Request with port=>"+port)
   	url:="http://"+ip+":"+port+"/"
	req, _ := http.NewRequest("GET",url , nil)
	client := &http.Client{}
        resp, err := client.Do(req)
	if(err == nil){
		if resp.StatusCode == 200 {
			return true
        } else {
            return false
		}
	}
       return false
} // function end

func main() {

    go UpdateClientLists()
	http.HandleFunc("/EntryToServer/", request_response_middle)
	http.HandleFunc("/middleToServer/", request_response_exit)
	http.HandleFunc("/saveConfig/", saveClients)
	fmt.Println("Tor Server Is Up and Running... ")
	http.ListenAndServe(":8081", nil)
}
