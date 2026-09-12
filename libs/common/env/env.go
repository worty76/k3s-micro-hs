package env

type Environment string

const (
	Dev   Environment = "dev"
	Prod  Environment = "prod"
	Stage Environment = "stage"
)

func (e Environment) String() string {
	return string(e)
}
