package env

import (
	"flag"
	"os"
	"path"
	"strings"

	"github.com/pelletier/go-toml"
)

var WorkPath = func() string {
	cpath, _ := os.Getwd()
	cpath = strings.Replace(cpath, "\\", "/", strings.Count(cpath, "\\")) //这可能是将连接的其他符号，全部转换成我们需要的格式

	rpath := cpath
	for {
		_, err := os.Stat(path.Join(rpath, "conf"))
		if err == nil || os.IsExist(err) {
			return rpath
		}
		if rpath == path.Dir(rpath) {
			return cpath
		}
		rpath = path.Dir(rpath)
	}
}()

var (
	tomlsSetPath = path.Join(WorkPath, "conf")
	cpatchMap    = map[string]interface{}{}
	config       *toml.Tree
)

var (
	commandC      = flag.String("c", "", "config name")
	commandCPath  = flag.String("cpath", "", "choose absolute config path to use")
	commandCPatch = flag.String("cpatch", "", "an easy monkey patch for config")
)

func init() {
	os.Chdir(WorkPath)
	flag.Parse()

	var configPath string
	switch true {
	case *commandCPath != "":
		configPath = *commandCPath
	case *commandC != "":
		configPath = path.Join(tomlsSetPath, *commandC+".toml")
	default:
		configPath = path.Join(tomlsSetPath, "SilverBussinessServer.toml")
	}

	var err error
	config, err = toml.LoadFile(configPath)
	if err != nil {
		//panic(err.Error())
	}

	if *commandCPatch != "" {
		escapeCPatchString(*commandCPatch, &cpatchMap)
	}
}

func escapeCPatchString(cpatch string, m *map[string]interface{}) {
	sslice := strings.Split(cpatch, ";")
	for _, i := range sslice {
		kvslice := strings.Split(i, "=")
		(*m)[kvslice[0]] = kvslice[1]
	}
}

// Get : return value from config.toml
func Get(key string) interface{} {
	value, ok := cpatchMap[key]
	if ok == true {
		return value
	}
	value = config.Get(key)
	cpatchMap[key] = value
	return value
}
