package zip1

import "testing"

func TestZip(t *testing.T) {
	zipDir("/home/chenchong/gopath/src/awesomeProject/test", "/home/chenchong/gopath/src/awesomeProject/test/zip/test.zip")
}

func TestUnzip(t *testing.T) {
	unzip("/home/chenchong/gopath/src/awesomeProject/test/zip/test.zip", "/home/chenchong/gopath/src/awesomeProject/test/zip/test")
}
