package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goproject/src/collection"
	"goproject/src/domain"
	"goproject/src/optional"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Factory[T domain.Entity] struct {
	Create func() T
}

type Manager[T domain.Entity] struct {
	Database *sql.DB
	Table    string
	Factory  Factory[T]
	context  sync.Map
}

func NewManager[T domain.Entity](
	database *sql.DB,
	factory Factory[T],
) *Manager[T] {
	return &Manager[T]{
		Database: database,
		Table:    factory.Create().Table(),
		Factory:  factory,
		context:  sync.Map{},
	}
}

func (m *Manager[T]) Persist(entity T) {
	tempID := uuid.New().String()
	m.context.Store(tempID, &entity)
}

func (m *Manager[T]) Flush() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := m.Database.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		log.Printf("Failed to begin transaction: %v \n\n", err)
		return err
	}

	var flushErr error = nil

	m.context.Range(func(key any, value any) bool {
		entity := value.(*T)

		log.Println(*entity)
		log.Println()

		var o *optional.Optional[T]

		if (*entity).GetId() == 0 {
			o = m.doInsert(entity)
		} else {
			o = m.doUpdate(entity)
		}

		if o.IsEmpty() {
			flushErr = errors.New("failed flush")
			return false
		}

		return true
	})

	log.Printf("flushErr: %v \n\n", flushErr)

	if flushErr != nil {
		err := tx.Commit()
		if err != nil {
			log.Printf("Failed commit: %s \n\n", err)
			return err
		}
	} else {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Failed rollback: %s \n\n", err)
			return err
		}
	}

	return flushErr
}

func (m *Manager[T]) Find(id int64) *optional.Optional[T] {
	if cachedEntity, ok := m.context.Load(id); ok {
		return optional.Of(*(cachedEntity.(*T)))
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", m.Table)
	row := m.Database.QueryRow(query, id)

	entity := m.Factory.Create()
	err := entity.FlushRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return optional.OfNil[T]()
		}

		log.Printf("Failed to flush row: %v \n\n", err)
		return optional.OfNil[T]()
	}

	m.context.Store(id, &entity)
	return optional.Of(entity)
}

func (m *Manager[T]) Count() int64 {
	query := fmt.Sprintf("SELECT count(*) FROM %s", m.Table)
	row := m.Database.QueryRow(query)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		log.Printf("Failed to execute count query: %v \n\n", err)
		return 0
	}

	return count
}

func (m *Manager[T]) FindAll() *collection.Collection[T] {
	rows, err := m.Database.Query(fmt.Sprintf("SELECT * FROM %s", m.Table))
	if err != nil {
		log.Fatalf("Query failed: %v \n\n", err)
	}

	defer rows.Close()

	c := &collection.Collection[T]{}
	for rows.Next() {
		entity := m.Factory.Create()
		_ = entity.FlushRows(rows)
		c.Append(entity)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error occurred during rows iteration: %v \n\n", err)
	}

	return c
}

func (m *Manager[T]) Delete(entity T) bool {
	var err error
	table := m.Table

	id := entity.GetId()

	if id != 0 {
		query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
		_, err = m.Database.Exec(query, id)
		m.context.Delete(id)
	}

	if err != nil {
		return false
	}

	return true
}

func (m *Manager[T]) doUpdate(entity *T) *optional.Optional[T] {
	e := *entity
	mappings := e.Mappings()

	var (
		updateFields []string
		values       []interface{}
	)

	entityValue := reflect.ValueOf(e).Elem()

	for dbField, structField := range mappings {
		field := entityValue.FieldByName(structField)
		if field.IsValid() && field.CanInterface() {
			if dbField == "id" {
				continue
			}

			updateFields = append(updateFields, fmt.Sprintf("%s = $%d", dbField, len(values)+1))
			values = append(values, field.Interface())
		}
	}

	if len(updateFields) == 0 {
		log.Println("No fields to update")
		return optional.OfNil[T]()
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d RETURNING id",
		e.Table(),
		strings.Join(updateFields, ", "),
		len(values)+1,
	)

	log.Printf("Query in doUpdate method: %s \n\n", query)

	values = append(values, e.GetId())
	var id int64
	err := m.Database.QueryRow(query, values...).Scan(&id)
	if err != nil {
		log.Printf("Failed to execute query: %v \n\n", err)
		return optional.OfNil[T]()
	}

	return optional.Of(e)
}

func (m *Manager[T]) doInsert(entity *T) *optional.Optional[T] {
	e := *entity
	mappings := e.Mappings()

	var (
		fields       []string
		placeholders []string
		updateFields []string
		values       []interface{}
	)

	entityValue := reflect.ValueOf(e).Elem()

	for dbField, structField := range mappings {
		field := entityValue.FieldByName(structField)
		if field.IsValid() && field.CanInterface() {
			if dbField == "id" {
				continue
			}

			fields = append(fields, dbField)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))
			values = append(values, field.Interface())
			updateFields = append(updateFields, fmt.Sprintf("%s = EXCLUDED.%s", dbField, dbField))
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (id) DO UPDATE SET %s RETURNING id",
		e.Table(),
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updateFields, ", "),
	)

	log.Printf("Query in doInsert method: %s \n\n", query)

	var id int64
	err := m.Database.QueryRow(query, values...).Scan(&id)
	if err != nil {
		log.Printf("Failed to execute query: %v \n\n", err)
		return optional.OfNil[T]()
	}

	e.SetId(id)
	return optional.Of(e)
}

func (m *Manager[T]) Exists(criteria string) bool {
	query := fmt.Sprintf(
		"SELECT exists (SELECT TRUE FROM %s WHERE %s)",
		m.Table,
		criteria,
	)

	row := m.Database.QueryRow(query)

	result := false
	_ = row.Scan(&result)

	return result
}

func formatValue(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return fmt.Sprintf("'%s'", value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", value.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", value.Float())
	case reflect.Bool:
		if value.Bool() {
			return "TRUE"
		}
		return "FALSE"
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			return fmt.Sprintf("'%s'", string(value.Bytes()))
		}
	default:
	}

	if value.Type() == reflect.TypeOf(time.Time{}) {
		return fmt.Sprintf("'%s'", value.Interface().(time.Time).Format("2006-01-02 15:04:05"))
	}

	return fmt.Sprintf("'%v'", value.Interface())
}
