package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		<html>
			<head>
				<title>Chat</title>
			</head>
			<body>
				Let's Chat!
			</body>
		</html>
		`))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
