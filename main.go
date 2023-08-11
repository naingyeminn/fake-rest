package main

import (
	"fmt"
 	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log access request details to console
		log.Printf("Received request: %s %s, Source: %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Retrieve the desired HTTP status code from an environment variable
		statusCodeEnv := os.Getenv("HTTP_STATUS_CODE")
		statusCode, err := strconv.Atoi(statusCodeEnv)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Printf("Error parsing HTTP_STATUS_CODE: %v", err)
			return
		}

		// Retrieve the desired error percentage from an environment variable
		errorPercentageEnv := os.Getenv("ERROR_PERCENTAGE")
		errorPercentage, err := strconv.Atoi(errorPercentageEnv)
		if err != nil {
			errorPercentage = 0 // Default to 0 if environment variable is not set or invalid
		}

		// Determine if an error should be simulated based on the error percentage
		simulateError := false
		if errorPercentage > 0 && errorPercentage <= 100 {
			randomNum := rand.Intn(100) + 1
			if randomNum <= errorPercentage {
				simulateError = true
			}
		}

		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Error reading request body: %v", err)
				return
			}
			defer r.Body.Close()
			log.Printf("Received POST data: %s", string(body))
    }
      
		delay := 0 // Default delay of 0 seconds
		delayEnv := os.Getenv("DELAY_SECONDS")
		if delayEnv != "" {
			delayInt, err := strconv.Atoi(delayEnv)
			if err == nil {
				delay = delayInt
			}
		}

		time.Sleep(time.Duration(delay) * time.Second)

		if simulateError {
			http.Error(w, http.StatusText(statusCode), statusCode)
			log.Printf("Simulated error due to error percentage: %d%%, HTTP_STATUS_CODE: %d", errorPercentage, statusCode)
		} else {
			responseMessage := "Success"
			responseMessageEnv := os.Getenv("HTTP_RESPONSE_MESSAGE")
			if responseMessageEnv != "" {
				responseMessage = responseMessageEnv
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, responseMessage)
		}
	})

	port := 8080
	fmt.Printf("Server is running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

