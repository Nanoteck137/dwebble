package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
)

type Player struct {
	Id string `db:"id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

type Queue struct {
	Id       string `db:"id"`
	PlayerId string `db:"player_id"`
	UserId   string `db:"user_id"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

type DefaultQueue struct {
	PlayerId string `db:"player_id"`
	UserId   string `db:"user_id"`
	QueueId  string `db:"queue_id"`
}

func PlayerQuery() *goqu.SelectDataset {
	query := dialect.From("players").
		Select(
			"players.id",

			"players.created",
			"players.updated",
		).
		Prepared(true)

	return query
}

func QueueQuery() *goqu.SelectDataset {
	query := dialect.From("queues").
		Select(
			"queues.id",
			"queues.player_id",
			"queues.user_id",

			"queues.name",

			"queues.created",
			"queues.updated",
		).
		Prepared(true)

	return query
}

func DefaultQueueQuery() *goqu.SelectDataset {
	query := dialect.From("default_queues").
		Select(
			"default_queues.player_id",
			"default_queues.user_id",
			"default_queues.queue_id",
		).
		Prepared(true)

	return query
}

func (db *Database) GetPlayerById(ctx context.Context, id string) (Player, error) {
	query := PlayerQuery().
		Where(goqu.I("players.id").Eq(id))

	var item Player
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Player{}, ErrItemNotFound
		}

		return Player{}, err
	}

	return item, nil
}

type CreatePlayerParams struct {
	Id string

	Created int64
	Updated int64
}

func (db *Database) CreatePlayer(ctx context.Context, params CreatePlayerParams) error {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	if params.Id == "" {
		return errors.New("id cannot be empty")
	}

	query := dialect.Insert("players").
		Rows(goqu.Record{
			"id": params.Id,

			"created": created,
			"updated": updated,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetQueue(ctx context.Context, id string) (Queue, error) {
	query := QueueQuery().
		Where(goqu.I("queues.id").Eq(id))

	var item Queue
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Queue{}, ErrItemNotFound
		}

		return Queue{}, err
	}

	return item, nil
}

type CreateQueueParams struct {
	Id       string
	PlayerId string
	UserId   string

	Name string

	Created int64
	Updated int64
}

func (db *Database) CreateQueue(ctx context.Context, params CreateQueueParams) (Queue, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		// TODO(patrik): Create: utils.CreateQueueId
		id = utils.CreateId()
	}

	query := dialect.Insert("queues").
		Rows(goqu.Record{
			"id":        id,
			"player_id": params.PlayerId,
			"user_id":   params.UserId,

			"name": params.Name,

			"created": created,
			"updated": updated,
		}).
		Returning(
			"queues.id",
			"queues.player_id",
			"queues.user_id",

			"queues.name",

			"queues.created",
			"queues.updated",
		)

	var res Queue
	err := db.Get(&res, query)
	if err != nil {
		return Queue{}, err
	}

	return res, nil
}

func (db *Database) GetDefaultQueue(ctx context.Context, playerId, userId string) (DefaultQueue, error) {
	query := DefaultQueueQuery().
		Where(
			goqu.I("default_queues.player_id").Eq(playerId),
			goqu.I("default_queues.user_id").Eq(userId),
		)

	var item DefaultQueue
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DefaultQueue{}, ErrItemNotFound
		}

		return DefaultQueue{}, err
	}

	return item, nil
}

type CreateDefaultQueueParams struct {
	PlayerId string
	UserId   string
	QueueId  string
}

func (db *Database) CreateDefaultQueue(ctx context.Context, params CreateDefaultQueueParams) error {
	query := dialect.Insert("default_queues").
		Rows(goqu.Record{
			"player_id": params.PlayerId,
			"user_id":   params.UserId,
			"queue_id":  params.QueueId,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
