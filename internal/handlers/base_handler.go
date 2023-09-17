// Package handlers ...
package handlers

import "github.com/SerjRamone/metrius/internal/storage"

// baseHandler ...
type baseHandler struct {
	Storage storage.Storage
}

// NewBaseHandler creates new BaseHandler
func NewBaseHandler(storage storage.Storage) baseHandler {
	return baseHandler{Storage: storage}
}
