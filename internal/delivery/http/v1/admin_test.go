package v1

import (
	"bytes"
	"errors"
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"

	mock_service "github.com/TakoB222/postingAds-api/internal/service/mocks"
)

func TestAdminSignIn(t *testing.T){
	type mockBehavior func(s *mock_service.MockAdmin, input service.SignInInput)

	testTable := []struct{
		name string
		inputBody string
		inputSignIn adminSignInInput
		mockBehavior mockBehavior
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "ok",
			inputBody: `{"email":"example@gmail.com","password":"somePassword"}`,
			inputSignIn: adminSignInInput{
				Email: "example@gmail.com",
				Password: "somePassword",
			},
			mockBehavior: func(s *mock_service.MockAdmin, input service.SignInInput) {
				s.EXPECT().AdminSignIn(input).Return(service.Tokens{AccessToken: "AccessToken", RefreshToken: "RefreshToken"}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"access_token":"AccessToken","refresh_token":"RefreshToken"}`,
		},
		{
			name: "empty input field",
			inputBody: `{"email":"example@gmail.com"}`,
			inputSignIn: adminSignInInput{
				Email: "example@gmail.com",
			},
			mockBehavior: func(s *mock_service.MockAdmin, input service.SignInInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "ok",
			inputBody: `{"email":"example@gmail.com","password":"somePassword"}`,
			inputSignIn: adminSignInInput{
				Email: "example@gmail.com",
				Password: "somePassword",
			},
			mockBehavior: func(s *mock_service.MockAdmin, input service.SignInInput) {
				s.EXPECT().AdminSignIn(input).Return(service.Tokens{}, errors.New("service failure"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admin := mock_service.NewMockAdmin(c)
			testCase.mockBehavior(admin, service.SignInInput{
				Email: testCase.inputSignIn.Email,
				Password: testCase.inputSignIn.Password,
			})

			services := &service.Service{Admin:admin}
			handler := Handler{services:services}

			r := gin.New()
			r.POST("/adminSignIn", handler.adminSignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/adminSignIn", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestAdminRefreshTokens(t *testing.T){
	type mockBehavior func(s *mock_service.MockAdmin, input service.RefreshInput)

	testTable := []struct{
		name string
		inputBody string
		inputRefresh refreshTokensInput
		mockBehavior mockBehavior
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "ok",
			inputBody: `{"RefreshToken":"token"}`,
			inputRefresh: refreshTokensInput{
				RefreshToken: "token",
			},
			mockBehavior: func(s *mock_service.MockAdmin, input service.RefreshInput) {
				s.EXPECT().AdminRefreshSession(input).Return(service.Tokens{AccessToken: "AccessToken", RefreshToken: "RefreshToken"}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"access_token":"AccessToken","refresh_token":"RefreshToken"}`,
		},
		{
			name: "invalid input body",
			inputBody: `{"RefreshToken":"token}`,
			inputRefresh: refreshTokensInput{},
			mockBehavior: func(s *mock_service.MockAdmin, input service.RefreshInput) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "service error",
			inputBody: `{"RefreshToken":"token"}`,
			inputRefresh: refreshTokensInput{
				RefreshToken: "token",
			},
			mockBehavior: func(s *mock_service.MockAdmin, input service.RefreshInput) {
				s.EXPECT().AdminRefreshSession(input).Return(service.Tokens{}, errors.New("service failure"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admin := mock_service.NewMockAdmin(c)
			testCase.mockBehavior(admin, service.RefreshInput{
				RefreshToken: testCase.inputRefresh.RefreshToken,
			})

			services := &service.Service{Admin:admin}
			handler := Handler{services:services}

			r := gin.New()
			r.POST("/adminRefreshTokens", handler.adminRefreshTokens)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/adminRefreshTokens", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestAdminGetAllAds(t *testing.T){
	type mockBehavior func(s *mock_service.MockAdmin)

	testTable := []struct{
		name string
		mockBehavior mockBehavior
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(s *mock_service.MockAdmin) {
				s.EXPECT().AdminGetAllAdsByAdmin().Return([]domain.Ad{}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[]`,
		},
		{
			name: "service error",
			mockBehavior: func(s *mock_service.MockAdmin) {
				s.EXPECT().AdminGetAllAdsByAdmin().Return([]domain.Ad{}, errors.New("service failure"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admin := mock_service.NewMockAdmin(c)
			testCase.mockBehavior(admin)

			services := &service.Service{Admin:admin}
			handler := Handler{services:services}

			r := gin.New()
			r.GET("/adminGetAllAds", handler.adminGetAllAds)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/adminGetAllAds", bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestAdminGetAd(t *testing.T){
	type mockBehavior func(s *mock_service.MockAdmin)

	testTable := []struct{
		name string
		mockBehavior mockBehavior
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(s *mock_service.MockAdmin) {
				s.EXPECT().AdminGetAd("1").Return(domain.Ad{}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"Id":0,"UserId":"","Title":"","Category":"","Description":"","Price":0,"Contacts":"","Published":false,"ImagesURL":null}`,
		},
		{
			name: "service error",
			mockBehavior: func(s *mock_service.MockAdmin) {
				s.EXPECT().AdminGetAd("1").Return(domain.Ad{}, errors.New("service failure"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admin := mock_service.NewMockAdmin(c)
			testCase.mockBehavior(admin)

			services := &service.Service{Admin:admin}
			handler := Handler{services:services}

			r := gin.New()
			r.GET("/adminGetAd/:id", handler.adminGetAd)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/adminGetAd/1", bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestAdminUpdateAd(t *testing.T){
	type mockBehavior func(s *mock_service.MockAdmin, ad service.Ads)

	testTable := []struct{
		name string
		inputBody string
		inputAd adminUpdateAdInput
		mockBehavior mockBehavior
		expectedStatusCode int
		expectedResponseBody string
	}{
		{
			name: "ok",
			inputBody: `{"title":"title","category":"category","description":"description","price":100,"contacts":{"name":"name", "phone_number":"number", "email":"email", "location":"location"},"published":true,"images_url":["url"]}`,
			inputAd: adminUpdateAdInput{
				Title: "title",
				Category: "category",
				Description: "description",
				Price: 100,
				Contacts: adminInputUpdateContacts{
					Name: "name",
					Phone_number: "number",
					Email: "email",
					Location: "location",
				},
				Published: true,
				ImagesURL: []string{"url"},
			},
			mockBehavior: func(s *mock_service.MockAdmin, ad service.Ads) {
				s.EXPECT().AdminUpdateAd("1", ad).Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `"updated"`,
		},
		{
			name: "Empty input field",
			inputBody: `{"title":"title","description":"description","price":100,"contacts":{"name":"name", "phone_number":"number", "email":"email", "location":"location"},"published":true,"images_url":["url"]}`,
			inputAd: adminUpdateAdInput{},
			mockBehavior: func(s *mock_service.MockAdmin, ad service.Ads) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "service failure",
			inputBody: `{"title":"title","category":"category","description":"description","price":100,"contacts":{"name":"name", "phone_number":"number", "email":"email", "location":"location"},"published":true,"images_url":["url"]}`,
			inputAd: adminUpdateAdInput{
				Title: "title",
				Category: "category",
				Description: "description",
				Price: 100,
				Contacts: adminInputUpdateContacts{
					Name: "name",
					Phone_number: "number",
					Email: "email",
					Location: "location",
				},
				Published: true,
				ImagesURL: []string{"url"},
			},
			mockBehavior: func(s *mock_service.MockAdmin, ad service.Ads) {
				s.EXPECT().AdminUpdateAd("1", ad).Return(errors.New("service failure"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admin := mock_service.NewMockAdmin(c)
			testCase.mockBehavior(admin, service.Ads{
				Title: testCase.inputAd.Title,
				Category: testCase.inputAd.Category,
				Description: testCase.inputAd.Description,
				Price: testCase.inputAd.Price,
				Contacts: service.Contacts(testCase.inputAd.Contacts),
				Published: testCase.inputAd.Published,
				ImagesURL: testCase.inputAd.ImagesURL,
			})

			services := &service.Service{Admin:admin}
			handler := Handler{services:services}

			r := gin.New()
			r.PUT("/adminUpdateAd/:id", handler.adminUpdateAd)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/adminUpdateAd/1", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
