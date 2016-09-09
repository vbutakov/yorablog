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

	// BasePhotosPath is the path for photos
	BasePhotosPath string

	// BaseTemplatesPath is the path for templates files
	BaseTemplatesPath string

	// BaseDSN is the dsn for db connect
	BaseDSN string

	//SMTP params
	SMTPServer string
	SMTPPort   string
	SMTPUser   string
	SMTPPass   string
	SMTPFrom   string
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

	BasePhotosPath = os.Getenv("BASEPHOTOSPATH")
	if len(BasePhotosPath) == 0 {
		return errors.New("Error! Env variable BASEPHOTOSPATH is not defined.\n")
	}

	BaseTemplatesPath = os.Getenv("BASETEMPLATESPATH")
	if len(BaseTemplatesPath) == 0 {
		return errors.New("Error! Env variable BASETEMPLATESPATH is not defined.\n")
	}

	BaseDSN = os.Getenv("BASEDSN")
	if len(BaseDSN) == 0 {
		return errors.New("Error! Env variable BASEDSN is not defined.\n")
	}

	SMTPServer = os.Getenv("SMTPSERVER")
	if len(SMTPServer) == 0 {
		return errors.New("Error! Env variable SMTPSERVER is not defined.\n")
	}

	SMTPPort = os.Getenv("SMTPPORT")
	if len(SMTPPort) == 0 {
		return errors.New("Error! Env variable SMTPPORT is not defined.\n")
	}

	SMTPUser = os.Getenv("SMTPUSER")
	if len(SMTPUser) == 0 {
		return errors.New("Error! Env variable SMTPUSER is not defined.\n")
	}

	SMTPPass = os.Getenv("SMTPPASS")
	if len(SMTPPass) == 0 {
		return errors.New("Error! Env variable SMTPPASS is not defined.\n")
	}

	SMTPFrom = os.Getenv("SMTPFROM")
	if len(SMTPFrom) == 0 {
		return errors.New("Error! Env variable SMTPFrom is not defined.\n")
	}

	return nil
}
