package data

import (
	"context"
	"database/sql"
	"slices"
	"time"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"

)

type Role struct {
	ID        int64     `json:"id"`  
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"-"`
}

//function that performs the validation check
func ValidateRole(v *validator.Validator, role *Role) {
	//check if the Role field is empty
	v.Check(role.Role != "", "content", "cannot be left blank")
	v.Check(len(role.Role) <= 50, "content", "must no be more than 50 bytes long")
}

// Setup model
type RoleModel struct {
	DB *sql.DB
}


//Insert a new row in the role table
func (r RoleModel) Insert(role *Role) error{

   // the SQL query to be executed against the database table
    query := `INSERT INTO role (role)
        	VALUES ($1)
        	RETURNING id, created_at`

   args := []any{role.Role}
  
//Create a context with a 3-second timeout. NO database operation should take more than 3 secs or we will quit it
ctx, cancel := context.WithTimeout(context.Background(), 3  * time.Second)
defer cancel()

// to update the Role struct later on 
return r.DB.QueryRowContext(ctx, query, args...).Scan(&role.ID, &role.CreatedAt)

}

// Get a specific Role from the role table
func (r RoleModel) Get(id int64) (*Role, error) {
	// check if the id is valid
	 if id < 1 {
		 return nil, ErrRecordNotFound
	 }
	// the SQL query to be executed against the database table
	 query := `
		 SELECT id, role, created_at
		 FROM role
		 WHERE id = $1 `

	 // declare a variable of type Role to store the returned role
	 var role Role

	 // Set a 3-second context/timer
	 ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	 defer cancel()
	 
	 err := r.DB.QueryRowContext(ctx, query, id).Scan (&role.ID,&role.Role, &role.CreatedAt)
	
	// check for which type of error
	if err != nil {
    	switch {
        	case errors.Is(err, sql.ErrNoRows):
            	return nil, ErrRecordNotFound
        	default:
            	return nil, err
        	}
    	}
	return &role, nil
}

func (r RoleModel) Update(role *Role) error {
    // SQL query to update the role name
    query := `
        UPDATE role
        SET role = $1
        WHERE id = $2
    `

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Execute the update
    _, err := r.DB.ExecContext(ctx, query, role.Role, role.ID)
    return err
}
		

// Delete a specific Role from the role table
func (r RoleModel) Delete(id int64) error {

    // check if the id is valid
    if id < 1 {
        return ErrRecordNotFound
    }
   // the SQL query to be executed against the database table
    query := `
        DELETE FROM role
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	 
	 // ExecContext does not return any rows unlike QueryRowContext. 
	 // It only returns  information about the the query execution
	 // such as how many rows were affected
		result, err := r.DB.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}

	// Were any rows  delete?
    rowsAffected, err := result.RowsAffected()
    if err != nil {
       return err
   }
	// Probably a wrong id was provided or the client is trying to
	// delete an already deleted role
   if rowsAffected == 0 {
       return ErrRecordNotFound
   }

   return nil
}

