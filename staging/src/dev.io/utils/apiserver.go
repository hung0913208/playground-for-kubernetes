package utils

import (
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
)

type Handler func(http.ResponseWriter, *http.Request)

type Version struct {
  endpoints map[string]*Api
  code string
}

type Alias struct {
  methods map[string]*Api
  enable bool
}

type Api struct {
  methods map[string]Handler

  level int
  owner *ApiServer
  enable bool
  name, code string
}

type ApiServer struct {
  versions map[string]*Version
  aliases map[string]*Alias
  router *mux.Router

  base, currentVersion string
}

const (
  PUBLIC    = 0
  PRIVATE   = 1
  PROTECTED = 2
)

/*! \brief Make an alias path to specific endpoint
 *
 *  This method is used to create a new alias to specific endpoint which
 * is the convention way to split api into several version
 *
 *  \param path: the absolute path of this alias
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) alias(path string) *Api {
  var endpoint *Alias

  if tmp, ok := self.owner.aliases[path]; ok {
    endpoint = tmp
  } else {
    endpoint = &Alias{}

    endpoint.methods = make(map[string]*Api)
    endpoint.enable = true
  }

  if endpoint.methods == nil {
    endpoint.methods = make(map[string]*Api)
  }

  for k, _ := range(self.methods) {
    endpoint.methods[k] = self
  }

  self.owner.aliases[path] = endpoint
  self.owner.router.HandleFunc(path,
    func(w http.ResponseWriter, r *http.Request) {
      if link, ok := self.owner.aliases[path]; ! ok {
        self.nok(w)(404, "not found")
      } else if ! link.enable {
        self.nok(w)(404, "not found")
      } else if api, ok := link.methods[r.Method]; ! ok {
        self.nok(w)(404, "not found")
      } else {
        self.owner.reorder(api.name, api.code)(w, r)
      }
    })

  return self
}

/*! \brief Check if the endpoint is allowed to handle requests
 *
 *  This method is used to check and return what if the endpoint could be used
 * to handle requests
 *
 *  \param r: the user request
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) isAllowed(r *http.Request) bool {
  if ! self.enable {
    return false
  }

  switch(self.level) {
    case PUBLIC:
      return true

    case PRIVATE:
      return self.owner.isLocal(r)

    case PROTECTED:
      return self.owner.isInternal(r)

    default:
      return false
  }
}

/*! \brief Switch main version for configuring handlers
 *
 *  This method is used to switch main version, which is used when we would like
 * to configure multiple version for single RESTful API
 *
 *  \param code: the code version
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) version(code string) *ApiServer {
  return self.owner.version(code)
}

/*! \brief Set handler to resolve specific endpoint's method
 *
 *  This method is used to assign a handler to solve specific endpoint's method
 *
 *  \param method: the method we would like to resolve
 *  \param handler: the handler
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) handle(method string, handler Handler) *Api {
  self.methods[method] = handler
  return self
}

/*! \brief Access an endpoint object
 *
 *  This method is used to access an endpoint object using ApiServer, if the
 * endpoint is non-existing, this will create and return the new one
 *
 *  \param endpoint: the endpoint name
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) endpoint(endpoint string) *Api {
  return self.owner.endpoint(endpoint)
}

/*! \brief Mock a specific path to this endpoint
 *
 *  This method is used to link a path to specific endpoint in order to handle
 * requests which are send directly to this path
 *
 *  \param path: the path which will receive requests
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *Api) mock(path string) *Api {
  var dest string

  if len(self.owner.base) > 0 {
    dest = fmt.Sprintf("/%s/%s%s", self.owner.base, self.code, path)
  } else {
    dest = fmt.Sprintf("/%s%s", self.code, path)
  }

  self.owner.router.HandleFunc(dest,
    self.owner.reorder(self.name, self.code))

  if len(self.owner.base) > 0 {
    path = fmt.Sprintf("/%s%s", self.owner.base, path)
  }

  return self.alias(path)
}

/*! \brief Send ok code and message to client
 *
 *  This function is used to produce a lambda which is used to write an ok
 * message to client
 *
 *  \param w: the response writer
 *  \return func(string): a lambda which is used to pack message and code
 *                        into an json object
 */
func (self *Api) ok(w http.ResponseWriter) func(string) {
  return self.owner.ok(w)
}

/*! \brief Send nok code and message to client
 *
 *  This function is used to produce a lambda which is used to write a self.nok
 * message to client
 *
 *  \param w: the response writer
 *  \return func(int, string): a lambda which is used to pack message and code
 *                             into an json object
 */
