package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/engmtcdrm/minno/utils/env"
	pp "github.com/engmtcdrm/minno/utils/prettyprint"
	"github.com/spf13/cobra"
)

type Spice int

const (
	Mild Spice = iota + 1
	Medium
	Hot
)

func (s Spice) String() string {
	switch s {
	case Mild:
		return "Mild "
	case Medium:
		return "Medium-Spicy "
	case Hot:
		return "Spicy-Hot "
	default:
		return ""
	}
}

type Order struct {
	Burger       Burger
	Side         string
	Name         string
	Instructions string
	Discount     bool
}

type Burger struct {
	Type     string
	Toppings []string
	Spice    Spice
}

func init() {
	rootCmd.AddCommand(initialize())
}

func initialize() *cobra.Command {
	init := &cobra.Command{
		Use:     "initialize",
		Short:   "init the rkl cfg.",
		Long:    "init provision the rkl configuration file.",
		Example: "rkl init",
		Aliases: []string{"i", "init"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// p := tea.NewProgram(tui.InitialModel())

			// if _, err := p.Run(); err != nil {
			// 	return err
			// }
			var burger Burger
			var order = Order{Burger: burger}
			order.Discount = true

			dir, _ := env.AppHomeDir()

			form := huh.NewForm(
				huh.NewGroup(
					// huh.NewNote().
					// 	Title("Welcome to "+env.AppNm+"!").
					// 	Description("This is the first time you're running rkl. Let's get you set up."),
					huh.NewSelect[string]().
						Options(
							huh.NewOption("Username/Password", "username-password"),
							huh.NewOption("Token", "token"),
							huh.NewOption("Client Credentials", "client-credentials"),
						).
						Title("What kind of credential(s) would you like to create?").
						Description("At Charm we truly have a burger for everyone.").
						Value(&order.Burger.Type),
					huh.NewInput().
						Title("What's your name?").
						Value(&order.Name).
						// EchoMode(huh.EchoModeNone).),
						// Prompt(" "),
						Inline(true),
					huh.NewConfirm().
						Title("Would you like 15% off?").
						Value(&order.Discount).
						Affirmative("Yes!").
						Negative("No.").
						Inline(true),
					huh.NewInput().
						Title("Your name").
						Value(&order.Name),
					huh.NewFilePicker().
						Description("Select a file").
						CurrentDirectory(dir).
						ShowPermissions(false).
						DirAllowed(false).
						AllowedTypes([]string{"cred"}),
				),
			)

			// println(order.Burger.Type)
			// println(order.Discount)

			err := form.
				WithTheme(pp.ThemeMinno()).
				Run()
			if err != nil {
				if err.Error() == "user aborted" {
					fmt.Println("User aborted")
					os.Exit(0)
				} else {
					return err
				}
			}

			return nil
		},
	}
	return init
}
