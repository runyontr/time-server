package main

import (
	t "time"
)


type TimeService interface {
	GetTime() (time, error)
}


func New() (TimeService, error) {

	return &timeServiceImplementation{}, nil

}


//timeServiceImplementation is the implementation of the TimeService interface.
type timeServiceImplementation struct {
}

//GetAppInfo returns the app info of the running application
func (s *timeServiceImplementation) GetTime() (time, error) {

	info := time{}
	info.CurrentTimeMillis=  t.Now().UnixNano() / int64(t.Millisecond)

	return info, nil
}

