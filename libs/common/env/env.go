package env

type Environment string

const (
	Dev   Environment = "development"
	Prod  Environment = "production"
	Stage Environment = "staging"
)

func (e Environment) String() string {
	return string(e)
}

func (e Environment) IsValid() bool {
	switch e {
	case Dev, Stage, Prod:
		return true
	default:
		return false
	}
}
