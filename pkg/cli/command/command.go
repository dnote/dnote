package command

import (
	"os"
	"strings"

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

// flagHasDefaultOption checks whether the given flag has a default value
// in the given flag set.
func flagHasDefaultOption(flag *flag.Flag, flagSet *flag.FlagSet) bool {
	if flag == nil {
		return false
	}

	return flag.NoOptDefVal != ""
}

// isFlagWithAdditionalArg checks whether the given string represents a flag
// which is followed by an argument. For instance, suppose that there is no
// default value for a flag "--flag". Then, given "--flag", this function returns
// true because "--flag" must be followed by an argument. Similar logic applies to
// a short-hand flag using a single dash. (e.g. "-f")
func isFlagWithAdditionalArg(flagStr string, flagSet *flag.FlagSet) bool {
	// --flag arg
	if strings.HasPrefix(flagStr, "--") && !strings.Contains(flagStr, "=") {
		flagKey := flagStr[2:]
		flag := flagSet.Lookup(flagKey)

		return !flagHasDefaultOption(flag, flagSet)
	}

	// -f arg
	if strings.HasPrefix(flagStr, "-") && !strings.Contains(flagStr, "=") && len(flagStr) == 2 {
		flagKey := flagStr[1:]
		flag := flagSet.ShorthandLookup(flagKey)

		return !flagHasDefaultOption(flag, flagSet)
	}

	return false
}

// filterFlags removes flags and their values from the given arguments. It returns
// a filtered slice which contains only non-flag arguments.
func filterFlags(args []string, c *Command) []string {
	ret := []string{}
	flags := c.Flags()

	idx := 0

	for idx < len(args) {
		currentArg := args[idx]

		// "--" signifies the end of command line flags
		if currentArg == "--" {
			ret = append(ret, args[idx+1:]...)
			break
		}

		if isFlagWithAdditionalArg(currentArg, flags) {
			idx = idx + 2
			continue
		}

		if currentArg != "" && !strings.HasPrefix(currentArg, "-") {
			ret = append(ret, currentArg)
		}

		idx = idx + 1
	}

	return ret
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

	cmd, flags := c.findSubCommand(args)
	if cmd == nil {
		// not found. show suggestion
		return nil
	}

	cmd.execute(flags)

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
func (c *Command) findSubCommand(args []string) (*Command, []string) {
	nonFlagArgs := filterFlags(args, c)
	log.Debug("non-flag arguments: %s\n", nonFlagArgs)

	subCommand := nonFlagArgs[0]
	log.Debug("sub-command: '%s'\n", subCommand)

	for _, cmd := range c.commands {
		if cmd.Name == subCommand || cmd.hasAlias(subCommand) {
			return cmd, c.argsWithout(args, subCommand)
		}
	}

	return nil, []string{}
}

// argsWithout removes the fist non-flag occurrence of the given key in the
// given slice of arguments. It does not modify the given slice.
func (c *Command) argsWithout(args []string, key string) []string {
	flags := c.Flags()

	idx := -1
	pos := 0

	for pos < len(args) {
		currentArg := args[pos]

		// "--" signifies the end of command line flags
		if currentArg == "--" {
			break
		}

		if isFlagWithAdditionalArg(currentArg, flags) {
			pos = pos + 2
			continue
		}

		if !strings.HasPrefix(currentArg, "-") {
			if currentArg == key {
				idx = pos
				break
			}
		}

		pos = pos + 1
	}

	if idx == -1 {
		return args
	}

	ret := []string{}
	ret = append(ret, args[:pos]...)
	ret = append(ret, args[pos+1:]...)

	return ret
}

// AddCommand adds the given command as a subcommand.
func (c *Command) AddCommand(cmd *Command) {
	cmd.Parent = c

	c.commands = append(c.commands, cmd)
}
