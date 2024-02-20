package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Pet struct {
	Health int
	Img    int
}

type Resp struct {
	IsHealthier bool `json:"is_healthier"`
	Scale       int  `json:"scale"`
}

func main() {
	// Serve static files from the 'assets' directory
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	resp, err := http.Get("http://172.16.24.166:8000/health")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	var respTia Resp
	err = json.Unmarshal(body, &respTia)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the pet name from the URL query parameter "name"
		petName := r.URL.Query().Get("name")
		if petName == "" {
			petName = "stu" // Default name if none is provided
		}

		pet := Pet{Health: respTia.Scale, Img: respTia.Scale} // Initialize pet with dynamic name

		if respTia.Scale == 1 {
			pet.Health = 20
		} else if respTia.Scale == 2 {
			pet.Health = 40
		} else if respTia.Scale == 3 {
			pet.Health = 60
		} else if respTia.Scale == 4 {
			pet.Health = 80
		} else if respTia.Scale == 5 {
			pet.Health = 100
		}

		// Generate HTML with embedded pet name and health value
		fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	   <script>
	   // Refresh the page after a delay of 3 seconds
	   setTimeout(function(){
		   location.reload();
	   }, 30000); // 3000 milliseconds = 3 seconds 
	   </script>
    	<meta charset="UTF-8">
    <style>
			body, html {
				height: 100;
				margin: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				flex-direction: column;
				background-color: white; /* Set background color to white */
			}
			.pet-image {
				width: 250px; /* Fixed width for all images */
				height: auto; /* Maintain aspect ratio */
				margin: 20px 0;
			}
			#health-bar-container {
				width: 200px;
				height: 20px;
				border: 1px solid #000;
				display: flex; /* Use flexbox for alignment */
				justify-content: flex-start; /* Align health bar to the start */
				background-color: #eee; /* Background of the health bar */
			}
			#health-bar {
				width: %d%%;
				background-color: purple;
				height: 100;
			}
			.pet-name {
				margin-top: 100px;
				font-size: 24px;
			}
			.btn {
				margin: 10px;
			}
		</style>
	</head>
	<body>
		<div class="pet-name">stuiegotchi</div>
		<img src="/assets/%d.png" alt="Pet Image" class="pet-image" />
		<b>Health: %d%%</b>
		<div id="health-bar-container">
			<div id="health-bar"></div>
		</div>
		<button class="btn" onclick="window.location.href = '/?name=st'">Feed</button>
	  </div>
	</body>
	</html>
	`, pet.Health, pet.Img, pet.Health)
	})

	http.ListenAndServe(":8080", nil)
}
