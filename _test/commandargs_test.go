package _test

import (
    "testing"

    "github.com/alexj212/webmgmt"
)

func TestCommandArgsEmptyParse(t *testing.T) {

    cmdArgs, _ := webmgmt.NewCommandArgs("", nil)

    if cmdArgs != nil {
        t.Errorf("Command Args should error on zero len string")
    }
}

func TestCommandArgsEmptySpacesParse(t *testing.T) {

    cmdArgs, _ := webmgmt.NewCommandArgs("   ", nil)

    if cmdArgs != nil {
        t.Errorf("Command Args should error on zero len string")
    }
}

func TestCommandArgsEmptyTabParse(t *testing.T) {

    cmdArgs, _ := webmgmt.NewCommandArgs("\t", nil)

    if cmdArgs != nil {
        t.Errorf("Command Args should error on zero len string")
    }
}

func TestCommandArgsEmptySpaceTabParse(t *testing.T) {

    cmdArgs, _ := webmgmt.NewCommandArgs(" \t", nil)

    if cmdArgs != nil {
        t.Errorf("Command Args should error on zero len string")
    }
}

func TestCommandArgsSingleCmdParse(t *testing.T) {

    cmdArgs, err := webmgmt.NewCommandArgs("hello", nil)

    if err != nil {
        t.Errorf("Command Args should have parsed")
    }

    if cmdArgs.CmdName != "hello" {
        t.Errorf("Command should have parsed cmd name to 'hello'")
    }

    if len(cmdArgs.Args) != 0 {
        t.Errorf("Command args should be 0")
    }

}

func TestCommandArgsCmdParse(t *testing.T) {

    cmdArgs, err := webmgmt.NewCommandArgs("hello world", nil)

    if err != nil {
        t.Errorf("Command Args should have parsed")
    }

    if cmdArgs.CmdName != "hello" {
        t.Errorf("Command should have parsed cmd name to 'hello'")
    }

    if len(cmdArgs.Args) != 1 {
        t.Errorf("Command args should be 1")
    }

    if cmdArgs.Args[0] != "world" {
        t.Errorf("Command Args[0] should be 'world'")
    }
}

// topic Hello World Stinky
// topic "Hello World" Stinky
// group 156 topic Hello World Stinky

func TestCommandArgsTopicParse(t *testing.T) {

    cmdArgs, err := webmgmt.NewCommandArgs("topic hello world stinky", nil)

    if err != nil {
        t.Errorf("Command Args should have parsed")
    }

    if cmdArgs.CmdName != "topic" {
        t.Errorf("Command should have parsed cmd name to 'hello'")
    }

    t.Logf("cmdArgs.String(): %v\n", cmdArgs.String())
    t.Logf("Args: %v\n", cmdArgs.Args)
    t.Logf("PealOff(0): %v\n", cmdArgs.PealOff(0))
    t.Logf("PealOff(1): %v\n", cmdArgs.PealOff(1))
    t.Logf("PealOff(2): %v\n", cmdArgs.PealOff(2))
    t.Logf("PealOff(3): %v\n", cmdArgs.PealOff(3))
    t.Logf("PealOff(4): %v\n", cmdArgs.PealOff(4))

}

func TestCommandArgsTopicQuotedParse(t *testing.T) {

    cmdArgs, err := webmgmt.NewCommandArgs("topic \"hello world\" stinky", nil)

    if err != nil {
        t.Errorf("Command Args should have parsed")
    }

    if cmdArgs.CmdName != "topic" {
        t.Errorf("Command should have parsed cmd name to 'hello'")
    }

    t.Logf("cmdArgs.String(): %v\n", cmdArgs.String())
    t.Logf("Args: %v\n", cmdArgs.Args)
    t.Logf("PealOff(0): %v\n", cmdArgs.PealOff(0))
    t.Logf("PealOff(1): %v\n", cmdArgs.PealOff(1))
    t.Logf("PealOff(2): %v\n", cmdArgs.PealOff(2))
    t.Logf("PealOff(3): %v\n", cmdArgs.PealOff(3))
    t.Logf("PealOff(4): %v\n", cmdArgs.PealOff(4))

}

func TestCommandArgsCmdShiftParse(t *testing.T) {

    cmdArgs, _ := webmgmt.NewCommandArgs("entry session mQn5owrjKEdKZFB7BgeWA6JB9ypOvL", nil)

    t.Logf("cmdArgs: %v\n", cmdArgs.String())
    t.Logf("cmdArgs: %v\n", cmdArgs.Args)
    newCmdLine, _ := cmdArgs.Shift()

    t.Logf("newCmdLine: %v\n", newCmdLine.String())
    t.Logf("newCmdLine: %v\n", newCmdLine.Args)

}

//
// func TestCommandArgsCmdQuotedArgParse(t *testing.T) {
//
//     cmdArgs, _ := webmgmt.NewCommandArgs("hello \"world, suckers\"", nil)
//
//     if cmdArgs != nil {
//         t.Errorf("Command Args should error on zero len string")
//     }
// }
// func TestCommandArgsCmdMixedArgParse(t *testing.T) {
//
//     cmdArgs, _ := webmgmt.NewCommandArgs("hello \"world, suckers\" alex", nil)
//
//     if cmdArgs != nil {
//         t.Errorf("Command Args should error on zero len string")
//     }
// }
//
//
//
// func TestCommandArgsCmdShiftMultiQuotedParse(t *testing.T) {
//
//     cmdArgs, _ := webmgmt.NewCommandArgs("entry debug a \"quoted variable\" ", nil)
//
//     newCmdLine, _ := cmdArgs.Shift()
//
//     if newCmdLine != nil {
//         t.Errorf("Command Args should error on zero len string")
//     }
// }
//
//
// func TestCommandArgsCmdMixedArgsShiftParse(t *testing.T) {
//
//     cmdArgs, _ := webmgmt.NewCommandArgs("hello \"world, suckers\" alex", nil)
//
//     newCmdLine, _ := cmdArgs.Shift()
//
//     if newCmdLine != nil {
//         t.Errorf("Command Args should error on zero len string")
//     }
// }
