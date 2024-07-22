package main

type multiflag []string

func (i *multiflag) String() string {
	return "my string representation"
}

func (i *multiflag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
