package webmgmt

import (
    "encoding/json"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "golang.org/x/crypto/bcrypt"
)

// Authr struct
type Authr struct {
    SecretKey string "secretkeyjwt"
}

// User struct
type User struct {
    gorm.Model
    Name     string `json:"name"`
    Email    string `gorm:"unique" json:"email"`
    Password string `json:"password"`
    Role     string `json:"role"`
}

// Authentication struct
type Authentication struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Token struct
type Token struct {
    Role        string `json:"role"`
    Email       string `json:"email"`
    TokenString string `json:"token"`
}

// Error struct
type Error struct {
    IsError bool   `json:"isError"`
    Message string `json:"message"`
}

//--------------HELPER FUNCTIONS---------------------

//SetError error message in Error struct
func SetError(err Error, message string) Error {
    err.IsError = true
    err.Message = message
    return err
}

//GenerateHashPassword password as input and generate new hash password from it
func GenerateHashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

//CheckPasswordHash plain password with hash password
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

//GenerateJWT JWT token
func (a *Authr) GenerateJWT(email, role string) (string, error) {
    var mySigningKey = []byte(a.SecretKey)
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["authorized"] = true
    claims["email"] = email
    claims["role"] = role
    claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

    tokenString, err := token.SignedString(mySigningKey)
    if err != nil {
        fmt.Errorf("something went Wrong: %s", err.Error())
        return "", err
    }

    return tokenString, nil
}

//---------------------MIDDLEWARE FUNCTION-----------------------

//IsAuthorized whether user is authorized or not
func (a *Authr) IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        if r.Header["Token"] == nil {
            var err Error
            err = SetError(err, "No Token Found")
            json.NewEncoder(w).Encode(err)
            return
        }

        var mySigningKey = []byte(a.SecretKey)

        token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("There was an error in parsing token.")
            }
            return mySigningKey, nil
        })

        if err != nil {
            var err Error
            err = SetError(err, "Your Token has been expired.")
            json.NewEncoder(w).Encode(err)
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            if claims["role"] == "admin" {
                r.Header.Set("Role", "admin")
                handler.ServeHTTP(w, r)
                return

            } else if claims["role"] == "user" {
                r.Header.Set("Role", "user")
                handler.ServeHTTP(w, r)
                return

            }
        }
        var resErr Error
        resErr = SetError(resErr, "Not Authorized.")
        json.NewEncoder(w).Encode(err)
    }
}

//InitializeAuthMux all auth routes
func InitializeAuthMux(prefix string, router *mux.Router) (a *Authr) {
    a = &Authr{}
    a.InitialMigration()
    router.HandleFunc(fmt.Sprintf("%s/signup", prefix), a.SignUp).Methods("POST")
    router.HandleFunc(fmt.Sprintf("%s/signin", prefix), a.SignIn).Methods("POST")
    router.HandleFunc(fmt.Sprintf("%s/admin", prefix), a.IsAuthorized(a.AdminIndex)).Methods("GET")
    router.HandleFunc(fmt.Sprintf("%s/user", prefix), a.IsAuthorized(a.UserIndex)).Methods("GET")
    router.HandleFunc(fmt.Sprintf("%s/", prefix), a.Index).Methods("GET")
    router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
    })
    return a
}

//InitializeAuthGin all auth routes for gin router
func InitializeAuthGin(prefix string, router *gin.Engine) (a *Authr) {
    a = &Authr{}
    a.InitialMigration()
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    // loge.Info("app.Config.WebPath: %v\n", app.webPath)
    router.POST(fmt.Sprintf("%s/signup", prefix), func(c *gin.Context) {
        a.SignUp(c.Writer, c.Request)
    })
    router.POST(fmt.Sprintf("%s/signin", prefix), func(c *gin.Context) {
        a.SignIn(c.Writer, c.Request)
    })

    router.GET(fmt.Sprintf("%s/admin", prefix), func(c *gin.Context) {
        a.IsAuthorized(a.AdminIndex)(c.Writer, c.Request)
    })
    router.GET(fmt.Sprintf("%s/user", prefix), func(c *gin.Context) {
        a.IsAuthorized(a.UserIndex)(c.Writer, c.Request)
    })

    router.GET(fmt.Sprintf("%s/", prefix), func(c *gin.Context) {
        a.Index(c.Writer, c.Request)
    })

    return a
}

// SignUp ROUTE HANDLER
func (a *Authr) SignUp(w http.ResponseWriter, r *http.Request) {
    connection := a.GetDatabase()
    defer a.CloseDatabase(connection)

    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        var err Error
        err = SetError(err, "Error in reading payload.")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(err)
        return
    }

    var dbuser User
    connection.Where("email = ?", user.Email).First(&dbuser)

    //check email is already registered or not
    if dbuser.Email != "" {
        var err Error
        err = SetError(err, "Email already in use")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(err)
        return
    }

    user.Password, err = GenerateHashPassword(user.Password)
    if err != nil {
        log.Fatalln("Error in password hashing.")
    }

    //insert user details in database
    connection.Create(&user)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// SignIn ROUTE HANDLER
func (a *Authr) SignIn(w http.ResponseWriter, r *http.Request) {
    connection := a.GetDatabase()
    defer a.CloseDatabase(connection)

    var authDetails Authentication

    err := json.NewDecoder(r.Body).Decode(&authDetails)
    if err != nil {
        var err Error
        err = SetError(err, "Error in reading payload.")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(err)
        return
    }

    var authUser User
    connection.Where("email = 	?", authDetails.Email).First(&authUser)

    if authUser.Email == "" {
        var err Error
        err = SetError(err, "Username or Password is incorrect")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(err)
        return
    }

    check := CheckPasswordHash(authDetails.Password, authUser.Password)

    if !check {
        var err Error
        err = SetError(err, "Username or Password is incorrect")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(err)
        return
    }

    validToken, err := a.GenerateJWT(authUser.Email, authUser.Role)
    if err != nil {
        var err Error
        err = SetError(err, "Failed to generate token")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(err)
        return
    }

    var token Token
    token.Email = authUser.Email
    token.Role = authUser.Role
    token.TokenString = validToken
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(token)
}

// Index ROUTE HANDLER
func (a *Authr) Index(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("HOME PUBLIC INDEX PAGE"))
}

// AdminIndex ROUTE HANDLER
func (a *Authr) AdminIndex(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Role") != "admin" {
        w.Write([]byte("Not authorized."))
        return
    }
    w.Write([]byte("Welcome, Admin."))
}

// UserIndex ROUTE HANDLER
func (a *Authr) UserIndex(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Role") != "user" {
        w.Write([]byte("Not Authorized."))
        return
    }
    w.Write([]byte("Welcome, User."))
}
