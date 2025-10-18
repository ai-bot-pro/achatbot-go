package params

type SerperApiArgs struct {
	Query string `json:"query"`
	GL    string `json:"gl"`   // Country code
	HL    string `json:"hl"`   // Language code
	Page  int    `json:"page"` // Page number
	Num   int    `json:"num"`  // Number of results
}

type SearchApiArgs struct {
	Query  string `json:"query"`
	Engine string `json:"engine"` // Search engine
	GL     string `json:"gl"`     // Country code
	HL     string `json:"hl"`     // Language code
	Page   int    `json:"page"`   // Page number
	Num    int    `json:"num"`    // Number of results
}
