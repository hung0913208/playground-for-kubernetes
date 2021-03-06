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
  apis map[string]*Api

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
func (self *Api) Alias(path string) *Api {
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
        self.Nok(w)(404, "not found")
      } else if ! link.enable {
        self.Nok(w)(404, "not found")
      } else if api, ok := link.methods[r.Method]; ! ok {
        self.Nok(w)(404, "not found")
      } else {
        self.owner.reorder(api.name, api.code)(w, r)
      }
    })

  return self
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
func (self *Api) Version(code string) *ApiServer {
  return self.owner.Version(code)
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
func (self *Api) Handle(method string, handler Handler) *Api {
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
func (self *Api) Endpoint(endpoint string) *Api {
  return self.owner.Endpoint(endpoint)
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
func (self *Api) Mock(path string) *Api {
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

  return self.Alias(path)
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
func (self *Api) Ok(w http.ResponseWriter) func(string) {
  return self.owner.Ok(w)
}

/*! \brief Send nok code and message to client
 *
 *  This function is used to produce a lambda which is used to write a self.Nok
 * message to client
 *
 *  \param w: the response writer
 *  \return func(int, string): a lambda which is used to pack message and code
 *                             into an json object
 */
func (self *Api) Nok(w http.ResponseWriter) func(int, string) {
  return self.owner.Nok(w)
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

/* ------------------------- ApiServer ---------------------------- */

/*! \brief Access a specific endpoint
 *
 *  \param endpoint: the endpoint name which is used to separate APIs
 *                   inside a same version
 *  \return *Api: to make a chain call, we will return itself to
 *                make calling next function easily
 */
func (self *ApiServer) Endpoint(endpoint string) *Api {
  if len(self.currentVersion) == 0 {
    return nil
  } else {
    if ver, ok := self.versions[self.currentVersion]; ! ok {
      return nil
    } else {
      if _, ok := ver.endpoints[endpoint]; ! ok {
        ver.endpoints[endpoint] = self.newApi(endpoint)
      }

      return ver.endpoints[endpoint]
    }
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
func (self *ApiServer) Version(code string) *ApiServer {
  if _, ok := self.versions[code]; ! ok {
    self.versions[code] = &Version{}

    self.versions[code].code = code
    self.versions[code].endpoints = make(map[string]*Api)
  }

  self.currentVersion = code
  return self
}

/*! \brief Send nok code and message to client
 *
 *  This function is used to produce a lambda which is used to write a self.Nok
 * message to client
 *
 *  \param w: the response writer
 *  \return func(int, string): a lambda which is used to pack message and code
 *                             into an json object
 */
func (self *ApiServer) Nok(w http.ResponseWriter) func(int, string) {
  return func(code int, message string) {
    Pack(w)(code, message)
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
func (self *ApiServer) Ok(w http.ResponseWriter) func(string) {
  return func(message string) {
    Pack(w)(200, message)
  }
}

func (self *ApiServer) GetMuxer() *mux.Router {
  return self.router
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
      self.Nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if api, ok := ver.endpoints[endpoint]; ! ok {
      self.Nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if handler, ok := api.methods[r.Method]; ! ok {
      self.Nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if api.isAllowed(r) {
      handler(w, r)
    } else {
      self.Nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    }
  }
}

/*! \brief Create a new API
 *
 *  This method is used to create a new API object and store it to our database
 * for serving and keeping track
 *
 *  \param name: the api name
 *  \return *Api: to make a chain call, we will return itself to make calling
 *                next function easily
 */
func (self *ApiServer) newApi(name string) *Api {

  if ret, ok := self.apis[name]; ! ok {
    ret = &Api{}

    ret.code = self.currentVersion
    ret.name = name
    ret.owner = self
    ret.enable = true
    ret.methods = make(map[string]Handler)

    return ret
  } else {
    return ret
  }
}

/*! \brief Check if request is local or not
 *
 *  This method is used to check if the request is produced by this itself
 * or not
 *
 *  \param r: the request
 *  \return bool: return if the request is created by itself or not
 */
func (self *ApiServer) isLocal(r *http.Request) bool {
  return false
}

/*! \brief Check if request is internal or not
 *
 *  This method is used to check if the request is produced by our cluster
 * or from outside
 *
 *  \param r: the request
 *  \return bool: return if the request is created by cluster or not
 */
func (self *ApiServer) isInternal(r *http.Request) bool {
  return false
}

/*! \brief Snift in comming requests before redirect it to correct service
 *
 *  This method is used to listen request from everywhere and redirect them
 * to correct placement
 *
 *  \param next: the handler which is registered
 *  \return http.Handler: the actual handler which server will do
 */
func (self *ApiServer) handleMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
func Pack(w http.ResponseWriter) func(int, string) {
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

/*! \brief Create Api server
 *
 *  This function is used to generate ApiServer which is used to build
 * complicated RESTful APIs
 *
 *  \return *ApiServer: the server object
 */
func NewApiServer() *ApiServer {
  ret := &ApiServer{}

  ret.router = mux.NewRouter()
  ret.versions = make(map[string]*Version)
  ret.aliases = make(map[string]*Alias)

  ret.router.Use(ret.handleMiddleware)
  return ret
}
