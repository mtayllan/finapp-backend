package main

import (
	"log"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnRecordAfterCreateRequest().Add(handleCreateTransaction(app))
	app.OnRecordAfterUpdateRequest().Add(handleUpdateTransaction(app))
	app.OnRecordAfterDeleteRequest().Add(handleDeleteTransaction(app))

	app.OnRecordAfterUpdateRequest().Add(handleUpdateAccount(app))

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func handleCreateTransaction(app core.App) func(*core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		if e.Record.TableName() != "transactions" {
			return nil
		}

		updateAccountBalance(app, e.Record.GetString("account"))

		return nil
	}
}

func handleUpdateTransaction(app core.App) func(*core.RecordUpdateEvent) error {
	return func(e *core.RecordUpdateEvent) error {
		record := e.Record
		if record.TableName() != "transactions" {
			return nil
		}

		oldRecord := e.Record.OriginalCopy()

		if record.GetString("account") != oldRecord.GetString("account") {
			updateAccountBalance(app, oldRecord.GetString("account"))
		}

		updateAccountBalance(app, record.GetString("account"))

		return nil
	}
}

func handleDeleteTransaction(app core.App) func(*core.RecordDeleteEvent) error {
	return func(e *core.RecordDeleteEvent) error {
		record := e.Record
		if record.TableName() != "transactions" {
			return nil
		}

		updateAccountBalance(app, record.GetString("account"))

		return nil
	}
}

func handleUpdateAccount(app core.App) func(*core.RecordUpdateEvent) error {
	return func(e *core.RecordUpdateEvent) error {
		record := e.Record
		if record.TableName() != "accounts" {
			return nil
		}

		updateAccountBalance(app, record.GetId())

		return nil
	}
}

func updateAccountBalance(app core.App, accountId string) error {
	record, _ := app.Dao().FindRecordById("accounts", accountId)

	var balance float64

	app.Dao().DB().
		Select("COALESCE(SUM(CASE WHEN type = 'inbound' THEN value ELSE -value END), 0)").
		From("transactions").
		Where(dbx.HashExp{"account": accountId}).
		Row(&balance)

	record.Set("current_balance", balance+record.GetFloat("initial_balance"))
	app.Dao().SaveRecord(record)

	return nil
}
