package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/arizdn234/crud-users-with-login-system/internal/models"
	"github.com/arizdn234/crud-users-with-login-system/internal/repository"
	"github.com/gorilla/mux"

	myUtils "github.com/arizdn234/crud-users-with-login-system/internal/utils"
)

type UserHandler struct {
	UserRepository *repository.UserRepository
}

func NewUserHandler(ur *repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepository: ur}
}

func (uh *UserHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	info := `
	Simple Login & Register system with CRUD on users data!

	Routes Available:
	- GET    /                  : Welcome message
	- GET    /users             : Get all users
	- POST   /users             : Create a new user (using create method)
	- GET    /users/{id}        : Get user by ID
	- PUT    /users/{id}        : Update user by ID
	- DELETE /users/{id}        : Delete user by ID
	- POST   /login             : User login
	- POST   /register          : Register new user (using register method)
	`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(info))
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	hashedPassword, err := myUtils.HashPassword(user.Password + user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	if err := uh.UserRepository.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updateReq struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := uh.UserRepository.FindByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if updateReq.Name != nil {
		user.Name = *updateReq.Name
	}

	if updateReq.Email != nil {
		user.Email = *updateReq.Email
	}

	if updateReq.Password != nil && *updateReq.Password != "" {
		hashedPassword, err := myUtils.HashPassword(*updateReq.Password + user.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.Password = hashedPassword
	}

	if err := uh.UserRepository.Update(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := uh.UserRepository.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("deletion was successful"))
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var user []models.User

	if err := uh.UserRepository.FindAll(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uh.UserRepository.FindByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// validate user data
	if err := uh.validateUser(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// email check
	existingUser, err := uh.UserRepository.FindByEmail(newUser.Email)
	if err == nil && existingUser != nil {
		http.Error(w, "email already registered", http.StatusBadRequest)
		return
	}

	// hash
	hashed, err := myUtils.HashPassword(newUser.Password + newUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newUser.Password = hashed

	// fmt.Printf("newUser: %v\n", newUser)
	// fmt.Printf("existingUser: %v\n", existingUser)

	if err := uh.UserRepository.Create(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (uh *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	existingUser, err := uh.UserRepository.FindByEmail(user.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	hashedRequestBodyPassword, _ := myUtils.HashPassword(user.Password + user.Email)

	if hashedRequestBodyPassword != existingUser.Password {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// user credentials have been verified
	// generate jwt token
	token, err := myUtils.CreateToken(existingUser.ID, existingUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Printf("token: %v\n", token)
	// fmt.Println(myUtils.VerifyToken(token))

	// set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true, // set HttpOnly to true for better security
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("login successful"))
}

func (uh *UserHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logout successful"))
}

func (uh *UserHandler) validateUser(user models.User) error {
	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	// Validate minimum password length
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Validate password composition (at least one uppercase, one lowercase, and one number)
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(user.Password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(user.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(user.Password)

	if !hasUppercase || !hasLowercase || !hasNumber {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}
