package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"unicode"

	"example.com/basicWebApp/storage"
	"golang.org/x/crypto/bcrypt"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("******rootHandler running ******")
	http.Redirect(w, r, "/login", http.StatusFound)
}

// registerHandler serves form for registering new users
func RegisterHandler(w http.ResponseWriter, r *http.Request) { 
	fmt.Println("******registerHandler running ******")
	templates.ExecuteTemplate(w, "register.html", nil)

}

// registerAuthHandler creates new user in database
/* 
		1. check username criteria.
		2. check passsword criteria.
		3. check first and last name criteria.
		4. check if username already exists in database.
		5. create bcrypt hash from password.
	*/
func RegisterAuthHandler(w http.ResponseWriter, r *http.Request) { 
	fmt.Println("******registerAuthHandler running ******")
	r.ParseForm()
	username := r.FormValue("userName")
	firstname := r.FormValue("firstName")
	lastname := r.FormValue("lastName")
	password := r.FormValue("password")
	
	// check username for only alphaNumric characters.
	var userNameAlphaNumeric = true
	// range through each char of username and check if char is letter and or number. 
	for _, char := range username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		// if !unicode.IsLetter(char) && !unicode.IsNumber(char) 
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			userNameAlphaNumeric = false
		}
	}
	var firstNameAlpha = true
	for _, char := range firstname {
		if unicode.IsLetter(char) == false {
			firstNameAlpha = false
		}
	}
	var lastNameAlpha = true
	for _, char := range lastname {
		if unicode.IsLetter(char) == false  {
			lastNameAlpha = false
		}
	}
	// check length username, first name, last name.
	var userNameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		userNameLength = true
	}
	var firstNameLength bool
	if 1 <= len(username) && len(username) <= 25 {
		firstNameLength = true
	}
	var lastNameLength bool
	if 1 <= len(username) && len(username) <= 25 {
		lastNameLength = true
	}
	// check passsword criteria
	fmt.Println("password:", password, "\npswdLength:", len(password))
	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true
		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", userNameAlphaNumeric, "\nuserNameLength:", userNameLength, "\nfirstNameAlpha:", firstNameAlpha, "\nuserNameLength:", firstNameLength, "\nlastNameAlpha:", lastNameAlpha, "\nlastNameLength:", lastNameLength,)
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !userNameAlphaNumeric || !userNameLength || !firstNameAlpha || !firstNameLength || !lastNameAlpha || !lastNameLength {
		templates.ExecuteTemplate(w, "register.html", "please check username and password criteria")
		return
	}
	err := storage.IsUser(username)
	if err != sql.ErrNoRows {
		fmt.Println("usernmae already exists, err:", err)
		templates.ExecuteTemplate(w, "register.html", "username taken")
		 return	
	}
	//create hash from password
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		templates.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	fmt.Println("hash:", hash)
	fmt.Println("string(hash):", string(hash))
	var stringHash string = string(hash)
	// finally insert user into db with validated credentials and generated hash.
	insertStmt, err := storage.PrepareInsert(&username, &lastname, &firstname, &stringHash)
	if err != nil {
		fmt.Println("error preparing statement:", err)
		templates.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	err = storage.InsertUser(insertStmt, &username, &lastname, &firstname, &stringHash)
	if err != nil {
		fmt.Println("error inserting new user")
		templates.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	fmt.Fprint(w, "congrats, your account has been successfully created")
	
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	templates.ExecuteTemplate(w, "login.html", nil)
}

func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("username:", username, "password:", password)
	// retrieve password from db to compare (hash) with user supplied password's hash
	hash, err := storage.GetHash(username)
	fmt.Println("hash from db:", hash)
	if err != nil {
		fmt.Println("error selecting Hash in db by Username")
		templates.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	// func CompareHashAndPassword(hashedPassword, password []byte) error
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// returns nill on succcess
	if err == nil {
		fmt.Fprint(w, "You have successfully logged in :)")
		return
	}
	fmt.Println("incorrect password")
	templates.ExecuteTemplate(w, "login.html", "check username and password")
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****AccountHandler running*****")
	templates.ExecuteTemplate(w, "account.html", nil)
}