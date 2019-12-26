package survey

import (
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

func SetOutoutDir(dir string) {
	workDir = dir
}

func Quit() {
	famLocker.Lock()
	defer famLocker.Unlock()

	famsNum := len(fams)
	if famsNum == 0 {
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(famsNum)
	for name, fam := range fams {
		fam.quitChan <- wg
		delete(fams, name)
	}
	wg.Wait()
}

func expose(v periodism, family, name string, param *VarParam) bool {
	family = normalizeName(family)
	name = normalizeName(name)

	fam := getFamily(family)
	if fam == nil {
		return false
	}
	v.setName(name)
	if !param.NoDump {
		exproter := fam.recordfile.AddRecordValue(name)
		v.setDumper(exproter)
	}
	fam.vars.append(v)
	return true
}

type family struct {
	name       string
	filename   string
	recordfile *RecordFile
	quitChan   chan *sync.WaitGroup
	vars       *periodNode
}

var (
	famLocker sync.Mutex
	fams      = map[string]*family{}
	workDir   string
	nilvar    NilVar
)

func init() {
	if wd, err := os.Getwd(); err != nil {
		workDir = "./"
	} else {
		workDir = wd
	}
}

func getFamily(name string) *family {
	famLocker.Lock()
	defer famLocker.Unlock()
	fam, ok := fams[name]
	if !ok {
		fam = addFamily(name)
	}
	return fam
}

// lock by caller
func addFamily(name string) *family {
	filename := path.Join(workDir, name+".txt")
	rf, err := OpenRecordFile(filename, 128)
	if err != nil {
		return nil
	}
	fam := &family{
		name:       name,
		filename:   filename,
		recordfile: rf,
		quitChan:   make(chan *sync.WaitGroup),
		vars:       &periodNode{nil, nil},
	}
	go fam.bgrun()
	fams[name] = fam
	return fam
}

func (fam *family) bgrun() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer fam.recordfile.Close()

	for {
		select {
		case <-ticker.C:
			fam.tickPer1s()
		case wg := <-fam.quitChan:
			fam.flushvars()
			wg.Done()
			return
		}
	}
}

func (fam *family) tickPer1s() {
	fam.vars.traverse(func(period periodism) {
		period.tick()
	})
}

func (fam *family) flushvars() {
	fam.vars.traverse(func(period periodism) {
		if v, ok := period.(Var); ok {
			v.Flush()
		}
	})
}

func normalizeName(s string) string {
	return strings.ToLower(s)
}
