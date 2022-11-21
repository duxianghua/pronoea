package contact

import (
	"regexp"
	"time"

	"github.com/duxianghua/pronoea/internal/config"
	"github.com/go-pg/pg/v10"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
)

var DEFAULT_EMAIL = []string{"xianghua.du2@homepartners.tech"}

type ServiceInfo struct {
	tableName      struct{} `pg:"service_info"`
	Id             int64
	Service        string
	Server         string
	Contact_type   string
	Contact_first  string
	Contact_second string
	Contact_third  string
}

type ServiceInfoList struct {
}

type Contact struct {
	cache *cache.Cache
}

var (
	contact *Contact
)

func NewContactStore(cfg config.DBConfig) Contact {
	return Contact{
		cache: cache.New(30*time.Minute, 10*time.Minute),
	}
}

func Store() *Contact {
	if contact == nil {
		contact = &Contact{
			cache: cache.New(30*time.Minute, 10*time.Minute),
		}
	}
	return contact
}

func (c *Contact) Sync(cfg config.DBConfig) error {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.Host,
		User:     cfg.Username,
		Password: cfg.Password,
		Database: cfg.Database,
	})
	defer db.Close()

	// select all
	var serviceInfoList []ServiceInfo
	err := db.Model(&serviceInfoList).Select()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	c.cache.Set("serviceInfo", serviceInfoList, cache.NoExpiration)
	// for _, i := range serviceInfoList {
	// 	c.cache.Set(i.Service, i, cache.NoExpiration)
	// 	log.Info().Str("service", i.Service).Interface("service_info", i).Msg("add contact info to cache")
	// }
	log.Info().Int("count", len(serviceInfoList)).Msg("sync contact success")
	return nil
}

func (c *Contact) Get(serivce string) []string {
	var emails []string
	if len(serivce) == 0 {
		return emails
	}
	data, ok := c.cache.Get("serviceInfo")
	if !ok {
		return emails
	}
	items := data.([]ServiceInfo)
	for _, i := range items {
		ok, err := regexp.MatchString(serivce, i.Service)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		if ok {
			if ok, _ := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, i.Contact_first); ok {
				emails = append(emails, i.Contact_first)
			}
			if ok, _ := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, i.Contact_second); ok {
				emails = append(emails, i.Contact_second)
			}
			if ok, _ := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, i.Contact_third); ok {
				emails = append(emails, i.Contact_third)
			}
			return emails
		}
	}
	return emails
}
