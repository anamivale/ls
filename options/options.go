package options

type Flags struct {
	Long      bool
	Recursive bool
	All       bool
	Reverse   bool
	Time      bool
}

func Options(flags string) Flags {
	var flagsStruct Flags

	for _, flag := range flags {
		switch flag {
		case 'R':
			flagsStruct.Recursive = true
		case 'l':
			flagsStruct.Long = true
		case 'r':
			flagsStruct.Reverse = true
		case 't':
			flagsStruct.Time = true
		case 'a':
			flagsStruct.All = true

		}
	}
	return flagsStruct
}
