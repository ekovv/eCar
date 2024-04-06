package handler

import (
	"bytes"
	"eCar/config"
	"eCar/internal/constants"
	"eCar/internal/domains/mocks"
	"eCar/internal/shema"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type serviceMock func(c *mocks.Service)

func TestHandler_GetNewCars(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cars := []shema.Car{
		{
			RegNum: "DDDDDD",
			Mark:   "Honda",
			Model:  "Civic",
			Year:   2015,
			Owner: shema.People{
				Name:       "Petr",
				Surname:    "Petrov",
				Patronymic: "Petrovich",
			},
		},
	}

	tests := []struct {
		name        string
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK#1",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("AddCar", mock.Anything, cars).Return(nil).Times(1)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "BAD",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("AddCar", mock.Anything, mock.Anything).Return(errors.New("can't add car")).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			service := mocks.NewService(t)
			h := NewHandler(service, config.Config{})
			tt.serviceMock(service)

			httpmock.RegisterResponder("GET", "http://localhost:63342/info",
				httpmock.NewStringResponder(200, `{"regNum":"DDDDDD","mark":"Honda","model":"Civic","year":2015,"owner":{"name":"Petr","surname":"Petrov","patronymic":"Petrovich"}}`))

			path := "/api/add"
			g.POST(path, h.GetNewCars)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer([]byte(`{"regNums":["DDDDDD"]}`)))
			request.Header.Set("Content-Type", "application/json")

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_GetData(t *testing.T) {
	filter := shema.Filter{
		RegNum: "DDDDDD",
		Mark:   "Honda",
		Page:   1,
		Limit:  2,
	}

	tests := []struct {
		name        string
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK#1",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetData", mock.Anything, filter.RegNum, filter.Mark, filter.Model, filter.Year, filter.OwnerName, filter.OwnerSurname,
					filter.OwnerPatronymic, filter.Page, filter.Limit).Return([]shema.Car{{
					RegNum: "DDDDDD",
					Mark:   "Honda",
					Model:  "Civic",
					Year:   2015,
					Owner: shema.People{
						Name:       "Petr",
						Surname:    "Petrov",
						Patronymic: "Petrovich",
					},
				},
				}, nil).Times(1)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "BAD",
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetData", mock.Anything, filter.RegNum, filter.Mark, filter.Model, filter.Year, filter.OwnerName, filter.OwnerSurname,
					filter.OwnerPatronymic, filter.Page, filter.Limit).Return(nil, constants.ErrInvalidData).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			service := mocks.NewService(t)
			h := NewHandler(service, config.Config{})
			tt.serviceMock(service)

			path := "/api/all"
			g.POST(path, h.GetData)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer([]byte(`{"regNum":"DDDDDD", "mark": "Honda", "page": 1, "limit": 2}`)))
			request.Header.Set("Content-Type", "application/json")

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_DeleteData(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK#1",
			id:   2,
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("DeleteCar", mock.Anything, 2).Return(nil).Times(1)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "BAD",
			id:   1231,
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("DeleteCar", mock.Anything, 1231).Return(constants.ErrInvalidData).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			service := mocks.NewService(t)
			h := NewHandler(service, config.Config{})
			tt.serviceMock(service)

			path := "/api/delete/:id"
			g.DELETE(path, h.DeleteData)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/delete/%d", tt.id), nil)
			request.Header.Set("Content-Type", "application/json")

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_UpdateData(t *testing.T) {
	filter := shema.Filter{
		RegNum: "DDDDDD",
		Mark:   "Honda",
	}
	tests := []struct {
		name        string
		id          int
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK#1",
			id:   2,
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("UpdateCar", mock.Anything, 2, filter).Return(nil).Times(1)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "BAD",
			id:   1231,
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("UpdateCar", mock.Anything, 1231, filter).Return(constants.ErrInvalidData).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			service := mocks.NewService(t)
			h := NewHandler(service, config.Config{})
			tt.serviceMock(service)

			path := "/api/update/:id"
			g.PUT(path, h.UpdateData)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/update/%d", tt.id), bytes.NewBuffer([]byte(`{"regNum":"DDDDDD", "mark": "Honda"}`)))
			request.Header.Set("Content-Type", "application/json")

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}
