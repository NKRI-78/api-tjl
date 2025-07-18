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

	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.JwtAuthentication)

	errMkidr := os.MkdirAll("public", os.ModePerm)
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

	// Translate
	router.HandleFunc("/api/v1/extract/translate-lang", controllers.TranslateLang).Methods("POST")

	// Administration
	router.HandleFunc("/api/v1/province", controllers.Province).Methods("GET")
	router.HandleFunc("/api/v1/city/{province_id}", controllers.City).Methods("GET")
	router.HandleFunc("/api/v1/district/{city_id}", controllers.District).Methods("GET")
	router.HandleFunc("/api/v1/subdistrict/{district_id}", controllers.Subdistrict).Methods("GET")
	router.HandleFunc("/api/v1/country", controllers.Country).Methods("GET")
	router.HandleFunc("/api/v1/country-store", controllers.CountryStore).Methods("POST")
	router.HandleFunc("/api/v1/country-delete", controllers.CountryDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/country-update", controllers.CountryUpdate).Methods("PUT")

	// Register User Branch
	router.HandleFunc("/api/v1/admin/register-user-branch", controllers.RegisterUserBranch).Methods("POST")
	router.HandleFunc("/api/v1/admin/update-user-branch", controllers.UpdateUserBranch).Methods("PUT")

	// Import User Create
	router.HandleFunc("/api/v1/import/user/create", controllers.ImportUserCreate).Methods("POST")
	router.HandleFunc("/api/v1/import/user/list", controllers.AdminListCandidateImport).Methods("GET")

	// List Candidate Passes
	router.HandleFunc("/api/v1/candidate-passes/list", controllers.CandidatePassesList).Methods("GET")
	router.HandleFunc("/api/v1/candidate-passes/form/list", controllers.CandidatePassesFormList).Methods("GET")
	router.HandleFunc("/api/v1/candidate-passes/info/departure", controllers.CandidateInfoDeparture).Methods("GET")
	router.HandleFunc("/api/v1/candidate-passes/form", controllers.CandidatePassesForm).Methods("POST")

	// View Candidate PDF Departure
	router.HandleFunc("/api/v1/view/pdf/departure", controllers.ViewPdfDeparture).Methods("GET")

	// View Candidate PDF Apply Job Offline
	router.HandleFunc("/api/v1/view/pdf/apply-job-offline", controllers.ViewPdfApplyJobOffline).Methods("GET")

	// Inbox
	router.HandleFunc("/api/v1/inbox/list", controllers.InboxList).Methods("GET")
	router.HandleFunc("/api/v1/inbox/detail/{id}", controllers.InboxDetail).Methods("GET")
	router.HandleFunc("/api/v1/inbox/badge", controllers.InboxBadge).Methods("GET")

	// Auth
	router.Handle("/api/v1/login", rateLimiter.LimitMiddleware(http.HandlerFunc(controllers.Login))).Methods("POST")
	router.Handle("/api/v1/register", rateLimiter.LimitMiddleware(http.HandlerFunc(controllers.Register))).Methods("POST")
	router.Handle("/api/v1/forgot-password", rateLimiter.LimitMiddleware(http.HandlerFunc(controllers.ForgotPassword))).Methods("PUT")
	router.Handle("/api/v1/update-email", rateLimiter.LimitMiddleware(http.HandlerFunc(controllers.UpdateEmail))).Methods("PUT")

	// Admin
	router.HandleFunc("/api/v1/admin/list/apply/job", controllers.AdminListApplyJob).Methods("GET")
	router.HandleFunc("/api/v1/admin/detail/apply/job/{id}", controllers.AdminDetailApplyJob).Methods("GET")

	// Admin Auth
	router.Handle("/api/v1/login-admin", rateLimiter.LimitMiddleware(http.HandlerFunc(controllers.LoginAdmin))).Methods("POST")

	// Admin User
	router.HandleFunc("/api/v1/admin/list/user/{type}", controllers.AdminListUser).Methods("GET")
	router.HandleFunc("/api/v1/admin/update/user", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/v1/admin/delete/user", controllers.DeleteUser).Methods("DELETE")

	// Admin Chart
	router.HandleFunc("/api/v1/admin/chart/summary", controllers.Summary).Methods("GET")

	// Branch
	router.HandleFunc("/api/v1/branch", controllers.Branch).Methods("GET")

	// Banner
	router.HandleFunc("/api/v1/banner", controllers.BannerList).Methods("GET")
	router.HandleFunc("/api/v1/banner-store", controllers.BannerStore).Methods("POST")
	router.HandleFunc("/api/v1/banner-update", controllers.BannerUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/banner-delete", controllers.BannerDelete).Methods("DELETE")

	// News
	router.HandleFunc("/api/v1/news", controllers.NewsList).Methods("GET")
	router.HandleFunc("/api/v1/news-detail/{id}", controllers.NewsDetail).Methods("GET")
	router.HandleFunc("/api/v1/news-store", controllers.NewsStore).Methods("POST")
	router.HandleFunc("/api/v1/news-update", controllers.NewsUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/news-store-image", controllers.NewsStoreImage).Methods("POST")
	router.HandleFunc("/api/v1/news-update-image", controllers.NewsUpdateImage).Methods("PUT")
	router.HandleFunc("/api/v1/news-delete-image", controllers.NewsDeleteImage).Methods("DELETE")
	router.HandleFunc("/api/v1/news-delete", controllers.NewsDelete).Methods("DELETE")

	// Event
	router.HandleFunc("/api/v1/event", controllers.EventList).Methods("GET")
	router.HandleFunc("/api/v1/event-detail/{id}", controllers.EventDetail).Methods("GET")
	router.HandleFunc("/api/v1/event-store", controllers.EventStore).Methods("POST")
	router.HandleFunc("/api/v1/event-update", controllers.EventUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/event-store-image", controllers.EventStoreImage).Methods("POST")
	router.HandleFunc("/api/v1/event-update-image", controllers.EventUpdateImage).Methods("PUT")
	router.HandleFunc("/api/v1/event-delete-image", controllers.EventDeleteImage).Methods("DELETE")
	router.HandleFunc("/api/v1/event-delete", controllers.EventDelete).Methods("DELETE")

	// Document
	router.HandleFunc("/api/v1/document", controllers.DocumentList).Methods("GET")
	router.HandleFunc("/api/v1/document-store", controllers.DocumentStore).Methods("POST")
	router.HandleFunc("/api/v1/get-document-additional/{type}", controllers.GetDocumentAdditional).Methods("GET")
	router.HandleFunc("/api/v1/documment-additional-store", controllers.DocumentAdditionalStore).Methods("POST")
	router.HandleFunc("/api/v1/document-additional-update", controllers.DocumentAdditionalUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/document-additional-delete", controllers.DocumentAdditionalDelete).Methods("DELETE")

	// Profile
	router.HandleFunc("/api/v1/profile", controllers.Profile).Methods("GET")
	router.HandleFunc("/api/v1/profile-update", controllers.ProfileUpdate).Methods("PUT")

	// Otp
	router.HandleFunc("/api/v1/resend-otp", controllers.ResendOtp).Methods("POST")
	router.HandleFunc("/api/v1/verify-otp", controllers.VerifyOtp).Methods("POST")

	// Apply Job
	router.HandleFunc("/api/v1/apply/job", controllers.ApplyJob).Methods("POST")

	// Assign
	router.HandleFunc("/api/v1/assign/document/apply/job", controllers.AssignDocumentApplyJob).Methods("POST")
	router.HandleFunc("/api/v1/assign/branch", controllers.AssignBranch).Methods("POST")
	router.HandleFunc("/api/v1/assign/education", controllers.AssignEducation).Methods("POST")
	router.HandleFunc("/api/v1/assign/work", controllers.AssignWork).Methods("POST")
	router.HandleFunc("/api/v1/assign/address", controllers.AssignAddress).Methods("POST")
	router.HandleFunc("/api/v1/assign/exercise", controllers.AssignExercise).Methods("POST")
	router.HandleFunc("/api/v1/assign/language", controllers.AssignLanguage).Methods("POST")
	router.HandleFunc("/api/v1/assign/apply-job", controllers.AssignApplyJob).Methods("POST")

	// Badges
	router.HandleFunc("/api/v1/apply-job-badges", controllers.ApplyJobBadges).Methods("GET")
	router.HandleFunc("/api/v1/admin/apply-job-badges", controllers.AdminApplyJobBadges).Methods("GET")
	router.HandleFunc("/api/v1/admin/candidate-passes-badges", controllers.AdminCandidatePassesBadges).Methods("GET")

	// Update Apply Job
	router.HandleFunc("/api/v1/update/apply/job", controllers.UpdateApplyJob).Methods("PUT")

	// List Apply Job
	router.HandleFunc("/api/v1/list/apply/job", controllers.ListApplyJob).Methods("GET")
	router.HandleFunc("/api/v1/list/apply/job/history", controllers.ListApplyJobHistory).Methods("GET")

	// Info Apply Job
	router.HandleFunc("/api/v1/info/apply/job/{id}", controllers.InfoApplyJob).Methods("GET")

	// Fcm
	router.HandleFunc("/api/v1/firebase/init-fcm", controllers.InitFcm).Methods("POST")

	// Company
	router.HandleFunc("/api/v1/company/list", controllers.CompanyList).Methods("GET")
	router.HandleFunc("/api/v1/company/store", controllers.CompanyStore).Methods("POST")
	router.HandleFunc("/api/v1/company/update", controllers.CompanyUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/company/delete", controllers.CompanyDelete).Methods("DELETE")

	// Jobs
	router.HandleFunc("/api/v1/job", controllers.JobList).Methods("GET")
	router.HandleFunc("/api/v1/job-detail/{id}", controllers.JobDetail).Methods("GET")
	router.HandleFunc("/api/v1/job-store", controllers.JobStore).Methods("POST")
	router.HandleFunc("/api/v1/job-skills-category-delete", controllers.JobSkillCategoryDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/job-skills-category-list", controllers.JobSkillCategoryList).Methods("GET")
	router.HandleFunc("/api/v1/job-skills-category-store", controllers.JobSkillCategoryStore).Methods("POST")
	router.HandleFunc("/api/v1/job-update", controllers.JobUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/job-delete", controllers.JobDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/job-favourite", controllers.JobFavourite).Methods("POST")
	router.HandleFunc("/api/v1/job-categories", controllers.JobCategory).Methods("GET")
	router.HandleFunc("/api/v1/job-places", controllers.JobPlace).Methods("GET")

	// Job Category
	router.HandleFunc("/api/v1/job-category-store", controllers.JobCategoryStore).Methods("POST")
	router.HandleFunc("/api/v1/job-category-update", controllers.JobCategoryUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/job-category-delete", controllers.JobCategoryDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/job-category-count", controllers.JobCategoryCount).Methods("GET")

	// Type Jobs
	router.HandleFunc("/api/v1/type-job-list", controllers.TypeJobList).Methods("GET")
	router.HandleFunc("/api/v1/type-job-store", controllers.TypeJobStore).Methods("POST")
	router.HandleFunc("/api/v1/type-job-delete", controllers.TypeJobDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/type-job-update", controllers.TypeJobUpdate).Methods("PUT")

	// Role
	router.HandleFunc("/api/v1/role/list", controllers.RoleList).Methods("GET")

	// Icons
	router.HandleFunc("/api/v1/icons", controllers.IconList).Methods("GET")
	router.HandleFunc("/api/v1/icons/store", controllers.IconStore).Methods("POST")
	router.HandleFunc("/api/v1/icons/update", controllers.IconUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/icons/delete", controllers.IconDelete).Methods("DELETE")

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
	router.HandleFunc("/api/v1/comment-store", controllers.CommentStore).Methods("POST")
	router.HandleFunc("/api/v1/comment-delete", controllers.CommentDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/reply-store", controllers.ReplyStore).Methods("POST")
	router.HandleFunc("/api/v1/reply-delete", controllers.ReplyDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/forum-store", controllers.ForumStore).Methods("POST")
	router.HandleFunc("/api/v1/forum-like", controllers.ForumLike).Methods("POST")
	router.HandleFunc("/api/v1/forum-delete", controllers.ForumDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/forum-list", controllers.ForumList).Methods("GET")
	router.HandleFunc("/api/v1/forum-type", controllers.ForumCategory).Methods("GET")
	router.HandleFunc("/api/v1/forum-detail/{id}", controllers.ForumDetail).Methods("GET")

	// Media
	router.HandleFunc("/api/v1/media/upload", controllers.Upload).Methods("POST")

	// Client
	router.HandleFunc("/api/v1/media/client/list", controllers.ClientList).Methods("GET")
	router.HandleFunc("/api/v1/media/client/store", controllers.ClientStore).Methods("POST")
	router.HandleFunc("/api/v1/media/client/update", controllers.ClientUpdate).Methods("PUT")
	router.HandleFunc("/api/v1/media/client/delete", controllers.ClientDelete).Methods("DELETE")

	// Test Send E-mail
	router.HandleFunc("/api/v1/email/send", controllers.SendEmail).Methods("POST")

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
