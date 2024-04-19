package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/arizdn234/crud-users-with-login-system/internal/handlers"
	"github.com/arizdn234/crud-users-with-login-system/internal/models"
	"github.com/arizdn234/crud-users-with-login-system/internal/repository"
)

var _ = Describe("User Handlers", func() {
	var (
		userHandler      *handlers.UserHandler
		userRepository   *repository.UserRepository
		server           *httptest.Server
		mockUserResponse models.User
	)

	BeforeEach(func() {
		// Setup mocked repository
		userRepository = &repository.UserRepository{}

		// Setup user handler with mocked repository
		userHandler = handlers.NewUserHandler(userRepository)

		// Setup test server
		server = httptest.NewServer(http.HandlerFunc(userHandler.CreateUser))

		// Setup mocked user response
		mockUserResponse = models.User{
			Name:     "Khalid Khasmiri",
			Email:    "khal1denby@example.com",
			Password: "nosecret",
		}
	})

	AfterEach(func() {
		server.Close()
	})

	// Test CreateUser
	It("should create a new user (indicated via Email data)", func() {
		userJSON, _ := json.Marshal(mockUserResponse)

		resp, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewBuffer(userJSON))
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		var createdUser models.User
		err = json.NewDecoder(resp.Body).Decode(&createdUser)
		Expect(err).NotTo(HaveOccurred())
		Expect(createdUser.Email).To(Equal(mockUserResponse.Email))
	})

	// Test AuthenticateUser (success)
	It("should authenticate a user with valid credentials", func() {
		loginData := models.User{
			Email:    mockUserResponse.Email,
			Password: mockUserResponse.Password,
		}

		loginDataJSON, _ := json.Marshal(loginData)
		resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(loginDataJSON))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	// Test AuthenticateUser (invalid)
	It("should not authenticate a user with invalid credentials", func() {
		invalidLoginData := models.User{
			Email:    mockUserResponse.Email,
			Password: "wrongpassword",
		}

		invalidLoginDataJSON, _ := json.Marshal(invalidLoginData)
		resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(invalidLoginDataJSON))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})
})
