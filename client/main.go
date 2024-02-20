package main

import (
	"fmt"
	"net/http"
)

type Pet struct {
	Name   string
	Health int
	Img    int
}

func main() {
	// Serve static files from the 'assets' directory
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the pet name from the URL query parameter "name"
		petName := r.URL.Query().Get("name")
		if petName == "" {
			petName = "nina" // Default name if none is provided
		}

		pet := Pet{Name: petName, Health: 80, Img: 4} // Initialize pet with dynamic name

		// Generate HTML with embedded pet name and health value
		fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
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
				font-size: 24px;
				margin-bottom: 20px;
			}
		</style>
	</head>
	<body>
		<div class="pet-name">tiagochi: %s</div>
		<img src="/assets/%d.png" alt="Pet Image" class="pet-image" />
		<b>Health: %d%%</b>
		<div id="health-bar-container">
			<div id="health-bar"></div>
		</div>
	</body>
	</html>
	`, pet.Health, pet.Name, pet.Img, pet.Health)
	})

	http.ListenAndServe(":8080", nil)
}
