package database

func AddToken(token string) error {
	conn, err := InitDB()
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec("INSERT INTO config (token) VALUES (?)", token)
	if err != nil {
		panic(err)
	}
	return nil
}

func GetToken() (string, error) {
	conn, err := InitDB()
	if err != nil {
		panic(err)
	}
	var token string
	err = conn.QueryRow("SELECT token FROM config").Scan(&token)
	if err != nil {
		panic(err)
	}
	return token, nil
}

func UpdateToken(token string) error {
	conn, err := InitDB()
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec("UPDATE config SET token =?", token)
	if err != nil {
		panic(err)
	}
	return nil
}

func DeleteToken() error {
	conn, err := InitDB()
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec("DELETE FROM config")
	if err != nil {
		panic(err)
	}
	return nil
}
