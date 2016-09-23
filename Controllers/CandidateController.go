package Controllers

import (
	"fmt"
	"net/http"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"

	config "github.com/SivaShhankar/CMS_NEW/Database"
	models "github.com/SivaShhankar/CMS_NEW/Models"
)

func GetAllApplicantsInfo(session *mgo.Session) []models.ApplicantInfo {

	var b []models.ApplicantInfo

	dataStore := config.NewDataStore()

	// Close the session.
	defer dataStore.Close()
	col := dataStore.Collection("JobCandidates")
	candidates := Candidates{C: col}

	iter := candidates.C.Find(nil).Sort("Name").Iter()
	result := models.ApplicantInfo{}
	for iter.Next(&result) {
		b = append(b, result)
	}
	return b
}

func GetApplicantByMobileNumber(session *mgo.Session, MobileNumber int) []models.ApplicantInfo {
	var b []models.ApplicantInfo

	dataStore := config.NewDataStore()

	// Close the session.
	defer dataStore.Close()
	col := dataStore.Collection("JobCandidates")
	candidates := Candidates{C: col}

	//iter := candidates.C.Find(bson.M{"mobile": &bson.RegEx{Pattern: MobileNumber, Options: "i"}}).Sort("name").Iter()
	iter := candidates.C.Find(bson.M{"mobile": MobileNumber}).Iter()
	result := models.ApplicantInfo{}

	for iter.Next(&result) {

		b = append(b, result)
	}

	return b
}

type Candidates struct {
	C *mgo.Collection
}

func SearchCandidates(session *mgo.Session, searchType, searchValue string) []models.ApplicantInfo {

	var b []models.ApplicantInfo

	dataStore := config.NewDataStore()

	// Close the session.
	defer dataStore.Close()
	col := dataStore.Collection("JobCandidates")
	candidates := Candidates{C: col}

	iter := candidates.C.Find(bson.M{searchType: &bson.RegEx{Pattern: searchValue, Options: "i"}}).Sort("name").Iter()
	result := models.ApplicantInfo{}
	for iter.Next(&result) {
		b = append(b, result)
	}
	return b
}

func FilterCandidates(session *mgo.Session, filterType string, sFrom, sTo string) []models.ApplicantInfo {

	var b []models.ApplicantInfo

	dataStore := config.NewDataStore()

	from, _ := strconv.Atoi(sFrom)
	to, _ := strconv.Atoi(sTo)

	// Close the session.
	defer dataStore.Close()
	col := dataStore.Collection("JobCandidates")
	candidates := Candidates{C: col}

	iter := candidates.C.Find(bson.M{filterType: bson.M{"$gt": from - 1, "$lt": to + 1}}).Sort(filterType).Iter()

	result := models.ApplicantInfo{}

	for iter.Next(&result) {
		b = append(b, result)
	}

	fmt.Println(b)

	return b
}

func SaveInfo(session *mgo.Session, r *http.Request, mode string) {

	var err error
	name := r.FormValue("name")
	sage := r.FormValue("age")
	gender := r.FormValue("gender")
	smobile := r.FormValue("mobile")
	email := r.FormValue("email")
	location := r.FormValue("location")

	qualification := r.FormValue("qualification")
	specialization := r.FormValue("specialization")
	department := r.FormValue("department")
	position := r.FormValue("position")
	sExpMonth := r.FormValue("expMonth")
	sExpYear := r.FormValue("expYear")

	//fmt.Println("ex", sExpMonth)

	age, _ := strconv.Atoi(sage)
	mobile, _ := strconv.Atoi(smobile)
	expMonth, _ := strconv.Atoi(sExpMonth)
	expYear, _ := strconv.Atoi(sExpYear)

	experience := (float32)(expYear*12+expMonth) / 12

	fmt.Println("no ex", experience)

	_, handler, err := r.FormFile("file")
	var cvpath string

	// If no file has selected in the Form, it will throw an error
	// Cond 1 : if mode  is update, then retreive file value from hidden text box
	// Cond 2 : if mode is Insert, then retreive file value from file field

	if err != nil && mode == "Update" {
		cvpath = r.FormValue("uploadedFile")
	} else {
		cvpath = r.FormValue("name") + "-" + r.FormValue("mobile") + "-" + handler.Filename
	}

	dataStore := config.NewDataStore()

	// Close the session.
	defer dataStore.Close()

	col := dataStore.Collection("JobCandidates")
	candidates := Candidates{C: col}

	if mode == "Insert" {
		err = candidates.C.Insert(&models.ApplicantInfo{
			Name:           name,
			Age:            age,
			Gender:         gender,
			Mobile:         mobile,
			Email:          email,
			Location:       location,
			Qualification:  qualification,
			Specialization: specialization,
			Department:     department,
			Position:       position,
			Experience:     experience,
			CvPath:         cvpath,
		})

	} else if mode == "Update" {
		err = candidates.C.Update(bson.M{"mobile": mobile}, &models.ApplicantInfo{
			Name:           name,
			Age:            age,
			Gender:         gender,
			Mobile:         mobile,
			Email:          email,
			Location:       location,
			Qualification:  qualification,
			Specialization: specialization,
			Department:     department,
			Position:       position,
			Experience:     experience,
			CvPath:         cvpath,
		})
	}

	if err != nil {
		panic(err)
	}
}
