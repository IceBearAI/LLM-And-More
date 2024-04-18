package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

type TokenSource string

func (c TokenSource) String() string {
	return string(c)
}

const (
	TokenSourceQw TokenSource = "qw" //企业微信验证登录，生成token
	TokenSourceNd TokenSource = "nd" //宁盾验证登录，生成token
)

// ArithmeticCustomClaims 自定义声明
type ArithmeticCustomClaims struct {
	Source   TokenSource `json:"source"`
	QwUserid string      `json:"qwUserid"`
	Email    string      `json:"email"`
	UserId   uint        `json:"userId"`
	IsAdmin  bool        `json:"isAdmin"`
	//jwt.StandardClaims
	jwt.RegisteredClaims
}

type ArithmeticTerminalClaims struct {
	UserId      uint   `json:"userId"`
	Cluster     string `json:"cluster"`
	Namespace   string `json:"namespace"`
	ServiceName string `json:"serviceName"`
	PodName     string `json:"podName"`
	Container   string `json:"container"`
	//jwt.StandardClaims
	jwt.RegisteredClaims
}

func (s ArithmeticTerminalClaims) Valid() error {
	return nil
}

// ArithmeticGptClaims token
type ArithmeticGptClaims struct {
	Email string `json:"email"`
	//jwt.StandardClaims
	jwt.RegisteredClaims
}

//func Sign(email string, uid string, sessionTimeout int64) (string, error) {
//
//	expAt := time.Now().Add(time.Duration(sessionTimeout)).Unix()
//
//	fmt.Println("expAt", expAt)
//
//	// 创建声明
//	claims := ArithmeticCustomClaims{
//		UserId: uid,
//		Name:   email,
//		StandardClaims: jwt.StandardClaims{
//			ExpiresAt: expAt,
//			Issuer:    "sys",
//		},
//	}
//
//	//创建token，指定加密算法为HS256
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
//	//生成token
//	return token.SignedString([]byte(GetJwtKey()))
//}
