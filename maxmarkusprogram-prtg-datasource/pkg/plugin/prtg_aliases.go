package plugin

import (
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtg"
)

type Api = prtg.Api
type ApiInterface = prtg.ApiInterface
type PRTGAPI = prtg.PRTGAPI

type Annotation = prtg.Annotation
type AnnotationQuery = prtg.AnnotationQuery
type AnnotationResponse = prtg.AnnotationResponse
type Device = prtg.Device
type Group = prtg.Group
type KeyValue = prtg.KeyValue
type PrtgAnnotation = prtg.PrtgAnnotation
type PrtgAnnotationResponse = prtg.PrtgAnnotationResponse
type PrtgChannelValueStruct = prtg.PrtgChannelValueStruct
type PrtgChannelsListResponse = prtg.PrtgChannelsListResponse
type PrtgDeviceListItemStruct = prtg.PrtgDeviceListItemStruct
type PrtgDevicesListResponse = prtg.PrtgDevicesListResponse
type PrtgGroupListItemStruct = prtg.PrtgGroupListItemStruct
type PrtgGroupListResponse = prtg.PrtgGroupListResponse
type PrtgHistoricalDataResponse = prtg.PrtgHistoricalDataResponse
type PrtgManualMethodResponse = prtg.PrtgManualMethodResponse
type PrtgSensorListItemStruct = prtg.PrtgSensorListItemStruct
type PrtgSensorsListResponse = prtg.PrtgSensorsListResponse
type PrtgStatusListResponse = prtg.PrtgStatusListResponse
type PrtgValues = prtg.PrtgValues
type Sensor = prtg.Sensor
type StringOrNumber = prtg.StringOrNumber

func NewApi(baseURL, apiKey string, cacheTime, requestTimeout time.Duration) *Api {
	return prtg.NewApi(baseURL, apiKey, cacheTime, requestTimeout)
}
