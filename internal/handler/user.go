package handler

import (
	"fmt"
	"net/http"

	"forum/internal/models"
	"forum/internal/validator"
)

// REGISTER
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := models.UserData{IsAuth: false} // save data input from input

	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form.Get("username") // username
		email := r.Form.Get("email")       // email
		password := r.Form.Get("password") // password

		v := validator.NewValidator()

		user := &models.User{}

		user.UserName = username
		user.Email = email
		user.Password = password

		// VALIDATE DATA FROM USER
		if models.ValidateUser(v, user); !v.Valid() {
			if v.Errors["username"] != "" {
				data.Errors.UserName = v.Errors["username"]
			}

			if v.Errors["email"] != "" {
				data.Errors.Email = v.Errors["email"]
			}

			if v.Errors["password"] != "" {
				data.Errors.Password = v.Errors["password"]
			}

		}

		// set data from user
		data.UserName = username
		data.Email = email
		data.Password = password

		// check user is exist in db
		existUser, _ := h.service.UserService.GetUserByEmail(data.Email)

		// if email and username alredy exist
		if existUser != nil {
			if existUser.Email == email {
				data.Errors.Email = "THIS EMAIL ALREADY EXIST!!"
			}

			if existUser.UserName == username {
				data.Errors.Email = "USERNAME ALREADY EXIST"
			}
		}

		// if errors is empty create user
		if data.Errors.UserName == "" && data.Errors.Email == "" && data.Errors.Password == "" {
			err := h.service.UserService.CreateUser(user)
			if err != nil {
				http.Error(w, "INTERNAL SERVER ERROR", http.StatusBadRequest)
				return
			}

			http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // redirect to sign-in page
		}

	}
	// render sign-up page
	h.Render(w, "sign_up.html", data)
}

// LOGIN
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	data := models.LoginData{IsAuth: false}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		data.User.Email = r.Form.Get("email")       // email
		data.User.Password = r.Form.Get("password") // password

		// check email and password for correctness
		user, err := h.service.UserService.Login(data.User.Email, data.User.Password)
		if err != nil {
			if err == models.ErrUserNotFound {
				data.Errors.Email = "Invalid email"
				data.Errors.Password = "Wrong email or password"
			} else if err == models.ErrWrongPassword {
				data.Errors.Password = "Wrong password!Try again:))"
			}
		} else {
			data.IsAuth = true

			// delete session
			// only one active session it's has
			err = h.service.SessionService.DeleteSessionByUser(user.ID)
			if err != nil {
				fmt.Println("Err,", err)
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}

			// after delete old session
			// create new seesion
			session, err := h.service.SessionService.CreateSession(user.ID, data.Errors.Email)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// set cokkie
			cokkie := &http.Cookie{
				Name:     "session",
				Value:    session.Token,
				Path:     "/",
				Expires:  session.ExpireTime,
				HttpOnly: true,
				MaxAge:   7200,
			}

			http.SetCookie(w, cokkie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	}
	h.Render(w, "sign_in.html", data)
}

// LOGOUT
func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set cookie default value
	expiredCookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}

	http.SetCookie(w, expiredCookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
