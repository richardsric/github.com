package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"html/template"

	_ "github.com/lib/pq"
)

var db *sql.DB
var temp *template.Template

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://amobi:iyochu@ipaytsa.com/ipaytsa_gov?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to database succesfully.")

	//temp = template.Must(template.ParseFiles("home.gohtml"))
}

type mfbs struct {
	// This holds single Service ID and a slice of mfb names
	SerID   string
	MfbList []string
}

type general struct {
	//This holds List Of Commercial Bank names
	CommercialList []string
}
type tempdata struct {
	Gen general
	Mf  []mfbs
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8085", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	servIDRows, errSid := db.Query("select DISTINCT service_id from _paypoints where service_id is NOT NULL")
	comRows, errCr := db.Query("select paypoint_name from _paypoints where paypoint_type = '0'")
	fmt.Println("Executing Service ID retrieval Query.")
	if errSid != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Println("Executing Commercial Bank retrieval Query.")
	if errCr != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	defer servIDRows.Close()
	defer comRows.Close()

	commNames := make([]string, 0) // this is slice of strings to hold Commercial Bank List
	var commName string            //This holds one commercial bank value
	var comm general
	var mfbName string //This holds one MFB bank value
	MfbList := make([]mfbs, 0)
	var mfbconstruct mfbs         //holds one mfb construct
	mfbNames := make([]string, 0) // this is slice of strings to hold MFBank List

	fmt.Println("Looping Through Commercial Bank Rows.")
	for comRows.Next() {

		err1 := comRows.Scan(&commName)
		if err1 != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		commNames = append(commNames, commName)
		fmt.Println(commNames)

		comm = general{commNames} // this now builds the general struct
	}
	fmt.Println(comm)

	//the loop that fetches each service_id from db starts here
	var ServId string //use this to hold the service ID
	fmt.Println("looping through Service ID Rows.")
	for servIDRows.Next() {

		err1 := servIDRows.Scan(&ServId) //store on the ServId field on each fetch
		if err1 != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		//Use ServId to retrieve MFBs in the service ID.
		mfbRows, errmfb := db.Query("SELECT paypoint_name FROM _paypoints WHERE service_id = '" + ServId + "' ORDER BY paypoint_name")
		//use what is fetched to query for the paypoint_names
		if errmfb != nil {
			println("the second query didn't run successfully")
			//http.Error(w, http.StatusText(500), 500)
			return
		}
		fmt.Println("Executing MFB retrieval Query using " + ServId)
		//loop through the paypoint_names gotten and append to slice of MfbList field in my Mfbs Struct
		fmt.Println("Looping through MFB rows.")
		for mfbRows.Next() {

			err3 := mfbRows.Scan(&mfbName)
			if err3 != nil {
				println("the third query didn't run successfully")
				//http.Error(w, http.StatusText(500), 500)
				return
			}
			mfbNames = append(mfbNames, mfbName) //the append happends here

		}
		//End for mfbNames
		fmt.Fprintln(os.Stdout, "MFBs for Service ID", ServId, "\n", mfbNames)
		fmt.Println("Constructing MFB data type")
		mfbconstruct = mfbs{
			ServId,
			mfbNames}
		MfbList = append(MfbList, mfbconstruct)
		fmt.Fprintln(os.Stdout, MfbList)

	} // the end  for the service id loop.

	data := tempdata{
		comm,
		MfbList,
	}
	//Printing data to be parsed to template
	fmt.Fprintln(os.Stdout, "Template Data is:\n", data)
	//Write to http ResponseWriter
	fmt.Fprintln(w, "Template Data is:\n", data)
	jsondata, _ := json.Marshal(data)
	fmt.Fprintln(w, "Template JSON Data is:\n", string(jsondata))

	fmt.Fprintln(w, "System Data is:\n", os.Stdout)
}
