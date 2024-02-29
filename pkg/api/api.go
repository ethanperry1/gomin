package api

type Option func()

type ArgType int

const (
	Any ArgType = iota
	Fltr
	Name
)

type CommandArgument interface {
	Type() ArgType
	Value() string
}

type CommandSurface interface {
	Parent() (CommandSurface, bool)
	Command() CommandArgument
}

type PackageInstance interface {
	FileParent
	FunctionParent
	File(name string) FileInstance
}

type FileInstance interface {
	FunctionParent
	Method(receiver string, name string) CommandSurface
	Function(name string) CommandSurface
	Literal(index int) CommandSurface
}

type FunctionParent interface {
	Functions() Functions
}

type FileParent interface {
	Files() Files
}

type PackageContext interface {
	FileParent
	CommandSurface
}

type FileContext interface {
	FunctionParent
	CommandSurface
}

type Packages Context[PackageContext]

type Files Context[FileContext]

type Functions Context[CommandSurface]

type Context[T CommandSurface] interface {
	Filter[T]
	CommandSurface
}

type Filter[T CommandSurface] interface {
	Filter(string) T
}