func (r RoleModel) GetAll(role string, filters Filters) ([]*Role, Metadata, error) {

    // Dynamic ORDER BY — make sure filters.sortColumn() and filters.sortDirection() are safe
    query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), id, role, created_at
        FROM role
        WHERE (to_tsvector('simple', role) @@ plainto_tsquery('simple', $1) OR $1 = '')
        ORDER BY %s %s, id ASC
        LIMIT $2 OFFSET $3
    `, filters.sortColumn(), filters.sortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := r.DB.QueryContext(ctx, query, role, filters.limit(), filters.offset())
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    roles := []*Role{}

    for rows.Next() {
        var role Role
        err := rows.Scan(&totalRecords, &role.ID, &role.Role, &role.CreatedAt)
        if err != nil {
            return nil, Metadata{}, err
        }
        roles = append(roles, &role)
    }

    if err = rows.Err(); err != nil {
        return nil, Metadata{}, err
    }

    metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
    return roles, metadata, nil
}

/*----------------------------------------------users_role------------------------------------------------------*/

// We will have the role in a slice which we will be able to search
type Roles []string

// Is the role  found for the Role slice
func (r Roles) Include(role string) bool {
	return slices.Contains(r, role)
}

// // What are all the roles associated with the user
// func (r RoleModel) GetAllForUser(userID int64) (Roles, error) {
// 	query := `
//             SELECT role.role
//             FROM role 
//         	INNER JOIN users_role ON 
//             users_role.role_id = role.id
// 			INNER JOIN users ON users_role.user_id = users.id
//             WHERE users.id = $1
//           `
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	rows, err := r.DB.QueryContext(ctx, query, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Ensure to release resources after use
// 	defer rows.Close()
// 	// Store the role for the user in our slice
// 	var roles Roles
// 	for rows.Next() {
// 		var role string

// 		err := rows.Scan(&role)
// 		if err != nil {
// 			return nil, err
// 		}
// 		roles = append(roles, role)
// 	} // end of for loop

// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return roles, nil

// }

func (r RoleModel) AddForUserRole(userID int64, roleIDs ...int) error {
    query := `
        INSERT INTO users_role (user_id, role_id)
        VALUES ($1, unnest($2::int[]))
    `
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := r.DB.ExecContext(ctx, query, userID, pq.Array(roleIDs))
    return err
}

// retrieves all roles associated with a specific user,
// including the user's name.
func (r RoleModel) GetForUserRole(userID int64) (string, []string, error) {
    query := `
    SELECT u.fname || ' ' || u.lname AS full_name,
           array_agg(ro.role ORDER BY ro.id)
    FROM users u
    JOIN users_role ur ON u.id = ur.user_id
    JOIN role ro ON ur.role_id = ro.id
    WHERE u.id = $1
    GROUP BY u.fname, u.lname
    `

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    var userName string
    var roles []string

    err := r.DB.QueryRowContext(ctx, query, userID).Scan(&userName, pq.Array(&roles))
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", nil, fmt.Errorf("no roles found for user %d", userID)
        }
        return "", nil, err
    }

    return userName, roles, nil
}


// UpdateForUser replaces a specific role for a user
func (r RoleModel) UpdateForUserRole(userID, oldRoleID, newRoleID int) error {
    query := `
        UPDATE users_role
        SET role_id = $1
        WHERE user_id = $2 AND role_id = $3
    `
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    result, err := r.DB.ExecContext(ctx, query, newRoleID, userID, oldRoleID)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no role found to update for user %d", userID)
    }

    return nil
}

// DeleteForUser removes a specific role from a user
func (r RoleModel) DeleteForUserRole(userID, roleID int) error {
    query := `
        DELETE FROM users_role
        WHERE user_id = $1 AND role_id = $2
    `
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    result, err := r.DB.ExecContext(ctx, query, userID, roleID)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("user %d does not have role %d", userID, roleID)
    }

    return nil
}

// Exists checks if a user already has the specified role.
func (m *RoleModel) Exists(userID, roleID int) (bool, string, error) {
    var roleName string
    query := `
        SELECT r.role
        FROM users_role ur
        JOIN role r ON ur.role_id = r.id
        WHERE ur.user_id = $1 AND ur.role_id = $2
        LIMIT 1
    `
    err := m.DB.QueryRow(query, userID, roleID).Scan(&roleName)
    if err == sql.ErrNoRows {
        return false, "", nil // no duplicate found
    }
    if err != nil {
        return false, "", err
    }
    return true, roleName, nil
}


//retrieves all users with their associated roles.
func (r RoleModel) GetAllUsersWithRoles() ([]map[string]any, error) {
    query := `
        SELECT u.id, u.fname || ' ' || u.lname AS full_name, array_agg(ro.role)
        FROM users u
        LEFT JOIN users_role ur ON u.id = ur.user_id
        LEFT JOIN role ro ON ur.role_id = ro.id
        GROUP BY u.id, u.fname, u.lname
    `

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := r.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []map[string]any

    for rows.Next() {
        var id int64
        var name string
        var roles []string

        err := rows.Scan(&id, &name, pq.Array(&roles))
        if err != nil {
            return nil, err
        }

        results = append(results, map[string]any{
            "id":    id,
            "name":  name,
            "roles": roles,
        })
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return results, nil
}