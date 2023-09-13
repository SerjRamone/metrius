// Package handlers ...
package handlers

import "github.com/SerjRamone/metrius/internal/storage"

// BaseHandler ...
type BaseHandler struct {
	Storage storage.Storage
}

// NewBaseHandler creates new BaseHandler
func NewBaseHandler(storage storage.Storage) BaseHandler {
	return BaseHandler{Storage: storage}
}
