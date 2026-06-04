package plugin

import (
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/models"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtgtime"
)

func SetDefaultTimezone(timezone string) {
	prtgtime.SetDefaultTimezone(timezone)
}

func ParseTimeInit(s *models.PluginSettings) {
	prtgtime.ParseTimeInit(s)
}

func parsePRTGDateTime(datetime string) (time.Time, string, error) {
	return prtgtime.ParseDateTime(datetime)
}
