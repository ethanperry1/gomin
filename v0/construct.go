package v0

type commandArguments struct {
	argType ArgType
	value   any
}

func (arg *commandArguments) Type() ArgType {
	return arg.argType
}

func (arg *commandArguments) Value() any {
	return arg.value
}

type commandSurface struct {
	parent  CommandSurface
	command CommandArgument
}

func (surface *commandSurface) Parent() (CommandSurface, bool) {
	return surface.parent, surface.parent != nil
}

func (surface *commandSurface) Command() CommandArgument {
	return surface.command
}

type packagesRuleBuilder struct {
	CommandSurface
}

func (builder *packagesRuleBuilder) Filter(s string) PackageContext {
	return &packageContextRuleBuilder{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Fltr,
				value:   s,
			},
		},
	}
}

type packageContextRuleBuilder struct {
	CommandSurface
}

func (builder *packageContextRuleBuilder) Files() Files {
	return &filesRuleBuilder{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Any,
			},
		},
	}
}

type filesRuleBuilder struct {
	CommandSurface
}

func (builder *filesRuleBuilder) Filter(s string) FileContext {
	return &fileContextRuleBuilder{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Fltr,
				value:   s,
			},
		},
	}
}

type fileContextRuleBuilder struct {
	CommandSurface
}

func (builder *fileContextRuleBuilder) Functions() Functions {
	return &functionsRuleBuilder{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Any,
			},
		},
	}
}

type functionsRuleBuilder struct {
	CommandSurface
}

func (builder *functionsRuleBuilder) Filter(s string) CommandSurface {
	return &functionCommandSurface{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Fltr,
				value:   s,
			},
		},
	}
}

type functionCommandSurface struct {
	CommandSurface
}

type packageInstanceRuleBuilder struct {
	CommandSurface
}

func (builder *packageInstanceRuleBuilder) Files() Files {
	return &filesRuleBuilder{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Any,
			},
		},
	}
}

func (builder *packageInstanceRuleBuilder) Functions() Functions {
	return &functionsRuleBuilder{
		CommandSurface: &commandSurface{
			parent: &commandSurface{
				parent: builder.CommandSurface,
				command: &commandArguments{
					argType: Any,
				},
			},
			command: &commandArguments{
				argType: Any,
			},
		},
	}
}

func (builder *packageInstanceRuleBuilder) File(name string) FileInstance {
	return &fileInstanceRuleBuilder{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Name,
				value:   name,
			},
		},
	}
}

type fileInstanceRuleBuilder struct {
	CommandSurface
}

func (builder *fileInstanceRuleBuilder) Functions() Functions {
	return &functionsRuleBuilder{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Any,
			},
		},
	}
}

func (builder *fileInstanceRuleBuilder) Method(receiver string, name string) CommandSurface {
	return &functionCommandSurface{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Pair,
				value:   NamePair{receiver, name},
			},
		},
	}
}

func (builder *fileInstanceRuleBuilder) Function(name string) CommandSurface {
	return &functionCommandSurface{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Name,
				value:   name,
			},
		},
	}
}

func (builder *fileInstanceRuleBuilder) Literal(index int) CommandSurface {
	return &functionCommandSurface{
		CommandSurface: &commandSurface{
			parent: builder.CommandSurface,
			command: &commandArguments{
				argType: Index,
				value:   index,
			},
		},
	}
}

func AllPackages() Packages {
	return &packagesRuleBuilder{
		CommandSurface: &commandSurface{
			parent: nil,
			command: &commandArguments{
				argType: Any,
			},
		},
	}
}

func AllFiles() Files {
	return &filesRuleBuilder{
		CommandSurface: &commandSurface{
			parent: &commandSurface{
				parent: nil,
				command: &commandArguments{
					argType: Any,
				},
			},
			command: &commandArguments{
				argType: Any,
			},
		},
	}
}

func AllFunctions() Functions {
	return &functionsRuleBuilder{
		CommandSurface: &commandSurface{
			parent: &commandSurface{
				parent: &commandSurface{
					parent: nil,
					command: &commandArguments{
						argType: Any,
					},
				},
				command: &commandArguments{
					argType: Any,
				},
			},
			command: &commandArguments{
				argType: Any,
			},
		},
	}
}

func Package(name string) PackageInstance {
	return &packageInstanceRuleBuilder{
		CommandSurface: &commandSurface{
			parent: nil,
			command: &commandArguments{
				argType: Name,
				value:   name,
			},
		},
	}
}
