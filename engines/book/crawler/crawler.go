package crawler

import (
	"fmt"
	"reflect"
	"sync"
)

// Command ...
type Command struct {
	Parser string      // 解析器名称
	Data   interface{} // 数据
}

// SiteParser 网页解析器
type SiteParser interface {
	Command(data interface{}) Command // 将解析器转成字符器
	Parser(cmd *Command) bool         // 解析
}

// TaskGroup ...
type TaskGroup struct {
	Sites          []SiteParser
	cmdChan        chan Command
	stop           chan bool
	wg             sync.WaitGroup
	parserCallback map[string]func(cmd Command) // 解析器回调处理函数表
}

// NewTaskGroup ...
func NewTaskGroup(taskNum int) *TaskGroup {
	tasks := &TaskGroup{
		cmdChan:        make(chan Command, 100),
		stop:           make(chan bool),
		wg:             sync.WaitGroup{},
		parserCallback: make(map[string]func(cmd Command)),
	}
	for i := 0; i < taskNum; i++ {
		go tasks.routine()
	}

	return tasks
}

// SetParser ...
func (c *TaskGroup) SetParser(site SiteParser, call func(cmd Command)) SiteParser {
	name := fmt.Sprint(reflect.TypeOf(site))
	c.Sites = append(c.Sites, site)
	c.parserCallback[name[1:]] = call

	return site
}

// AddCommand ...
func (c *TaskGroup) AddCommand(cmd Command) {
	c.wg.Add(1)
	c.cmdChan <- cmd
}

// GetParser ...
func (c *TaskGroup) GetParser(name string) SiteParser {
	for _, parser := range c.Sites {
		parserName := fmt.Sprint(reflect.TypeOf(parser))[1:]
		if parserName == name {
			return parser
		}
	}

	return nil
}

// Wait ...
func (c *TaskGroup) Wait() {
	c.wg.Wait()
	c.stop <- true
}

func (c *TaskGroup) routine() {
	for {
		select {
		case cmd := <-c.cmdChan:
			if parser := c.GetParser(cmd.Parser); parser != nil {
				if parser.Parser(&cmd) == true {
					call := c.parserCallback[cmd.Parser]
					if call != nil {
						call(cmd)
					}
				}
			}
			c.wg.Done()
		case stop := <-c.stop:
			if stop {
				c.stop <- true
				return
			}
		}
	}
}
