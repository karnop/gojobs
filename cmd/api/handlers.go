package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/karnop/gojobs/internal/data"
	"github.com/karnop/gojobs/internal/validator"
	"net/http"
	"os"
	"strconv"
	"time"
)

// USER HANDLERS

// registerUserHandler handles user signup
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// struct to store incoming json
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// decoding the request
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// create the user struct
	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
		Role:  "candidate",
	}

	// setting the password
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// validating the data
	v := validator.New()
	data.ValidateUser(v, user)

	if !v.Valid() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		json.NewEncoder(w).Encode(v.Errors)
		return
	}

	// insert into Database
	err = app.Users.Insert(user)
	if err != nil {
		// checking for duplicates
		if errors.Is(err, data.ErrDuplicateEmail) {
			v.AddError("email", "a user with this email address already exists")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict) // 409
			json.NewEncoder(w).Encode(v.Errors)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// loginUserHandler used for user login
func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// decoding request in a struct
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validating user inputs
	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email and Password required", http.StatusBadRequest)
		return
	}

	// finding user
	user, err := app.Users.GetByEmail(input.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// check password
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if !match {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generating JWT
	claims := jwt.MapClaims{
		"sub":  user.Id,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signing the token with secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// send token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// JOB HANDLERS

// createJobHandler handles POST request to add a new job
func (app *application) createJobHandler(w http.ResponseWriter, r *http.Request) {
	// variable to hold the incoming data
	var job data.Job

	// decoding json body from request
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// get user id from context
	userId := r.Context().Value("userId").(int)
	job.UserId = userId

	// RBAC check
	// only recruiters can post job
	user, err := app.Users.Get(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
        return
	}

	if user.Role != "recruiter" {
		http.Error(w, "Forbidden: Only recruiters can post jobs", http.StatusForbidden) // 403
        return
	}

	// validation logic
	v := validator.New()
	data.ValidateJob(v, &job) // v is already passed by reference

	if !v.Valid() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(v.Errors)
		return
	}

	// SQL query (postgres)
	// RETURNING id allows us to get the auto-generated ID back from SQLSERVER.
	// $1, $2, $3, $4 are placeholders for the data to prevent SQL Injection.
	query := `
		INSERT INTO jobs (title, description, company, salary, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	// QueryRow executes a query that returns exactly one row (the ID).
	err = app.DB.QueryRow(query, job.Title, job.Description, job.Company, job.Salary, job.UserId).Scan(&job.Id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// structured log
	app.Logger.Info("Job created successfully", 
        "job_id", job.Id, 
        "user_id", job.UserId,
        "title", job.Title,
    )

	// responding to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

// listjobshandler handles GET request to show all jobs
func (app *application) listJobsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string
		Company string
		data.Filters
	}

	v := validator.New()

	// parse Query Parameters
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Company = app.readString(qs, "company", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id") // Default sort by Id

	// validating filters
	input.Filters.SortSafelist = []string{"id", "title", "company", "salary", "-id", "-title", "-company", "-salary"}
	
	// calling db
	jobs, err := app.Jobs.GetAll(input.Title, input.Company, input.Filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// sending Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jobs": jobs,
	})

}

// getJobHandler fetches a single job by its Id
func (app *application) getJobHandler(w http.ResponseWriter, r *http.Request) {
	// extracting id from url path
	idStr := r.PathValue("id")

	// converting string id to integer
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid job Id", http.StatusBadRequest)
		return
	}

	// sql query
	query := `SELECT id, title, description, company, salary FROM jobs WHERE id = $1`
	var job data.Job

	err = app.DB.QueryRow(query, id).Scan(
		&job.Id,
		&job.Title,
		&job.Description,
		&job.Company,
		&job.Salary,
	)

	// handle errors
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Job not found", http.StatusNotFound)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// applyjobhandler POST req applies for a job 
func (app *application) applyJobHandler(w http.ResponseWriter, r *http.Request) {
	// get jobId from URL
	idStr := r.PathValue("id")
	jobId, err := strconv.Atoi(idStr)
	if err != nil || jobId < 1 {
		http.Error(w, "Invalid Job ID", http.StatusBadRequest)
		return
	}

	// getting user id from Context
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// create the application struct
	jobApp := &data.JobApplication{
		JobId : jobId,
		UserId : userId,
	}

	err = app.Applications.Insert(jobApp) 
	if err != nil {
		if errors.Is(err, data.ErrDuplicateApplication) {
			app.clientError(w, http.StatusConflict) // 409 Conflict
			// Ideally, send a JSON message: {"error": "You have already applied"}
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jobApp)
}
