package auth

import (
	"testing"
)

func Test_strMd5(t *testing.T) {
	t.Log(strMd5("2323"))
}

func Test_GetUserPassword(t *testing.T) {
	pass := GetUserPassword("123456", "12341")
	t.Log(strMd5(pass))
}

func Test_Encrypt(t *testing.T) {
	pass, err := Encrypt("123456")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pass)
}
