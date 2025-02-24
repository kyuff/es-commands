package commands_test

type TestCommand struct {
	Value string
}

func (cmd TestCommand) CommandName() string {
	return "TestCommand"
}

type TestDoubleCommand struct {
	Value string
}

func (cmd TestDoubleCommand) CommandName() string {
	return "TestCommand"
}

type TestPointerCommand struct {
	Value string
}

func (cmd *TestPointerCommand) CommandName() string {
	return "TestPointerCommand"
}

type TestPanicCommand struct {
	Value string
}

func (cmd TestPanicCommand) CommandName() string {
	panic("TestPanicCommand")
}

type TestMyCommand struct {
	Value string
}

func (cmd TestMyCommand) CommandName() string {
	return "TestMyCommand"
}

func (cmd TestMyCommand) CustomMethod() bool {
	return true
}
