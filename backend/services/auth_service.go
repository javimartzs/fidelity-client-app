package services

import (
	"errors"
	"fidelity-client-app/config"
	"fidelity-client-app/models"
	"regexp"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

// RegisterNewUser (logica de negocio del registro de nuevos usuarios)
func (s *AuthService) RegisterNewUser(user *models.User, passConfirm string) error {

	// Validamos que los campos obligatorios estan llenos
	if user.FirstName == "" {
		return errors.New("el nombre es obligatorio")
	}
	if user.LastName == "" {
		return errors.New("el apellido es obligatorio")
	}
	if user.BirthDate == "" {
		return errors.New("la fecha de nacimiento es obligatoria")
	}
	if user.Gender == "" {
		return errors.New("el sexo es obligatorio")
	}

	// Comprobamos que el correo tiene formato de correo electronico
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("el formato del correo no es valido")
	}

	// Comprobamos si el correo YA esta registrado
	var existingEmail models.User
	if err := s.DB.Where("email = ?", user.Email).First(&existingEmail).Error; err == nil {
		return errors.New("el usuario ya está registrado")
	}

	// Comprobamos que la contraseña tiene minimo 8 caracteres
	if len(user.Password) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	// Comprobamos que la contraseña tiene al menos un numero y una mayuscula
	var hasUpper bool
	var hasNumber bool

	for _, char := range user.Password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper {
		return errors.New("la contraseña debe tener al menos una mayuscula")
	}
	if !hasNumber {
		return errors.New("la contraseña debe tener al menos un numero")
	}

	// Comprobamos si la contraseña y su confirmacion son iguales
	if user.Password != passConfirm {
		return errors.New("ambas contraseñas no coinciden")
	}

	// Hasheamos la contraseña
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error al procesar la contraseña")
	}

	// Guardamos la contraseña hasheada
	user.Password = string(hashedPass)
	// Generamos el UUID
	user.ID = uuid.NewString()
	// Forzamos el rol del usuario
	user.Role = "customer-client"

	// Guardamos al nuevo usuario en la tabla de User
	if err := s.DB.Save(&user).Error; err != nil {
		return errors.New("error al guardar usuario, intentelo de nuevo")
	}

	return nil
}

// LoginUser (logica de negovio del login de usuarios)
func (s *AuthService) LoginUser(email, password string) (*models.User, error) {
	var user models.User

	// Buscamos al usuario en la tabla con el correo
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("el usuario no esta registrado")
	}

	// Si el correo existe confirmamos la contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("el usuario o contraseña son incorrectos")
	}
	// Devolvemos los datos del usuario
	return &user, nil

}

// GenerateJWT maneja la generacion de un Json Web Token del usuario cuando se logea
func (s *AuthService) GenerateJWT(userID, email, role string) (string, error) {

	// Definimos los claims del JWT
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iss":     "fidelity-client-app",
	}

	// Generamos un nuevo Json Web Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmamos el Json Web Token generado
	tokenString, err := token.SignedString([]byte(config.Vars.JwtKey))
	if err != nil {
		return "", errors.New("error al acceder")
	}

	return tokenString, nil
}
