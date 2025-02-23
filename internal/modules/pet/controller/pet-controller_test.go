package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/infrastructure/responder"
	"pet-store/internal/models"
	"pet-store/internal/modules/pet/service"
	"pet-store/internal/modules/pet/service/mocks"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockPetService struct{}

func TestUpdatePetForm(t *testing.T) {
	cases := []struct {
		nameTest  string
		petId     string
		namePet   string
		statusPet string
		respError string
		mockError int
	}{
		{
			nameTest:  "Success",
			petId:     "12",
			namePet:   "Jora",
			statusPet: "available",
		},
		{
			nameTest:  "Error in Mck",
			petId:     "10",
			namePet:   "Egorka",
			statusPet: "available",
			mockError: 2022,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.nameTest, func(t *testing.T) {
			t.Parallel()

			serviceUpdateSaveFormMock := mocks.NewPeter(t)

			if serviceUpdateSaveFormMock == nil {
				t.Fatal("mock service is nil")
			}

			if tc.respError == "" && tc.mockError == 0 {
				serviceUpdateSaveFormMock.On("UpdatePetForm", mock.Anything, tc.namePet, tc.statusPet, tc.petId).
					Return(service.RequestOut{
						Status:    true,
						ErrorCode: errors.NoError,
					}).
					Once()
			}
			if tc.mockError == 2022 {
				serviceUpdateSaveFormMock.On("UpdatePetForm", mock.Anything, tc.namePet, tc.statusPet, tc.petId).
					Return(service.RequestOut{
						Status:    false,
						ErrorCode: errors.UpdatePetFormError,
					}).
					Once()
			}

			decoder := godecoder.NewDecoder(jsoniter.Config{
				EscapeHTML:             true,
				SortMapKeys:            true,
				ValidateJsonRawMessage: true,
				DisallowUnknownFields:  true,
			})
			logger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			responseManager := responder.NewResponder(decoder, logger)

			petHandler := &Pet{
				service:   serviceUpdateSaveFormMock,
				Responder: responseManager,
				Decoder:   decoder,
			}
			var reqFrom []string
			if tc.namePet != "" {
				reqFrom = append(reqFrom, fmt.Sprintf("name=%s", tc.namePet))
			}
			if tc.statusPet != "" {
				reqFrom = append(reqFrom, fmt.Sprintf("status=%s", tc.statusPet))
			}

			reqBody := []byte(strings.Join(reqFrom, "&"))

			httpReq, err := http.NewRequest(http.MethodPost, "/pet/"+tc.petId, bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			// Устанавливаем заголовок Content-Type
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Запускаем HTTP-сервер
			rr := httptest.NewRecorder()

			// Настроим маршрут
			router := chi.NewRouter()
			router.Post("/pet/{petId}", petHandler.UpdatePetForm)

			router.ServeHTTP(rr, httpReq)

			// Проверяем, что мок был вызван
			serviceUpdateSaveFormMock.AssertExpectations(t)

			// Проверка ответа
			if tc.mockError != 0 {
				// Если произошла ошибка, проверяем ответ с ошибкой
				if tc.mockError == 2022 {
					assert.Equal(t, http.StatusInternalServerError, rr.Code)
				}
			} else {
				// Если успешный ответ
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestFindPetbyID(t *testing.T) {
	cases := []struct {
		nameTest  string
		petId     string
		pet       models.Pet
		mockError int
	}{
		{
			nameTest: "Success",
			petId:    "12",
			pet: models.Pet{
				ID: 12,
				Category: models.Category{
					ID:   1,
					Name: "Dog",
				},
				Name:      "Juchka",
				PhotoUrls: []string{"url1", "url2"},
				Tags: []models.Tag{
					{ID: 1, Name: "frendly"},
				},
				Status: "available",
			},
		},
		{
			nameTest:  "Error Not Found",
			petId:     "12",
			mockError: errors.PetServiceErrPetNotFound,
		},
		{
			nameTest:  "Error on server",
			petId:     "12",
			mockError: errors.PetServiceFindPetbyID,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.nameTest, func(t *testing.T) {
			t.Parallel()

			serviceMock := mocks.NewPeter(t)
			if serviceMock == nil {
				t.Fatal("mock service is nil")
			}

			if tc.mockError == 0 {
				serviceMock.On("FindPetbyID", mock.Anything, tc.petId).
					Return(service.RequestOutWithPet{
						Pet: models.Pet{
							ID: 12,
							Category: models.Category{
								ID:   1,
								Name: "Dog",
							},
							Name:      "Juchka",
							PhotoUrls: []string{"url1", "url2"},
							Tags: []models.Tag{
								{ID: 1, Name: "frendly"},
							},
							Status: "available",
						},
						Status:    true,
						ErrorCode: errors.NoError,
					}).
					Once()
			}
			if tc.mockError != 0 {
				if tc.mockError == errors.PetServiceFindPetbyID {
					serviceMock.On("FindPetbyID", mock.Anything, tc.petId).
						Return(service.RequestOutWithPet{
							Status:    false,
							ErrorCode: errors.PetServiceFindPetbyID,
						}).
						Once()
				}
				if tc.mockError == errors.PetServiceErrPetNotFound {
					serviceMock.On("FindPetbyID", mock.Anything, tc.petId).
						Return(service.RequestOutWithPet{
							Status:    false,
							ErrorCode: errors.PetServiceErrPetNotFound,
						}).
						Once()
				}
			}
			decoder := godecoder.NewDecoder(jsoniter.Config{
				EscapeHTML:             true,
				SortMapKeys:            true,
				ValidateJsonRawMessage: true,
				DisallowUnknownFields:  true,
			})
			logger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			responseManager := responder.NewResponder(decoder, logger)

			petHandler := &Pet{
				service:   serviceMock,
				Responder: responseManager,
				Decoder:   decoder,
			}

			httpReq, err := http.NewRequest(http.MethodGet, "/pet/"+tc.petId, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Запускаем HTTP-сервер
			rr := httptest.NewRecorder()

			// Настроим маршрут
			router := chi.NewRouter()
			router.Get("/pet/{petId}", petHandler.FindPetbyID)

			router.ServeHTTP(rr, httpReq)

			// Проверяем, что мок был вызван
			serviceMock.AssertExpectations(t)

			// Проверка ответа
			if tc.mockError != 0 {
				// Если произошла ошибка, проверяем ответ с ошибкой
				if tc.mockError == errors.PetServiceErrPetNotFound {
					assert.Equal(t, http.StatusNotFound, rr.Code)
				}
				if tc.mockError == errors.PetServiceFindPetbyID {
					assert.Equal(t, http.StatusInternalServerError, rr.Code)
				}
			} else {
				// Если успешный ответ
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestAddPet(t *testing.T) {
	cases := []struct {
		nameTest  string
		req       service.PetAddRequest
		mockError int
	}{
		{
			nameTest: "Success",
			req: service.PetAddRequest{
				Category:  service.Category{Name: "Dog"},
				Name:      "Rex",
				PhotoUrls: []string{"http://example.com/photo1.jpg", "http://example.com/photo2.jpg"},
				Tags: []service.Tag{
					{Name: "Friendly"},
					{Name: "Loyal"},
				},
				Status: "available",
			},
		},
		{
			nameTest: "Error",
			req: service.PetAddRequest{
				Category:  service.Category{Name: "Dog"},
				Name:      "Rex",
				PhotoUrls: []string{"http://example.com/photo1.jpg", "http://example.com/photo2.jpg"},
				Tags: []service.Tag{
					{Name: "Friendly"},
					{Name: "Loyal"},
				},
				Status: "available",
			},
			mockError: errors.AddPetErr,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.nameTest, func(t *testing.T) {
			t.Parallel()

			serviceMock := mocks.NewPeter(t)
			if serviceMock == nil {
				t.Fatal("mock service is nil")
			}

			if tc.mockError == 0 {
				serviceMock.On("AddPet", mock.Anything, tc.req).
					Return(service.RequestOut{
						Status:    true,
						ErrorCode: errors.NoError,
					}).
					Once()
			}
			if tc.mockError != 0 {
				if tc.mockError == errors.AddPetErr {
					serviceMock.On("AddPet", mock.Anything, tc.req).
						Return(service.RequestOut{
							Status:    false,
							ErrorCode: errors.AddPetErr,
						}).
						Once()
				}
			}

			decoder := godecoder.NewDecoder(jsoniter.Config{
				EscapeHTML:             true,
				SortMapKeys:            true,
				ValidateJsonRawMessage: true,
				DisallowUnknownFields:  true,
			})
			logger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			responseManager := responder.NewResponder(decoder, logger)

			petHandler := &Pet{
				service:   serviceMock,
				Responder: responseManager,
				Decoder:   decoder,
			}

			petRequest := service.PetAddRequest{
				Category:  service.Category{Name: "Dog"},
				Name:      "Rex",
				PhotoUrls: []string{"http://example.com/photo1.jpg", "http://example.com/photo2.jpg"},
				Tags: []service.Tag{
					{Name: "Friendly"},
					{Name: "Loyal"},
				},
				Status: "available",
			}

			// Сериализуем в JSON
			jsonData, err := json.Marshal(petRequest)
			if err != nil {
				t.Fatalf("Error marshalling to JSON: %s", err)
			}

			httpReq, err := http.NewRequest(http.MethodPost, "/pet", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatal(err)
			}
			httpReq.Header.Set("Content-Type", "application/json")

			// Запускаем HTTP-сервер
			rr := httptest.NewRecorder()

			// Настроим маршрут
			router := chi.NewRouter()
			router.Post("/pet", petHandler.AddPet)

			router.ServeHTTP(rr, httpReq)

			// Проверяем, что мок был вызван
			serviceMock.AssertExpectations(t)

			// Проверка ответа
			if tc.mockError != 0 {
				// Если произошла ошибка, проверяем ответ с ошибкой
				if tc.mockError == errors.AddPetErr {
					assert.Equal(t, http.StatusInternalServerError, rr.Code)
				}
			} else {
				// Если успешный ответ
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestFindPetbyStatus(t *testing.T) {
	cases := []struct {
		nameTest  string
		statusPet []string
		pets      []models.Pet
		mockError int
	}{
		{
			nameTest:  "Success",
			statusPet: []string{"available"},
			pets: []models.Pet{
				{
					ID: 12,
					Category: models.Category{
						ID:   1,
						Name: "Dog",
					},
					Name:      "Juchka",
					PhotoUrls: []string{"url1", "url2"},
					Tags: []models.Tag{
						{ID: 1, Name: "frendly"},
					},
					Status: "available",
				},
				{
					ID: 14,
					Category: models.Category{
						ID:   1,
						Name: "Cat",
					},
					Name:      "Fedor",
					PhotoUrls: []string{"url1", "url2"},
					Tags: []models.Tag{
						{ID: 1, Name: "frendly"},
					},
					Status: "available",
				},
			},
		},
		{
			nameTest:  "Error",
			statusPet: []string{"available", "pending"},
			mockError: errors.PetServiceFindPetbyStatus,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.nameTest, func(t *testing.T) {
			t.Parallel()

			serviceMock := mocks.NewPeter(t)
			if serviceMock == nil {
				t.Fatal("mock service is nil")
			}

			if tc.mockError == 0 {
				serviceMock.On("FindPetbyStatus", mock.Anything, tc.statusPet).
					Return(service.RequestOutWithPets{
						Status:    true,
						ErrorCode: errors.NoError,
						Pets: []models.Pet{
							{
								ID: 12,
								Category: models.Category{
									ID:   1,
									Name: "Dog",
								},
								Name:      "Juchka",
								PhotoUrls: []string{"url1", "url2"},
								Tags: []models.Tag{
									{ID: 1, Name: "frendly"},
								},
								Status: "available",
							},
							{
								ID: 14,
								Category: models.Category{
									ID:   1,
									Name: "Cat",
								},
								Name:      "Fedor",
								PhotoUrls: []string{"url1", "url2"},
								Tags: []models.Tag{
									{ID: 1, Name: "frendly"},
								},
								Status: "available",
							},
						},
					}).
					Once()
			}
			if tc.mockError != 0 {
				if tc.mockError == errors.PetServiceFindPetbyStatus {
					serviceMock.On("FindPetbyStatus", mock.Anything, tc.statusPet).
						Return(service.RequestOutWithPets{
							Status:    false,
							ErrorCode: errors.PetServiceFindPetbyStatus,
						}).
						Once()
				}
			}

			decoder := godecoder.NewDecoder(jsoniter.Config{
				EscapeHTML:             true,
				SortMapKeys:            true,
				ValidateJsonRawMessage: true,
				DisallowUnknownFields:  true,
			})
			logger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			responseManager := responder.NewResponder(decoder, logger)

			petHandler := &Pet{
				service:   serviceMock,
				Responder: responseManager,
				Decoder:   decoder,
			}
			var strUrl string
			for i := range tc.statusPet {
				strUrl += fmt.Sprintf("status=%s&", tc.statusPet[i])
			}
			strUrl = strUrl[:len(strUrl)-1]

			httpReq, err := http.NewRequest(http.MethodGet, "/pet/findByStatus"+"?"+strUrl, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Запускаем HTTP-сервер
			rr := httptest.NewRecorder()

			// Настроим маршрут
			router := chi.NewRouter()
			router.Get("/pet/findByStatus", petHandler.FindPetbyStatus)

			router.ServeHTTP(rr, httpReq)

			// Проверяем, что мок был вызван
			serviceMock.AssertExpectations(t)

			// Проверка ответа
			if tc.mockError != 0 {
				// Если произошла ошибка, проверяем ответ с ошибкой
				if tc.mockError == errors.PetServiceFindPetbyStatus {
					assert.Equal(t, http.StatusInternalServerError, rr.Code)
				}
			} else {
				// Если успешный ответ
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestUpdatePet(t *testing.T) {
	cases := []struct {
		nameTest     string
		updatePetReq service.PetUpdateRequest
		mockError    int
	}{
		{
			nameTest: "Success",
			updatePetReq: service.PetUpdateRequest{
				Category:  service.Category{Name: "Dog"},
				Name:      "Rex",
				PhotoUrls: []string{"http://example.com/photo1.jpg", "http://example.com/photo2.jpg"},
				Tags: []service.Tag{
					{Name: "Friendly"},
					{Name: "Loyal"},
				},
				Status: "available",
			},
		},
		{
			nameTest: "Error",
			updatePetReq: service.PetUpdateRequest{
				Category:  service.Category{Name: "Dog"},
				Name:      "Rex",
				PhotoUrls: []string{"http://example.com/photo1.jpg", "http://example.com/photo2.jpg"},
				Tags: []service.Tag{
					{Name: "Friendly"},
					{Name: "Loyal"},
				},
				Status: "available",
			},
			mockError: errors.PetServiceUpdateErr,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.nameTest, func(t *testing.T) {
			t.Parallel()

			serviceMock := mocks.NewPeter(t)
			if serviceMock == nil {
				t.Fatal("mock service is nil")
			}

			if tc.mockError == 0 {
				serviceMock.On("UpdatePet", mock.Anything, tc.updatePetReq).
					Return(service.RequestOut{
						Status:    true,
						ErrorCode: errors.NoError,
					}).
					Once()
			}
			if tc.mockError != 0 {
				if tc.mockError == errors.PetServiceUpdateErr {
					serviceMock.On("UpdatePet", mock.Anything, tc.updatePetReq).
						Return(service.RequestOut{
							Status:    false,
							ErrorCode: errors.PetServiceUpdateErr,
						}).
						Once()
				}
			}

			decoder := godecoder.NewDecoder(jsoniter.Config{
				EscapeHTML:             true,
				SortMapKeys:            true,
				ValidateJsonRawMessage: true,
				DisallowUnknownFields:  true,
			})
			logger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			responseManager := responder.NewResponder(decoder, logger)

			petHandler := &Pet{
				service:   serviceMock,
				Responder: responseManager,
				Decoder:   decoder,
			}

			petRequest := service.PetUpdateRequest{
				Category:  service.Category{Name: "Dog"},
				Name:      "Rex",
				PhotoUrls: []string{"http://example.com/photo1.jpg", "http://example.com/photo2.jpg"},
				Tags: []service.Tag{
					{Name: "Friendly"},
					{Name: "Loyal"},
				},
				Status: "available",
			}

			// Сериализуем в JSON
			jsonData, err := json.Marshal(petRequest)
			if err != nil {
				t.Fatalf("Error marshalling to JSON: %s", err)
			}

			httpReq, err := http.NewRequest(http.MethodPut, "/pet", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatal(err)
			}
			httpReq.Header.Set("Content-Type", "application/json")

			// Запускаем HTTP-сервер
			rr := httptest.NewRecorder()

			// Настроим маршрут
			router := chi.NewRouter()
			router.Put("/pet", petHandler.UpdatePet)

			router.ServeHTTP(rr, httpReq)

			// Проверяем, что мок был вызван
			serviceMock.AssertExpectations(t)

			// Проверка ответа
			if tc.mockError != 0 {
				// Если произошла ошибка, проверяем ответ с ошибкой
				if tc.mockError == errors.PetServiceUpdateErr {
					assert.Equal(t, http.StatusInternalServerError, rr.Code)
				}
			} else {
				// Если успешный ответ
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}
