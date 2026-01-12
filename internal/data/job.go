package data

// job represents a job posting in the application
type Job struct {
	// omitempty means if id is empty, hide it in JSON
	Id string `json:"id,omitempty"` 
	Title string `json:"title"`
	Description string `json:"description"`
	Company string `json:"company"`
	Salary string `json:"salary"`
}