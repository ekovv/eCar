package handler

import (
	"eCar/config"
	"eCar/internal/domains"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestHandler_GetNewCars(t *testing.T) {
	type fields struct {
		service domains.Service
		engine  *gin.Engine
		config  config.Config
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Handler{
				service: tt.fields.service,
				engine:  tt.fields.engine,
				config:  tt.fields.config,
			}
			s.GetNewCars(tt.args.c)
		})
	}
}
