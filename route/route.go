package route

import (
	"../controller"
	"../module"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
)

/**
@brief: IniRoute will declare our "routes" in backend server
*/
func IniRoute(app *mux.Router) {

	app.Use(CommonMiddleware)
	app.Methods("POST").Path("/register").HandlerFunc(controller.Register) // call register function when the url is ...:8080/register
	app.Methods("POST").Path("/login").HandlerFunc(controller.Login)
	app.Methods("GET").Path("/verify-email/{token}").HandlerFunc(controller.VerifyEmail)

	api := app.PathPrefix("/api").Subrouter() // create a "locked path" (add a fixed part of the url)
	v1 := api.PathPrefix("/v1").Subrouter()

	v1.Methods("DELETE").Path("/delete-user").HandlerFunc(IsAuthorized(controller.DeleteAccount))
	rPost := v1.PathPrefix("/post").Subrouter()
	rPost.Methods("GET").Path("").HandlerFunc(controller.GetPosts)
	rPost.Methods("GET").Path("/{id}").HandlerFunc(controller.GetPost)
	rPost.Methods("GET").Path("/user/{id}").HandlerFunc(controller.GetPostByUserId)
	rPost.Methods("POST").Path("/like/{id}").HandlerFunc(IsAuthorized(controller.LikePost))
	rPost.Methods("POST").Path("/dislike/{id}").HandlerFunc(IsAuthorized(controller.DisLikePost))
	rPost.Methods("POST").Path("").HandlerFunc(IsAuthorized(controller.NewPost))
	rPost.Methods("DELETE").Path("/{id}").HandlerFunc(IsAuthorized(controller.DeletePost))

	rMessage := v1.PathPrefix("/message").Subrouter()
	rMessage.Methods("GET").Path("").HandlerFunc(controller.GetMessages)
	rMessage.Methods("GET").Path("/{id}").HandlerFunc(controller.GetMessage)
	rMessage.Methods("GET").Path("/user/{id}").HandlerFunc(controller.GetMessageByUserId)
	rMessage.Methods("GET").Path("/post/{id}").HandlerFunc(controller.GetMessageByPostId)
	rMessage.Methods("POST").Path("").HandlerFunc(IsAuthorized(controller.NewMessage))
	rMessage.Methods("DELETE").Path("/{id}").HandlerFunc(IsAuthorized(controller.DeleteMessage))
}

/**
@brief: CommonMiddleware define header's response
*/
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}

/**
@brief: IsAuthorized will check if our token is right or wrong to grant acces
*/
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			var err module.Error
			err = module.SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte("mysupersecretsecuritytoken")

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token.")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err module.Error
			err = module.SetError(err, "Your Token has been expired.")
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { // token claims will check if the token ir valid or invalid
			if claims["role"] == "admin" {
				r.Header.Set("Role", "admin")
				r.Header.Set("Email", fmt.Sprintf("%v", claims["email"]))
				handler.ServeHTTP(w, r)
				return

			} else if claims["role"] == "user" {
				r.Header.Set("Role", "user")
				r.Header.Set("Email", fmt.Sprintf("%v", claims["email"]))
				handler.ServeHTTP(w, r)
				return

			}
		}
		var reserr module.Error
		reserr = module.SetError(reserr, "Not Authorized.")
		json.NewEncoder(w).Encode(err)
	}
}
