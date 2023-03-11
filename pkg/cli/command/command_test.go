package command

import (
	"fmt"
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

	subcommandResult1, subcommandArgResult1 := cmd.findSubCommand([]string{"bar"})
	assert.Equal(t, subcommandResult1, &subCommand1, "Subcommand 'bar' mismatch")
	assert.DeepEqual(t, subcommandArgResult1, []string{}, "Subcommand arg for 'bar' mismatch")

	subcommandResult2, subcommandArgResult2 := cmd.findSubCommand([]string{"baz", "echo"})
	assert.Equal(t, subcommandResult2, &subCommand2, "Subcommand 'baz' mismatch")
	assert.DeepEqual(t, subcommandArgResult2, []string{"echo"}, "Subcommand arg for 'baz' mismatch")

	// Should match an alias
	subcommandResult3, subcommandArgResult3 := cmd.findSubCommand([]string{"quz"})
	assert.Equal(t, subcommandResult3, &subCommand1, "Subcommand 'quz' mismatch")
	assert.DeepEqual(t, subcommandArgResult3, []string{}, "Subcommand arg for 'quz' mismatch")

	// Should not match if not exists
	subcommandResult4, subcommandArgResult4 := cmd.findSubCommand([]string{"qux"})
	assert.Equal(t, subcommandResult4, (*Command)(nil), "Subcommand 'qux' mismatch")
	assert.DeepEqual(t, subcommandArgResult4, []string{}, "Subcommand arg for 'qux' mismatch")

	// Should not match itself
	subcommandResult5, subcommandArgResult5 := subCommand1.findSubCommand([]string{"bar"})
	assert.Equal(t, subcommandResult5, (*Command)(nil), "Subcommand 'bar' should not exist on 'bar'")
	assert.DeepEqual(t, subcommandArgResult5, []string{}, "Subcommand arg for 'bar' should be empty when there is no match")
}

func executeCommand(root *Command, args ...string) error {
	return root.Execute()
}

func TestFilterFlags(t *testing.T) {
	testCases := []struct {
		arguments []string
		expected  []string
	}{
		{
			arguments: []string{"a", "b", "c"},
			expected:  []string{"a", "b", "c"},
		},
		{
			arguments: []string{"-b", "cmd"},
			expected:  []string{"cmd"},
		},
		{
			arguments: []string{"foo", "-b", "--str"},
			expected:  []string{"foo"},
		},
		{
			arguments: []string{"foo", "-b", "baz"},
			expected:  []string{"foo", "baz"},
		},
		{
			arguments: []string{"foo", "-i", "123"},
			expected:  []string{"foo"},
		},
		{
			arguments: []string{"-b", "bar", "-i", "1", "a"},
			expected:  []string{"bar", "a"},
		},
		{
			arguments: []string{"a", "-b", "bar", "-i", "1"},
			expected:  []string{"a", "bar"},
		},
		{
			arguments: []string{"--foo", "b", "-baz"},
			expected:  []string{},
		},
		{
			arguments: []string{"--s=hello", "foo"},
			expected:  []string{"foo"},
		},
		{
			arguments: []string{"-shello", "foo"},
			expected:  []string{"foo"},
		},
		{
			arguments: []string{"-shello", "foo", "-i1"},
			expected:  []string{"foo"},
		},
		{
			arguments: []string{"-si", "foo"},
			expected:  []string{"foo"},
		},
	}

	cmd := Command{}
	cmd.Flags().StringP("str", "s", "", "")
	cmd.Flags().IntP("int", "i", 1, "")
	cmd.Flags().BoolP("bool", "b", false, "")
	// Default 'bool' to "true"
	cmd.Flags().Lookup("bool").NoOptDefVal = "true"

	for idx, testCase := range testCases {
		got := filterFlags(testCase.arguments, &cmd)

		assert.DeepEqual(t, got, testCase.expected, fmt.Sprintf("result mismatch for test case %d", idx))
	}
}
