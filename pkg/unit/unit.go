package unit

type Unit interface {
	Name() string
	Children() []Unit
}

type Project struct {
	name string
}

func (project *Project) Name() string {
	return project.name
}

type Package struct {
	name string
}

func (pack *Package) Name() string {
	return pack.name
}

type File struct {
	name string
}

func (file *File) Name() string {
	return file.name
}

type Block struct {
	name string
}

func (block *Block) Name() string {
	return block.name
}