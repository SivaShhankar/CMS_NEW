package Models

import "gopkg.in/mgo.v2/bson"

type ApplicantInfo struct {
	ID             bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name           string        `json:"name"`
	Age            int           `json:"age"`
	Gender         string        `json:"gender"`
	Mobile         int           `json:"mobile"`
	Email          string        `json:"email"`
	Location       string        `json:"location"`
	Qualification  string        `json:"qualification"`
	Specialization string        `json:"specialization"`
	Department     string        `json:"department"`
	JobCode        string        `json:"jobcode"`
	Position       string        `json:"position"`
	Experience     float64       `json:"experience"`
	CvPath         string        `json:"cvpath"`
	SourceFrom     string        `json:"sourcefrom"`
}
