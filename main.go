package main

import (
        "bytes"
        "context"
        "fmt"
        "github.com/tsocial/goproxy/utils"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "io/ioutil"
        "net/http"
        "net/http/httputil"
        "net/url"
)

const (
        clientID = ""
        clientSecret = ""
        //oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
        userInfoScope = "https://www.googleapis.com/auth/userinfo.email"
        redirectURL =  ""
)
var conf *oauth2.Config

func main() {
        remote, err := url.Parse("http://localhost:5000")
        if err != nil {
                panic(err)
        }

        proxy := httputil.NewSingleHostReverseProxy(remote)
        http.HandleFunc("/", handler(proxy))
        http.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
                if r.URL.Path != "/oauth" {
                        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                        return
                }

                bucket := r.URL.Query().Get("state")
                token, err := conf.Exchange(context.Background(), r.URL.Query().Get("code"))
                if err != nil {
                        panic(err)
                }
                cookie    :=    http.Cookie{Name: "token",Value:token.AccessToken}
                http.SetCookie(w, &cookie)
                re := utils.PopRequest(bucket)
                method := re.Method
                url := re.Url
                body := re.Body
                req1, err := http.NewRequest(method, remote.String()+url, bytes.NewReader(body))
                if err != nil {
                        panic(err)
                }

                client := http.Client{}
                resp, err3 := client.Do(req1)
                if err3 != nil {
                        panic(err3)
                }

                bodyb, _ := ioutil.ReadAll(resp.Body)
                w.Write(bodyb)
        })
        err = http.ListenAndServe(":8080", nil)
        if err != nil {
                panic(err)
        }
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
        return func(w http.ResponseWriter, r *http.Request) {
                _, err := r.Cookie("token")
                if err != nil{
                        fmt.Println("Not authorized")
                        conf = &oauth2.Config{
                                ClientID:     clientID,
                                ClientSecret: clientSecret,
                                RedirectURL:  redirectURL,
                                Scopes: []string{
                                        userInfoScope,
                                },
                                Endpoint: google.Endpoint,
                        }
                        reqByte, _ := ioutil.ReadAll(r.Body)
                        bucket := utils.CacheRequest(&utils.HttpRequest{
                                r.Method, r.RequestURI, reqByte,
                        })
                        fmt.Println("login bucket no. ", bucket)


                        fmt.Println("url :", conf.AuthCodeURL(bucket))
                        req, err := http.NewRequest("GET", conf.AuthCodeURL(bucket), nil)
                        if err != nil {
                                panic(err)
                        }

                        client := http.Client{}
                        resp, err3 := client.Do(req)
                        if err3 != nil {
                                panic(err3)
                        }

                        bodyb, _ := ioutil.ReadAll(resp.Body)
                        w.Write(bodyb)
                }else{
                        fmt.Println(r)
                        p.ServeHTTP(w, r)
                }

        }
}