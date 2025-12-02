package database

import (
	"context"
	"fmt"

	"github.com/braunkc/todo-db/config"
	"github.com/braunkc/todo-db/internal/application/repository"
	"github.com/braunkc/todo-db/internal/domain/entities"
	valueobjects "github.com/braunkc/todo-db/internal/domain/value_objects/query"
	"github.com/braunkc/todo-db/internal/infra/database/postgres/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type databaseRepository struct {
	db     *gorm.DB
	mapper Mapper
}

func NewDatabaseService(cfg *config.Config, mapper Mapper) (repository.Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate user: %w", err)
	}
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		return nil, fmt.Errorf("failed to migrate task: %w", err)
	}

	return &databaseRepository{
		db:     db,
		mapper: mapper,
	}, nil
}

func (r *databaseRepository) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	u, err := r.mapper.UserToModel(user)
	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return nil, err
	}

	return r.mapper.UserToDomain(u), nil
}

func (r *databaseRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return r.mapper.UserToDomain(&user), nil
}

func (r *databaseRepository) DeleteUserByID(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *databaseRepository) CreateTask(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	t, err := r.mapper.TaskToModel(task)
	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return nil, err
	}

	return r.mapper.TaskToDomain(t), nil
}

func (r *databaseRepository) GetTaskByID(ctx context.Context, ID string) (*entities.Task, error) {
	var t models.Task
	if err := r.db.WithContext(ctx).Where("id = ?", ID).First(&t).Error; err != nil {
		return nil, err
	}

	return r.mapper.TaskToDomain(&t), nil
}

func (r *databaseRepository) GetTasks(ctx context.Context, query *valueobjects.GetTasksQuery) ([]*entities.Task, error) {
	q := r.db.Model(&models.Task{}).Where("user_id = ?", query.UserID())

	if len(query.Filters().Statuses) > 0 {
		q = q.Where("status IN ?", query.Filters().Statuses)
	}

	if len(query.Filters().Priorities) > 0 {
		q = q.Where("priority IN ?", query.Filters().Priorities)
	}

	if query.Title() != "" {
		// ILIKE for postgres
		// can be replace to LOWER(title) LIKE LOWER(?)
		q = q.Where("title ILIKE ?", "%"+query.Title()+"%")
	}

	var orderField string
	switch query.OrderBy().Field {
	case valueobjects.SortByPriority:
		orderField = "priority"
	case valueobjects.SortByDueDate:
		orderField = "due_date"
	case valueobjects.SortByCreatedAt:
		orderField = "created_at"
	default:
		orderField = "priority"
	}

	dir := "ASC"
	if query.OrderBy().Direction == valueobjects.SortDesc {
		dir = "DESC"
	}

	q = q.Order(orderField + " " + dir)

	offset := (query.PageNumber() - 1) * query.PageSize()
	q = q.Limit(int(query.PageSize())).Offset(int(offset))

	var t []models.Task
	if err := q.WithContext(ctx).Find(&t).Error; err != nil {
		return nil, err
	}

	var tasks []*entities.Task
	for _, task := range t {
		tasks = append(tasks, r.mapper.TaskToDomain(&task))
	}

	return tasks, nil
}

func (r *databaseRepository) UpdateTask(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	t, err := r.mapper.TaskToModel(task)
	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Save(t).Error; err != nil {
		return nil, err
	}

	return r.mapper.TaskToDomain(t), nil
}

func (r *databaseRepository) DeleteTasks(ctx context.Context, IDs []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", IDs).Delete(&models.Task{}).Error
}
