// Allows sending of emails to users
package email

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"io/ioutil"

	"strings"

	"github.com/pkg/errors"
)

const MAILGUN_URL string = "https://api.mailgun.net/v3"

var MailgunFailure = errors.New("email - Mailgun returned server error status code")
var MailgunUnauthorized = errors.New("email - Mailgun incorrect auth key")

// Emailer defines methods for our mailgun implementation
type Emailer interface {
	SendConfirmation(to, name, id string) error
}

type Mailgun struct {
	key           string
	domain        string
	confirmDomain string
}

func InitMailgunService(key, domain, confirmDomain string) *Mailgun {
	return &Mailgun{key: key, domain: domain, confirmDomain: confirmDomain}
}

// SendConfirmation sends a confirmation email to users email address
func (m *Mailgun) SendConfirmation(to, name, id string) error {
	c := http.Client{Timeout: time.Second * time.Duration(30)}

	from := fmt.Sprintf("Delayed NZ <confirm@%v>", m.domain)
	text := fmt.Sprintf("Hey %v\r\n\r\nThank you for signing up to Delayed. Before you can login, please confirm your email.\r\nClick here https://%v/confirm/%v\r\n\r\nThe Delayed Team", name, m.confirmDomain, id)

	form := url.Values{}
	form.Set("from", from)
	form.Set("to", to)
	form.Set("subject", "Confirm Your Email")
	form.Set("text", text)

	mgURL := fmt.Sprintf("%v/%v/messages", MAILGUN_URL, m.domain)
	req, err := http.NewRequest("POST", mgURL, strings.NewReader(form.Encode()))

	if err != nil {
		return fmt.Errorf("email - SendConfirmationEmail: failed to create request: %v", err)
	}

	req.SetBasicAuth("api", m.key)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)

	defer resp.Body.Close()
	bd, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bd))

	if resp.StatusCode == 401 {
		return MailgunUnauthorized
	}

	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		return MailgunFailure
	}

	return nil
}
