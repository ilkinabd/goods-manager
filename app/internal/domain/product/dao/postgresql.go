package dao

import (
	"context"
	filter2 "github.com/ilkinabd/goods-manager/app/internal/domain/product/filter"

	sq "github.com/Masterminds/squirrel"
	db "github.com/ilkinabd/goods-manager/app/pkg/client/postgresql/model"
	"github.com/ilkinabd/goods-manager/app/pkg/errors"
	"github.com/ilkinabd/goods-manager/app/pkg/logging"
)

type productDAOPostgres struct {
	queryBuilder sq.StatementBuilderType
	client       PostgreSQLClient
}

func NewProductDAOPostgres(client PostgreSQLClient) ProductDAO {
	return &productDAOPostgres{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
	}
}

const (
	scheme      = "public"
	table       = "product"
	tableScheme = scheme + "." + table
)

func (s *productDAOPostgres) All(ctx context.Context, filtering []filter2.Criteria, sorting filter2.Sortable) ([]*Product, error) {
	query := s.queryBuilder.
		Select("id").
		Columns(
			"name",
			"description",
			"image_id",
			"price",
			"currency_id",
			"rating",
			"category_id",
			"specification",
			"created_at",
			"updated_at",
		).
		From(tableScheme)

	for _, filter := range filtering {
		query = filter.MeetCriteria(query)
	}

	query = sorting.Sort(query)

	sql, args, err := query.ToSql()
	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args,
	})
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	rows, err := s.client.Query(ctx, sql, args...)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, err
	}

	defer rows.Close()

	list := make([]*Product, 0)

	for rows.Next() {
		ps := Product{}
		if err = rows.Scan(
			&ps.ID,
			&ps.Name,
			&ps.Description,
			&ps.ImageID,
			&ps.Price,
			&ps.CurrencyID,
			&ps.Rating,
			&ps.CategoryID,
			&ps.Specification,
			&ps.CreatedAt,
			&ps.UpdatedAt,
		); err != nil {
			err = db.ErrScan(err)
			logger.Error(err)
			return nil, err
		}

		list = append(list, &ps)
	}

	return list, nil
}

func (s *productDAOPostgres) Create(ctx context.Context, m map[string]interface{}) error {
	sql, args, buildErr := s.queryBuilder.
		Insert(tableScheme).
		SetMap(m).
		PlaceholderFormat(sq.Dollar).ToSql()

	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args,
	})
	if buildErr != nil {
		buildErr = db.ErrCreateQuery(buildErr)
		logger.Error(buildErr)
		return buildErr
	}

	if exec, execErr := s.client.Exec(ctx, sql, args...); execErr != nil {
		execErr = db.ErrDoQuery(execErr)
		logger.Error(execErr)
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Insert() {
		execErr = db.ErrDoQuery(errors.New("product was not created. 0 rows were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}

func (s *productDAOPostgres) One(ctx context.Context, id string) (*Product, error) {
	sql, args, buildErr := s.queryBuilder.
		Select("id").
		Columns(
			"name",
			"description",
			"image_id",
			"price",
			"currency_id",
			"rating",
			"category_id",
			"specification",
			"created_at",
			"updated_at",
		).
		From(tableScheme).
		Where(sq.Eq{"id": id}).ToSql()

	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args,
	})
	if buildErr != nil {
		buildErr = db.ErrCreateQuery(buildErr)
		logger.Error(buildErr)
		return nil, buildErr
	}

	var ps Product

	err := s.client.QueryRow(ctx, sql, args...).Scan(
		&ps.ID,
		&ps.Name,
		&ps.Description,
		&ps.ImageID,
		&ps.Price,
		&ps.CurrencyID,
		&ps.Rating,
		&ps.CategoryID,
		&ps.Specification,
		&ps.CreatedAt,
		&ps.UpdatedAt,
	)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, err
	}

	return &ps, nil
}

func (s *productDAOPostgres) Update(ctx context.Context, id string, m map[string]interface{}) error {
	sql, args, buildErr := s.queryBuilder.
		Update(tableScheme).
		SetMap(m).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args,
	})
	if buildErr != nil {
		buildErr = db.ErrCreateQuery(buildErr)
		logger.Error(buildErr)
		return buildErr
	}

	if exec, execErr := s.client.Exec(ctx, sql, args...); execErr != nil {
		execErr = db.ErrDoQuery(execErr)
		logger.Error(execErr)
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Update() {
		execErr = db.ErrDoQuery(errors.New("product was not updated. 0 rows were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}

func (s *productDAOPostgres) Delete(ctx context.Context, id string) error {
	sql, args, buildErr := s.queryBuilder.
		Delete(tableScheme).
		Where(sq.Eq{"id": id}).
		ToSql()

	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args,
	})
	if buildErr != nil {
		buildErr = db.ErrCreateQuery(buildErr)
		logger.Error(buildErr)
		return buildErr
	}

	if exec, execErr := s.client.Exec(ctx, sql, args...); execErr != nil {
		execErr = db.ErrDoQuery(execErr)
		logger.Error(execErr)
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Delete() {
		execErr = db.ErrDoQuery(errors.New("product was not deleted. 0 rows were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}
