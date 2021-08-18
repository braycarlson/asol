package asol

type Login struct {
	AccountId      float64
	Connected      bool
	Error          bool
	GasToken       string
	IdToken        string
	IsInLoginQueue bool
	IsNewPlayer    bool
	Puuid          string
	State          string
	SummonerId     float64
	UserAuthToken  string
	Username       string
}
