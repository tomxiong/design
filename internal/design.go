package design

import (
	"context"
	"design/internal/conf"
	"design/internal/dao"
	"design/internal/model"
)

type Design struct {
	c   *conf.Config
	dao *dao.Dao
}

// New init
func New(c *conf.Config) (d *Design) {
	d = &Design{
		c:   c,
		dao: dao.New(c),
	}
	return d
}

func (d *Design) Close() error {
	// TODO: implement Close
	return nil
}

func (d *Design) Register(reg model.Member) (bool, error) {
	return d.dao.Register(context.Background(), reg)
}

func (d *Design) ListMember(role, status string) ([]model.Member, error) {
	return d.dao.ListMember(context.Background(), role, status)
}

func (d *Design) SetEmail(mem model.Member) (bool, error) {
	return d.dao.SetEmail(context.Background(), mem.Id, mem.Email)
}

func (d *Design) GetMember(token string) (model.Member, error) {
	return d.dao.GetMember(context.Background(), token)
}
