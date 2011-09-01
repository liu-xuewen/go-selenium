/* Remote selenium client */
package selenium

import (
	"bytes"
	"fmt"
	"http"
	"io/ioutil"
	"json"
	_ "log"
	"os"
	"strings"
)

var errors = map[int]string{
	7:  "no such element",
	8:  "no such frame",
	9:  "unknown command",
	10: "stale element reference",
	11: "element not visible",
	12: "invalid element state",
	13: "unknown error",
	15: "element is not selectable",
	17: "javascript error",
	19: "xpath lookup error",
	21: "timeout",
	23: "no such window",
	24: "invalid cookie domain",
	25: "unable to set cookie",
	26: "unexpected alert open",
	27: "no alert open",
	28: "script timeout",
	29: "invalid element coordinates",
	32: "invalid selector",
}

const (
	SUCCESS          = 0
	DEFAULT_EXECUTOR = "http://127.0.0.1:4444/wd/hub"
	JSON_TYPE        = "application/json"
)

type BrowserProfile *map[string]string

type BrowserCapabilities struct {
	BrowserName, Version, Platform string
	JavaScriptEnabled, takesScreenShots, HandlesAlerts, DatabaseEnabled,
	LocationContextEnabled, ApplicationCacheEnabled,
	BrowserConnectionEnabled, CSSEnabled, WebStorageEnabled, Rotatable,
	AcceptSSLCerts, NativeElements bool
}

type WebDriver struct {
	SessionId, Executor string
	Capabilities        *BrowserCapabilities
	profile             BrowserProfile
}

type serverReply struct {
	SessionId *string
	Status    int
}

func isMimeType(response *http.Response, mtype string) bool {
	if ctype, ok := response.Header["Content-Type"]; ok {
		return strings.HasPrefix(ctype[0], mtype)
	}

	return false
}

func (wd *WebDriver) makeRequest(cmd *Command, params *Params) (*http.Request, os.Error) {
	path, err := cmd.URL(params)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	url := wd.Executor + path
	request, err := http.NewRequest(cmd.Method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", JSON_TYPE)

	return request, nil
}

func cleanNils(buf []byte) {
	for i, b := range buf {
		if b == 0 {
			buf[i] = ' '
		}
	}
}

func (wd *WebDriver) execute(cmd *Command, params *Params) ([]byte, os.Error) {
	request, err := wd.makeRequest(cmd, params)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		buf = []byte(response.Status)
	}

	if (err != nil) || (response.StatusCode >= 400) {
		return nil, os.NewError(string(buf))
	}

	cleanNils(buf)

	reply := new(serverReply)
	if isMimeType(response, JSON_TYPE) {
		err := json.Unmarshal(buf, reply)
		if err != nil {
			return nil, err
		}

		if reply.Status != 0 {
			message, ok := errors[reply.Status]
			if !ok {
				message = fmt.Sprintf("unknown error - %d", reply.Status)
			}

			return nil, os.NewError(message)
		}
		return buf, err
	} else if isMimeType(response, "image/png") {
		// FIXME
	}

	message := "no content type in reply"
	ctype, ok := response.Header["Content-Type"]
	if ok {
		message = fmt.Sprintf("unknown reply content type: %v", ctype)
	}
	return nil, os.NewError(message)

}

type Build struct {
	Version, Revision, Time string
}

type OS struct {
	Arch, Name, Version string
}

type Status struct {
	Build
	OS
}

type statusReply struct {
	Value Status
}

func (wd *WebDriver) Status() (*Status, os.Error) {
	reply, err := wd.execute(CMD_STATUS, nil)
	if err != nil {
		return nil, err
	}

	status := new(statusReply)
	err = json.Unmarshal(reply, status)
	if err != nil {
		return nil, err
	}

	return &status.Value, nil
}

func (wd *WebDriver) NewSession() (string, os.Error) {
	_, err := wd.execute(CMD_NEW_SESSION, nil)
	if err != nil {
		return "", err
	}

	return wd.SessionId, nil
}


func New(capabilities *BrowserCapabilities, executor string,
profile BrowserProfile) (*WebDriver, os.Error) {

	if len(executor) == 0 {
		executor = DEFAULT_EXECUTOR
	}

	wd := &WebDriver{Executor: executor,
		Capabilities: capabilities,
		profile:      profile}

	// FIXME: Handle profile

	return wd, nil
}