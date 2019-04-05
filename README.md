# FASTTOR	

A web application built in golang that can be used for anonymous browsing. This application is based on peer 2 peer architecture.


## How to run the application
This application can be run using the steps given below:

### Steps

Run the cdc_v1.go file using the "go run cdc_v1.go". 

- It is the main server file. 
- It will maintains a registry for IPs for those clients that will connect to it. 
- Those IPs will be used for requesting web pages.

-  Run the client_server.go file using the "go run client_server.go". 

-  It will used for the user preferences.
-  Preferences like whether a user wants to act as a server of not.

-  Once application starts to run , open the browser and enter the url of the webpage that you want to fetch in this way:

- "ip:port/fastor/webPageURL. 

- It will display the webpage anonymously.

