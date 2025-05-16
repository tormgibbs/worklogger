package scratches


// cmd/login.go

// p := tea.NewProgram(tui.NewLoginModel())
// m, err := p.Run()
// if err != nil {
// 	cmd.PrintErrln("Error running Login TUI", err)
// 	return
// }

// lm := m.(tui.LoginModel)

// switch lm.Selected {
// case "Github":
// 	fmt.Println("Logging in with GitHub...")
// case "Local Authentication":
// 	fmt.Println("Logging in locally...")
// default:
// 	fmt.Println("No valid login method selected")
// }


// tui/login.go

// func RunLoginUI() (string, error) {
// 	m := NewLoginModel()
// 	p := tea.NewProgram(m)

// 	finalModel, err := p.Run()
// 	if err != nil {
// 		return "", err
// 	}

// 	model := finalModel.(LoginModel)
// 	return model.Selected, nil
// }

