package controller

import (
	"../model/token"
	"../model/user"
	"../module"
	"../module/password"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

/**
@brief: register log our user
*/
func Login(w http.ResponseWriter, r *http.Request) {
	var authDetails user.Authentication

	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		var err module.Error
		err = module.SetError(err, "Error in reading payload.")
		json.NewEncoder(w).Encode(err)
		return
	}

	authUser, _ := user.GetUserByEmail(authDetails.Email) //autehntificate him by email

	if authUser.Email == "" || !authUser.EmailValidate {
		var err module.Error
		err = module.SetError(err, "Username or Password is incorrect!")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	check := password.CheckPasswordHash(authDetails.Password, authUser.Password)

	if !check {
		var err module.Error
		err = module.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := token.GenerateJWT(authUser.Email, authUser.Role)
	if err != nil {
		var err module.Error
		err = module.SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var tToken token.Token
	tToken.Email = authUser.Email
	tToken.Role = authUser.Role
	tToken.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tToken)

}

/**
 *	@brief: Register create user
 */
func Register(w http.ResponseWriter, r *http.Request) {
	register := new(user.Register)
	User := new(user.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&register)

	if register.Email != "" && register.Password != "" && register.Pseudo != "" {
		pass, _ := password.HashPassword(register.Password)
		User.Password = pass
		User.Email = register.Email
		User.Pseudo = register.Pseudo
		User.Role = "user"
		User.EmailValidate = true
		User.TokenEmail = user.TokenGenerator(10)

		newUser, err := user.NewUser(User)
		if err != nil {
			var err module.Error
			err = module.SetError(err, "can't register")
			json.NewEncoder(w).Encode(err)
			return
		}
		postsJson, _ := json.Marshal(newUser)
		w.Write(postsJson)
		return
	}
	var err module.Error
	err = module.SetError(err, "can't register")
	json.NewEncoder(w).Encode(err)
	return

}

/**
 * 	@brief: VerifyEmail activate user's account
 */
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	User, err := user.GetUserByToken(token)

	if err == false {
		var err module.Error
		err = module.SetError(err, "token not found")
		json.NewEncoder(w).Encode(err)
		return
	}
	user.UpdateUser(User, user.User{EmailValidate: true})
	userJson, _ := json.Marshal(User)
	w.Write(userJson)
}

/**
@brief: work in progress
*/
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	User, _ := user.GetUserByEmail(r.Header.Get("Email"))
	user.DeleteUser(User)
}
