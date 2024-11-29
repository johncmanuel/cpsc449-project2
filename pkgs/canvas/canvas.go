package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CanvasClient struct {
	BaseURL    string
	AuthToken  string
	HTTPClient *http.Client
}

type Assignment struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CourseID    int    `json:"course_id"`
	Description string `json:"description"`
	DueAt       string `json:"due_at"`
}

type Course struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
	Term    struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		StartAt       string `json:"start_at"`
		EndAt         string `json:"end_at"`
		CreatedAt     string `json:"created_at"`
		WorkflowState string `json:"workflow_state"`
	} `json:"term"`
}

func NewCanvasClient(baseURL, authToken string) *CanvasClient {
	return &CanvasClient{
		BaseURL:    baseURL,
		AuthToken:  authToken,
		HTTPClient: &http.Client{},
	}
}

// https://canvas.instructure.com/doc/api/courses.html#method.courses.index
func (c *CanvasClient) GetCurrentTermCourses() ([]Course, error) {
	// NOTE: There might be a way to dynamically get the current term (as of 11/29/24, it'll be Fall 2024)
	// Need to add enrollment_term_id to the url queries tho, but for sake of time, we can hard code it
	// https://canvas.instructure.com/doc/api/enrollment_terms.html#method.terms_api.index
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/courses?published=true&include[]=term", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var courses []Course
	err = json.Unmarshal(body, &courses)
	if err != nil {
		return nil, err
	}

	// Print courses, course term and id
	// TODO: ID isn't being printed. Maybe we don't have the proper permissions?
	for _, course := range courses {
		fmt.Printf("Course: %s (ID: %d), Term: %s\n (ID: %d)", course.Name, course.ID, course.Term.Name, course.Term.ID)
	}

	return courses, nil
}

func (c *CanvasClient) GetAssignmentsForCourse(courseID int) ([]Assignment, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/courses/%d/assignments", c.BaseURL, courseID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var assignments []Assignment
	err = json.Unmarshal(body, &assignments)
	if err != nil {
		return nil, err
	}

	return assignments, nil
}

func (c *CanvasClient) GetAllAssignmentsForCurrentTerm() (map[int]map[string][]Assignment, error) {
	courses, err := c.GetCurrentTermCourses()
	if err != nil {
		return nil, err
	}

	allAssignments := make(map[int]map[string][]Assignment)

	for _, course := range courses {
		assignments, err := c.GetAssignmentsForCourse(course.ID)
		if err != nil {
			fmt.Printf("Error fetching assignments for course %d: %v\n", course.ID, err)
			continue
		}
		allAssignments[course.ID] = map[string][]Assignment{
			course.Name: assignments,
		}
	}

	return allAssignments, nil
}
