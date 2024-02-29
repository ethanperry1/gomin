package v0

type GlobalOption func() ([]*ruleSet, error)
type Option func() ([]*ruleSet, error)

type ArgType int

const (
	Any ArgType = iota
	Fltr
	Name
	Pair
	Index
)

type CommandArgument interface {
	Type() ArgType
	Value() any
}

type CommandSurface interface {
	Parent() (CommandSurface, bool)
	Command() CommandArgument
}

type PackageInstance interface {
	CommandSurface
	FileParent
	FunctionParent
	File(name string) FileInstance
}

type FileInstance interface {
	CommandSurface
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
