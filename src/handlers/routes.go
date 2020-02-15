package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type EncodeOrErrorInterface struct {
	Write     http.ResponseWriter
	Error     error
	ErrorCode int
	Encode    interface{}
}

func EncodeOrError(opt EncodeOrErrorInterface) {
	log.Print("########### 1")
	if opt.Error != nil {
		log.Print("########### 2")
		EncodeMessage(opt.Write, opt.Error, opt.ErrorCode)
	} else {
		log.Print("########### 3")
		err := json.NewEncoder(opt.Write).Encode(opt.Encode)
		if err != nil {
			log.Print("########### 4")
			EncodeMessage(opt.Write, errors.New("Error encoding data to JSON"), 400)
		}

	}
}

type messageError struct {
	Msg string `json:"msg"`
}

func EncodeMessage(w http.ResponseWriter, error error, httpCode int) { //httpCode ...int) {
	log.Print(error.Error())
	_msg := messageError{}
	_msg.Msg = error.Error()
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(_msg)
}
