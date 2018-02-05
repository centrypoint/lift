package lift

import (
	"net/http"
	"time"
	"net/url"
	"log"
	"errors"
	"encoding/json"
)

type Instance instance

type Route route

type Params struct {
	QueryParams map[string]string
}

func (p Params) New() Params {
	return Params{
		QueryParams: make(map[string]string),
	}
}

type instance struct {
	routes map[string]Route
}

type route struct {
	Path     string
	Method   string
	Params   Params
	Resolver func(params Params) (status int, response interface{}, err error)
}

func New() Instance {
	return Instance{routes: make(map[string]Route)}
}

func (ro *Route) prepare() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var (
			err      error
			response interface{}
			res      []byte
		)

		params := Params.New(Params{})
		responseStatus := 500
		start := time.Now()

		defer func(writer *http.ResponseWriter, method string, url *url.URL, s *time.Time, status *int) {
			if (*status) != 200 {
				(*writer).WriteHeader(*status)
			}
		}(&rw, r.Method, r.URL, &start, &responseStatus)

		defer func(e *error) {
			if _e := recover();
				_e != nil {
				log.Println(_e)
			}
			if (*e) != nil {
				log.Println(*e)
			}
		}(&err)

		defer r.Body.Close()

		if r.Method != ro.Method {
			responseStatus = http.StatusMethodNotAllowed
			return
		}

		if len(ro.Params.QueryParams) > 0 {
			for v := range ro.Params.QueryParams {
				value := r.URL.Query().Get(v)
				if len(value) < 1 {
					err = errors.New("not enough query params")
					return
				}
				params.QueryParams[v] = value
			}
		}

		responseStatus, response, err = ro.Resolver(params)

		if err != nil {
			return
		}

		if response == nil {
			responseStatus = 204
			return
		}

		if res, err = json.Marshal(response);
			err != nil {
			responseStatus = 500
			return
		}

		rw.Write(res)
	})
}

func (i *Instance) Register(r Route) {
	i.routes[r.Path] = r
}

func (i *Instance) Kindle() http.Handler {
	mux := http.NewServeMux()
	for p, v := range i.routes {
		mux.Handle(p, v.prepare())
	}
	return mux
}
