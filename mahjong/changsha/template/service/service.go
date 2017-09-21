package service

import "context"

type ChangShaRoomTemplateService interface {
	GetAll() *ChangShaRoomTemplateConfig
	GetRoundList() []*ChangShaRoomRoundTemplate
	GetPeopleList() []*ChangShaRoomPeopleTemplate
	GetZhuaNiaoList() []*ChangShaRoomZhuaNiaoTemplate
	GetPeopleTemplateById(id int) *ChangShaRoomPeopleTemplate
	GetZhuaNiaoTemplateById(id int) *ChangShaRoomZhuaNiaoTemplate
	GetRoundTemplateById(id int) *ChangShaRoomRoundTemplate
}

type ChangShaRoomRoundTemplate struct {
	Id    int `json:"id"`
	Round int `json:"round"`
	Cost  int `json:"cost"`
}

type ChangShaRoomPeopleTemplate struct {
	Id     int `json:"id"`
	People int `json:"people"`
}

type ChangShaRoomZhuaNiaoTemplate struct {
	Id       int `json:"id"`
	ZhuaNiao int `json:"zhuaNiao"`
}

type ChangShaRoomTemplateConfig struct {
	RoundList    []*ChangShaRoomRoundTemplate    `json:"round"`
	PeopleList   []*ChangShaRoomPeopleTemplate   `json:"people"`
	ZhuaNiaoList []*ChangShaRoomZhuaNiaoTemplate `json:"zhuaNiao"`
}

type changShaRoomTemplateService struct {
	config *ChangShaRoomTemplateConfig
}

func (csrts *changShaRoomTemplateService) GetRoundList() []*ChangShaRoomRoundTemplate {
	return csrts.config.RoundList
}

func (csrts *changShaRoomTemplateService) GetPeopleList() []*ChangShaRoomPeopleTemplate {
	return csrts.config.PeopleList
}
func (csrts *changShaRoomTemplateService) GetZhuaNiaoList() []*ChangShaRoomZhuaNiaoTemplate {
	return csrts.config.ZhuaNiaoList
}

func (csrts *changShaRoomTemplateService) GetAll() *ChangShaRoomTemplateConfig {
	return csrts.config
}

func (csrts *changShaRoomTemplateService) GetPeopleTemplateById(id int) *ChangShaRoomPeopleTemplate {
	for _, template := range csrts.GetPeopleList() {
		if template.Id == id {
			return template
		}
	}
	return nil
}

func (csrts *changShaRoomTemplateService) GetZhuaNiaoTemplateById(id int) *ChangShaRoomZhuaNiaoTemplate {
	for _, template := range csrts.GetZhuaNiaoList() {
		if template.Id == id {
			return template
		}
	}
	return nil
}

func (csrts *changShaRoomTemplateService) GetRoundTemplateById(id int) *ChangShaRoomRoundTemplate {
	for _, template := range csrts.GetRoundList() {
		if template.Id == id {
			return template
		}
	}
	return nil
}

func NewChangShaRoomTemplateService(c *ChangShaRoomTemplateConfig) ChangShaRoomTemplateService {
	return &changShaRoomTemplateService{
		config: c,
	}
}

const (
	key = "changsha_template_service"
)

func WithChangShaTemplateService(ctx context.Context, csrts ChangShaRoomTemplateService) context.Context {
	return context.WithValue(ctx, key, csrts)
}

func ChangShaTemplateServiceInContext(ctx context.Context) ChangShaRoomTemplateService {
	us, ok := ctx.Value(key).(ChangShaRoomTemplateService)
	if !ok {
		return nil
	}
	return us
}
