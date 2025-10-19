package main

import (
	//   "encoding/json"
	"fmt"
	"net/http"
	"errors"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

func (a *application)createRoleHandler(w http.ResponseWriter, r *http.Request) { 
	var incomingData struct {

	 Role  string  `json:"role"`
	}
	
	// err := json.NewDecoder(r.Body).Decode(&incomingData)
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		// a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	
	role := &data.Role {
		Role: incomingData.Role,
	}

	//INitialize a validator instance
	v := validator.New()

	//do the validation
	data.ValidateRole(v, role)
	if !v.IsEmpty() {
		// a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Add the role to the database table
	err = a.roleModel.Insert(role)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
 
	// Set a Location header. The path to the newly created role
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/roles/%d", role.ID))

	// Send a JSON response with 201 (new resource created) status code
	data := envelope{
		"role": role,
	  }
 	err = a.writeJSON(w, http.StatusCreated, data, headers)
 	if err != nil {
	  	a.serverErrorResponse(w, r, err)
	  	return
  	}
	
}

func (a *application) displayRoleHandler (w http.ResponseWriter, r *http.Request) {
   id, err := a.readIDParam(r)
   if err != nil {
       a.notFoundResponse(w, r)
       return 
   }

   // Call Get() to retrieve the quotte with the specified id
   role, err := a.roleModel.Get(id)
   if err != nil {
       switch {
           case errors.Is(err, data.ErrRecordNotFound):
              a.notFoundResponse(w, r)
           default:
              a.serverErrorResponse(w, r, err)
       }
       return 
   }

   // display the role
   data := envelope {
	"role": role,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
	a.serverErrorResponse(w, r, err)
	return 
	}
}

func (a *application) updateRoleHandler (w http.ResponseWriter, r *http.Request) {

	// Get the id from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return 
	}

	// Call Get() to retrieve the role with the specified id
	role, err := a.roleModel.Get(id)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
			   a.notFoundResponse(w, r)
			default:
			   a.serverErrorResponse(w, r, err)
		}
		return 
	}

	
 	var incomingData struct {
        Role  *string  `json:"role"`
    }  

	// perform the decoding
	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
 	// We need to now check the fields to see which ones need updating
 	// if incomingData.Role is nil, no update was provided
	if incomingData.Role != nil {
		role.Role = *incomingData.Role
	}
 
 	// Before we write the updates to the DB let's validate
	v := validator.New()
	data.ValidateRole(v, role)
	if !v.IsEmpty() {
		 a.failedValidationResponse(w, r, v.Errors)  
		 return
	}

	// perform the update
    err = a.roleModel.Update(role)
    if err != nil {
       a.serverErrorResponse(w, r, err)
       return 
   }
   data := envelope {
                "role": role,
          }
   err = a.writeJSON(w, http.StatusOK, data, nil)
   if err != nil {
       a.serverErrorResponse(w, r, err)
       return 
   }

}

func (a *application) deleteRoleHandler (w http.ResponseWriter, r *http.Request) {

	id, err := a.readIDParam(r)
   if err != nil {
       a.notFoundResponse(w, r)
       return 
   }

   err = a.roleModel.Delete(id)

   if err != nil {
       switch {
           case errors.Is(err, data.ErrRecordNotFound):
              a.notFoundResponse(w, r)
           default:
              a.serverErrorResponse(w, r, err)
       }
       return 
   }

   // display the role
   data := envelope {
	"message": "role successfully deleted",
	}
		err = a.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}


func (a *application) listRoleHandler (w http.ResponseWriter, r *http.Request) {

	var queryParametersData struct {
		Role string
		data.Filters
	}
	// get the query parameters from the URL
	queryParameters := r.URL.Query()

	// Load the query parameters into our struct
    queryParametersData.Role = a.getSingleQueryParameter(
		queryParameters,
		"role",
		"")      

	//create a new validator instance
	v := validator.New()

	queryParametersData.Filters.Page = a.getSingleIntegerParameter(queryParameters, "page", 1, v)

	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(queryParameters, "page_size", 10,v)

	queryParametersData.Filters.Sort = a.getSingleQueryParameter(queryParameters, "sort", "id")
	
	queryParametersData.Filters.SortSafeList = []string {"id", "role","-id", "-role"}


	//check if our filters are valid
	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return

	}

	roles, metadata, err := a.roleModel.GetAll(queryParametersData.Role, queryParametersData.Filters)
	if err != nil {
    	a.serverErrorResponse(w, r, err)
    	return
  	}

	data := envelope {
    	"roles": roles,
		"@metadata": metadata,
   	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
    	a.serverErrorResponse(w, r, err)
  	}
}






// getUserRolesHandler retrieves all roles for a specific user.
func (a *application) getUserRolesHandler(w http.ResponseWriter, r *http.Request) {
   
	userID, err := a.readIDParam(r)
	if err != nil {
    	a.notFoundResponse(w, r)
    	return
	}

    // Fetch the user and their roles
    userName, roles, err := a.roleModel.GetForUserRole(userID)
    if err != nil {
        if errors.Is(err, data.ErrRecordNotFound) {
            a.notFoundResponse(w, r)
            return
        }
        a.serverErrorResponse(w, r, err)
        return
    }

    // Return as JSON
    data := envelope{
        "user": map[string]any{
            "id":    userID,
            "name":  userName,
            "roles": roles,
        },
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

func (a *application) updateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := a.readIDParam(r) // get user ID from URL
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Decode JSON body
    var input struct {
        OldRoleID int `json:"old_role_id"`
        NewRoleID int `json:"new_role_id"`
    }

    err = a.readJSON(w, r, &input)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    // Update the role for the user
    err = a.roleModel.UpdateForUserRole(int(userID), input.OldRoleID, input.NewRoleID)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{
        "message": fmt.Sprintf("user %d role updated successfully", userID),
    }
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

func (a *application) deleteUserRoleHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := a.readIDParam(r) // get user ID from URL
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Decode JSON body
    var input struct {
        RoleID int `json:"role_id"`
    }

    err = a.readJSON(w, r, &input)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    // Delete the role for the user
    err = a.roleModel.DeleteForUserRole(int(userID), input.RoleID)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{
        "message": fmt.Sprintf("role %d removed from user %d", input.RoleID, userID),
    }
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

//  retrieves all users and their associated roles.
func (a *application) listUsersWithRolesHandler(w http.ResponseWriter, r *http.Request) {
    usersWithRoles, err := a.roleModel.GetAllUsersWithRoles()
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{
        "users": usersWithRoles,
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}