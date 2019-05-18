package strpair

type StrPair struct {
	str1, str2 string
}

func New(str1, str2 string) StrPair {
	return StrPair{str1, str2}
}

func (f StrPair) Str1() string {
	return f.str1
}

func (f StrPair) Str2() string {
	return f.str2
}

func (f StrPair) Get() (string, string) {
	return f.str1, f.str2
}
