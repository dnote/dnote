package command

import (
	"os"

	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

type Command struct {
	// Name is the name of the command. This is used to find the subcommand
	// and is case sensitive.
	Name string

	// flags is a FlagSet which can parse and hold flags and their values.
	flags *flag.FlagSet

	// args is a list of arments
	args []string

	// RunE is a function that contains the logic for the command
	RunE func(cmd *Command, args []string) error

	// commands is the list of subcommands that this command has.
	commands []*Command

	// Parent is a pointer to the parent command, if any, of which the command
	// is a subcommand.
	Parent *Command

	Use           string
	Short         string
	SilenceErrors bool
	SilenceUsage  bool
	Example       string
	Aliases       []string
	PreRunE       func(cmd *Command, args []string) error
	Deprecated    string
	Long          string
	Run           func(cmd *Command, args []string)
}

// Flags returns a flag set for the command. If not initialized yet, it initializes
// one and returns the result.
func (c *Command) Flags() *flag.FlagSet {
	if c.flags == nil {
		c.flags = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	}

	return c.flags
}

// ParseFlags parses the given slice of arguments using the flag set of the command.
func (c *Command) ParseFlags(args []string) error {
	err := c.Flags().Parse(args)
	if err != nil {
		return errors.Wrap(err, "Error parsing flags")
	}

	return nil
}

// Root returns the root command of the given command.
func (c *Command) Root() *Command {
	if c.Parent != nil {
		return c.Parent.Root()
	}

	return c
}

// setArgs sets the arguments for the command. It is useful while writing tests.
func (c *Command) setArgs(args []string) {
	c.args = args
}

// Args returns the argument for the command. By default, os.Args[1:] is used.
func (c *Command) Args() []string {
	args := c.args
	if c.args == nil {
		args = os.Args[1:]
	}

	return args

}

// Execute runs the root command. It is meant to be called on the root command.
func (c *Command) Execute() error {
	// Call Execute on the root command
	if c.Parent != nil {
		return c.Root().Execute()
	}

	args := c.Args()
	log.Debug("root command received arguments: %s\n", args)

	cmd := c.findSubCommand(args[0])
	if cmd == nil {
		// not found. show suggestion
		return nil
	}

	cmd.execute(args[1:])

	return nil
}

// execute runs the command.
func (c *Command) execute(args []string) error {
	log.Debug("command '%s' called with arguments: %s\n", c.Name, args)

	if err := c.ParseFlags(args); err != nil {
		return err
	}

	nonFlagArgs := c.Flags().Args()
	log.Debug("command '%s' called with non-flag arguments: %s\n", c.Name, nonFlagArgs)

	if err := c.RunE(c, nonFlagArgs); err != nil {
		return err
	}

	return nil
}

// hasAlias checks whether the command has the given alias.
func (c *Command) hasAlias(targetAlias string) bool {
	for _, alias := range c.Aliases {
		if alias == targetAlias {
			return true
		}
	}

	return false
}

// findSubCommand finds and returns an appropriate subcommand to be called, based
// on the given slice of arguments. It also returns a slice of arguments with which
// the subcommand should be called.
func (c *Command) findSubCommand(name string) *Command {
	log.Debug("sub-command: '%s'\n", name)

	for _, cmd := range c.commands {
		if cmd.Name == name || cmd.hasAlias(name) {
			return cmd
		}
	}

	return nil
}

// AddCommand adds the given command as a subcommand.
func (c *Command) AddCommand(cmd *Command) {
	cmd.Parent = c

	c.commands = append(c.commands, cmd)
}
