package a

var pv = 1 // want `variable name "pv" is shorter than 3 characters`

var longEnough = 2

const cx = 3 // want `constant name "cx" is shorter than 3 characters`

const (
	ab = iota // want `constant name "ab" is shorter than 3 characters`
	abc
)

type thing struct {
	ab int // fields are not checked
}

func (t *thing) method(ab int, cd int) (rv int, err error) {
	x := 1 // want `variable name "x" is shorter than 3 characters`
	_ = x
	var y = 2 // want `variable name "y" is shorter than 3 characters`
	_ = y
	var ab2 int // three characters is long enough
	_ = ab2
	for k, v := range []int{ab, cd} { // want `variable name "k" is shorter than 3 characters` `variable name "v" is shorter than 3 characters`
		_, _ = k, v
	}
	for i := 0; i < 3; i++ {
	}
	for j := 0; j < 3; j++ { // want `variable name "j" is shorter than 3 characters`
	}
	if ok := t != nil; ok {
	}
	mp := map[string]int{}    // want `variable name "mp" is shorter than 3 characters`
	if n, ok := mp["k"]; ok { // want `variable name "n" is shorter than 3 characters`
		_ = n
	}
	var any1 interface{} = 1
	switch s := any1.(type) { // want `variable name "s" is shorter than 3 characters`
	case int:
		_ = s
	}
	switch any1.(type) {
	case int:
	}
	fn := func(q int) int { return q } // want `variable name "fn" is shorter than 3 characters`
	_ = fn
	tx, wg := 1, 2
	_, _ = tx, wg
	_, z := 1, 2 // want `variable name "z" is shorter than 3 characters`
	_ = z
	ch := make(chan int) // want `variable name "ch" is shorter than 3 characters`
	select {
	case cv := <-ch: // want `variable name "cv" is shorter than 3 characters`
		_ = cv
	default:
	}
	var _ = 1
	rv = ab + cd
	return rv, err
}

func gen[T any](z T) T { return z }
