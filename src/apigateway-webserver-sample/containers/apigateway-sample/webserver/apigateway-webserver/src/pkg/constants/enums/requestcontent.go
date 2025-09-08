package enums

type RequestContentType int

const (
	None RequestContentType = iota + 1
	Urlencoded
	Raw
	Json
)

func (r RequestContentType) String() string {
	return [...]string{"", "application/x-www-form-urlencoded", "text/plain", "application/json"}[r-1]
}

func (r RequestContentType) EnumIndex() int {
	return int(r)
}
