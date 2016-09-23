package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/gorilla/mux"
	//"time"
	"strconv"

	config "github.com/SivaShhankar/CMS_NEW/Database"
	models "github.com/SivaShhankar/CMS_NEW/Models"

	model "github.com/SivaShhankar/CMS_NEW/Controllers"
)

var message = ""

// Info - struct for handling forms data
type Info struct {
	BEditMode bool
	Details   []models.ApplicantInfo
	UserMsg   string
	Operation string
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		message = ""
		t, _ := template.ParseFiles("Templates/Upload.html")

		d := Info{
			BEditMode: false,
			Details:   nil,
			UserMsg:   "",
			Operation: "Insert",
		}
		t.Execute(w, d)

	} else {

		mode := r.FormValue("mode")
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")

		if err != nil && mode == "Insert" {
			fmt.Println(err)
			return
		}

		if handler != nil {
			defer file.Close()

			if mode == "Update" {
				os.Remove("Templates/" + config.AppConfig.CVLocation + r.FormValue("uploadedFile"))
			}
			f, err := os.Create("Templates/" + config.AppConfig.CVLocation + r.FormValue("name") + "-" + r.FormValue("mobile") + "-" + handler.Filename)

			if err != nil {
				fmt.Println(err)
				return
			}

			defer f.Close()

			io.Copy(f, file)
		}

		fmt.Println("TEST", mode)
		model.SaveInfo(config.Session, r, mode)

		//message = "File Uploaded Successfully !"

		d := Info{

			BEditMode: false,
			Details:   nil,
			UserMsg:   "File Uploaded Successfully !",
			Operation: "Insert",
		}
		//t.Execute(w, d)

		t, _ := template.ParseFiles("Templates/Upload.html")
		t.Execute(w, d)
	}
}

func view(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("Templates/ViewCandidates.html")
	a := model.GetAllApplicantsInfo(config.Session)
	t.Execute(w, a)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("Templates/Index.html")
	t.Execute(w, nil)

}

func search(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	searchType := r.FormValue("searchType")
	searchValue := r.FormValue("searchBox")

	fmt.Println(searchType, searchValue)

	t, _ := template.ParseFiles("Templates/ViewCandidates.html")
	candidateDetails := model.SearchCandidates(config.Session, searchType, searchValue)
	t.Execute(w, candidateDetails)
}

func filter(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	filterType := r.FormValue("filterType")
	from := r.FormValue("from")
	to := r.FormValue("to")

	t, _ := template.ParseFiles("Templates/ViewCandidates.html")
	candidateDetails := model.FilterCandidates(config.Session, filterType, from, to)

	t.Execute(w, candidateDetails)
}

// EditData - Edit user data
func EditData(h http.ResponseWriter, r *http.Request) {
	//id, err := strconv.ParseInt(mux.Vars(r)["mobile"], 10, 64)
	mobileno, _ := strconv.Atoi(r.FormValue("mobileNumber"))
	r.URL.RawQuery = ""

	//fmt.Println(r.FormValue("mobileNumber"))
	t, _ := template.ParseFiles("Templates/Upload.html")

	candidateDetails := model.GetApplicantByMobileNumber(config.Session, mobileno)
	//fmt.Println(candidateDetails[0].Name)
	//t.Execute(h, candidateDetails)
	//var d
	if candidateDetails == nil {
		d := struct {
			BEditMode bool
			UserMsg   string
			Operation string
		}{

			BEditMode: false,
			UserMsg:   "No details found",
			Operation: "Insert",
		}
		t.Execute(h, d)

	} else {
		d := struct {
			BEditMode bool
			Details   models.ApplicantInfo
			UserMsg   string
			Operation string
		}{

			BEditMode: true,
			Details:   candidateDetails[0],
			UserMsg:   "",
			Operation: "Update",
		}
		t.Execute(h, d)
		//fmt.Println(candidateDetails[0].Name)
	}
	//fmt.Println(message)

}

// GetPort -- get the Port from the Dynamic environment
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	fmt.Println("Running Port is" + port)
	return ":" + port
}

func main() {

	config.LoadAppConfig()
	config.CreateDBSession()
	config.AddIndexes()

	mux := http.NewServeMux()

	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("Templates/css"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("Templates/images"))))

	mux.Handle("/JS/", http.StripPrefix("/JS/", http.FileServer(http.Dir("Templates/JS"))))
	mux.Handle("/Files/", http.StripPrefix("/Files/", http.FileServer(http.Dir("Templates/Files"))))
	mux.Handle("/Index", http.HandlerFunc(index))
	mux.Handle("/Upload", http.HandlerFunc(upload))
	mux.Handle("/View", http.HandlerFunc(view))
	mux.Handle("/Search", http.HandlerFunc(search))
	mux.Handle("/Filter", http.HandlerFunc(filter))
	mux.Handle("/", http.HandlerFunc(index))
	mux.Handle("/EditData", http.HandlerFunc(EditData)) //.Methods("GET")

	log.Println("Listening...")

	http.ListenAndServe(GetPort(), mux)
}
