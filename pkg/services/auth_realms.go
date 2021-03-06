// Copyright Red Hat

package services

import (
	"context"
	"encoding/json"
	errs "errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	chi "github.com/go-chi/chi/v5"
	"github.com/identitatem/idp-configs-api/pkg/common"
	"github.com/identitatem/idp-configs-api/pkg/db"
	"github.com/identitatem/idp-configs-api/pkg/errors"
	"github.com/identitatem/idp-configs-api/pkg/models"
	"gorm.io/gorm"
)

func GetAuthRealmsForAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authRealms []models.AuthRealm

	// Get account from request header
	account, err := common.GetAccount(r)

	if (err != nil) {        
		errors.RespondWithBadRequest(err.Error(), w)
        return		
	}
	
	// Fetch Auth Realms for specific account from the DB
	result := db.DB.Where("Account = ?", account).Find(&authRealms)

	if result.Error != nil {		
		errors.RespondWithBadRequest(result.Error.Error(), w)
		return
	}

	// TODO: support filtering and searching by name (query param)

	// Respond with auth realms for the account
	json.NewEncoder(w).Encode(&authRealms)	
}

func CreateAuthRealmForAccount(w http.ResponseWriter, r *http.Request) {
    var authRealm *models.AuthRealm

	w.Header().Set("Content-Type", "application/json")

	// Get Account from request context
	account, err := common.GetAccount(r)
	if (err != nil) {        
		errors.RespondWithBadRequest(err.Error(), w)
        return		
	}	

	// Parse the request payload
	authRealm, err = authRealmFromRequestBody(r.Body)

	if err != nil {
		errors.RespondWithBadRequest(err.Error(), w)
		return
	}	

	// Request body must contain the auth-realm name and custom_resource
	if (authRealm.Name == "" || authRealm.CustomResource == nil) {		
		errors.RespondWithBadRequest("The request body must contain 'name' and 'custom_resource'", w)
		return	
	}

	// TODO: Additional validations on Custom Resource/ evaluating checksum

	// If the request body contains an Account number, it should match the requestor's Account number retrieved from the reqeust context
	if (authRealm.Account != "") {
		err = validateAccount(authRealm, account)
		if (err != nil) {
			errors.RespondWithBadRequest("Account in the request body does not match account for the authenticated user", w)
			return			
		}		
	}	
	// Set account for the auth realm record
	authRealm.Account = account

	// Create record for auth realm in the DB		
	tx := db.DB.Create(&authRealm)
	if tx.Error != nil {
		errorMessage := tx.Error.Error()
		if (strings.Contains(strings.ToLower(errorMessage), "unique constraint")) {	// The error message looks a little different between sqlite and postgres 
			// Unique constraint violated (return 409)
			errors.RespondWithConflict("Error creating record in the DB: " + tx.Error.Error(), w)			
		} else {
			// Error updating the DB		
			errors.RespondWithInternalServerError("Error creating record in the DB: " + tx.Error.Error(), w)
		}
		return			
	}

	// Return ID for created record (Temporily responding with the complete record)
	authRealmJSON, _ := json.Marshal(authRealm)	
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(authRealmJSON))
}

type key int

const AuthRealmKey key = 0

// AuthRealmCtx is a handler for Auth Realm requests
func AuthRealmCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var authRealm models.AuthRealm
		// Get account from request header
		account, err := common.GetAccount(r)

		if (err != nil) {        
			errors.RespondWithBadRequest(err.Error(), w)
			return		
		}
		if authRealmID := chi.URLParam(r, "id"); authRealmID != "" {
			// Fetch record based on Auth Realm ID
			result := db.DB.First(&authRealm, authRealmID)

			if (errs.Is(result.Error, gorm.ErrRecordNotFound)) {
				errors.RespondWithNotFound(result.Error.Error(), w)	// Record not found
				return					
			} else if (result.Error != nil) {	// other error
				errors.RespondWithInternalServerError(result.Error.Error(), w)
				return
			}
			
			if authRealm.Account != "" {
				// Check that the requestor's account matches the account in the DB
				err = validateAccount(&authRealm, account)
				if (err != nil) {        
					errors.RespondWithForbidden("Requestor's account does not match the Auth Realm account", w)
					return		
				}				
			}
			ctx := context.WithValue(r.Context(), AuthRealmKey, &authRealm)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})	
}

func GetAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	// Respond with auth realm
	if authRealm := getAuthRealm(w, r); authRealm != nil {
		json.NewEncoder(w).Encode(authRealm)
	}	
}

func UpdateAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	authRealm := getAuthRealm(w, r)	// Auth realm looked up from DB based on auth realm ID from the request
	if authRealm == nil {
		return
	}	

	incoming, err := authRealmFromRequestBody(r.Body)	// Auth realm from the request body
	if err != nil {
		errors.RespondWithBadRequest(err.Error(), w)
		return
	}
	
	// TODO: Validations on payload

	if (incoming.Account != "" && incoming.Account != authRealm.Account) {
		errors.RespondWithBadRequest("Account number in request body does not match the auth realm account", w)
		return
	}

	// If both name and custom_resource are missing in the request body, there is nothing to update
	if (incoming.Name == "" || incoming.CustomResource == nil) {		
		errors.RespondWithBadRequest("The request body must contain 'name' and 'custom_resource' for update", w)
		return	
	}
	
	incoming.ID = authRealm.ID
	incoming.Account = authRealm.Account
	incoming.CreatedAt = authRealm.CreatedAt

	// Save updates to the DB
	if err := db.DB.Save(&incoming).Error; err != nil {
		if (strings.Contains(strings.ToLower(err.Error()), "unique constraint")) {	// The error message looks a little different between sqlite and postgres 
			// Unique constraint violated (return 409)
			errors.RespondWithConflict("Error updating record in the DB: " + err.Error(), w)			
		} else {
			// Error updating the DB		
			errors.RespondWithInternalServerError("Error updating record in the DB: " + err.Error(), w)
		}
		return
	}

	json.NewEncoder(w).Encode(incoming)
}

func DeleteAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	authRealm := getAuthRealm(w, r)	// Auth realm looked up from DB based on auth realm ID from the request
	if authRealm == nil {
		return
	}

	// Implementing a soft delete for now (The records will continue to exist in the DB with the "DeletedAt" field set. They will not be retrievable by regular queries, but using the )
	if err := db.DB.Delete(&authRealm).Error; err != nil {
		errors.RespondWithInternalServerError(err.Error(), w)
		return
	}		

	fmt.Fprintf(w, "Auth realm with ID %d was successfully deleted", authRealm.ID)	
}

func validateAccount (authRealm *models.AuthRealm, account string) (error) {	
	if (authRealm.Account != "" && authRealm.Account != account) {
		// Account in the request body must match the account of the authenticated user
		return fmt.Errorf("mismatch in account")
	}

	return nil
}

func getAuthRealm(w http.ResponseWriter, r *http.Request) *models.AuthRealm {
	ctx := r.Context()
	authRealm, ok := ctx.Value(AuthRealmKey).(*models.AuthRealm)
	if !ok {
		errors.RespondWithBadRequest("The request must include an auth realm id", w)
		return nil
	}	
	return authRealm
}

func authRealmFromRequestBody(rc io.ReadCloser) (*models.AuthRealm, error) {
	defer rc.Close()
	var authRealm models.AuthRealm
	err := json.NewDecoder(rc).Decode(&authRealm)

	switch {
	case err == io.EOF:
		// empty request body
		return nil, fmt.Errorf("request body must not be empty")		
	case err != nil:
		// other error		
		return nil, err
	}	

	return &authRealm, err
}
