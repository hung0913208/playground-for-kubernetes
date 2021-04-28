package utils

import (
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
)

type Handler func(http.ResponseWriter, *http.Request)

type Version struct {
  endpoints map[string]*Api
}

type Alias struct {
  code map[string]string
  enable bool
}

type Api struct {
  methods map[string]Handler

  level int
  owner *ApiServer
  enable bool
  name, version string
}

type ApiServer struct {
  versions map[string]*Version
  alias map[string]*Alias
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
  if , ok := self.owner.aliases[path]; ! ok {
    self.owner.aliases[path] = &Alias{}
  }

  self.aliases[path] = &Alias{}
  self.aliases[path].code = make(map[string]string)

  for key, val := range self.mainlines {
    self.aliases[path].code[key] = val
  }

  self.owner.router.HandleFunc(path,
    func(w http.ResponseWriter, r *http.Request){
      code := self.aliases[path].code[r.Method]

      if ver, ok := self.versions[code]; ! ok {
        self.nok(w)(404, fmt.Sprintf("Not found %s", path))
      } else if handler, ok := ver.methods[r.Method]; ! ok {
        self.nok(w)(404, fmt.Sprintf("Not found %s", path))
      } else if self.isAllowed(r) {
        handler(w, r)
      } else {
        self.nok(w)(404, fmt.Sprintf("Not found %s", path))
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
  if ! self.owner.versions[self.version].enable {
    return false
  } else if ! self.enable {
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
  if len(self.main) == 0 {
    panic("Please specifiy version before doing anything")
  }

  if _, ok := self.versions[self.main]; ! ok {
    panic("There is something wrong with creating new version")
  }

  if _, ok := self.mainlines[method]; ! ok {
    self.mainlines[method] = self.main
  }

  self.versions[self.main].methods[method] = handler
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
    dest = fmt.Sprintf("/%s/%s%s", self.owner.base, self.version, path)
  } else {
    dest = fmt.Sprintf("/%s%s", self.version, path)
  }

  self.owner.router.HandleFunc(dest,
    self.owner.reorder(self.name, self.version))

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
  if ver := self.getCurrentVersion(); ver == nil {
    return nil
  } else {
    if _, ok := ver.endpoints[endpoint]; ! ok {
      ver.endpoints[endpoint] = self.newApi(endpoint)
    }

    return ver.endpoints[endpoint]
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
    if api, ok := self.endpoints[endpoint]; ! ok {
      self.nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if ver, ok := api.versions[code]; ! ok {
      self.nok(w)(404, fmt.Sprintf("Not found %s", endpoint))
    } else if handler, ok := ver.methods[r.Method]; ! ok {
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

  ret.main = ""
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
  ret.endpoints = make(map[string]*Api)

  return ret
}