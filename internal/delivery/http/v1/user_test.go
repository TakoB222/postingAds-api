package v1

import (
	"bytes"
	"errors"
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/internal/service"
	mock_service "github.com/TakoB222/postingAds-api/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

//------------------Test functions for Authorization implementation------------------

func TestSignUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, input service.UserSignUpInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            signUpInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"firstName":"artem", "lastName":"yevsuikov", "email":"example1@gmail.com", "password":"dfghjk1503"}`,
			inputUser: signUpInput{
				FirstName: "artem",
				LastName:  "yevsuikov",
				Email:     "example1@gmail.com",
				Password:  "dfghjk1503",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input service.UserSignUpInput) {
				s.EXPECT().SignUp(input).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Empty field",
			inputBody: `{"lastName":"yevsuikov", "email":"example1@gmail.com", "password":"dfghjk1503"}`,
			inputUser: signUpInput{
				LastName: "yevsuikov",
				Email:    "example1@gmail.com",
				Password: "dfghjk1503",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, input service.UserSignUpInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "service error",
			inputBody: `{"firstName":"artem", "lastName":"yevsuikov", "email":"example1@gmail.com", "password":"dfghjk1503"}`,
			inputUser: signUpInput{
				FirstName: "artem",
				LastName:  "yevsuikov",
				Email:     "example1@gmail.com",
				Password:  "dfghjk1503",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input service.UserSignUpInput) {
				s.EXPECT().SignUp(input).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, service.UserSignUpInput{
				FirsName: testCase.inputUser.FirstName,
				LastName: testCase.inputUser.LastName,
				Email:    testCase.inputUser.Email,
				Password: testCase.inputUser.Password,
			})

			services := &service.Service{Authorization: auth}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/signUp", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signUp", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestRefreshTokens(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, input service.RefreshInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputRefresh         refreshTokensInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "ok",
			inputBody:    `{"RefreshToken":"someToken"}`,
			inputRefresh: refreshTokensInput{RefreshToken: "someToken"},
			mockBehavior: func(s *mock_service.MockAuthorization, input service.RefreshInput) {
				s.EXPECT().RefreshSession(input).Return(service.Tokens{RefreshToken: "someRefreshToken", AccessToken: "someAccessToken"}, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"access_token":"someAccessToken","refresh_token":"someRefreshToken"}`,
		},
		{
			name:                 "Empty input token field",
			inputBody:            `{"RefreshToken":""}`,
			inputRefresh:         refreshTokensInput{RefreshToken: ""},
			mockBehavior:         func(s *mock_service.MockAuthorization, input service.RefreshInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:         "service error",
			inputBody:    `{"RefreshToken":"someToken"}`,
			inputRefresh: refreshTokensInput{RefreshToken: "someToken"},
			mockBehavior: func(s *mock_service.MockAuthorization, input service.RefreshInput) {
				s.EXPECT().RefreshSession(input).Return(service.Tokens{}, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, service.RefreshInput{
				RefreshToken: testCase.inputRefresh.RefreshToken,
			})

			services := &service.Service{Authorization: auth}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/refreshTokens", handler.refreshTokens)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refreshTokens", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestSignIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, input service.SignInInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputSignIn          signInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			inputBody:   `{"email":"example@gmail.com","password":"somePassword"}`,
			inputSignIn: signInInput{Email: "example@gmail.com", Password: "somePassword"},
			mockBehavior: func(s *mock_service.MockAuthorization, input service.SignInInput) {
				s.EXPECT().SignIn(input).Return(service.Tokens{RefreshToken: "someRefreshToken", AccessToken: "someAccessToken"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"access_token":"someAccessToken","refresh_token":"someRefreshToken"}`,
		},
		{
			name:                 "Empty field",
			inputBody:            `{"email":"example@gmail.com"}`,
			inputSignIn:          signInInput{Email: "example@gmail.com"},
			mockBehavior:         func(s *mock_service.MockAuthorization, input service.SignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:        "service error",
			inputBody:   `{"email":"example@gmail.com","password":"somePassword"}`,
			inputSignIn: signInInput{Email: "example@gmail.com", Password: "somePassword"},
			mockBehavior: func(s *mock_service.MockAuthorization, input service.SignInInput) {
				s.EXPECT().SignIn(input).Return(service.Tokens{}, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, service.SignInInput{
				Email:    testCase.inputSignIn.Email,
				Password: testCase.inputSignIn.Password,
			})

			services := &service.Service{Authorization: auth}
			handler := &Handler{services: services}

			r := gin.New()
			r.POST("/signIn", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signIn", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

//-------------------------------------------------------------------------
//тестирование вспомагательной функции - getUserId
func TestGetUserId(t *testing.T) {
	testTable := []struct {
		name           string
		userId         interface{}
		setUserContext bool
		expectedBody   string
		expectedError  error
	}{
		{
			name:           "ok",
			userId:         "1",
			setUserContext: true,
			expectedBody:   "1",
			expectedError:  nil,
		},
		{
			name:           "Empty body user context",
			userId:         "",
			setUserContext: true,
			expectedBody:   "",
			expectedError:  errors.New("empty body of user context"),
		},
		{
			name:           "Wrong type of user context",
			userId:         1,
			setUserContext: true,
			expectedBody:   "",
			expectedError:  errors.New("invalid type of userId from context"),
		},
		{
			name:           "Empty user context",
			setUserContext: false,
			expectedBody:   "",
			expectedError:  errors.New("empty user context"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			ctx := new(gin.Context)

			if testCase.setUserContext {
				ctx.Set(userContext, testCase.userId)
			}

			userId, err := getUserId(ctx)

			assert.Equal(t, testCase.expectedBody, userId)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

//------------------Test functions for Ads implementation------------------

func TestGetAllAds(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAd, userId interface{})

	testTable := []struct {
		name                 string
		userId               interface{}
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Ok",
			userId: "1",
			mockBehavior: func(s *mock_service.MockAd, userId interface{}) {
				s.EXPECT().GetAllAds(userId).Return([]domain.Ad{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[]`,
		},
		{
			name:   "Service error",
			userId: "1",
			mockBehavior: func(s *mock_service.MockAd, userId interface{}) {
				s.EXPECT().GetAllAds(userId).Return([]domain.Ad{}, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ad := mock_service.NewMockAd(c)
			testCase.mockBehavior(ad, testCase.userId)

			services := &service.Service{Ad: ad}
			handler := &Handler{services: services}

			r := gin.New()
			r.GET("/getAllAds", func(ctx *gin.Context) {
				ctx.Set(userContext, testCase.userId)
			}, handler.getAllAds)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/getAllAds", bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestCreateAd(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAd, userId string, ad service.Ads)

	testTable := []struct {
		name                 string
		setUserContext       bool
		inputBody            string
		inputAd              inputAd
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:           "ok",
			setUserContext: true,
			inputBody:      `{"title": "someTitle","category": "category/category","description": "someDescription","price": 100,"contacts": {"name":"someName","phone_number":"somePhoneNumber","email":"someEmail","location":"someLocation"}, "published": true, "images_url": ["someImageURL"]}`,
			inputAd: inputAd{
				Title:       "someTitle",
				Category:    "category",
				Description: "someDescription",
				Price:       100,
				Contacts: inputContacts{
					Name:         "someName",
					Phone_number: "somePhoneNumber",
					Email:        "someEmail",
					Location:     "someLocation",
				},
				Published: true,
				ImagesURL: []string{"someImageURL"},
			},
			mockBehavior: func(s *mock_service.MockAd, userId string, ad service.Ads) {
				s.EXPECT().CreateAd(userId, ad).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Empty input field",
			setUserContext:       false,
			inputBody:            `{"category": "category/category","description": "someDescription","price": 100,"contacts": {"name":"someName","phone_number":"somePhoneNumber","email":"someEmail","location":"someLocation"}, "published": true, "images_url": ["someImageURL"]}`,
			mockBehavior:         func(s *mock_service.MockAd, userId string, ad service.Ads) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:           "service fail",
			setUserContext: true,
			inputBody:      `{"title": "someTitle","category": "category/category","description": "someDescription","price": 100,"contacts": {"name":"someName","phone_number":"somePhoneNumber","email":"someEmail","location":"someLocation"}, "published": true, "images_url": ["someImageURL"]}`,
			inputAd: inputAd{
				Title:       "someTitle",
				Category:    "category",
				Description: "someDescription",
				Price:       100,
				Contacts: inputContacts{
					Name:         "someName",
					Phone_number: "somePhoneNumber",
					Email:        "someEmail",
					Location:     "someLocation",
				},
				Published: true,
				ImagesURL: []string{"someImageURL"},
			},
			mockBehavior: func(s *mock_service.MockAd, userId string, ad service.Ads) {
				s.EXPECT().CreateAd(userId, ad).Return(0, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ad := mock_service.NewMockAd(c)
			testCase.mockBehavior(ad, "1", service.Ads{
				Title:       testCase.inputAd.Title,
				Category:    testCase.inputAd.Category,
				Description: testCase.inputAd.Description,
				Price:       testCase.inputAd.Price,
				Contacts:    service.Contacts(testCase.inputAd.Contacts),
				Published:   testCase.inputAd.Published,
				ImagesURL:   testCase.inputAd.ImagesURL,
			})

			services := &service.Service{Ad: ad}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/createAd", func(ctx *gin.Context) {
				if testCase.setUserContext {
					ctx.Set(userContext, "1")
				}
			}, handler.createAd)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/createAd", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestDeleteAd(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAd, userId string, adId string)

	testTable := []struct {
		name                 string
		setUserContext       bool
		adId                 string
		userId               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:           "ok",
			setUserContext: true,
			adId:           "1",
			userId:         "1",
			mockBehavior: func(s *mock_service.MockAd, userId string, adId string) {
				s.EXPECT().DeleteAd(userId, adId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"deleted"`,
		},
		{
			name:           "service error",
			setUserContext: true,
			adId:           "1",
			userId:         "1",
			mockBehavior: func(s *mock_service.MockAd, userId string, adId string) {
				s.EXPECT().DeleteAd(userId, adId).Return(errors.New("service error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ad := mock_service.NewMockAd(c)
			testCase.mockBehavior(ad, testCase.userId, testCase.adId)

			services := &service.Service{Ad: ad}
			handler := Handler{services: services}

			r := gin.New()
			r.DELETE("/deleteAd/:id", func(ctx *gin.Context) {
				if testCase.setUserContext {
					ctx.Set(userContext, "1")
				}
			}, handler.deleteAd)

			req := httptest.NewRequest("DELETE", "/deleteAd/"+testCase.adId, bytes.NewBufferString(""))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}
