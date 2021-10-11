package storage

import "time"

type Device struct {
	Did             int32     `gorm:"column:did;primary_key:true;AUTO_INCREMENT;comment:'宠端设备ID'" json:"did"`                // 设备唯一id
	Type            int8      `gorm:"column:type;type:varchar(20);NOT NULL;default:'';comment:'设备类型 0-宠端 1-人端'" json:"type"` // 设备类型
	Mac             string    `gorm:"column:mac;type:varchar(25);default:'';comment:'设备mac地址'" json:"mac"`
	Sn              string    `gorm:"column:sn;type:varchar(30);primary_key;default:'';comment:'设备序列号'"  json:"sn"`
	DeviceModel     string    `gorm:"column:device_model;type:varchar(120);default:'';comment:'设备型号'" json:"device_model"`        // 版本型号
	DeviceName      string    `gorm:"column:device_name;type:varchar(30);default:'';comment:'设备名称'" json:"device_name"`           // 设备名称
	DeviceVersion   string    `gorm:"column:device_version;type:varchar(25);default:'';comment:'设备版本'" json:"device_version"`     // 设备版本号
	SoftwareVersion string    `gorm:"column:software_version;type:varchar(25);default:'';comment:'软件版本'" json:"software_version"` // 软件名称
	DeviceNickname  string    `gorm:"column:device_nickname;type:varchar(30);default:'';comment:'设备昵称'" json:"device_nickname"`   // 设备昵称
	DeviceAvatar    string    `gorm:"column:device_avatar;type:varchar(120);default:'';comment:'设备头像'" json:"device_avatar"`      // 设备头像
	AudioId         int32     `gorm:"column:audio_id;type:int(11);default:'1';comment:'录音ID'" json:"audio_id"`                    // 宠物播放录音的ID
	AudioName       string    `gorm:"column:audio_name;type:varchar(25);default:'';comment:'录音名称'" json:"audio_name"`             // 宠物播放录音的名称
	AudioUrl        string    `gorm:"column:audio_url;type:varchar(120);default:'';comment:'录音URL'" json:"audio_url"`             // 宠物播放录音的Url
	WorkModel       uint8     `gorm:"column:work_model;type:int(11);comment:'设备工作模式'" json:"work_model"`                          // 工作模式： 1：发送位置信息 2：发送和接收用户信息 3：中继模式
	Interval        int32     `gorm:"column:interval_time;type:int(11);comment:'设备上报时间间隔'" json:"interval"`                       // 设备上报时间间隔
	SearchSwitch    int32     `gorm:"column:search_switch;type:int(11);comment:'社区寻找开关 1-开 2-关'" json:"search_switch"`            // 社区寻找开关 1-开 2-关
	LocLightSw      int32     `gorm:"column:loc_light_sw;type:tinyint(4);comment:'位指示灯开关 1-开 2-关'" json:"loc_light_sw"`           // 定位指示灯开关 1-开 2-关
	AlarmLightSw    int32     `gorm:"column:alarm_light_sw;type:tinyint(4);comment:'报警指示灯开关 1-开 2-关'" json:"alarm_light_sw"`      // 报警指示灯开关 1-开 2-关
	SoundSw         int32     `gorm:"column:sound_sw;type:tinyint(4);comment:'声音开关 1-开 2-关'" json:"sound_sw"`                     // 声音开关 1-开 2-关
	Sensor          int32     `gorm:"column:sensor;type:tinyint(4);comment:'传感器灵敏度'" json:"sensor"`                               // 传感器灵敏度
	DataState       uint8     `gorm:"column:data_state;type:int(11);NOT NULL;default:'1';comment:'数据状态'" json:"-"`
	CreateTime      time.Time `gorm:"column:create_time;type:timestamp;default:'1970-01-01 21:00:01'" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:timestamp;default:'1970-01-01 21:00:01'" json:"update_time"` // 修改时间
	Source          string    `gorm:"column:source;type:varchar(15);default:'';comment:'source'" json:"-"`

	LoraExtend string `gorm:"column:lora_extend;type:varchar(255);default:'';comment:'loraWAN模式下的参数'" json:"lora_extend"`
	JoinState  int    `gorm:"-" json:"join_state"`
}

// gorm自动建表
func (d *Device) TableName() string {
	return "device"
}

// 表名
func (d *Device) table() string {
	return "device"
}
