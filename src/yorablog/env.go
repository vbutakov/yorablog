package main

import (
	"errors"
	"os"
)

var (
	// BaseServeAddr is the addr for ListenAndServe
	BaseServeAddr string

	// BaseStaticPath is the path for static files
	BaseStaticPath string

	// BaseTemplatesPath is the path for templates files
	BaseTemplatesPath string

	// BaseDSN is the dsn for db connect
	BaseDSN string
)

// InitEnv initialize all global params from env variables
func InitEnv() error {
	BaseServeAddr = os.Getenv("BASESERVEADDR")
	if len(BaseServeAddr) == 0 {
		BaseServeAddr = ":8080"
	}

	BaseStaticPath = os.Getenv("BASESTATICPATH")
	if len(BaseStaticPath) == 0 {
		return errors.New("Error! Env variable BASESTATICPATH is not defined.\n")
	}

	BaseTemplatesPath = os.Getenv("BASETEMPLATESPATH")
	if len(BaseTemplatesPath) == 0 {
		return errors.New("Error! Env variable BASETEMPLATESPATH is not defined.\n")
	}

	BaseDSN = os.Getenv("BASEDSN")
	if len(BaseDSN) == 0 {
		return errors.New("Error! Env variable BASEDSN is not defined.\n")
	}

	return nil
}