func (self *Api) nok(w http.ResponseWriter) func(int, string) {
  return self.owner.nok(w)
}

/* ------------------------- ApiServer ---------------------------- */

/*! \brief Access a specific endpoint
 *
 *  \param endpoint: the endpoint name which is used to separate APIs
 *                   inside a same version
 *  \return *Api: to make a chain call, we will return itself to
 *                make calling next function easily
 */
func (self *ApiServer) endpoint(endpoint string) *Api {
  if len(self.currentVersion) == 0 {
    return nil
  } else {
    if _, ok := self.endpoints[endpoint]; ! ok {
      self.endpoints[endpoint] = self.newApi(endpoint)
    }

    return self.endpoints[endpoint]
  }
}

/*! \brief Create a new version which is used to group APIs
 *
 *  This method is used to switch to another version for storing new API
 * and create new version if it didn't exist
 *
 *  \param code: the version code
 *  \return *ApiServer: to make a chain call, we will return itself
 *                      to make calling next function easily
 */
func (self *ApiServer) version(code string) *ApiServer {
  if _, ok := self.versions[code]; ! ok {
    self.versions[code] = &Version{}

    self.versions[code].code = code
    self.versions[code].endpoints = make(map[string]*Api)
  }

  self.currentVersion = code
  return self
}

/*! \brief Order a handler to redirect request to specific endpoint's version
 *
 *  This method is used to link a path to specific endpoint's version, in order
 * to build a large scale version of single one endpoint
 *
 *  \param path: the path which will receive requests
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *ApiServer) reorder(endpoint, code string) Handler {
  return func(w http.ResponseWriter, r *http.Request) {
    if ver, ok := self.versions[code]; ! ok {
      self.nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if api, ok := ver.endpoints[endpoint]; ! ok {
      self.nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if handler, ok := api.methods[r.Method]; ! ok {
      self.nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if api.isAllowed(r) {
      handler(w, r)
    } else {
      self.nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    }
  }
}

func (self *ApiServer) newApi(name string) *Api {
  ret := &Api{}

  ret.code = self.currentVersion
  ret.name = name
  ret.owner = self
  ret.enable = true
  ret.methods = make(map[string]Handler)
  return ret
}

func (self *ApiServer) isLocal(r *http.Request) bool {
  fmt.Println(r.RemoteAddr)
  return false
}

func (self *ApiServer) isInternal(r *http.Request) bool {
  fmt.Println(r.RemoteAddr)
  return false
}


func (self *ApiServer) handleMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.RequestURI)
    next.ServeHTTP(w, r)
  })
}

/* --------------------------- helper ----------------------------- */

/*! \brief Pack code and message into an json object and write back to client
 *
 *  This function is used to produce a lambda which is used to write a message
 * as response to client in a form way
 *
 *  \param w: the response writer
 *  \return func(int, string): a lambda which is used to pack message and code
 *                             into an json object
 */
func pack(w http.ResponseWriter) func(int, string) {
  return func(code int, message string) {
    if message[0] == '{' && message[len(message) - 1] == '}' {
      fmt.Fprintf(w, "{\"code\": %d, \"data\": %s}", code, message)
    } else if message[0] == '[' && message[len(message) - 1] == ']' {
      fmt.Fprintf(w, "{\"code\": %d, \"data\": %s}", code, message)
    } else {
      fmt.Fprintf(w, "{\"code\": %d, \"data\": \"%s\"}", code, message)
    }
  }
}

/*! \brief Send nok code and message to client
 *
 *  This function is used to produce a lambda which is used to write a self.nok
 * message to client
 *
 *  \param w: the response writer
 *  \return func(int, string): a lambda which is used to pack message and code
 *                             into an json object
 */
func (self *ApiServer) nok(w http.ResponseWriter) func(int, string) {
  return func(code int, message string) {
    pack(w)(code, message)
  }
}

/*! \brief Send ok code and message to client
 *
 *  This function is used to produce a lambda which is used to write an ok
 * message to client
 *
 *  \param w: the response writer
 *  \return func(string): a lambda which is used to pack message and code
 *                        into an json object
 */
func (self *ApiServer) ok(w http.ResponseWriter) func(string) {
  return func(message string) {
    pack(w)(200, message)
  }
}

/* --------------------------- public ----------------------------- */

func (self *ApiServer) GetMuxer() *mux.Router {
  return self.router
}

func NewApiServer() *ApiServer {
  ret := &ApiServer{}

  ret.router = mux.NewRouter()
  ret.versions = make(map[string]*Version)
  ret.aliases = make(map[string]*Alias)

  ret.router.Use(ret.handleMiddleware)
  return ret
}
