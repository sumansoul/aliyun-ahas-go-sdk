package meta

type Meta struct {
	license   string
	namespace string
	deployEnv string

	inVpc        bool
	regionId     string
	vpcId        string
	privateIp    string
	hostIp       string
	hostName     string
	pid          string
	cid          string
	tid          string
	instanceId   string
	deviceType   int
	uid          string
	version      string
	ahasEndpoint string

	tidChan chan string

	debugging bool
}

func (m *Meta) TidChan() chan string {
	return m.tidChan
}

func (m *Meta) SetUid(uid string) {
	m.uid = uid
}

func (m *Meta) SetTid(tid string) {
	m.tid = tid
	m.tidChan <- tid
}

func (m *Meta) SetCid(cid string) {
	m.cid = cid
}

func (m *Meta) Cid() string {
	return m.cid
}

func (m *Meta) Version() string {
	return m.version
}

func (m *Meta) InstanceId() string {
	return m.instanceId
}

func (m *Meta) Tid() string {
	return m.tid
}

func (m *Meta) Uid() string {
	return m.uid
}

func (m *Meta) Pid() string {
	return m.pid
}

func (m *Meta) AhasEndpoint() string {
	return m.ahasEndpoint
}

func (m *Meta) HostName() string {
	return m.hostName
}

func (m *Meta) PrivateIp() string {
	return m.privateIp
}

func (m *Meta) HostIp() string {
	return m.hostIp
}

func (m *Meta) VpcId() string {
	return m.vpcId
}

func (m *Meta) RegionId() string {
	return m.regionId
}

func (m *Meta) InVpc() bool {
	return m.inVpc
}

func (m *Meta) DeviceType() int {
	return m.deviceType
}
