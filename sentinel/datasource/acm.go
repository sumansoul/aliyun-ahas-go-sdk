package datasource

import (
	"encoding/json"
	"time"

	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	sentinelConf "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/aliyun/aliyun-ahas-go-sdk/logger"
	"github.com/aliyun/aliyun-ahas-go-sdk/meta"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
)

const (
	AcmGroupId = "ahas-sentinel"

	FlowRuleDataIdPrefix            = "flow-rule-"
	SystemRuleDataIdPrefix          = "system-rule-"
	CircuitBreakingRuleDataIdPrefix = "degrade-rule-"
	ParamFlowRuleDataIdPrefix       = "param-flow-rule-"
)

func formFlowRuleDataId(userId, namespace, appName string) string {
	return FlowRuleDataIdPrefix + userId + "-" + namespace + "-" + appName
}

func formSystemRuleDataId(userId, namespace, appName string) string {
	return SystemRuleDataIdPrefix + userId + "-" + namespace + "-" + appName
}

func formCircuitBreakingRuleDataId(userId, namespace, appName string) string {
	return CircuitBreakingRuleDataIdPrefix + userId + "-" + namespace + "-" + appName
}

func formParamFlowRuleDataId(userId, namespace, appName string) string {
	return ParamFlowRuleDataIdPrefix + userId + "-" + namespace + "-" + appName
}

func InitAcm(acmHost string, conf Config, m *meta.Meta) error {
	ch := m.TidChan()
	select {
	case <-ch:
		break
	case <-time.After(30 * time.Second):
		return errors.New("wait AHAS transport timeout")
	}

	clientConfig := constant.ClientConfig{
		TimeoutMs:      conf.TimeoutMs,
		ListenInterval: conf.ListenIntervalMs,
		NamespaceId:    m.Tid(),
		Endpoint:       acmHost + ":8080",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})
	if err != nil {
		return err
	}

	// Add flow/isolation rule config listener.
	flowRuleDataId := formFlowRuleDataId(m.Uid(), meta.Namespace(), sentinelConf.AppName())
	err = registerRuleDataSource(flowRuleDataId, onFlowRuleChange, configClient)
	if err != nil {
		return err
	}
	// Add system rule config listener.
	systemRuleDataId := formSystemRuleDataId(m.Uid(), meta.Namespace(), sentinelConf.AppName())
	err = registerRuleDataSource(systemRuleDataId, onSystemRuleChange, configClient)
	if err != nil {
		return err
	}
	// Add circuit breaking rule config listener.
	circuitBreakerRuleDataId := formCircuitBreakingRuleDataId(m.Uid(), meta.Namespace(), sentinelConf.AppName())
	err = registerRuleDataSource(circuitBreakerRuleDataId, onCircuitBreakingRuleChange, configClient)
	if err != nil {
		return err
	}
	// Add param flow rule config listener.
	paramFlowRuleDataId := formParamFlowRuleDataId(m.Uid(), meta.Namespace(), sentinelConf.AppName())
	err = registerRuleDataSource(paramFlowRuleDataId, onParamFlowRuleChange, configClient)
	if err != nil {
		return err
	}

	logging.Info("ACM data source initialized successfully", "flowRuleDataId", flowRuleDataId)
	logger.Infof("ACM data source initialized successfully, flow dataId: %s", flowRuleDataId)
	return nil
}

func registerRuleDataSource(dataId string, handler func(string), nacosClient config_client.IConfigClient) error {
	nacosConfig := vo.ConfigParam{
		Group:  AcmGroupId,
		DataId: dataId,
		OnChange: func(namespace, group, dataId, data string) {
			handler(data)
		},
	}
	go func() {
		data, err := nacosClient.GetConfig(nacosConfig)
		if err != nil {
			logging.Error(err, "Failed to getConfig from ACM", "dataId", dataId)
		} else if len(data) > 0 {
			handler(data)
		}
	}()
	return nacosClient.ListenConfig(nacosConfig)
}

func onFlowRuleChange(data string) {
	logging.Info("ACM data received for flow rules", "data", data)
	d := &struct {
		Version string
		Data    []*LegacyFlowRule
	}{}
	err := json.Unmarshal([]byte(data), d)
	if err != nil {
		logging.Error(err, "Failed to parse flow rules")
		return
	}
	flowRules := make([]*flow.Rule, 0)
	isolationRules := make([]*isolation.Rule, 0)
	for _, r := range d.Data {
		if r.MetricType == 0 {
			if rule := r.ToIsolationRule(); rule != nil {
				isolationRules = append(isolationRules, rule)
			} else {
				logging.Warn("Cannot convert received rule to isolation rule, ignoring", "ruleId", r.ID)
			}
		} else {
			if rule := r.ToFlowRule(); rule != nil {
				flowRules = append(flowRules, rule)
			} else {
				logging.Warn("Cannot convert received rule to flow rule, ignoring", "ruleId", r.ID)
			}
		}
	}
	// separate flow/isolation rules here
	_, err = flow.LoadRules(flowRules)
	if err != nil {
		logging.Error(err, "Failed to load flow rules")
		return
	}
	if len(isolation.GetRules()) == 0 && len(isolationRules) == 0 {
		// If both current and received isolation rules are empty, then do not update
		return
	}
	_, err = isolation.LoadRules(isolationRules)
	if err != nil {
		logging.Error(err, "Failed to load isolation rules")
		return
	}
}

func onSystemRuleChange(data string) {
	logging.Info("ACM data received for system rules", "data", data)
	d := &struct {
		Version string
		Data    []*LegacySystemRule
	}{}
	err := json.Unmarshal([]byte(data), d)
	if err != nil {
		logging.Error(err, "Failed to parse legacy system rules")
		return
	}
	arr := make([]*system.Rule, 0)
	for _, r := range d.Data {
		if rule := r.ToGoRule(); rule != nil {
			arr = append(arr, rule)
		}
	}
	_, err = system.LoadRules(arr)
	if err != nil {
		logging.Error(err, "Failed to load system rules")
		return
	}
}

func onCircuitBreakingRuleChange(data string) {
	logging.Info("ACM data received for circuit breaking rules", "data", data)
	d := &struct {
		Version string
		Data    []*LegacyDegradeRule
	}{}
	err := json.Unmarshal([]byte(data), d)
	if err != nil {
		logging.Error(err, "Failed to parse legacy degrade rules")
		return
	}
	arr := make([]*circuitbreaker.Rule, 0)
	for _, r := range d.Data {
		if rule := r.ToGoRule(); rule != nil {
			arr = append(arr, rule)
		}
	}
	_, err = circuitbreaker.LoadRules(arr)
	if err != nil {
		logging.Error(err, "Failed to load circuit breaking rules")
		return
	}
}

func onParamFlowRuleChange(data string) {
	logging.Info("ACM data received for hot-spot param flow rules", "data", data)
	d := &struct {
		Version string
		Data    []*LegacyParamFlowRule
	}{}
	err := json.Unmarshal([]byte(data), d)
	if err != nil {
		logging.Error(err, "Failed to parse legacy param flow rules")
		return
	}
	arr := make([]*hotspot.Rule, 0)
	for _, r := range d.Data {
		if rule := r.ToGoRule(); rule != nil {
			arr = append(arr, rule)
		}
	}
	_, err = hotspot.LoadRules(arr)
	if err != nil {
		logging.Error(err, "Failed to load param flow rules")
		return
	}
}
