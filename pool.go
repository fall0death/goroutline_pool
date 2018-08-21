package pool

type func_arg struct {
	f func([]interface{})

	args []interface{}
}
type Process_pool struct {
	running int32

	total int32

	ch chan *func_arg
}

func (p *Process_pool) run() {
	for sig := range p.ch {
		if sig == nil {
			return
		}
		sig.f(sig.args)

	}
}
