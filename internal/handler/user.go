package handler

import (
	"net/http"

	"forum/internal/models"
	"forum/internal/validator"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := models.SignupRequest{IsAuth: false}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			h.logger.Error("parse form signup", err)
			return
		}

		data.UserName = r.Form.Get("username")
		data.Email = r.Form.Get("email")
		data.Password = r.Form.Get("password")

		v := validator.NewValidator()

		if models.ValidateSignupRequest(v, &data); !v.Valid() {
			if v.ErrorsMap["username"] != "" {
				data.ErrorMessages.UserName = v.ErrorsMap["username"]
			}
			if v.ErrorsMap["email"] != "" {
				data.ErrorMessages.Email = v.ErrorsMap["email"]
			}
			if v.ErrorsMap["password"] != "" {
				data.ErrorMessages.Password = v.ErrorsMap["password"]
			}

			h.Render(w, "sign_up.html", data)
		}

		err = h.service.UserService.SignUpUser(&data)
		if err != nil {
			if err == models.ErrUserExists {
				data.ErrorMessages.Email = "Email or Username already exists"
				h.Render(w, "sign_up.html", data)
				return
			}
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			h.logger.Error("signup user", err)
			return
		}

		// Redirect to sign-in page on success
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return

	}
	h.Render(w, "sign_up.html", data)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	data := models.LoginRequest{IsAuth: false}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			h.logger.Error("parse form signin", err)
			return
		}

		data.User.Email = r.Form.Get("email")
		data.User.Password = r.Form.Get("password")

		user, err := h.service.UserService.Login(data.User.Email, data.User.Password)
		if err != nil {
			h.handleLoginError(err, &data)
		} else {
			data.IsAuth = true

			session, err := h.service.SessionService.SetSession(user.ID)
			if err != nil {
				h.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
				h.logger.Error("set session", err)
				return
			}

			cookie := &http.Cookie{
				Name:     "session",
				Value:    session.Token,
				Path:     "/",
				Expires:  session.ExpireTime,
				HttpOnly: true,
				MaxAge:   7200,
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	h.Render(w, "sign_in.html", data)
}

func (h *Handler) handleLoginError(err error, data *models.LoginRequest) {
	if err == models.ErrUserNotFound {
		data.ErrorMessages.Email = "Invalid email"
		data.ErrorMessages.Password = "Wrong email or password"
	} else if err == models.ErrWrongPassword {
		data.ErrorMessages.Password = "Wrong password!Try again:))"
	}
}

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
