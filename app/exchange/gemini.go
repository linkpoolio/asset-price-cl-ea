package exchange

import (
	"fmt"
	"strings"
)

type Gemini struct {
	Exchange
	Pairs []*Pair
}

type GeminiPairs []string

type GeminiTicker struct {
	Last   string        `json:"last"`
	Volume GeminiVolume  `json:"volume"`
}

type GeminiVolume map[string]interface{}

func (exc *Gemini) GetResponse(base, quote string) (*Response, *Error) {
	var gt GeminiTicker
	config := exc.GetConfig()
	excErr := HttpGet(config, fmt.Sprintf("/pubticker/%s%s", base, quote), &gt)
	if excErr != nil {
		return nil, excErr
	}

	qv, ok := gt.Volume[strings.ToUpper(quote)]
	if !ok {
		return  nil, &Error{config.Name, "400", "Quote volume wasn't returned in the ticker"}
	}

	return &Response{
		Name:   config.Name,
		Price:  ToFloat64(gt.Last),
		Volume: ToFloat64(qv),
	}, nil
}

func (exc *Gemini) SetPairs() *Error {
	var gp GeminiPairs
	config := exc.GetConfig()

	err := HttpGet(config, "/symbols", &gp)
	if err != nil {
		return err
	}

	for _, p := range gp {
		if len(p) == 6 {
			exc.Pairs = append(exc.Pairs, &Pair{strings.ToUpper(p[:3]), strings.ToUpper(p[3:])})
		}
	}
	return nil
}

func (exc *Gemini) GetConfig() *Config {
	return &Config{
		Name:    "Gemini",
		BaseUrl: "https://api.gemini.com/v1",
		Pairs:   exc.Pairs,
	}
}
