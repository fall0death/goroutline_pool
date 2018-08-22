package pool

import (
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
	"time"
)

//要传的函数和参数
type funcArg struct {
	f func([]interface{}) //要执行的函数

	args []interface{} //要执行的函数参数
}

//线程池的结构体
type Processpool struct {
	running int32 //当前正在运行的线程数

	total int32 //总共允许的最大线程数

	ch chan *funcArg //要传的函数和其参数的通道

	close int32 //线程池是否关闭

	lock sync.Mutex //Mutex锁

	once sync.Once //只执行一次（用来执行关闭线程池的操作)
}

//线程池中子线程的运行
func (p *Processpool) run() {
	defer p.decRunning()
	p.incRunning()
	for {
		sig, ok := <-p.ch
		if !ok {
			return
		}
		if sig == nil {

			return
		}
		sig.f(sig.args)
	}
}

//错误信息显示
var (
	InvalidSize = errors.New("invalid size for pool")
	Poolclosed  = errors.New("Pool has been closed")
)

//开启子线程，给线程池的运行线程总数加1
func (p *Processpool) incRunning() {
	atomic.AddInt32(&p.running, 1)
}

//结束子线程后，给线程池的运行线程总数减1
func (p *Processpool) decRunning() {
	atomic.AddInt32(&p.running, -1)
}
func (p *Processpool) Running() int32 {
	return atomic.LoadInt32(&p.running)
}
func (p *Processpool) Size() int32 {
	return p.total
}

//建立一个新的线程池样例
func NewProcessPool(size int32) (*Processpool, error) {
	if size <= 0 {
		return nil, InvalidSize
	}
	p := &Processpool{
		running: 0,
		total:   size,
		close:   0,
		ch:      make(chan *funcArg, size),
	}
	//for i:=0;i<int(p.total/4);i++ {
	//	go p.run()
	//}
	go p.periodDect()
	return p, nil
}

//提交要运行的函数，交给子线程运行
func (p *Processpool) ProcessSubmit(function func([]interface{}), args []interface{}) error {
	if atomic.LoadInt32(&p.close) == 1 {
		p.lock.Unlock()
		return Poolclosed
	}
	p.ch <- &funcArg{function, args}
	return nil
}

//心跳运行，动态规划线程的数量
func (p *Processpool) periodDect() {
	heartbeat := time.NewTicker(time.Millisecond)
	var proinc int = 0
	var prodec int = 0
	for range heartbeat.C {
		if atomic.LoadInt32(&p.close) == 1 {
			return
		}
		if len(p.ch) > 0 {
			proinc++
			prodec = 0
			if proinc >= 2 && p.running < p.total {
				go p.run()
				proinc = 0
			}
		} else {
			proinc = 0
			prodec++
			if prodec >= 5 {
				p.ch <- nil
				prodec = 0
			}
		}
	}
}

//结束线程池的运行
func (p *Processpool) PoolClose() {
	p.once.Do(func() {
		p.lock.Lock()
		atomic.StoreInt32(&p.close, 1)
		close(p.ch)
		p.lock.Unlock()
	})
}
