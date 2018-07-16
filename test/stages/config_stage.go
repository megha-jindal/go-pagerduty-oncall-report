package stages

import (
	"testing"

	"bytes"

	"github.com/rogersole/go-pagerduty-oncall-report/configuration"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type ConfigStage struct {
	t *testing.T

	configRaw            []byte
	config               *configuration.Configuration
	configError          error
	configUnmarshalError error

	mapValue interface{}
	mapError error
}

func ConfigTest(t *testing.T) (*ConfigStage, *ConfigStage, *ConfigStage) {
	stage := &ConfigStage{
		t: t,
	}

	return stage, stage, stage
}

func (s *ConfigStage) And() *ConfigStage {
	return s
}

func (s *ConfigStage) A_valid_configuration() *ConfigStage {
	s.configRaw = []byte(`
pdAuthToken: abcdefghijklm
rotationStartHour: 08:00:00
currency: £
rotationPrices:
  - type: weekday
    price: 1
  - type: weekend
    price: 1
  - type: bankholiday
    price: 2
rotationUsers:
  - name: "User 1"
    holidaysCalendar: uk
    userId: ABCDEF1
  - name: "User 2"
    holidaysCalendar: uk
    userId: ABCDEF2
schedulesToIgnore:
  - SCHED_1
  - SCHED_2
  - SCHED_3
`)
	return s
}

func (s *ConfigStage) A_malformed_configuration() *ConfigStage {
	s.configRaw = []byte(`
pdAuthToken: abcdefghijklm
	rotationStartHour: 08:00:00
  currency: £
		rotationPrices:
  	- type: weekday
    	price: 1
  - type: weekend
    price: 1
  	- type: bankholiday
    price: 2
	rotationUsers:
  - 	name: "User 1"
    holidaysCalendar: uk
    userId: ABCDEF1
  - name: "User 2"
    holidaysCalendar: uk
    userId: ABCDEF2
`)
	return s
}

func (s *ConfigStage) A_valid_configuration_correctly_loaded() *ConfigStage {
	s.A_valid_configuration().And().It_is_loaded()
	assert.Nil(s.t, s.configError)
	return s
}

func (s *ConfigStage) It_is_loaded() *ConfigStage {
	viper.SetConfigType("yaml")
	s.configError = viper.ReadConfig(bytes.NewBuffer(s.configRaw))
	if s.configError == nil {
		s.config = configuration.New()
		s.configUnmarshalError = viper.Unmarshal(s.config)
	}
	return s
}

func (s *ConfigStage) An_existing_price_is_requested() *ConfigStage {
	s.mapValue, s.mapError = s.config.FindPriceByDay("weekday")
	return s
}

func (s *ConfigStage) A_non_existing_price_is_requested() *ConfigStage {
	s.mapValue, s.mapError = s.config.FindPriceByDay("wokday")
	return s
}

func (s *ConfigStage) An_existing_rotation_info_is_requested() *ConfigStage {
	s.mapValue, s.mapError = s.config.FindRotationUserInfoByID("ABCDEF1")
	return s
}

func (s *ConfigStage) A_non_existing_rotation_info_is_requested() *ConfigStage {
	s.mapValue, s.mapError = s.config.FindRotationUserInfoByID("NONE")
	return s
}

func (s *ConfigStage) Value_is_found() *ConfigStage {
	assert.Nil(s.t, s.mapError)
	assert.NotNil(s.t, s.mapValue)
	return s
}

func (s *ConfigStage) Value_is_not_found() *ConfigStage {
	assert.NotNil(s.t, s.mapError)
	assert.Nil(s.t, s.mapValue)
	return s
}
func (s *ConfigStage) Config_error_is_created() *ConfigStage {
	assert.NotNil(s.t, s.configError)
	return s
}