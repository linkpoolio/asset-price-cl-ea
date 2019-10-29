package exchange

import (
	"fmt"
	"strings"
)

type Gemini struct {
	Exchange
}

type GeminiPairs []string

type GeminiTicker struct {
	Last   string        `json:"last"`
	Volume GeminiVolume  `json:"volume"`
}

type GeminiVolume map[string]interface{}

func (exc *Gemini) GetResponse(base, quote string) (*Response, error) {
	var gt GeminiTicker
	config := exc.GetConfig()
	excErr := exc.HttpGet(config, fmt.Sprintf("/pubticker/%s%s", base, quote), &gt)
	if excErr != nil {
		return nil, excErr
	}

	qv, ok := gt.Volume[strings.ToUpper(quote)]
	if !ok {
		return  nil, &Error{config.Name, "400", "Quote volume wasn't returned in the ticker"}
	}

	return &Response{
		Name:   config.Name,
		Price:  exc.toFloat64(gt.Last),
		Volume: exc.toFloat64(qv),
	}, nil
}

func (exc *Gemini) RefreshPairs() error {
	var gp GeminiPairs
	config := exc.GetConfig()

	err := exc.HttpGet(config, "/symbols", &gp)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for _, p := range gp {
		if len(p) == 6 {
			pairs = append(pairs, &Pair{strings.ToUpper(p[:3]), strings.ToUpper(p[3:])})
		}
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Gemini) GetConfig() *Config {
	return &Config{Name: "Gemini", BaseURL: "https://api.gemini.com/v1"}
}
