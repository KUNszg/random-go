package main

// import dependencies
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// random.org api key
var apiKey string

var apiURL = "https://api.random.org/json-rpc/4/invoke"

func getRes(requests string, length int, randomID int) (response string, err error) {

	// define structs for POST request payload
	// and let the json parser know how the key names should look like
	type Params struct {
		ApiKey                    string `json:"apiKey"`
		N                         int    `json:"n"`
		Min                       int    `json:"min"`
		Max                       int    `json:"max"`
		Replacement               bool   `json:"replacement"`
		Base                      int    `json:"base"`
		PregeneratedRandomization any    `json:"pregeneratedRandomization"`
	}

	type Payload struct {
		Jsonrpc string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  Params `json:"params"`
		Id      int    `json:"id"`
	}

	// assign values to struct
	payload := &Payload{
		Jsonrpc: "2.0",
		Method:  "generateIntegers",
		Params: Params{
			ApiKey:                    apiKey,
			N:                         length,
			Min:                       0,
			Max:                       10,
			Replacement:               true,
			Base:                      10,
			PregeneratedRandomization: nil,
		},
		Id: randomID,
	}

	// parse the struct data into byte data
	body, _ := json.Marshal(payload)

	// set the timeout timer for 5 seconds
	client := &http.Client{Timeout: 5 * time.Second}

	// send the POST request with the created payload
	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(body))

	// handle API timeout
	if err != nil {
		if os.IsTimeout(err) {
			verbose := "request timed out"
			return strconv.Itoa(http.StatusRequestTimeout) +
				":" + verbose, fmt.Errorf(verbose+" %s", err)
		}
		verbose := "failed to finish the request"
		return strconv.Itoa(http.StatusInternalServerError) +
			":" + verbose, fmt.Errorf(verbose+" %s", err)
	}

	// close the request at the end of function
	defer resp.Body.Close()

	// handle response code
	if resp.StatusCode == http.StatusOK {
		// read the response as []byte slice
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			verbose := "error reading data"
			return strconv.Itoa(http.StatusInternalServerError) +
				":" + verbose, fmt.Errorf(verbose+" %s", err)
		}

		// convert the []byte slice to string
		jsonStr := string(body)

		return jsonStr, nil
	} else {
		verbose := "failed to finish the request"
		return strconv.Itoa(http.StatusInternalServerError) +
			":" + verbose, fmt.Errorf(verbose+" %s", err)
	}
}

func standardDeviation(arr []int) (result float64) {
	var sum, mean, stddev float64

	// get the slice length
	var n = len(arr)

	// iterate to sum up all the values in provided int slice
	for i := 1; i <= n; i++ {
		sum += float64(arr[i-1])
	}

	// divide the summed up values by the number of values
	mean = sum / float64(n)

	// use the standard deviation formula part (x - s)^2 + (x - s)^2 (...)
	for j := 0; j < n; j++ {
		stddev += math.Pow(float64(arr[j])-mean, 2)
	}

	stddev = math.Sqrt(stddev / float64(n))

	// round up the result to 5 decimal points
	ratio := math.Pow(10, float64(5))
	stddev = math.Round(stddev*ratio) / ratio

	return stddev
}

func main() {
	// read the random.org API key from config file
	content, err := ioutil.ReadFile("config.txt")

	if err != nil {
		log.Fatal(err)
	}

	// stringify the file contents
	apiKey = string(content)

	// create a handler for REST listener
	router := mux.NewRouter()

	// listen on /random/mean path
	router.HandleFunc("/random/mean", func(w http.ResponseWriter, r *http.Request) {

		// set the format of response data
		w.Header().Set("Content-Type", "application/json")

		resp := make(map[string]any)

		length, err := strconv.Atoi(r.URL.Query().Get("length"))

		requests := r.URL.Query().Get("requests")
		_requests, err := strconv.Atoi(requests)

		// handle user errors
		if length < 3 {
			w.WriteHeader(http.StatusBadRequest)
			resp["status"] = strconv.Itoa(http.StatusBadRequest)
			resp["message"] = "length value has to be more than 2"
		}

		if _requests < 1 {
			w.WriteHeader(http.StatusBadRequest)
			resp["status"] = strconv.Itoa(http.StatusBadRequest)
			resp["message"] = "requests parameter value has to be more than 0"
		}

		if len(apiKey) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			resp["status"] = strconv.Itoa(http.StatusUnauthorized)
			resp["message"] = "no random.org API key was provided, check the app readme for proper configuration"
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["status"] = http.StatusInternalServerError
			resp["message"] = "error parsing the request parameters"
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["status"] = http.StatusInternalServerError
			resp["message"] = "error parsing the request parameters"
		}

		if len(resp) > 0 {
			json.NewEncoder(w).Encode(resp)
			log.Println(resp["message"])
			return
		}

		// get a random ID for random.org api requirements
		randomID := rand.Int()

		// create a struct for response from api
		type Res struct {
			Jsonrpc string                      `json:"jsonrpc"`
			Result  map[string]map[string][]int `json:"result"`
			Id      int                         `json:"id"`
		}

		resultMap := make(map[int][]int)

		// concurrent API calls using wait groups
		wg := sync.WaitGroup{}

		for i := 0; i < _requests; i++ {
			wg.Add(1)
			go func(i int) {
				result, err := getRes(requests, length, randomID)

				if err != nil {
					// get the error status codes and message
					_result := strings.Split(result, ":")
					if _result[0] == "500" || _result[0] == "408" {
						status, _ := strconv.Atoi(_result[0])

						w.WriteHeader(status)

						resp["status"] = status
						resp["message"] = _result[1]

						json.NewEncoder(w).Encode(resp)
						log.Println(resp["message"])
					}
					return
				}

				// parse the bytes into json
				var res *Res
				json.Unmarshal([]byte(result), &res)

				resultMap[i] = res.Result["random"]["data"]

				wg.Done()
			}(i)
		}
		wg.Wait()

		// create a struct for app's JSON response
		type Result struct {
			Stddev float64 `json:"stddev"`
			Data   []int   `json:"data"`
		}

		result := map[int]Result{}
		allElements := make(map[int][]int)

		for key, element := range resultMap {
			// put the received and calculated data into a slice
			result[key] = Result{standardDeviation(element), element}

			allElements[0] = append(allElements[0], element...)
		}

		// add the standard deviation of all elements received
		result[len(result)] = Result{standardDeviation(allElements[0]), allElements[0]}

		// send the data in JSON format
		json.NewEncoder(w).Encode(result)
	})

	// app listens on port 8080
	http.ListenAndServe(":8080", router)
}
