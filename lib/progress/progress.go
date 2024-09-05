package progress

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

const spinnerLen = 12
const cyclesPerSec = 0.5 

var spinnerChars = [spinnerLen]rune{'\u2581', '\u2583', '\u2584', '\u2585', '\u2586', '\u2587', '\u2588', '\u2587', '\u2586', '\u2585', '\u2584', '\u2583'}

const delayms = 1000 / (spinnerLen * cyclesPerSec)

const ansi_alt_buf_enable = "\x1b[?1049h"
const ansi_alt_buf_disable = "\x1b[?1049l"
const ansi_clear_screen = "\x1b[J"
const ansi_reset_cursor = "\x1b[H"

const completed_char = '\u2713'
const failed_char = '\u2715'
const skipped_char = '-'

const line_width = 60

type jobstat struct {
	line    string
	status  bool
	failed  bool
	spinner bool
	skipped bool
}

type Job struct {
	name        string
	completed   bool
	failed      bool
	jobs        []*Job
	level       int
	runner      func(*log.Logger, types.Target) bool
	configure   func() error
	target      types.Target
	canParallel bool
}

func (j *Job) NewChild(name string) *Job {
	nj := NewJob(name)
	nj.level = j.level + 1
	j.jobs = append(j.jobs, nj)
	return nj
}

func (j *Job) Run(ml *log.Logger) bool {
	var res bool
	if len(j.jobs) > 0 || j.runner == nil {
		res = true
		if j.canParallel {
			wg := sync.WaitGroup{}
			for _, cj := range j.jobs {
				wg.Add(1)
				go func() {
					cj.Run(ml)
					wg.Done()
				}()
			}
			wg.Wait()
			for _, cj := range j.jobs {
				if !cj.completed {
					res = false
				}
			}
		} else {
			for _, cj := range j.jobs {
				cres := cj.Run(ml)
				if !cres {
					res = false
					break
				}
			}
		}
	} else {
		if j.configure != nil {
			err := j.configure()
			if err != nil {
				res = false
				goto end
			}
		}
		res = j.runner(ml, j.target)
	}

end:
	if res {
		j.completed = true
	} else {
		j.failed = true
	}

	return res
}

func (j *Job) line() []*jobstat {
	var leader string
	if j.level > 0 {
		leader = "\u2514\u2500" + strings.Repeat("\u2500", (j.level-1)*2)
	}

	mj := jobstat{
		line:    leader + j.name,
		status:  j.completed,
		failed:  j.failed,
		spinner: true,
	}

	stats := []*jobstat{&mj}

	for _, js := range j.jobs {
		for _, lin := range js.line() {
			if j.failed && !lin.status && !lin.failed {
				lin.skipped = true
			}
			stats = append(stats, lin)
		}
	}

	return stats
}

func (j *Job) WithFunc(runner func(*log.Logger, types.Target) bool) *Job {
	j.runner = runner
	return j
}

func (j *Job) WithTarget(target types.Target) *Job {
	j.target = target
	return j
}

func (j *Job) WithConfigure(configure func() error) *Job {
	j.configure = configure
	return j
}

func (j *Job) WithParallel() *Job {
	j.canParallel = true
	return j
}

func NewJob(name string) *Job {
	return &Job{
		name: name,
	}
}

type Progress struct {
	jobs     []*Job
	stop     bool
	wg       sync.WaitGroup
	spinProg int
}

func (p *Progress) Render(name string, ml *log.Logger) bool {
	p.wg.Add(1)
	go p.render(name)
	res := true
	for _, jg := range p.jobs {
		r := jg.Run(ml)
		if !r {
			res = false
		}
	}
	p.wg.Wait()

	p.genFrame(name, os.Stdout)

	return res
}

func (p *Progress) render(name string) {
	fmt.Print(ansi_alt_buf_enable)

	buf := bytes.NewBuffer(nil)

	for {
		buf.Truncate(0)

		fmt.Print(ansi_clear_screen)
		fmt.Print(ansi_reset_cursor)

		stop := p.genFrame(name, buf)
		fmt.Println(buf.String())

		if stop {
			break
		}

		time.Sleep(time.Duration(math.Ceil(delayms)) * time.Millisecond)
	}

	fmt.Print(ansi_alt_buf_disable)
	p.wg.Done()
}

func (p *Progress) genFrame(top string, w io.Writer) bool {
	comp := true
	fmt.Fprintln(w, fmt.Sprintf("lbt: %s", top))
	fmt.Fprintln(w, strings.Repeat("\u2550", line_width+2))

	for _, jg := range p.jobs {
		stats := jg.line()
		for _, js := range stats {
			if js.spinner {
				var eChar rune
				if js.status {
					eChar = completed_char
				} else if js.failed {
					eChar = failed_char
				} else if js.skipped {
					eChar = skipped_char
				} else {
					eChar = spinnerChars[p.spinProg]
					comp = false
				}
				fmt.Fprintf(w, "%-*s %c\n", line_width, js.line, eChar)
			} else {
				fmt.Fprintf(w, "%-*s\n", line_width, js.line)
			}
		}
		fmt.Fprintln(w, " ")
	}

	if comp || p.stop {
		return true
	}

	p.spinProg = (p.spinProg + 1) % spinnerLen
	return false
}

func NewProgress(jobs ...*Job) *Progress {
	return &Progress{
		jobs: jobs,
	}
}
