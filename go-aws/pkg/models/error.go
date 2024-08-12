package models

type CaptionError struct {
	Stage        string // stage when this error occurred
	Context      string // more specific context for what we were attempting when this error happened
	ErrorMessage string // actual error message
}
