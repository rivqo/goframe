package db

import (
	// "errors"
	"reflect"
	// "strings"
	"time"
)

// Entity is the base struct for all ORM models
type Entity struct {
	ID        uint      `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Repository provides a generic way to interact with entities
type Repository[T any] struct {
	db *Database
}

// NewRepository creates a new repository for the given entity type
func NewRepository[T any](db *Database) *Repository[T] {
	return &Repository[T]{db: db}
}

// Create creates a new entity
func (r *Repository[T]) Create(entity *T) error {
	// Set created_at and updated_at
	setTimestamps(entity)
	
	return r.db.Create(entity)
}

// FindByID finds an entity by ID
func (r *Repository[T]) FindByID(id uint, dest *T) error {
	return r.db.Table(dest).Where("id = ?", id).First(dest)
}

func (r Repository[T]) FindByString(column string, dest *T, value string) error {
	return r.db.Table(dest).Where(column+" = ?", value).First(dest)
}

func (r *Repository[T]) FindByIDOrFail(id uint, dest *T) error {
	err := r.FindByID(id, dest)
	if err != nil {
		return err
	}

	return nil
}

// FindAll finds all entities matching the query
func (r *Repository[T]) FindAll(dest *[]T, conditions ...interface{}) error {
	query := r.db.Table(*new(T))
	
	if len(conditions) > 0 {
		if condition, ok := conditions[0].(string); ok {
			query = query.Where(condition, conditions[1:]...)
		}
	}
	
	return query.Find(dest)
}

// Update updates an entity
func (r *Repository[T]) Update(entity *T) error {
	// Set updated_at
	setUpdatedAt(entity)
	
	values := toMap(entity)
	delete(values, "id")       // Don't update ID
	delete(values, "created_at") // Don't update created_at
	
	return r.db.Table(entity).Where("id = ?", getID(entity)).Update(values)
}

// Delete deletes an entity
func (r *Repository[T]) Delete(entity *T) error {
	return r.db.Table(entity).Where("id = ?", getID(entity)).Delete()
}

// setTimestamps sets the created_at and updated_at fields
func setTimestamps(entity interface{}) {
	now := time.Now()
	
	val := reflect.ValueOf(entity).Elem()
	
	if createdField := val.FieldByName("CreatedAt"); createdField.IsValid() {
		createdField.Set(reflect.ValueOf(now))
	}
	
	if updatedField := val.FieldByName("UpdatedAt"); updatedField.IsValid() {
		updatedField.Set(reflect.ValueOf(now))
	}
}

// setUpdatedAt sets the updated_at field
func setUpdatedAt(entity interface{}) {
	now := time.Now()
	
	val := reflect.ValueOf(entity).Elem()
	
	if updatedField := val.FieldByName("UpdatedAt"); updatedField.IsValid() {
		updatedField.Set(reflect.ValueOf(now))
	}
}

// getID gets the ID field value
func getID(entity interface{}) uint {
	val := reflect.ValueOf(entity).Elem()
	
	if idField := val.FieldByName("ID"); idField.IsValid() {
		return uint(idField.Uint())
	}
	
	return 0
}

// toMap converts a struct to a map
func toMap(entity interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(entity).Elem()
	typ := val.Type()
	
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("db")
		
		if tag != "" && tag != "-" {
			result[tag] = val.Field(i).Interface()
		}
	}
	
	return result
}

