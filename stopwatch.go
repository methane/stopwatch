/*
stopwatch is simple manual profiler
*/

package stopwatch

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"sync"
	"time"
)

var Enabled = true

type Stopper interface {
	Stop()
}

type measure struct {
	name    string
	start   time.Time
	stopped bool
}

type counter struct {
	name  string
	count int64
	total time.Duration
}

var mm sync.Mutex
var metrics = make(map[string]*counter)

func (m *measure) Stop() {
	if !Enabled || m.stopped {
		return
	}

	d := time.Now().Sub(m.start)
	if d < 0 {
		d = 0
	}

	mm.Lock()
	ct := metrics[m.name]
	if ct == nil {
		ct = &counter{name: m.name}
		metrics[m.name] = ct
	}
	ct.count++
	ct.total += d
	mm.Unlock()
}

type dummy bool

func (d dummy) Stop() {
}

func Start(name string) Stopper {
	if !Enabled {
		return dummy(false)
	}
	t := time.Now()

	_, file, line, _ := runtime.Caller(1)
	name = fmt.Sprintf("%s %s:%d", name, file, line)

	return &measure{
		name:  name,
		start: t,
	}
}

type ccs []counter

func (cs ccs) Len() int {
	return len(cs)
}
func (cs ccs) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}
func (cs ccs) Less(i, j int) bool {
	return cs[i].total > cs[j].total
}

func Show() string {
	var cs ccs
	mm.Lock()
	for _, c := range metrics {
		cs = append(cs, *c)
	}
	mm.Unlock()
	sort.Sort(cs)

	buf := new(bytes.Buffer)
	buf.WriteString("name\tcount\tavg\ttotal\n")
	for _, c := range cs {
		avg := c.total / time.Duration(c.count)
		fmt.Fprintf(buf, "%s\t%d\t%s\t%s\n", c.name, c.count, avg, c.total)
	}
	return buf.String()
}

func Reset() {
	mm.Lock()
	metrics = make(map[string]*counter)
	mm.Unlock()
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := Show()
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(data))
}

func init() {
	http.HandleFunc("/debug/stopwatch", ServeHTTP)
}
