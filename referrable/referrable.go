package referrable

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chonla/yas/response"
	"github.com/stretchr/objx"
)

// Referrable is referrable items
type Referrable struct {
	values map[string][]string
	data   objx.Map
}

// NewReferrable creates an referrable object
func NewReferrable(resp *response.Response) (*Referrable, error) {
	values := map[string][]string{}

	values["statuscode"] = []string{fmt.Sprintf("%d", resp.StatusCode)}
	values["status"] = []string{resp.Status}

	for k, v := range resp.Header {
		key := strings.ToLower(fmt.Sprintf("header.%s", k))
		if values[key] == nil {
			values[key] = []string{}
		}
		for _, t := range v {
			values[key] = append(values[key], t)
		}
	}

	jsonObj, e := objx.FromJSON(resp.Body)
	if e != nil {
		return nil, e
	}

	return &Referrable{
		values: values,
		data:   jsonObj,
	}, nil
}

// Find to find a value of given key
func (a *Referrable) Find(k string) ([]string, bool) {
	re := regexp.MustCompile("(?i)^data\\.(.+)")
	match := re.FindStringSubmatch(k)
	if len(match) > 1 {
		if a.data.Has(match[1]) {
			return []string{a.data.Get(match[1]).String()}, true
		}
		return []string{}, false
	} else {
		val, ok := a.values[k]
		return val, ok
	}
}