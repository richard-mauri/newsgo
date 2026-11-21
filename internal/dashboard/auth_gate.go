package dashboard

import (
 "context"
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "regexp"
 "strings"
 "time"

 suiteauth "github.com/richard-mauri/homepage/shared/auth"
)

const newsgoProfileEndpoint = "https://home.themauris.org/api/session/profile"
const newsgoOrigin = "https://newsgo.themauris.org"
var newsgoBearerPattern = regexp.MustCompile(`^Bearer [a-f0-9]{64}$`)

var errNewsgoNoBearer=errors.New("verified Homepage bearer missing")

// NewsGoAuthGate validates every app API request with Homepage. It does not
// authorize based on client role labels, telemetry or a previous page load.
// The initial GET / shell and static assets must be public so a new browser
// can receive its child bearer in a fragment; their content is not authorization.
type NewsGoAuthGate struct {
 ProfileURL string
 Client *http.Client
}

func NewNewsGoAuthGate() *NewsGoAuthGate {
 return &NewsGoAuthGate{
  ProfileURL:newsgoProfileEndpoint,
  Client:&http.Client{Timeout:4*time.Second,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}},
 }
}

func (g *NewsGoAuthGate) verify(ctx context.Context,r *http.Request)(suiteauth.Identity,error) {
 authorization:=r.Header.Get("Authorization")
 if !newsgoBearerPattern.MatchString(authorization) {return suiteauth.Identity{},errNewsgoNoBearer}
 if g==nil || g.ProfileURL=="" || g.Client==nil{return suiteauth.Identity{},suiteauth.ErrVerificationUnavailable}
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,g.ProfileURL,nil)
 if err!=nil{return suiteauth.Identity{},err}
 req.Header.Set("Authorization",authorization)
 req.Header.Set("Origin",newsgoOrigin)
 req.Header.Set("Accept","application/json")
 response,err:=g.Client.Do(req)
 if err!=nil{return suiteauth.Identity{},err}
 defer response.Body.Close()
 if response.StatusCode==http.StatusUnauthorized || response.StatusCode==http.StatusForbidden{return suiteauth.Identity{},nil}
 if response.StatusCode!=http.StatusOK{return suiteauth.Identity{},suiteauth.ErrVerificationUnavailable}
 var profile struct{
  User string `json:"user"`
  Role suiteauth.Role `json:"role"`
  ExpiresAt string `json:"expiresAt"`
 }
 decoder:=json.NewDecoder(io.LimitReader(response.Body,8192))
 if err:=decoder.Decode(&profile);err!=nil{return suiteauth.Identity{},err}
 expiration,err:=time.Parse(time.RFC3339,profile.ExpiresAt)
 if err!=nil || !time.Now().Before(expiration){return suiteauth.Identity{},nil}
 identity:=suiteauth.Identity{User:strings.TrimSpace(profile.User),Role:profile.Role,Verified:true}
 if !suiteauth.Allowed(identity,suiteauth.UseApplication){return suiteauth.Identity{},nil}
 return identity,nil
}

// Protect applies the shared Homepage authorization guard at the outer HTTP
// boundary. Every /api/ route requires a fresh verified personal or Guest
// bearer. Unknown /api paths are guarded too, preventing accidental bypasses.
func (g *NewsGoAuthGate) Protect(next http.Handler) http.Handler {
 guarded:=suiteauth.Guard(func(ctx context.Context,r *http.Request)(suiteauth.Identity,error){
  identity,err:=g.verify(ctx,r)
  if errors.Is(err,errNewsgoNoBearer){return suiteauth.Identity{},nil}
  return identity,err
 },suiteauth.UseApplication,next)
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path=="/api"||strings.HasPrefix(r.URL.Path,"/api/"){
   guarded.ServeHTTP(w,r)
   return
  }
  next.ServeHTTP(w,r)
 })
}
