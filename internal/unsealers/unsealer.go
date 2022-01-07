package unsealers

type UnsealData struct {
	Keys      []string `yaml:"keys,omitempty"`
	Vaults    []string `yaml:"vaults,omitempty"`
	RootToken string   `yaml:"root_token,omitempty"`
}

type Unsealer interface {
	Prerequisites() error
	FetchUnsealData(map[string]interface{}) (*UnsealData, error)
	Unseal(UnsealData) error
}
