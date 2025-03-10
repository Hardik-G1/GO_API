package main

import (
	"log"
	"net/http"
	. "v1/Middleware"
	. "v1/MongoConnection"
	. "v1/Routes"
	. "v1/Routes/CreateRoutes"
	. "v1/Routes/DeleteRoutes"
	. "v1/Routes/ReadRoutes"
	. "v1/Routes/UpdateRoutes"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func main() {
	Connect()
	defer Close(Client, Ctx)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		// AllowedMethods: []string{"GET,POST,PUT,DELETE"}, // Allowing only get, just an example
	})
	r := mux.NewRouter()
	r.HandleFunc("/api/signup", AuthenticateUser(SignUp)).Methods("POST")
	r.HandleFunc("/api/verify/mail", VerifyMail).Methods("POST")
	r.HandleFunc("/api/resendVerification", AuthenticateUser(ResendVerificationMail)).Methods("POST")
	r.HandleFunc("/api/login", Login).Methods("POST")
	r.HandleFunc("/api/refresh", Refresh).Methods("GET")
	r.HandleFunc("/api/logout", Logout).Methods("GET")
	r.HandleFunc("/api/checkUsername/{name}", CheckUsernamePresent).Methods("GET")
	r.HandleFunc("/api/checkEmail/{mail}", CheckEmailPresent).Methods("GET")
	r.HandleFunc("/api/checkurl/{url}", CheckUrlisAvailable).Methods("GET")
	//guest allowed
	r.HandleFunc("/api/create", AuthenticateUser(CreateData)).Methods("POST")
	r.HandleFunc("/api/get/{url}", AuthenticateUser(GetData)).Methods("GET")
	r.HandleFunc("/api/search", SearchDataAll).Methods("POST")
	r.HandleFunc("/api/search/user", SearchUserAll).Methods("POST")
	r.HandleFunc("/api/get-user/{user}", AuthenticateUser(GetUserPage)).Methods("GET")
	//Authentication needed
	r.HandleFunc("/api/star/{url}/{set}", AuthenticateUser(StarData)).Methods("PUT")
	r.HandleFunc("/api/follow/{usertofollow}/{set}", AuthenticateUser(FollowUser)).Methods("PUT")
	r.HandleFunc("/api/fork/{url}", AuthenticateUser(ForkData)).Methods("PUT")
	r.HandleFunc("/api/edit/{url}", AuthenticateUser(EditData)).Methods("POST")
	r.HandleFunc("/api/editUser", AuthenticateUser(EditUser)).Methods("POST")
	r.HandleFunc("/api/delete/{url}", AuthenticateUser(DeleteData)).Methods("DELETE")
	r.HandleFunc("/api/report/{url}}", AuthenticateUser(ReportData)).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", c.Handler(r)))
}
