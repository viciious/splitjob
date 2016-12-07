package splitjob

import (
	"io"
	"bufio"
	"fmt"
	"testing"
	"os"
)

func splitPrintFile(f *os.File) error {
	var key uint32

	r := bufio.NewReader(f)

	readLine := func() (interface{}, uint32) {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			return nil, 0
		}
		key_ := key
		key = key + 1
		return line, key_
	}

	handleLine := func(line string) error {
		fmt.Print(line)
		return nil
	}

	job := New(&Options{
		Spawn: func() ThinkFn {
			return func(obj interface{}) error {
				return handleLine(obj.(string))
			}
		},
		Pull: func() (interface{}, uint32, bool) {
			o, s := readLine()
			if o == nil {
				return nil, 0, true
			}
			return o, s, false
		},
		NumSplits: 5,
		ChanSize:  2,
	})

	return job.Do()
}

func testFilePrint(t *testing.T, fName string) {
	f, err := os.Open(fName)
	if err != nil {
		t.Error(fmt.Sprintf("Could not open %s", fName))
		return
	}
	defer f.Close()

	err = splitPrintFile(f)
	if err != nil {
		t.Error(fmt.Sprintf("Printing %s failed with error %s", fName, err))
		return
	}

	fmt.Println(fmt.Sprintf("DONE with %s", fName))
}

func TestFilePrintJob(t *testing.T) {
	testFilePrint(t, "tests/file1.txt")
}
