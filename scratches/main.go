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

// loging in

// func LocalLogin() error {
// 	var username, password string

// 	fmt.Print("Enter your username: ")
// 	fmt.Scanln(&username)
// 	username = strings.TrimSpace(username)

// 	fmt.Print("Enter your password: ")
// 	bytePass, err := term.ReadPassword(uintptr(os.Stdin.Fd()))
// 	if err != nil {
// 		return fmt.Errorf("failed to read password: %w", err)
// 	}

// 	fmt.Println()

// 	password = string(bytePass)

// 	if err := VerifyLocalCredentials(username, password); err != nil {
// 		return fmt.Errorf("login failed: %w", err)
// 	}

// 	if err := SaveSession(Session{
// 		Method:        LocalAuth,
// 		Authenticated: true,
// 	}); err != nil {
// 		return fmt.Errorf("failed to save session: %w", err)
// 	}

// 	fmt.Println("Logged in successfully with local credentials.")
// 	return nil
// }

// init

// dbPath := filepath.Join(dir, "worklogger.db")
// if _, err := os.Stat(dbPath); os.IsNotExist(err) {
// Call your function to create the SQLite DB and tables
// err := SetupDB(dbPath)
// if err != nil {
// 	log.Fatalf("❌ Failed to initialize database: %v", err)
// }
// fmt.Println("🗃️  SQLite DB initialized")
// } else {
// 	fmt.Println("Database already exists, skipping DB setup")
// }

// cmd/log.go
// func runLogCmd() {
// 	db := getDB() // however you open your DB
// 	logs, err := data.LogModel{DB: db}.GetLogs()
// 	if err != nil {
// 		fmt.Println("error fetching logs:", err)
// 		return
// 	}

// 	model := tui.NewLogModel(logs)
// 	p := tea.NewProgram(model)
// 	if err := p.Start(); err != nil {
// 		fmt.Println("error running TUI:", err)
// 	}
// }
