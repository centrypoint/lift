package lift_test

import "testing"
import (
	"github.com/centrypoint/lift"
	"fmt"
	"net/http"
	"io/ioutil"
)

func TestHandlers(t *testing.T) {
	i := lift.New()
	i.Register(lift.Route{
		Params: lift.Params{QueryParams: map[string]string{"name": "name", "surname": "surname"}},
		Path: "/test",
		Method: "GET",
		Resolver: func(params lift.Params) (status int, response interface{}, err error) {
			status = 200
			err = nil
			response = fmt.Sprintf("hello, %s %s!\n", params.QueryParams["name"], params.QueryParams["surname"])
			return
		},
	})

	go http.ListenAndServe(":8080", i.Kindle())

	r, _ := http.Get("http://localhost:8080/test?name=tester&surname=testerov")
	defer r.Body.Close()
	res, _ := ioutil.ReadAll(r.Body)
	t.Logf("%s\n", string(res))
}
