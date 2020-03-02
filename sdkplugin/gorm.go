package sdkplugin

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hitechqb/gosdk"
	"github.com/hitechqb/gosdk/sdkcm"
	"github.com/jinzhu/gorm"
)

type gormPlugin struct {
	*plugin
	db *gorm.DB
}

func NewGormPlugin(options ...PluginOption) *gormPlugin {
	p := &plugin{}
	s := &gormPlugin{plugin: p}

	for _, o := range options {
		o(p)
	}

	return s
}

//func (s *gormPlugin) Name() string {
//	return s.name
//}
//
//func (s *gormPlugin) Order() int {
//	return s.order
//}

func (s *gormPlugin) GetDb() *gorm.DB {
	return s.db
}

func (s *gormPlugin) Run(service sdk.Service) error {
	logger := sdkcm.GetLogger()
	cf := sdkcm.GetConfig()

	connectionString := s.getConnectionString(cf)

	logger.Infof("🔔 Connecting to database on: %s\n", connectionString)
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	if err := db.DB().Ping(); err != nil {
		return err
	}

	logger.Infof(`🎉 Connected to "%s" database !`, cf.DbName())

	s.db = db
	return nil
}

func (s *gormPlugin) Stop() error {
	return s.db.Close()
}

func (s *gormPlugin) getConnectionString(cf sdkcm.SDKConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Asia%%2fSaigon",
		cf.DbUser(),
		cf.DbPassword(),
		cf.DbHost(),
		cf.DbPort(),
		cf.DbName())
}

//func WithName(name string) GormOption {
//	return func(p *gormPlugin) {
//		p.name = name
//	}
//}
//
//func WithOrder(order int) GormOption {
//	return func(p *gormPlugin) {
//		p.order = order
//	}
//}
