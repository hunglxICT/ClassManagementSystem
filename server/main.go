package main

import (
	"fmt"
	"os"
	"time"
	"strconv"
	"log"
	"strings"
	"io"
	"net/http"
	"html"
	"database/sql"
	"crypto/sha256"
	"./Config"
	"./Models"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
_	"github.com/go-sql-driver/mysql"
	"github.com/gin-contrib/cors"
)

var db *sql.DB
var STUDENT_ROLE int = 1
var TEACHER_ROLE int = 2
var RESPONSE_OK = 200
var RESPONSE_FAIL = 500
var STATUS_ON = 1

// -----------------------------jwt processor---------------------------

//jwt service
type JWTService interface {
	GenerateToken(username string, id, role int) string
	ValidateToken(token string) (*jwt.Token, error)
}
type authCustomClaims struct {
	Username 	string	`json:"name"`
	Id			int		`json:"id"`
	Role		int 	`json:"role"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issuer    string
}

//auth-jwt
func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issuer:    "Bikash",
	}
}

func getSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (service jwtServices) GenerateToken(username string, id, role int) string {
	claims := authCustomClaims{
		Username: username,
		Id: id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	//fmt.Printf("%d\n", claims.role)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})

}

func GetAccountFromCookie(token string) (jwt.MapClaims, error) {
	var test jwtServices
	result, err := test.ValidateToken(token)
	if result.Valid {
		claims := result.Claims.(jwt.MapClaims)
		return claims, nil
	}
	return nil, err
}

// -----------------------------end of jwt processor---------------------------


// --------------------------string verification-------------------------

func isUsername(s string) bool {
	for i := 0; i < len(s); i++ {
		if (s[i] >= 'A' && s[i] <= 'Z') {
			continue
		} else if (s[i] >= 'a' && s[i] <= 'z') {
			continue
		} else if (i > 0 && s[i] >= '0' && s[i] <= '9') {
			continue
		} else if (s[i] == '_') {
			continue
		} else {
			return false
		}
	}
	return true
}

func contain_dangerous_character(s string) bool {
	for i := 0; i < len(s); i++ {
		if (s[i] < 32 || s[i] > 126) {
			return true
		}
	}
	blacklist := "<>'\";\\/|?:"
	return strings.ContainsAny(s, blacklist)
}

func encode_string(s string) string {
	s = html.EscapeString(s)
	return s
}

func password_hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func upload(c *gin.Context) string {
	file, header, err := c.Request.FormFile("file") 
	if err != nil {
		fmt.Println(err)
		//c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return "error"
	}
	original_filename := header.Filename
	
	if (contain_dangerous_character(original_filename)) {
		return "error"
	}
	filename := password_hash(fmt.Sprintf("%lld",time.Now().Unix())) + original_filename
	out, err := os.Create("public/" + filename)
	if err != nil {
		//log.Fatal(err)
		return "error"
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		//log.Fatal(err)
		return "error"
	}
	filepath := "http://localhost:8080/file/" + filename
	return filepath
	//c.JSON(http.StatusOK, gin.H{"filepath": filepath})
}

// --------------------------end of string verification-------------------------

func main() {
	os.Setenv("SECRET", "9928829f03b6884bf62f8bebf61a0c089168c1fa")
	db, err := sql.Open("mysql",Config.DbURL(Config.BuildDBConfig()))
	//Config.DB, err = gorm.Open("mysql", Config.DbURL(Config.BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
		return
	}
	//defer Config.DB.Close()
	defer db.Close()
	//Config.DB.AutoMigrate(&Models.Account{})
	
	r := gin.Default()
	//----------------------------------------------------------------------------
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
	//----------------------------------------------------------------------------
	r.POST("/login", func(c *gin.Context) {
		acc := Models.Account {}
		c.Bind(&acc)
		//_username := c.PostForm("username")
		_username := acc.Username
		if (isUsername(_username) == false) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		//_password := c.PostForm("password")
		_password := acc.Password
		
		row := db.QueryRow("SELECT id,username,role,firstname,lastname FROM Account WHERE username = ? AND password = ? AND status = 1", encode_string(_username), password_hash(_password))
		
		err = row.Scan(&acc.Id, &acc.Username, &acc.Role, &acc.Firstname, &acc.Lastname)
		if err != nil {
			fmt.Println(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		var test jwtServices
		token := test.GenerateToken(acc.Username, acc.Id, acc.Role)
		c.JSON(RESPONSE_OK, gin.H{
			"id": acc.Id,
			"username": acc.Username,
			"firstName": acc.Firstname,
			"lastName": acc.Lastname,
			"authdata": token,
		})
	})
	//----------------------------------------------------------------------------
	r.GET("/list-student", func(c *gin.Context) {
		var accounts []Models.Account
		rows, err := db.Query("SELECT id, username, firstname, lastname, role, avatar FROM Account WHERE status=1 LIMIT 100")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var acc Models.Account
			rows.Scan(&acc.Id, &acc.Username,&acc.Firstname, &acc.Lastname, &acc.Role, &acc.Avatar)
			accounts = append(accounts, acc)
		}
		result := gin.H {
			"result": accounts,
		}
		c.JSON(RESPONSE_OK, result)
	})
	//----------------------------------------------------------------------------
	r.GET("/enrolled-students/:id", func(c *gin.Context) {
		id := c.Param("id")
		_Id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		var accounts []Models.Account
		rows, err := db.Query("SELECT Account.id, Account.username, Account.firstname, Account.lastname, Account.role, Account.avatar FROM Enroll, Account, Class WHERE Enroll.status=1 AND Enroll.classID = ? AND Enroll.studentID = Account.id AND Account.status=1 AND Enroll.classID = Class.id AND Class.status = 1 LIMIT 100", _Id)
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var acc Models.Account
			rows.Scan(&acc.Id, &acc.Username,&acc.Firstname, &acc.Lastname, &acc.Role, &acc.Avatar)
			accounts = append(accounts, acc)
		}
		result := gin.H {
			"result": accounts,
		}
		c.JSON(RESPONSE_OK, result)
	})
	//----------------------------------------------------------------------------
	r.GET("/list-class", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		var classes []Models.Class
		id := int(cookieinfo["id"].(float64))
		rows, err := db.Query("(SELECT id, teacherID, classname FROM Class WHERE status=1 AND teacherID = ? UNION SELECT Class.id, Class.teacherID, Class.classname FROM Class, Enroll WHERE Class.Status = 1 AND Enroll.Status = 1 AND Enroll.studentID = ? AND Enroll.classID = Class.id) LIMIT 100", id, id)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer rows.Close()
		
		for rows.Next() {
			var cla Models.Class
			rows.Scan(&cla.Id, &cla.Teacherid, &cla.Classname)
			classes = append(classes, cla)
		}
		result := gin.H {
			"result": classes,
		}
		c.JSON(RESPONSE_OK, result)
	})
	//----------------------------------------------------------------------------
	r.POST("/new-student", func(c *gin.Context) {
		/*
		_username := c.PostForm("username")
		_password := c.PostForm("password")
		_firstname := c.PostForm("firstname")
		_lastname := c.PostForm("lastname")
		_email := c.PostForm("email")
		_tel := c.PostForm("tel")
		_avatar := c.PostForm("avatar")
		_role := STUDENT_ROLE
		
		new_acc := Models.Account {
			Username: encode_string(_username),
			Password: encode_string(_password),
			Firstname: encode_string(_firstname),
			Lastname: encode_string(_lastname),
			Email: encode_string(_email),
			Tel: encode_string(_tel),
			Avatar: encode_string(_avatar),
			Role: _role,
		}
		*/
		
		new_acc := Models.Account{}
		c.Bind(&new_acc)
		new_acc.Role = STUDENT_ROLE
		if (isUsername(new_acc.Username) == false) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		new_acc.Password = password_hash(new_acc.Password)
		new_acc.Firstname = encode_string(new_acc.Firstname)
		new_acc.Lastname = encode_string(new_acc.Lastname)
		new_acc.Email = encode_string(new_acc.Email)
		new_acc.Tel = encode_string(new_acc.Tel)
		new_acc.Avatar = encode_string(new_acc.Avatar)
		//err = Config.DB.Create(&new_acc).Error;
		p, err := db.Prepare("INSERT INTO Account(username, password, firstname, lastname, email, tel, avatar, role, status) VALUES (?,?,?,?,?,?,?,?,1)")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(new_acc.Username,
						new_acc.Password,
						new_acc.Firstname,
						new_acc.Lastname,
						new_acc.Email,
						new_acc.Tel,
						new_acc.Avatar,
						new_acc.Role)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.POST("/new-class", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		if (int(cookieinfo["role"].(float64)) & TEACHER_ROLE == 0) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		_teacherid := int(cookieinfo["id"].(float64))
		
		new_class := Models.Class{}
		c.Bind(&new_class)
		
		new_class.Teacherid = _teacherid
		new_class.Classname = encode_string(new_class.Classname)
		
		p, err := db.Prepare("INSERT INTO Class(teacherID, classname, status) VALUES (?,?,1)")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(new_class.Teacherid,
						new_class.Classname)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.POST("/join-class", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		_adderrole := int(cookieinfo["role"].(float64))
		_adderid := int(cookieinfo["id"].(float64))
		
		new_join := Models.Enroll{}
		c.Bind(&new_join)
		//fmt.Printf("%d\n", _adderrole)
		if !((_adderrole & TEACHER_ROLE != 0) || ((_adderrole & STUDENT_ROLE != 0) && (_adderid == new_join.Studentid))) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		p, err := db.Prepare("INSERT INTO Enroll(studentID, classID, status) VALUES (?,?,1)")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(new_join.Studentid,
						new_join.Classid)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.POST("/dismiss-student", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		_adderrole := int(cookieinfo["role"].(float64))
		_adderid := int(cookieinfo["id"].(float64))
		
		new_join := Models.Enroll{}
		c.Bind(&new_join)
		//fmt.Printf("%d\n", _adderrole)
		if !((_adderrole & TEACHER_ROLE != 0) || ((_adderrole & STUDENT_ROLE != 0) && (_adderid == new_join.Studentid))) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		p, err := db.Prepare("UPDATE Enroll SET status=0 WHERE studentID=? AND classID=?")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(new_join.Studentid,
						new_join.Classid)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.POST("/add-exercise", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		if (int(cookieinfo["role"].(float64)) & TEACHER_ROLE == 0) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		new_exercise := Models.Exercises {}
		c.Bind(&new_exercise)
		new_exercise.Title = encode_string(new_exercise.Title)
		new_exercise.Description = encode_string(new_exercise.Description)
		new_exercise.Link = ""
		
		p, err := db.Prepare("INSERT INTO Exercises(classID, title, description, link, status) VALUES (?,?,?,?,1)")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(new_exercise.Classid,
						new_exercise.Title,
						new_exercise.Description,
						new_exercise.Link)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		var resultid int
		row := db.QueryRow("SELECT LAST_INSERT_ID()")
		err = row.Scan(&resultid)
		//fmt.Printf("%d\n",resultid)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		c.JSON(RESPONSE_OK, gin.H{
			"result": resultid,
		})
	})
	//----------------------------------------------------------------------------
	r.POST("/add-link-exercise/:id", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		if (int(cookieinfo["role"].(float64)) & TEACHER_ROLE == 0) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		id := c.Param("id")
		_Id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		var creatorid int
		row := db.QueryRow("SELECT Class.teacherID FROM Exercises, Class WHERE Exercises.id=? AND Exercises.classID=Class.id", _Id)
		err = row.Scan(&creatorid)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		if creatorid != int(cookieinfo["id"].(float64)) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		_link := upload(c)
		if (_link == "error") {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		p, err := db.Prepare("UPDATE Exercises SET link = ? WHERE id = ?")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(_link, _Id)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		c.JSON(RESPONSE_OK, gin.H{
			"message": RESPONSE_OK,
		})
	})
	//----------------------------------------------------------------------------
	r.GET("/list-exercises/:id", func(c *gin.Context) {
		var exercises []Models.Exercises
		id := c.Param("id")
		_Id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		rows, err := db.Query("SELECT id, classID, title, description, link FROM Exercises WHERE status=1 AND classID = ? LIMIT 100", _Id)
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var exer Models.Exercises
			rows.Scan(&exer.Id, &exer.Classid, &exer.Title, &exer.Description, &exer.Link)
			exercises = append(exercises, exer)
		}
		result := gin.H {
			"result": exercises,
		}
		c.JSON(RESPONSE_OK, result)
	})
	//----------------------------------------------------------------------------
	r.POST("/submit", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		_studentid := int(cookieinfo["id"].(float64))
		new_submission := Models.Submission {}
		c.Bind(&new_submission)

		new_submission.Studentid = _studentid
		new_submission.Description = encode_string(new_submission.Description)
		new_submission.Link = ""
		p, err := db.Prepare("INSERT INTO Submission(studentID, exerciseID, description, link, status) VALUES (?,?,?,?,1)")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(new_submission.Studentid,
						new_submission.Exerciseid,
						new_submission.Description,
						new_submission.Link)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		var resultid int
		row := db.QueryRow("SELECT LAST_INSERT_ID()")
		err = row.Scan(&resultid)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		c.JSON(RESPONSE_OK, gin.H{
			"result": resultid,
		})
	})
	//----------------------------------------------------------------------------
	r.POST("/add-link-submission/:id", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		id := c.Param("id")
		_Id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		_link := upload(c)
		if (_link == "error") {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		p, err := db.Prepare("UPDATE Submission SET link = ? WHERE id = ? AND studentID = ?")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(_link, _Id, int(cookieinfo["id"].(float64)))
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		c.JSON(RESPONSE_OK, gin.H{
			"message": RESPONSE_OK,
		})
	})
	//----------------------------------------------------------------------------
	r.GET("/get-submissions/:id", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		if (int(cookieinfo["role"].(float64)) & TEACHER_ROLE == 0) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		id := c.Param("id")
		_Id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		var submissions []Models.Submission
		
		rows, err := db.Query("SELECT id, studentID, exerciseID, description, link FROM Submission WHERE status=1 AND exerciseID = ? LIMIT 100", _Id)
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var submit Models.Submission
			rows.Scan(&submit.Id, &submit.Studentid, &submit.Exerciseid, &submit.Description, &submit.Link)
			submissions = append(submissions, submit)
		}
		result := gin.H {
			"result": submissions,
		}
		c.JSON(RESPONSE_OK, result)
	})
	//----------------------------------------------------------------------------
	r.GET("/profile/:id", func(c *gin.Context) {
		id := c.Param("id")
		_Id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		acc := Models.Account {
			Id: _Id,
		}
		row := db.QueryRow("SELECT username, firstname, lastname, email, tel, avatar, role FROM Account WHERE status = 1 AND id = ?", acc.Id)
		err = row.Scan(&acc.Username, &acc.Firstname, &acc.Lastname, &acc.Email, &acc.Tel, &acc.Avatar, &acc.Role)
		//fmt.Printf("%s\n%s\n",acc.Username,acc.Firstname)
		if err != nil {
			fmt.Println(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		} else {
			result := gin.H {
				"result": acc,
			}
			c.JSON(RESPONSE_OK, result)
		}
	})
	//----------------------------------------------------------------------------
	r.GET("/class-info/:id", func(c *gin.Context) {
		id := c.Param("id")
		_Id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		acc := Models.Class {
			Id: _Id,
		}
		row := db.QueryRow("SELECT teacherID, classname FROM Class WHERE status = 1 AND id = ?", acc.Id)
		err = row.Scan(&acc.Teacherid, &acc.Classname)
		
		if err != nil {
			fmt.Println(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		} else {
			result := gin.H {
				"result": acc,
			}
			c.JSON(RESPONSE_OK, result)
		}
	})
	//----------------------------------------------------------------------------
	r.GET("/exercise-detail/:id", func(c *gin.Context) {
		id := c.Param("id")
		_Id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		exer := Models.Exercises {
			Id: _Id,
		}
		row := db.QueryRow("SELECT classID, title, description, link FROM Exercises WHERE status = 1 AND id = ?", exer.Id)
		err = row.Scan(&exer.Classid, &exer.Title, &exer.Description, &exer.Link)
		
		if err != nil {
			fmt.Println(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		} else {
			result := gin.H {
				"result": exer,
			}
			c.JSON(RESPONSE_OK, result)
		}
	})
	//----------------------------------------------------------------------------
	r.POST("/edit-profile", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		acc := Models.Account {}
		c.Bind(&acc)
		edited_id := acc.Id
		
		editor_role := int(cookieinfo["role"].(float64))
		editor_id := int(cookieinfo["id"].(float64))
		if ((editor_role & TEACHER_ROLE == 0) && (edited_id != editor_id)) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		} else if (editor_role & TEACHER_ROLE != 0) {
			var edited_role int
			row := db.QueryRow("SELECT role FROM Account WHERE status = 1 AND id = ?", edited_id)
			err = row.Scan(&edited_role)
			if err != nil {
				fmt.Println(err)
				c.JSON(RESPONSE_FAIL, gin.H{
					"message": "error",
				})
				return
			} else {
				if ((edited_role & TEACHER_ROLE != 0) && (edited_id != editor_id)) {
					c.JSON(RESPONSE_FAIL, gin.H{
						"message": "error",
					})
					return
				}
			}
		}
		_password := acc.Password
		if _password != "" {
			_, err := db.Exec("UPDATE Account SET password=? WHERE id=?", password_hash(_password), acc.Id)
			if err != nil {
				fmt.Println(err)
				c.JSON(RESPONSE_FAIL, gin.H{
					"message": "error",
				})
			}
		}
		
		_email := acc.Email
		if _email != "" {
			_, err := db.Exec("UPDATE Account SET email=? WHERE id=?", encode_string(_email), acc.Id)
			if err != nil {
				fmt.Println(err)
				c.JSON(RESPONSE_FAIL, gin.H{
					"message": "error",
				})
			}
		}
		_tel := acc.Tel
		if _tel != "" {
			_, err := db.Exec("UPDATE Account SET tel=? WHERE id=?", encode_string(_tel), acc.Id)
			if err != nil {
				fmt.Println(err)
				c.JSON(RESPONSE_FAIL, gin.H{
					"message": "error",
				})
			}
		}
		_avatar := acc.Avatar
		if _avatar != "" {
			_, err := db.Exec("UPDATE Account SET avatar=? WHERE id=?", encode_string(_avatar), acc.Id)
			if err != nil {
				fmt.Println(err)
				c.JSON(RESPONSE_FAIL, gin.H{
					"message": "error",
				})
			}
		}
		
		if (editor_role & TEACHER_ROLE != 0) {
			_firstname := acc.Firstname
			if _firstname != "" {
				_, err := db.Exec("UPDATE Account SET firstname=? WHERE id=?", encode_string(_firstname), acc.Id)
				if err != nil {
					fmt.Println(err)
					c.JSON(RESPONSE_FAIL, gin.H{
						"message": "error",
					})
				}
			}
			_lastname := acc.Lastname
			if _lastname != "" {
				_, err := db.Exec("UPDATE Account SET lastname=? WHERE id=?", encode_string(_lastname), acc.Id)
				if err != nil {
					fmt.Println(err)
					c.JSON(RESPONSE_FAIL, gin.H{
						"message": "error",
					})
				}
			}
			_username := acc.Username
			if ((_username != "") && (isUsername(_username))) {
				_, err := db.Exec("UPDATE Account SET username=? WHERE id=?", encode_string(_username), acc.Id)
				if err != nil {
					fmt.Println(err)
					c.JSON(RESPONSE_FAIL, gin.H{
						"message": "error",
					})
				}
			}
		}
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.GET("/delete-account/:id", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}

		if (int(cookieinfo["role"].(float64)) & TEACHER_ROLE == 0) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		id := c.Param("id")
		_deletedid, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		var deleted_role int
		row := db.QueryRow("SELECT role FROM Account WHERE status = 1 AND id = ?", _deletedid)
		err = row.Scan(&deleted_role)
		if err != nil {
			fmt.Println(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		} else if (deleted_role & TEACHER_ROLE == 0 && deleted_role & STUDENT_ROLE != 0) {
			_, err := db.Exec("UPDATE Account SET status=0 WHERE id=?", _deletedid)
			if err != nil {
				fmt.Println(err)
				c.JSON(RESPONSE_FAIL, gin.H{
					"message": "error",
				})
				return
			}
		} else {
			fmt.Println(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.GET("/delete-class/:id", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}

		if (int(cookieinfo["role"].(float64)) & TEACHER_ROLE == 0) {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		id := c.Param("id")
		_deletedid, err := strconv.Atoi(id)
		_deletorid := int(cookieinfo["id"].(float64))
		if err != nil {
			log.Fatalln(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		var class_creator int
		row := db.QueryRow("SELECT teacherID FROM Class WHERE id = ?", _deletedid)
		err = row.Scan(&class_creator)
		if err != nil {
			fmt.Println(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		} else if (class_creator == _deletorid) {
			_, err := db.Exec("UPDATE Class SET status=0 WHERE id=?", _deletedid)
			if err != nil {
				fmt.Println(err)
				c.JSON(RESPONSE_FAIL, gin.H{
					"message": "error",
				})
				return
			}
		} else {
			fmt.Println(err)
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.POST("/send-messenger", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		_senderid := int(cookieinfo["id"].(float64))
		var new_messenger Models.Messenger
		c.Bind(&new_messenger)
		new_messenger.Senderid = _senderid
		new_messenger.Message = encode_string(new_messenger.Message)
		
		p, err := db.Prepare("INSERT INTO Messenger(senderID, receiverID, message, status) VALUES (?,?,?,1)")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(new_messenger.Senderid,
						new_messenger.Receiverid,
						new_messenger.Message)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.GET("/get-messenger/:id", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		_senderid := int(cookieinfo["id"].(float64))
		id := c.Param("id")
		_receiverid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		var messengers []Models.Messenger
		rows, err := db.Query("(SELECT id, senderID, receiverID, message, create_date FROM Messenger WHERE status=1 AND (((senderID = ?) AND (receiverID = ?)) OR ((senderID = ?) AND (receiverID = ?))) ORDER BY id DESC LIMIT 20) ORDER BY id ASC", _senderid, _receiverid, _receiverid, _senderid)
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var mess Models.Messenger
			rows.Scan(&mess.Id, &mess.Senderid, &mess.Receiverid, &mess.Message, &mess.Create_date)
			messengers = append(messengers, mess)
		}
		result := gin.H {
			"result": messengers,
		}
		c.JSON(RESPONSE_OK, result)
	})
	//----------------------------------------------------------------------------
	r.POST("/edit-messenger", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		_userid := int(cookieinfo["id"].(float64))
		
		var new_messenger Models.Messenger
		c.Bind(&new_messenger)
		_messengerid := new_messenger.Id
		fmt.Printf("%d",_messengerid)
		_message := encode_string(new_messenger.Message)
		
		row := db.QueryRow("SELECT senderID FROM Messenger WHERE id = ? AND status = 1", _messengerid)
		var _realsenderid int
		err = row.Scan(&_realsenderid)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		if _userid != _realsenderid {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		p, err := db.Prepare("UPDATE Messenger SET message = ? WHERE id = ?")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(encode_string(_message), _messengerid)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.POST("/delete-messenger", func(c *gin.Context) {
		cookie := c.Request.Header.Get("Authorization")
		cookieinfo, err := GetAccountFromCookie(cookie)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		
		var new_messenger Models.Messenger
		c.Bind(&new_messenger)
		_messengerid := new_messenger.Id
		
		_userid := int(cookieinfo["id"].(float64))
	
		row := db.QueryRow("SELECT senderID FROM Messenger WHERE id = ? AND status = 1", _messengerid)
		var _realsenderid int
		err = row.Scan(&_realsenderid)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		if _userid != _realsenderid {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		p, err := db.Prepare("UPDATE Messenger SET status = 0 WHERE id = ?")
		
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		defer p.Close()
		
		_, err = p.Exec(_messengerid)
		if err != nil {
			c.JSON(RESPONSE_FAIL, gin.H{
				"message": "error",
			})
			return
		}
		c.JSON(RESPONSE_OK, gin.H{
			"message": "ok",
		})
	})
	//----------------------------------------------------------------------------
	r.StaticFS("/file", http.Dir("public"))
	r.Run() // listen and serve on 0.0.0.0:9399 (for windows "localhost:8080")
}