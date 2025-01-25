package compute

type Query struct {
	command Command
	args    []string
}

func NewQuery(command Command, args []string) Query {
	return Query{
		command: command,
		args:    args,
	}
}

func (q Query) Command() Command {
	return q.command
}

func (q Query) Args() []string {
	return q.args
}
