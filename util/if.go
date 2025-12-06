package util

type IDisplay interface {
	String() string
}

type IMapForeach interface {
	Foreach(do func(key, value any) (isContinue bool)) // 遍历, isContinue:是否继续遍历
}
