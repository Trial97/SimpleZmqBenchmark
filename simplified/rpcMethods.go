package main

//Args used for the sum
type Args struct {
	A int
	B int
}

//MyServer used to sum the numbers
type MyServer struct{}

//Sum adds numbers
func (b *MyServer) Sum(a *Args, r *int) error {
	*r = a.A + a.B
	return nil
}
