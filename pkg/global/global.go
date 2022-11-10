package global

import (
	"log"
	"runtime"

	"github.com/yunling101/cayo/pkg/types"
)

const (
	Version   = "v1.1"
	VerifyKey = "e944c3373328376c09c189c55441d35b"
)

var (
	ClientGlobal *types.FlagConfig
	Cache        *types.MemoryCache
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
