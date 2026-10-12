package almanac

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/fatWill/agent-team-blog/backend/models"
	algo "github.com/fatWill/agent-team-blog/backend/pkg/almanac"
)

// Repo 黄历存储（单表原子操作）
type Repo struct {
	db *gorm.DB
}

// NewRepo 创建黄历存储
func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// GetByDate 查询单日；不存在返回 (nil, nil)。用 Find 而非 First，避免未命中时 GORM 打印 record not found 日志
func (r *Repo) GetByDate(date string) (*models.AlmanacDay, error) {
	var list []models.AlmanacDay
	if err := r.db.Where("date = ?", date).Limit(1).Find(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// ListRange 查询 [start, end] 闭区间（按日期升序）
func (r *Repo) ListRange(start, end string) ([]models.AlmanacDay, error) {
	var list []models.AlmanacDay
	err := r.db.Where("date >= ? AND date <= ?", start, end).Order("date ASC").Find(&list).Error
	return list, err
}

// SourceMap 返回 [start, end] 内已入库日期 → source，用于断点续爬
func (r *Repo) SourceMap(start, end string) (map[string]string, error) {
	type row struct {
		Date   string
		Source string
	}
	var rows []row
	err := r.db.Model(&models.AlmanacDay{}).Select("date, source").
		Where("date >= ? AND date <= ?", start, end).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(rows))
	for _, x := range rows {
		m[x.Date] = x.Source
	}
	return m, nil
}

// Upsert 按 date 幂等写入（存在则整行覆盖，保留 id 与 created_at）
func (r *Repo) Upsert(d *models.AlmanacDay) error {
	now := time.Now()
	d.UpdatedAt = now
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}},
		UpdateAll: true,
	}).Omit("id").Create(d).Error
}

// UpsertBatch 单事务批量按 date upsert，返回写入行数（任一失败整体回滚）
func (r *Repo) UpsertBatch(rows []models.AlmanacDay) (int, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		txr := &Repo{db: tx}
		for i := range rows {
			rows[i].ID = 0
			if err := txr.Upsert(&rows[i]); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return len(rows), nil
}

// GetOrCompute 优先读库，库里没有则用算法现算（不回写，避免请求路径上的写放大）
func (r *Repo) GetOrCompute(date string) (*models.AlmanacDay, error) {
	d, err := r.GetByDate(date)
	if err != nil || d != nil {
		return d, err
	}
	t, err := algo.ParseDate(date)
	if err != nil {
		return nil, err
	}
	return algo.Compute(t), nil
}
