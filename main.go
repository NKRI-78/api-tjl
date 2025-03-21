package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"superapps/controllers"
	helper "superapps/helpers"
	middleware "superapps/middlewares"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		helper.Logger("error", "Error getting env")
	}

	router := mux.NewRouter()

	router.Use(middleware.JwtAuthentication)

	// Check if the directory exists, create if it doesn't
	errMkidr := os.MkdirAll("public", os.ModePerm) // os.ModePerm ensures directory is created with the correct permissions
	if errMkidr != nil {
		log.Fatalf("Failed to create or access directory: %v", err)
	}

	// Open the public directory
	dir, err := os.Open("public")
	if err != nil {
		log.Fatalf("Failed to open public directory: %v", err)
	}
	defer dir.Close()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatalf("Failed to read directory contents: %v", err)
	}

	// Loop through each file in the directory
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			// Define static and public paths
			staticPath := "/" + fileInfo.Name() + "/"
			publicPath := "./public/" + fileInfo.Name() + "/"

			log.Printf("Serving static files from %s at %s", publicPath, staticPath)

			// Register (override if already exists) the route to serve static content
			router.PathPrefix(staticPath).Handler(http.StripPrefix(staticPath, http.FileServer(http.Dir(publicPath))))
		}
	}

	// Inisialisasi rate limiter: 2 permintaan per menit
	rateLimiter := middleware.NewRateLimiter(2, 1)

	// Administration
	router.HandleFunc("/api/v1/province", controllers.Province).Methods("GET")
	router.HandleFunc("/api/v1/city/{province_id}", controllers.City).Methods("GET")
	router.HandleFunc("/api/v1/district/{city_id}", controllers.District).Methods("GET")
	router.HandleFunc("/api/v1/subdistrict/{district_id}", controllers.Subdistrict).Methods("GET")
	router.HandleFunc("/api/v1/country", controllers.Country).Methods("GET")

	// Auth
	router.Handle("/api/v1/login", rateLimiter.LimitMiddleware(http.HandlerFunc(controllers.Login))).Methods("POST")
	router.Handle("/api/v1/register", rateLimiter.LimitMiddleware(http.HandlerFunc(controllers.Register))).Methods("POST")

	// Branch
	router.HandleFunc("/api/v1/branch", controllers.Branch).Methods("GET")

	// Banner
	router.HandleFunc("/api/v1/banner", controllers.BannerList).Methods("GET")
	router.HandleFunc("/api/v1/banner-store", controllers.BannerStore).Methods("POST")

	// Document
	router.HandleFunc("/api/v1/document", controllers.DocumentList).Methods("GET")
	router.HandleFunc("/api/v1/document-store", controllers.DocumentStore).Methods("POST")

	// Profile
	router.HandleFunc("/api/v1/profile", controllers.Profile).Methods("GET")
	router.HandleFunc("/api/v1/profile-update", controllers.ProfileUpdate).Methods("PUT")

	// Otp
	router.HandleFunc("/api/v1/resend-otp", controllers.ResendOtp).Methods("POST")
	router.HandleFunc("/api/v1/verify-otp", controllers.VerifyOtp).Methods("POST")

	// Admin
	router.HandleFunc("/api/v1/admin/job", controllers.AdminJobList).Methods("GET")

	// Apply Job
	router.HandleFunc("/api/v1/apply/job", controllers.ApplyJob).Methods("POST")

	// Update Apply Job
	router.HandleFunc("/api/v1/update/apply/job", controllers.UpdateApplyJob).Methods("PUT")

	// List Apply Job
	router.HandleFunc("/api/v1/list/apply/job", controllers.ListApplyJob).Methods("GET")

	// Info Apply Job
	router.HandleFunc("/api/v1/info/apply/job/{id}", controllers.InfoApplyJob).Methods("GET")

	// Jobs
	router.HandleFunc("/api/v1/job", controllers.JobList).Methods("GET")
	router.HandleFunc("/api/v1/job-detail/{id}", controllers.JobDetail).Methods("GET")
	router.HandleFunc("/api/v1/job-store", controllers.JobStore).Methods("POST")
	router.HandleFunc("/api/v1/job-favourite", controllers.JobFavourite).Methods("POST")
	router.HandleFunc("/api/v1/job-categories", controllers.JobCategory).Methods("GET")
	router.HandleFunc("/api/v1/job-places", controllers.JobPlace).Methods("GET")
	router.HandleFunc("/api/v1/job-category-count", controllers.JobCategoryCount).Methods("GET")

	// Language
	router.HandleFunc("/api/v1/language", controllers.Language).Methods("GET")

	// Form Biodata
	router.HandleFunc("/api/v1/form-biodata", controllers.FormBiodata).Methods("POST")
	router.HandleFunc("/api/v1/delete-form-biodata", controllers.DeleteFormBiodata).Methods("DELETE")

	// Form Address
	router.HandleFunc("/api/v1/form-address", controllers.FormAddress).Methods("POST")
	router.HandleFunc("/api/v1/delete-form-address", controllers.DeleteFormAddress).Methods("DELETE")

	// Form Education
	router.HandleFunc("/api/v1/form-education", controllers.FormEducation).Methods("POST")
	router.HandleFunc("/api/v1/delete-form-education", controllers.DeleteFormEducation).Methods("DELETE")
	router.HandleFunc("/api/v1/update-form-education", controllers.UpdateFormEducation).Methods("PUT")

	// Form Exercise
	router.HandleFunc("/api/v1/form-exercise", controllers.FormExercise).Methods("POST")
	router.HandleFunc("/api/v1/delete-form-exercise", controllers.DeleteFormExercise).Methods("DELETE")
	router.HandleFunc("/api/v1/update-form-exercise", controllers.UpdateFormExercise).Methods("PUT")

	// Form Work
	router.HandleFunc("/api/v1/form-work", controllers.FormWork).Methods("POST")
	router.HandleFunc("/api/v1/delete-form-work", controllers.DeleteFormWork).Methods("DELETE")
	router.HandleFunc("/api/v1/update-form-work", controllers.UpdateFormWork).Methods("PUT")

	// Form Language
	router.HandleFunc("/api/v1/form-language", controllers.FormLanguage).Methods("POST")
	router.HandleFunc("/api/v1/delete-form-language", controllers.DeleteFormLanguage).Methods("DELETE")
	router.HandleFunc("/api/v1/update-form-language", controllers.UpdateFormLanguage).Methods("PUT")

	// Forum
	router.HandleFunc("/api/v1/forum-store", controllers.ForumStore).Methods("POST")
	router.HandleFunc("/api/v1/forum-delete", controllers.ForumDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/forum-list", controllers.ForumList).Methods("GET")
	router.HandleFunc("/api/v1/forum-type", controllers.ForumCategory).Methods("GET")
	router.HandleFunc("/api/v1/forum-detail/{id}", controllers.ForumDetail).Methods("GET")

	// Content Comment
	router.HandleFunc("/api/v1/content/comment", controllers.CreateContentComment).Methods("POST")
	router.HandleFunc("/api/v1/content/comment/delete", controllers.DeleteContentComment).Methods("DELETE")

	// Content Like
	router.HandleFunc("/api/v1/content/like", controllers.CreateContentLike).Methods("POST")

	// Content Unlike
	router.HandleFunc("/api/v1/content/unlike", controllers.CreateContentUnlike).Methods("POST")

	// Media
	router.HandleFunc("/api/v1/media/upload", controllers.Upload).Methods("POST")

	portEnv := os.Getenv("PORT")
	port := ":" + portEnv

	// NOT SECURE FOR USE
	// server := new(http.Server)
	// server.Handler = router
	// server.Addr = ":" + port

	fmt.Println("Starting server at", port)

	server := &http.Server{
		Addr:              port,
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	errListenAndServe := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", errListenAndServe)
	}

	// errs := http.ListenAndServe(port, router)
	// if errs != nil {
	// 	fmt.Println("Error starting server:", errs)
	// }
}
