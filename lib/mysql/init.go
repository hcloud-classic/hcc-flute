package mysql

// Init : Prepare mysql connection
func Init() error {
	err := prepare()
	if err != nil {
		return err
	}

	return nil
}

// End : Close mysql connection
func End() {
	if Db != nil {
		_ = Db.Close()
	}
}
