package remoteok

import (
	"fmt"
	"gobot/pkg"
)

type SupabaseJob struct {
	ID          string `json:"id"`
	Epoch       int    `json:"epoch"`
	Slug        string `json:"slug"`
	Company     string `json:"company"`
	Position    string `json:"position"`
	Description string `json:"description"`
	Location    string `json:"location"`
	URL         string `json:"url"`
}

const dbName = "RemoteOk"

/*
GetAllSupabaseJob

Get All rows from RemoteOk Database, return arrays of SupabaseJob
*/
func GetAllSupabaseJob() []SupabaseJob {
	var result []SupabaseJob
	err := pkg.SupabaseClient.DB.From(dbName).Select("*").Execute(&result)
	if err != nil {
		fmt.Println("Error calling GetAllSupabaseJob", err)
		return nil
	}
	return result
}

/*
Insert

Insert a row into RemoteOk Database
*/
func Insert(job Job) {
	var results []SupabaseJob
	err := pkg.SupabaseClient.DB.From(dbName).Insert(SupabaseJob{
		ID:          job.ID,
		Epoch:       job.Epoch,
		Slug:        job.Slug,
		Company:     job.Company,
		Position:    job.Position,
		Description: job.Description,
		Location:    job.Location,
		URL:         job.URL,
	}).Execute(&results)
	if err != nil {
		fmt.Println("Error calling Insert", err)
		return
	}
}

/*
Delete

Delete a row in Games Database
*/
func Delete(post SupabaseJob) {
	var results []SupabaseJob
	err := pkg.SupabaseClient.DB.From(dbName).Delete().Eq("ID", post.ID).Execute(&results)
	if err != nil {
		fmt.Println("Error calling Delete", err)
		return
	}
}
