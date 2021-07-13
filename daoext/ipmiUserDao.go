package daoext

import (
	"errors"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/lib/passwordEncrypt"
)

func checkIfIPMIUserExist(bmcIP string) bool {
	sql := "select bmc_ip from ipmi_user where bmc_ip = ?"
	row := mysql.Db.QueryRow(sql, bmcIP)
	err := mysql.QueryRowScan(row, &bmcIP)
	if err != nil {
		return false
	}

	return true
}

// AddIPMIUser : Add IPMI user info of the node
func AddIPMIUser(bmcIP string, ipmiUserID string, ipmiUserPassword string) error {
	ipmiUserIDOk := len(ipmiUserID) != 0
	ipmiUserPasswordOk := len(ipmiUserPassword) != 0

	if !ipmiUserIDOk || !ipmiUserPasswordOk {
		return errors.New("AddIPMIUser(): ID or password is empty")
	}

	if checkIfIPMIUserExist(bmcIP) {
		logger.Logger.Println("AddIPMIUser(): IPMI user info is already exist for bmcIP=" + bmcIP)
		return nil
	}

	sql := "insert into ipmi_user(bmc_ip, id, password) values (?, ?, ?)"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		logger.Logger.Println("AddIPMIUser(): " + err.Error())
		return err
	}

	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(bmcIP, ipmiUserID, passwordEncrypt.EncryptPassword(ipmiUserPassword))
	if err != nil {
		logger.Logger.Println("AddIPMIUser(): " + err.Error())
		return err
	}

	return nil
}

// UpdateIPMIUser : Update IPMI user info of the node
func UpdateIPMIUser(bmcIP string, ipmiUserID string, ipmiUserPassword string) error {
	if !checkIfIPMIUserExist(bmcIP) {
		return AddIPMIUser(bmcIP, ipmiUserID, ipmiUserPassword)
	}

	ipmiUserIDOk := len(ipmiUserID) != 0
	ipmiUserPasswordOk := len(ipmiUserPassword) != 0

	var sql string

	if ipmiUserIDOk && ipmiUserPasswordOk {
		sql = "update ipmi_user set id = ?, password = ? where bmc_ip = ?"
	} else if ipmiUserIDOk {
		sql = "update ipmi_user set id = ? where bmc_ip = ?"
	} else if ipmiUserPasswordOk {
		sql = "update ipmi_user set password = ? where bmc_ip = ?"
	} else {
		return errors.New("UpdateIPMIUser(): ID or password is empty")
	}

	stmt, err := mysql.Prepare(sql)
	if err != nil {
		logger.Logger.Println("UpdateIPMIUser(): " + err.Error())
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var err2 error
	if ipmiUserIDOk && ipmiUserPasswordOk {
		_, err2 = stmt.Exec(ipmiUserID, passwordEncrypt.EncryptPassword(ipmiUserPassword), bmcIP)
	} else if ipmiUserIDOk {
		_, err2 = stmt.Exec(ipmiUserID, bmcIP)
	} else if ipmiUserPasswordOk {
		_, err2 = stmt.Exec(passwordEncrypt.EncryptPassword(ipmiUserPassword), bmcIP)
	}
	if err2 != nil {
		logger.Logger.Println("UpdateIPMIUser(): " + err2.Error())
		return err2
	}

	defer func() {
		_ = stmt.Close()
	}()

	return nil
}

// DeleteIPMIUser : Delete IPMI user info of the node
func DeleteIPMIUser(bmcIP string) error {
	sql := "delete from ipmi_user where bmc_ip = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		logger.Logger.Println("DeleteIPMIUser(): " + err.Error())
		return err
	}

	defer func() {
		_ = stmt.Close()
	}()

	_, err2 := stmt.Exec(bmcIP)
	if err2 != nil {
		logger.Logger.Println("DeleteIPMIUser(): " + err2.Error())
		return err
	}

	return nil
}
