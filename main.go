package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Model for course file
type Course struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

// fakeDB
var courses []Course

// middleware, helper - file
func (c *Course) IsEmpty() bool {
	return c.CourseName == ""
}

func main() {
	fmt.Println("We are building API in golang in MAIN")
	fmt.Println("Routers")

	r := mux.NewRouter()

	//	Seeding
	courses = append(courses, Course{CourseId: "1", CourseName: "MERN Stack", CoursePrice: 399, Author: &Author{Fullname: "Siddharth Priyadarsi", Website: "www.mernstack.com"}})
	courses = append(courses, Course{CourseId: "2", CourseName: "Golang", CoursePrice: 199, Author: &Author{Fullname: "Mr. Priyadarsi", Website: "www.python.com"}})

	//	routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateOnceCourse).Methods("PUT")
	r.HandleFunc("/course/{id}", deleteOneCourse).Methods("DELETE")

	//	listen to a port
	log.Fatal(http.ListenAndServe(":3000", r))
}

//Controllers-file

//serve home route

// r : reader is the value  which we read from route
// w	: writer is where we write the response
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to API</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)

}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	//grab it
	//	loop and check in db
	fmt.Println("Get one Course")
	w.Header().Set("Content-Type", "application/json")

	//grab it from request
	params := mux.Vars(r)

	//	loop through courses, find matching id and return matching res
	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("No Course found with given id")
	return

}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//	body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}

	//	what if {} : json but nothing iside it

	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("No data inside json")
		return
	}

	//	generate unique id, string
	//append new course into courses

	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return

}

func updateOnceCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//	first grab id from req
	params := mux.Vars(r)
	//	loop, id, remove, add with my ID

	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	//todo: send a response when id is not found
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one Course")
	w.Header().Set("Contetn-type", "application/json")

	//to get the params from the request
	params := mux.Vars(r)

	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			break
		}
	}
	return
	//	toDO return no course if CourseId is not available
}
