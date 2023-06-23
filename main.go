package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type StringData struct {
	Value string `json:"value"`
	ID    string `json:"id"`
}

type IntegerData struct {
	Value int    `json:"value"`
	ID    string `json:"id"`
}

var stringArr []StringData
var integerArr []IntegerData
var payloadArr []interface{}

func createData(w http.ResponseWriter, r *http.Request) {
	// Handler for creating string and integer data from the payload
	var payload struct {
		String  string `json:"string"`
		Integer int    `json:"integer"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request payload")
		return
	}

	// Generate unique IDs for string and integer data
	stringID := uuid.New().String()
	stringData := StringData{
		Value: payload.String,
		ID:    stringID,
	}
	stringArr = append(stringArr, stringData)

	integerID := uuid.New().String()
	integerData := IntegerData{
		Value: payload.Integer,
		ID:    integerID,
	}
	integerArr = append(integerArr, integerData)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Data created successfully")
}

func getData(w http.ResponseWriter, r *http.Request) {
	// Handler for retrieving saved string and integer data
	response := struct {
		Strings  []StringData `json:"strings"`
		Integers []IntegerData `json:"integers"`
	}{
		Strings:  stringArr,
		Integers: integerArr,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to retrieve data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func savePayload(w http.ResponseWriter, r *http.Request) {
	// Handler for saving incoming payloads as is
	var payload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request payload")
		return
	}

	// Generate unique ID for the payload
	payloadID := uuid.New().String()
	payload["id"] = payloadID
	payloadArr = append(payloadArr, payload)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Payload saved successfully")
}

func getPayloads(w http.ResponseWriter, r *http.Request) {
	// Handler for retrieving saved payloads
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(payloadArr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to retrieve payloads")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {

    http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createData(w, r)
		case http.MethodGet:
			getData(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

    http.HandleFunc("/payloads", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			savePayload(w, r)
		case http.MethodGet:
			getPayloads(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

/*
/data endpoint is used for saving string and integer values seperately by attaching unique id along 
with values
to save save string and integer data seperately use POST http://localhost:8080/data
to retrieve save data use GET http://localhost:8080/data

saved data will look like this

{
    "strings": [
        {
            "value": "string7",
            "id": "51b08548-bf7c-451c-906a-d9f8a79d57ac"
        },
        {
            "value": "string8",
            "id": "913ea9a3-44a0-457f-8d3e-6d1e1b07ee92"
        }
    ],
    "integers": [
        {
            "value": 55,
            "id": "23ac3735-9538-4df3-9bd6-b377a1583a90"
        },
        {
            "value": 56,
            "id": "6874b7f4-a01b-4c7b-b344-c9ebd26bec0c"
        }
    ]
}

to save string and data together with unique id payloads endpoint will be used

POST http://localhost:8080/payloads for saving payload with id attached to it
GET http://localhost:8080/payloads to fetch all saved payloads

saved data will look like

[
    {
        "id": "05ef4e10-ef8f-4bf3-8063-f83dcfcf4987",
        "integer": 55,
        "string": "string4"
    },
    {
        "id": "f75faaf9-e5a1-4ff2-a08f-5e03f343a5a3",
        "integer": 55,
        "string": "string5"
    },
    {
        "id": "eae3b489-f4be-4c2d-9d95-51c0b88b5846",
        "integer": 55,
        "string": "string7"
    }
]
*/