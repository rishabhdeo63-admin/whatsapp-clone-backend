package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/controller"
	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/dao"
	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/service"
	"github.com/stretchr/testify/assert"
)

// Mock ContactService implementing interface
type MockContactService struct{ service.ContactService }

func (m *MockContactService) AddContact(uid, cid int64) error { return nil }
func (m *MockContactService) GetContacts(uid int64) ([]*dao.ContactDAO, error) {
	return []*dao.ContactDAO{{UserID: 2, ContactID: 3}}, nil
}
func (m *MockContactService) RemoveContact(uid, cid int64) error { return nil }

func TestAddContactHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	ctrl := controller.NewContactController(&MockContactService{})
	router.POST("/contacts/add", ctrl.AddContact)

	body := map[string]int{"contact_id": 2}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/contacts/add", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestListContactsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	ctrl := controller.NewContactController(&MockContactService{})
	router.GET("/contacts/list", ctrl.GetContacts)

	req, _ := http.NewRequest("GET", "/contacts/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
