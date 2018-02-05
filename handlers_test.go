package lift_test

import "testing"
import (
	"../lift"
	"fmt"
	"net/http"
	"io/ioutil"
)

func TestHandlers(t *testing.T) {
	i := lift.New()
	i.Register(struct {
		Path     string
		Method   string
		QueryParams   []string
		Resolver func(params ...interface{}) (status int, response interface{}, err error)
	}{Path: "/test", Method: "GET", QueryParams: []string{"name"}, Resolver: func(params ...interface{}) (status int, response interface{}, err error) {
		status = 200
		response = fmt.Sprintf("hello, %s!\n", params[0])
		err = nil
		return
	}})
	

	go http.ListenAndServe(":8080", i.Kindle())

	r, _ := http.Get("http://localhost:8080/test?name=tester")
	defer r.Body.Close()
	res, _ := ioutil.ReadAll(r.Body)
	t.Logf("%s\n", string(res))
}
