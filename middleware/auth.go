package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lorezi/boolang/helpers"
	"github.com/lorezi/boolang/models"
)

// type Privileges struct {
// 	Role   bool
// 	Create bool
// 	Read   bool
// 	Update bool
// 	Delete bool
// }

func taxAuthorization(p models.PermissionGroup, w http.ResponseWriter, r *http.Request) map[string]bool {

	// var roles = Privileges{Role: false, Create: false, Read: false, Update: false, Delete: false}
	roles := make(map[string]bool, 5)

	for _, v := range p.Permission {
		if v.Role == "tax" {
			roles["role"] = true
			for _, v := range v.Actions {
				if v.Create {
					roles["create"] = true

				}
				if v.Read {
					roles["read"] = true
				}
				if v.Update {
					roles["update"] = true
				}
				if v.Delete {
					roles["delete"] = true
				}
			}
			return roles
		}
	}

	return roles

}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		permissions := ctx.Value("permissions").(models.PermissionGroup)
		fmt.Println(permissions)

		// business logic
		tr := taxAuthorization(permissions, w, r)
		_, ok := tr["role"]
		_, found := tr["create"]
		if ok && found {
			next.ServeHTTP(w, r)
			return
		}

		// pr := payAuthorization(permissions, w, r)
		// _, ok = pr["role"]
		// _, found = pr["create"]

		// if ok && found {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// fmt.Println(tr)

		msg := models.Result{
			Status:  "Authorization failure",
			Message: "Contact admin 😰😰😰",
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)

	})
}

// for k, v := range p {
// 	fmt.Printf("key: %v  value: %v", k, v)
// }
// xs := strings.Split(tk, ".")
// u := strings.TrimSpace(xs[1])

// Auth validate token and authorize users

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		clientToken := r.Header.Get("Authorization")

		if clientToken == "" {
			r := models.Result{
				Status:  "Authentication failure",
				Message: "No Authorization header provided 😰😰😰",
			}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(r)
			return
		}

		if !strings.HasPrefix(clientToken, "Bearer") {
			r := models.Result{
				Status:  "Authentication failure",
				Message: "Add 'Bearer' to the token 👍🏾👍🏾👍🏾",
			}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(r)
			return
		}

		if strings.HasPrefix(clientToken, "Bearer") {

			tk := strings.TrimPrefix(clientToken, "Bearer ")

			claims, err := helpers.ValidateToken(tk)

			if err != "" {
				r := models.Result{
					Status:  "Authentication failure",
					Message: err,
				}
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(r)
				return
			}

			// authorization(claims.Permissions)

			r.Header.Set("email", claims.Email)
			r.Header.Set("first_name", claims.FirstName)
			r.Header.Set("last_name", claims.LastName)
			r.Header.Set("uid", claims.UID)

			ctx := r.Context()
			ctx = context.WithValue(ctx, "permissions", claims.Permissions)
			r = r.WithContext(ctx)

			// r.Header.Set("permission", claims.Permissions)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)

		}

	})
}
