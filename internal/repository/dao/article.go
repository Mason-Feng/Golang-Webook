package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
}

func (a *ArticleGORMDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	res := a.db.WithContext(ctx).Model(&Article{}).Where("id=? AND author_id=?", art.Id, art.AuthorId).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"utime":   now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新失败，ID不会或者作者不对")
	}
	return nil
}

type ArticleGORMDAO struct {
	db *gorm.DB
}

func (a *ArticleGORMDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		dao := NewArticleGORMDAO(tx)
		if id > 0 {
			err = dao.UpdateById(ctx, art)
		} else {
			id, err = dao.Insert(ctx, art)
		}

		if err != nil {
			return err
		}
		art.Id = id
		now := time.Now().UnixMilli()
		pubArt := PublishedArticle(art)
		pubArt.Ctime = now
		pubArt.Utime = now
		err = tx.Clauses(clause.OnConflict{
			//对mysql不起效，但是可以兼容别的方言
			//INSERT XXX ON DUPLICATE KEY SET `title`=?
			//别的方言：sqlite INSERT XXX ON CONFLICT DO NOTHING
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   pubArt.Title,
				"content": pubArt.Content,
				"utime":   now,
			}),
		}).Create(&pubArt).Error
		return err
	})

	return id, err
}
func (a *ArticleGORMDAO) SyncV1(ctx context.Context, art Article) (int64, error) {
	tx := a.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()

	var (
		id  = art.Id
		err error
	)
	dao := NewArticleGORMDAO(tx)
	if id > 0 {
		err = dao.UpdateById(ctx, art)
	} else {
		id, err = dao.Insert(ctx, art)
	}

	if err != nil {
		return 0, err
	}
	art.Id = id
	now := time.Now().UnixMilli()
	pubArt := PublishedArticle(art)
	pubArt.Ctime = now
	pubArt.Utime = now
	err = tx.Clauses(clause.OnConflict{
		//对mysql不起效，但是可以兼容别的方言
		//INSERT XXX ON DUPLICATE KEY SET `title`=?
		//别的方言：sqlite INSERT XXX ON CONFLICT DO NOTHING
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   pubArt.Title,
			"content": pubArt.Content,
			"utime":   now,
		}),
	}).Create(&pubArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, err
}

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"type=varchar(4096)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`
	Ctime    int64  // 创建时间,时区 UTC 0毫秒数
	Utime    int64  // 更新时间
}

func NewArticleGORMDAO(db *gorm.DB) ArticleDAO {
	return &ArticleGORMDAO{
		db: db,
	}
}

func (a *ArticleGORMDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := a.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

type PublishedArticle Article
