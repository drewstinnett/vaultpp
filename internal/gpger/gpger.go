package gpger

type Gpger interface {
	Recipients([]string) error
	Encrypt([]byte) ([]byte, error)
}
