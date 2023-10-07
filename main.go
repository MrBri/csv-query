package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/surrealdb/surrealdb.go"
)

var tmpl *template.Template
var db *surrealdb.DB
var port = ":3333"

func init() {
	var err error
	tmpl, err = template.New("upload_template.html").ParseFiles("upload_template.html")
	if err != nil {
		panic(err)
	}
	db, err = surrealdb.New("ws://surrealdb:8000/rpc")
	if err != nil {
		panic(err)
	}

	_, err = db.Signin(map[string]interface{}{"user": "root", "pass": "root"}) // for production should be imported from a secure place, env, parameter store etc
	if err != nil {
		panic(err)
	}

	if _, err = db.Use("dev", "weather"); err != nil { // namespace then database
		panic(err)
	}
}

func main() {
	http.HandleFunc("/upload", uploadFileHandler)
	http.HandleFunc("/query", queryHandler)
	fmt.Println("Server started at", port)

	err := http.ListenAndServe(port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		start := time.Now()
		// Parse the form data to retrieve the file
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get the file from the form
		file, header, err := r.FormFile("csvfile")
		if err != nil {
			http.Error(w, "Unable to get the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		log.Printf("Uploaded header: %v\n", header.Header)

		// Stream and parse the CSV
		csvReader := csv.NewReader(file)
		_, err = csvReader.Read() // Headers
		if err != nil {
			http.Error(w, "Unable to read the CSV headers", http.StatusInternalServerError)
			return
		}

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "Error reading a CSV record", http.StatusInternalServerError)
				return
			}

			_, err = db.Create("weather", NewWeatherRecord(record))
			if err != nil {
				panic(err)
			}
		}
		elapsed := time.Since(start)

		output, err := db.Query("select count() from weather group all", map[string]interface{}{})
		log.Println(output)
		outSlice := output.([]interface{})
		var outMap map[string]interface{}
		for _, o := range outSlice {
			outMap = o.(map[string]interface{})
			fmt.Println(outMap)
		}
		outMap["file_db_insert_time"] = fmt.Sprint(elapsed)
		temp := outMap["time"]
		delete(outMap, "time")
		outMap["count_time"] = temp


		jsonOutput, err := json.MarshalIndent(outMap, "", "  ")
		if err != nil {
			http.Error(w, "Unable to convert to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonOutput)
	} else {
		// Render the upload form using the template
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	if d, ok := queryParams["date"]; ok {
		log.Println(queryParams, d)
		out, err := db.Select(fmt.Sprintf("weather:⟨%s⟩", d[0]))
		if err != nil {
			http.Error(w, "Unable to query the db", http.StatusInternalServerError)
			return
		}
		jsonOutput, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			http.Error(w, "Unable to convert to JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonOutput)
		return
	}
		
	limit, okl := queryParams["limit"] 	
	weather, okw := queryParams["weather"]	
	log.Println(limit, weather, okl, okw)
	queryStr := "select * from weather"
	if okw {
		queryStr = queryStr + fmt.Sprintf(" where weather = '%s'", weather[0])
	}
	log.Println(queryStr)
	if okl {
		queryStr = queryStr + fmt.Sprintf(" limit %s", limit[0])
	}
	log.Println(queryStr)
	out, err := db.Query(queryStr, map[string]interface{}{})
	if err != nil {
		http.Error(w, "Unable to query the db", http.StatusInternalServerError)
		return
	}

	jsonOutput, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		http.Error(w, "Unable to convert to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonOutput)
}
