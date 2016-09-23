package Handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	controllers "github.com/SivaShhankar/CMS_NEW/Controllers"
	config "github.com/SivaShhankar/CMS_NEW/Database"
	models "github.com/SivaShhankar/CMS_NEW/Models"
)

var message = ""

type Info struct {
	BEditMode  bool
	Details    []models.ApplicantInfo
	YearOfExp  string
	MonthOfExp string
	UserMsg    string
	Operation  string
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("Templates/Index.html")
	t.Execute(w, nil)
}

/* Upload the candidate information based on the mode.append
   If the mode is insert -> then add the new record into the list.
   If the mode is edit -> then update the existing record in the list.
*/
func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		message = ""
		t, _ := template.ParseFiles("Templates/Upload.html")

		d := Info{
			BEditMode:  false,
			Details:    nil,
			YearOfExp:  "",
			MonthOfExp: "",
			UserMsg:    "",
			Operation:  "Insert",
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

		controllers.SaveInfo(config.Session, r, mode)

		d := Info{
			BEditMode:  false,
			Details:    nil,
			YearOfExp:  "",
			MonthOfExp: "",
			UserMsg:    "Candidate Details Updated Successfully!",
			Operation:  "Insert",
		}

		t, _ := template.ParseFiles("Templates/Upload.html")
		t.Execute(w, d)
	}
}

// Delete the existing candidate based on the given mobile number in the candidate details.
func Delete(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	mobile := r.FormValue("mobileNumber")

	err := controllers.DeleteCandidateByMobileNumber(config.Session, mobile)

	if err != nil {
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, "/View", http.StatusSeeOther)
}

func View(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("Templates/ViewCandidates.html")
	a := controllers.GetAllApplicantsInfo(config.Session)
	t.Execute(w, a)
}

/* Search the candidates by different criteria with the help of some parameters.
Example the following:
1. Name => wants to search the candidate by name like "contains character or full name".
2. qualification => wants to search the candidate by qualification like "B.C.A or MCA, etc..."
......
*/
func Search(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	searchType := r.FormValue("searchType")
	searchValue := r.FormValue("searchBox")

	fmt.Println(searchType, searchValue)

	t, _ := template.ParseFiles("Templates/ViewCandidates.html")
	candidateDetails := controllers.SearchCandidatesByType(config.Session, searchType, searchValue)
	t.Execute(w, candidateDetails)
}

/* Filter the candidates by different criteria with the help of range parameters.
Example the following:
1. Age => wants to filter the candidate by age between 23 to 34.
2. Experience => wants to filter the candidate by the experience between 2 to 8.
*/

func Filter(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	filterType := r.FormValue("filterType")
	from := r.FormValue("from")
	to := r.FormValue("to")

	t, _ := template.ParseFiles("Templates/ViewCandidates.html")
	candidateDetails := controllers.FilterCandidatesByRange(config.Session, filterType, from, to)

	t.Execute(w, candidateDetails)
}

// Edit the existing candidate information
func Edit(h http.ResponseWriter, r *http.Request) {
	mobileno, _ := strconv.Atoi(r.FormValue("mobileNumber"))
	r.URL.RawQuery = ""
	t, _ := template.ParseFiles("Templates/Upload.html")

	candidateDetails := controllers.GetApplicantByMobileNumber(config.Session, mobileno)

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

		tempExp := candidateDetails[0].Experience

		strExp := strconv.FormatFloat(tempExp, 'f', 1, 64)

		experience := strings.Split(strExp, ".")

		d := struct {
			BEditMode  bool
			Details    models.ApplicantInfo
			YearOfExp  string
			MonthOfExp string
			UserMsg    string
			Operation  string
		}{

			BEditMode:  true,
			Details:    candidateDetails[0],
			YearOfExp:  experience[0],
			MonthOfExp: experience[1],
			UserMsg:    "",
			Operation:  "Update",
		}
		t.Execute(h, d)
	}
}
