package remoteok

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Job struct {
	Slug        string    `json:"slug,omitempty"`
	ID          string    `json:"id,omitempty"`
	Epoch       int       `json:"epoch,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Company     string    `json:"company,omitempty"`
	Position    string    `json:"position,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Description string    `json:"description,omitempty"`
	Location    string    `json:"location,omitempty"`
	SalaryMin   int       `json:"salary_min,omitempty"`
	SalaryMax   int       `json:"salary_max,omitempty"`
	URL         string    `json:"url,omitempty"`
}

// Get Jobs from remote Ok API. It will return array of Job struct
func getJobs() (*[]Job, error) {
	resp, err := http.Get("https://remoteok.com/api?api=1")
	if err != nil {
		fmt.Println("Error calling Call", err)
	}

	var data []Job
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Failed to decode JSON:", err)
		return nil, err
	}
	return &data, nil
}

/*
Filter function for remote Ok API.
Keep only job that fulfill these conditions:

- Position contains one of the position parameter
*/
func filter(job *[]Job, position ...string) *[]Job {
	var result []Job
	for _, s := range *job {
		for _, p := range position {
			if strings.Contains(strings.ToLower(s.Position), p) {
				result = append(result, s)
				break
			}
		}
	}
	return &result
}

/*
getJobsAndFilter
Get Jobs from remote Ok API, after that applied filter
*/
func getJobsAndFilter() []Job {
	result := make([]Job, 0)
	job, err := getJobs()
	if err != nil {
		fmt.Println("Error calling GetJobs", err)

	}
	result = *filter(job, "backend")
	return result
}
