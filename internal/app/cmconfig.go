package app

import (
	"fmt"
	"regexp"
	"strings"
)

//=================== ConfigMap ==================//
type CmConfig struct {
	Namespace   string
	Name        string
	RawData     map[string]string
	ServiceCall []*Service

	DBProperties map[string]string
}

func (c *CmConfig) GetKey() string {
	return fmt.Sprintf("%s|%s", c.Namespace, c.Name)
}

// TODO : work on pragmatic matching ..
func (c *CmConfig) GrabLinks() {

	for k, v := range c.RawData {
		firstmatch, _ := regexp.MatchString("http://[a-zA-Z\\-]+.[a-zA-Z\\-]+.svc.cluster.local:[0-9]+/.+", v)
		secondMatch, _ := regexp.MatchString("http://[a-zA-Z].+:[0-9].+", v)
		// fmt.Println(fmt.Sprintf("%s - %s", k, v))

		if c.isJdbdProperties(k, v) {
			c.DBProperties[k] = v
		} else {

			if firstmatch || secondMatch {
				fmt.Printf("%s - %s", "matched", v)
				s := &Service{
					Url: v,
				}
				s.parseUrl()

				if s.IsValidService() {
					fmt.Printf("===> %s\n", "valid")
					c.ServiceCall = append(c.ServiceCall, s)
				} else {
					fmt.Printf("###> %s\n", "invalid")
				}

			}

		}

	}
}

func (c *CmConfig) isJdbdProperties(key string, value string) bool {
	if strings.HasPrefix(key, "spring.datasource") {
		return true
	}

	if strings.Contains(key, ".hikari.") {
		return true
	}

	if strings.Contains(value, "jdbc:") {
		return true
	}

	return false

}

// pattern 1
// http://[a-zA-Z\-]+.[a-zA-Z\-]+.svc.cluster.local:[0-9]+/.+
// http://cash-account.owcsapi.svc.cluster.local:8443/cast-account/api

// pattern 2
// http://[a-zA-Z].+:[0-9].+
// http://cash-account:8443/cast-account/api

// ====================================================== //
type Service struct {
	Namespace string
	Name      string
	Url       string
}

func (s *Service) parseUrl() {
	//removal of http:// cash-account:8443/cast-account/api
	t := strings.Split(s.Url, "//")

	//cash-account  :   8443/cast-account/api
	// remove the port number if any ....
	t = strings.Split(t[1], ":")

	if len(t) > 0 {
		if !strings.HasSuffix(t[0], ".com") {

			t = strings.Split(t[0], ".")
			s.Name = t[0]
			if len(t) > 1 {
				// sanitize
				if t[1] != "svc" {
					s.Namespace = t[1]
				}
			}
		}
	}
}

func (s *Service) GetKey() string {
	return fmt.Sprintf("%s|%s", s.Namespace, s.Name)
}

func (s *Service) IsValidService() bool {
	return (len(s.Name) > 0)
}

func (s *Service) RequireNamespaceLookup() bool {
	return (len(s.Name) > 0) && (len(s.Namespace) == 0)
}
