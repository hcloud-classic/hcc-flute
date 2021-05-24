package ipmi

import (
	"errors"
	"hcc/flute/daoext"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"strconv"
	"strings"
	"time"
)

func getTodayByNumString() string {
	currentTime := time.Now()
	return currentTime.Format("060102")
}

func checkIfTodayNodeUptimeExist(nodeUUID string) (isExist bool, nodeUptimeMs int64) {
	var uptimeMs int64

	sql := "select node_uuid, uptime_ms from node_uptime where node_uuid = ? and day = ?"
	row := mysql.Db.QueryRow(sql, nodeUUID, getTodayByNumString())
	err := mysql.QueryRowScan(row, &nodeUUID, &uptimeMs)
	if err != nil {
		return false, 0
	}

	return true, uptimeMs
}

func insertTodayNodeUptime(launchedTime time.Time, nodeUUID string) error {
	node, errCode, errText := daoext.ReadNode(nodeUUID)
	if errCode != 0 {
		return errors.New(errText)
	}

	var uptimeMs int64
	if strings.ToLower(node.Status) == "on" {
		now := time.Now()
		timeDiff := now.Sub(launchedTime)
		uptimeMs = timeDiff.Milliseconds()
	} else {
		uptimeMs = 0
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("insertTodayNodeUptime(): Inserting new node uptime_ms (node_uuid=" + nodeUUID +
			", uptime_ms=" + strconv.Itoa(int(uptimeMs)) + ", day=" + getTodayByNumString() + ")")
	}

	sql := "insert into node_uptime(node_uuid, group_id, uptime_ms, day) values (?, ?, ?, ?)"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "insertTodayNodeUptime(): " + err.Error()
		logger.Logger.Println(errStr)
		return errors.New(errStr)
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(nodeUUID, node.GroupID, uptimeMs, getTodayByNumString())
	if err != nil {
		errStr := "insertTodayNodeUptime(): " + err.Error()
		logger.Logger.Println(errStr)
		return errors.New(errStr)
	}

	return nil
}

func updateTodayNodeUptime(launchedTime time.Time, nodeUUID string) error {
	exist, uptimeMS := checkIfTodayNodeUptimeExist(nodeUUID)
	if !exist {
		err := insertTodayNodeUptime(launchedTime, nodeUUID)
		if err != nil {
			return err
		}

		return nil
	}

	node, errCode, errText := daoext.ReadNode(nodeUUID)
	if errCode != 0 {
		return errors.New(errText)
	}

	if strings.ToLower(node.Status) == "off" {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("updateTodayNodeUptime(): Node is turned off (node_uuid=" + nodeUUID +
				", uptime_ms=" + strconv.Itoa(int(uptimeMS)) + ", day=" + getTodayByNumString() + ")")
		}

		return nil
	}

	now := time.Now()
	timeDiff := now.Sub(launchedTime)
	newUptimeMs := uptimeMS + timeDiff.Milliseconds()

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("updateTodayNodeUptime(): Updating node uptime_ms (node_uuid=" + nodeUUID +
			", new_uptime_ms=" + strconv.Itoa(int(newUptimeMs)) + ", before_uptime_ms=" + strconv.Itoa(int(uptimeMS)) +
			", day=" + getTodayByNumString() + ")")
	}

	sql := "update node_uptime set uptime_ms = ? where node_uuid = ? and day = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "updateTodayNodeUptime(): " + err.Error()
		logger.Logger.Println(errStr)
		return errors.New(errStr)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(newUptimeMs, nodeUUID, getTodayByNumString())
	if err != nil {
		errStr := "updateTodayNodeUptime(): " + err.Error()
		logger.Logger.Println(errStr)
		return errors.New(errStr)
	}

	return nil
}
