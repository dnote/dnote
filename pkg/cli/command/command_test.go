package command

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestAddCommand(t *testing.T) {
	cmd := Command{Name: "root command"}
	assert.Equal(t, len(cmd.commands), 0, "Commands length mismatch")

	subCommand1 := Command{Name: "foo"}
	cmd.AddCommand(&subCommand1)
	assert.Equal(t, subCommand1.Parent, &cmd, "subCommand1 Parent mismatch")
	assert.Equal(t, len(cmd.commands), 1, "Commands length mismatch")
	assert.Equal(t, cmd.commands[0], &subCommand1, "commands[0] mismatch")

	subCommand2 := Command{Name: "bar"}
	cmd.AddCommand(&subCommand2)
	assert.Equal(t, len(cmd.commands), 2, "Commands length mismatch")
	assert.Equal(t, subCommand2.Parent, &cmd, "subCommand2 Parent mismatch")
	assert.Equal(t, cmd.commands[0], &subCommand1, "commands[0] mismatch")
	assert.Equal(t, cmd.commands[1], &subCommand2, "commands[1] mismatch")
}

func TestHasAlias(t *testing.T) {
	cmd := Command{
		Name:    "foo",
		Aliases: []string{"f", "bar"},
	}

	assert.Equal(t, cmd.hasAlias("f"), true, "Command should have 'f' alias")
	assert.Equal(t, cmd.hasAlias("F"), false, "Command should not have 'F' alias")
	assert.Equal(t, cmd.hasAlias("bar"), true, "Command should have 'bar' alias")
	assert.Equal(t, cmd.hasAlias("BAR"), false, "Command should have 'BAR' alias")
	assert.Equal(t, cmd.hasAlias("baz"), false, "Command should not have 'baz' alias")
	assert.Equal(t, cmd.hasAlias(""), false, "Command should not have an empty alias")
}

func TestHasAlias_withoutAlias(t *testing.T) {
	cmd := Command{
		Name: "foo",
	}

	assert.Equal(t, cmd.hasAlias("f"), false, "Command should not have any alias")
	assert.Equal(t, cmd.hasAlias(""), false, "Command should not have any alias")
}

func TestCommandRoot(t *testing.T) {
	subCommand2 := Command{
		Name: "baz",
	}
	subCommand1 := Command{
		Name: "bar",
		commands: []*Command{
			&subCommand2,
		},
	}
	cmd := Command{
		Name: "foo",
		commands: []*Command{
			&subCommand1,
		},
	}

	subCommand1.Parent = &cmd
	subCommand2.Parent = &subCommand1

	assert.Equal(t, cmd.Root(), &cmd, "Command should already be a root")
	assert.Equal(t, subCommand1.Root(), &cmd, "subCommand1 root mismatch")
	assert.Equal(t, subCommand2.Root(), &cmd, "subCommand2 root mismatch")
}

func TestFindSubcommand(t *testing.T) {
	subCommand1 := Command{
		Name:    "bar",
		Aliases: []string{"quz"},
	}
	subCommand2 := Command{
		Name: "baz",
	}
	cmd := Command{
		Name: "foo",
		commands: []*Command{
			&subCommand1,
			&subCommand2,
		},
	}

	assert.Equal(t, cmd.findSubCommand("bar"), &subCommand1, "Subcommand 'bar' mismatch")
	assert.Equal(t, cmd.findSubCommand("baz"), &subCommand2, "Subcommand 'baz' mismatch")

	// Should match an alias
	assert.Equal(t, cmd.findSubCommand("quz"), &subCommand1, "Subcommand 'quz' mismatch")
	// Should not match if not exists
	assert.Equal(t, cmd.findSubCommand("qux"), (*Command)(nil), "Subcommand 'qux' mismatch")
	// Should not match itself
	assert.Equal(t, subCommand1.findSubCommand("bar"), (*Command)(nil), "Subcommand 'bar' mismatch")
}

func executeCommand(cmd *Command, args ...string) error {
	cmd.setArgs(args)
	return cmd.Execute()
}
