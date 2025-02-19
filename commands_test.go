package commands_test

type TestCommand struct {
	Value string
}

func (cmd TestCommand) Name() string {
	return "TestCommand"
}

type TestDoubleCommand struct {
	Value string
}

func (cmd TestDoubleCommand) Name() string {
	return "TestCommand"
}

type TestPointerCommand struct {
	Value string
}

func (cmd *TestPointerCommand) Name() string {
	return "TestPointerCommand"
}

type TestPanicCommand struct {
	Value string
}

func (cmd TestPanicCommand) Name() string {
	panic("TestPanicCommand")
}
