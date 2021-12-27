/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName="ms"
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`,description="The Actual status of Mongoshake"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// MongoShake is the Schema for the mongoshakes API
type MongoShake struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoShakeSpec   `json:"spec,omitempty"`
	Status MongoShakeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MongoShakeList contains a list of MongoShake
type MongoShakeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoShake `json:"items"`
}

// MongoShakeSpec defines the desired state of MongoShake
type MongoShakeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Image              string                 `json:"image,omitempty"`
	BackoffLimit       *int32                 `json:"backoffLimit,omitempty"`
	Pause              bool                   `json:"pause,omitempty"`
	ImagePullPolicy    corev1.PullPolicy      `json:"imagePullPolicy,omitempty"`
	Resources          *ResourcesSpec         `json:"resources,omitempty"`
	ReadinessProbe     *corev1.Probe          `json:"readinessProbe,omitempty"`
	LivenessProbe      *LivenessProbeExtended `json:"livenessProbe,omitempty"`
	Configuration      string                 `json:"configuration,omitempty"`
	Annotations        map[string]string      `json:"annotations,omitempty"`
	NodeSelector       map[string]string      `json:"nodeSelector,omitempty"`
	Tolerations        []corev1.Toleration    `json:"tolerations,omitempty"`
	VolumeSpec         *VolumeSpec            `json:"volumeSpec,omitempty"`
	Collector          *CollectorSpec         `json:"collector,omitempty"`
	HealthyCheckEnable bool                   `json:"healthy_check_enable,omitempty"`
}

type JobStatus struct {
	Size    int      `json:"size"`
	Ready   int      `json:"ready"`
	Status  AppState `json:"status,omitempty"`
	Message string   `json:"message,omitempty"`
}

// MongoShakeStatus defines the observed state of MongoShake
type MongoShakeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	State       AppState           `json:"state,omitempty"`
	JobStatus   *JobStatus         `json:"job_status,omitempty"`
	Message     string             `json:"message,omitempty"`
	Conditions  []ClusterCondition `json:"conditions,omitempty"`
	Host        string             `json:"host,omitempty"`
	HealthCheck HealthyCheckType   `json:"health_check,omitempty"`
}

// CollectorSpec defines the global configuration of MongoShake
type CollectorSpec struct {
	// current configuration version, do not modify.
	// 当前配置文件的版本号，请不要修改该值。
	ConfVersion int32 `json:"conf_version,omitempty"`

	// collector name
	// id用于输出pid文件等信息。
	Id string `json:"id,omitempty"`

	// 多个MongoShake同时拉取一个源端数据时，true为开启高可用设置。默认为false
	MasterQuorum bool `json:"master_quorum,omitempty"`

	// http api interface. Users can use this api to monitor mongoshake.
	// `curl 127.0.0.1:9100`.
	// We also provide a restful tool named "mongoshake-stat" to
	// print ack, lsn, checkpoint and qps information based on this api.
	// usage: `./mongoshake-stat --port=9100`
	// 全量和增量的restful监控端口，可以用curl查看内部监控metric统计情况。详见wiki。
	FullSyncHttpPort int32 `json:"full_sync_http_port,omitempty"`
	IncrSyncHttpPort int32 `json:"incr_sync_http_port,omitempty"`

	// profiling on net/http/profile
	// profiling端口，用于查看内部go堆栈。
	SystemProfilePort int32 `json:"system_profile_port,omitempty"`

	Log          *LogSpec         `json:"log,omitempty"`
	Mode         SyncMode         `json:"mode,omitempty"`
	Mongo        *MongoSpec       `json:"mongo,omitempty"`
	Tunnel       *TunnelSpec      `json:"tunnel,omitempty"`
	MongoConnect MongoConnectMode `json:"mongo_connect,omitempty"`
	Filter       *FilterSpec      `json:"filter,omitempty"`
	CheckPoint   *CheckPointSpec  `json:"check_point,omitempty"`

	// transform from source db or collection namespace to dest db or collection namespace.
	// at most one of these two parameters can be given.
	// transform: fromDbName1.fromCollectionName1:toDbName1.toCollectionName1;fromDbName2:toDbName2
	// 转换命名空间，比如a.b同步后变成c.d，谨慎建议开启，比较耗性能。
	TransformNamespace string        `json:"transform_namespace"`
	FullSync           *FullSyncSpec `json:"full_sync,omitempty"`
	IncrSync           *IncrSyncSpec `json:"incr_sync,omitempty"`

	// 特殊字段，标识源端类型，默认为空。阿里云MongoDB serverless集群请配置aliyun_serverless
	SpecialSourceDbFlag string `json:"special_source_db_flag,omitempty"`
}

// FetchMethod
// oplog: fetch oplog from source mongodb (default)
// change_stream: use change to receive change event from source mongodb, support MongoDB >= 4.0.
// we recommand to use change_stream if possible.
type FetchMethod string

const (
	FetchMethodOplog        FetchMethod = "oplog"
	FetchMethodChangeStream FetchMethod = "change_stream"
)

type ShardKeyType string

const (
	ShardKeyAuto       ShardKeyType = "auto"
	ShardKeyId         ShardKeyType = "id"
	ShardKeyCollection ShardKeyType = "collection"
)

// IncrSyncSpec defines the incr sync configuration
type IncrSyncSpec struct {
	MongoFetchMethod FetchMethod `json:"mongo_fetch_method,omitempty"`
	// After the document is updated, the fields that only need to be updated are set to false,
	// and the contents of all documents are set to true
	// 更新文档后,只需要更新的字段则设为false,需要全部文档内容则设为true
	// 只在mongo_fetch_method = change_stream 模式下生效，且性能有所下降
	ChangeStreamWatchFullDocument bool `json:"change_stream_watch_full_document,omitempty"`

	// global id. used in active-active replication.
	// this parameter is not supported on current open-source version.
	// gid用于双活防止环形复制，目前只用于阿里云云上MongoDB，如果是阿里云云上实例互相同步
	// 希望开启gid，请联系阿里云售后，sharding的有多个gid请以分号(;)分隔。
	OplogGids string `json:"oplog_gids,omitempty"`

	// distribute data to different worker by hash key to run in parallel.
	// [auto] 		decide by if there has unique index in collections.
	// 		 		use `collection` if has unique index otherwise use `id`.
	// [id] 			shard by ObjectId. handle oplogs in sequence by unique _id
	// [collection] 	shard by ns. handle oplogs in sequence by unique ns
	// hash的方式，id表示按文档hash，collection表示按表hash，auto表示自动选择hash类型。
	// 如果没有索引建议选择id达到非常高的同步性能，反之请选择collection。
	ShardKey ShardKeyType `json:"shard_key,omitempty"`

	// if shard_key is collection, and users want to improve performance when some collections
	// do not have unique key.
	// 对于按collection哈希，如果某些表不具有唯一索引，则可以设置按_id哈希以提高并发度。
	// 用户需要确认该表不会创建唯一索引，一旦检测发现存在唯一索引，则会立刻crash退出。
	// 例如，db1.collection1;db2.collection2，不支持仅指定db
	ShardByObjectIdWhitelist string `json:"shard_by_object_id_whitelist,omitempty"`

	// oplog transmit worker concurrent
	// if the source is sharding, worker number must equal to shard numbers.
	// 内部发送的worker数目，如果机器性能足够，可以提高worker个数。
	Worker int32 `json:"worker,omitempty"`

	// how many writing threads will be used in one worker.
	// 对于目的端是kafka等非direct tunnel，启用多少个序列化线程，必须为"incr_sync.worker"的倍数。
	// 默认为"incr_sync.worker"的值。
	TunnelWriteThread int32 `json:"tunnel_write_thread,omitempty"`

	// set the sync delay just like mongodb secondary slaveDelay parameter. unit second.
	// 设置目的端的延迟，比如延迟源端20分钟，类似MongoDB本身主从同步slaveDelay参数，单位：秒
	// 0表示不启用
	TargetDelay int32 `json:"target_delay,omitempty"`

	// memory queue configuration, plz visit FAQ document to see more details.
	// do not modify these variables if the performance and resource usage can
	// meet your needs.
	// 内部队列的配置参数，如果目前性能足够不建议修改，详细信息参考FAQ。
	WorkerBatchQueueSize    int32 `json:"worker_batch_queue_size,omitempty"`
	AdaptiveBatchingMaxSize int32 `json:"adaptive_batching_max_size,omitempty"`
	FetcherBufferCapacity   int32 `json:"fetcher_buffer_capacity,omitempty"`

	// --- direct tunnel only begin ---
	// if tunnel type is direct, all the below variable should be set
	// 下列参数仅用于tunnel为direct的情况。

	// oplog changes to Insert while Update found non-exist (_id or unique-index)
	// 如果_id不存在在目的库，是否将update语句修改为insert语句。
	ExecutorUpsert bool `json:"executor_upsert,omitempty"`

	// oplog changes to Update while Insert found duplicated key (_id or unique-index)
	// 如果_id存在在目的库，是否将insert语句修改为update语句。
	ExecutorInsertOnDupUpdate bool `json:"executor_insert_on_dup_update,omitempty"`

	// db. write duplicated logs to mongoshake_conflict
	// 如果写入存在冲突，记录冲突的文档。
	ConflictWriteTo string `json:"conflict_write_to,omitempty"`

	// enable majority write in incrmental sync.
	// the performance will degrade if enable.
	// 增量阶段写入端是否启用majority write
	ExecutorMajorityEnable bool `json:"executor_majority_enable,omitempty"`
	// --- direct tunnel only end ---
}

// CreateIndexOption  create index option.
type CreateIndexOption string

const (
	CreateIndexOptionNone       CreateIndexOption = "none"
	CreateIndexOptionForeground CreateIndexOption = "foreground"
	CreateIndexOptionBackground CreateIndexOption = "background"
)

// FullSyncSpec defines the full sync configuration
// 全量同步参数配置
type FullSyncSpec struct {
	// the number of collection concurrence
	// 并发最大拉取的表个数，例如，6表示同一时刻shake最多拉取6个表。
	ReaderCollectionParallel int32 `json:"reader_collection_parallel,omitempty"`
	// the number of document writer thread in each collection.
	// 同一个表内并发写的线程数，例如，8表示对于同一个表，将会有8个写线程进行并发写入。
	ReaderWriteDocumentParallel int32 `json:"reader_write_document_parallel,omitempty"`
	// number of documents in a batch insert in a document concurrence
	// 目的端写入的batch大小，例如，128表示一个线程将会一次聚合128个文档然后再写入。
	ReaderDocumentBatchSize int32 `json:"reader_document_batch_size,omitempty"`
	// max number of fetching thread per table. default is 1
	// 单个表最大拉取的线程数，默认是单线程拉取。需要具备splitVector权限。
	// 注意：对单个表来说，仅支持索引对应的value是同种类型，如果有不同类型请勿启用该配置项！
	ReaderParallelThread int32 `json:"reader_parallel_thread,omitempty"`
	// the parallel query index if set full_sync.reader.parallel_thread. index should only has 1 field.
	// 如果设置了full_sync.reader.parallel_thread，还需要设置该参数，并行拉取所扫描的index，value
	// 必须是同种类型。对于副本集，建议设置_id；对于集群版，建议设置shard_key。key只能有1个field。
	ReaderParallelIndex string `json:"reader_parallel_index,omitempty"`
	// drop the same name of collection in dest mongodb in full synchronization
	// 同步时如果目的库存在，是否先删除目的库再进行同步，true表示先删除再同步，false表示不删除。
	CollectionExistDrop bool `json:"collection_exist_drop,omitempty"`

	// none: do not create indexes.
	// foreground: create indexes when data sync finish in full sync stage.
	// background: create indexes when starting.
	// 全量期间数据同步完毕后，是否需要创建索引，none表示不创建，foreground表示创建前台索引，
	// background表示创建后台索引。
	CreateIndex CreateIndexOption `json:"create_index,omitempty"`
	// convert insert to update when duplicate key found
	// 如果_id存在在目的库，是否将insert语句修改为update语句。
	ExecutorInsertOnDupUpdate bool `json:"executor_insert_on_dup_update,omitempty"`
	// filter orphan document for source type is sharding.
	// 源端是sharding，是否需要过滤orphan文档
	ExecutorFilterOrphanDocument bool `json:"executor_filter_orphan_document,omitempty"`
	// enable majority write in full sync.
	// the performance will degrade if enable.
	// 全量阶段写入端是否启用majority write
	ExecutorMajorityEnable bool `json:"executor_majority_enable,omitempty"`
}

// CheckPointSpec defines the checkpoint info ,used in resuming from break point
// checkpoint存储信息，用于支持断点续传。
type CheckPointSpec struct {
	// context.storage.url is used to mark the checkpoint store database. E.g., mongodb://127.0.0.1:20070
	// if not set, checkpoint will be written into source mongodb(db=mongoshake)
	// checkpoint的具体写入的MongoDB地址，如果不配置，对于副本集和分片集群都将写入源库(db=mongoshake)
	// 2.4版本以后不需要配置为源端cs的地址。
	StorageUrl string `json:"storage_url,omitempty"`
	// checkpoint db's name.
	// checkpoint存储的db的名字
	StorageDb string `json:"storage_db,omitempty"`
	// checkpoint collection's name.
	// checkpoint存储的表的名字，如果启动多个mongoshake拉取同一个源可以修改这个表名以防止冲突。
	StorageCollection string `json:"storage_collection,omitempty"`
	// set if enable ssl
	StorageUrlMongoSslRootCaFile string `json:"storage_url_mongo_ssl_root_ca_file,omitempty"`
	// real checkpoint: the fetching oplog position.
	// pay attention: this is UTC time which is 8 hours latter than CST time. this
	// variable will only be used when checkpoint is not exist.
	// 本次开始拉取的位置，如果checkpoint已经存在（位于上述存储位置）则该参数无效，
	// 如果需要强制该位置开始拉取，需要先删除原来的checkpoint，详见FAQ。
	// 若checkpoint不存在，且该值为1970-01-01T00:00:00Z，则会拉取源端现有的所有oplog。
	// 若checkpoint不存在，且该值不为1970-01-01T00:00:00Z，则会先检查源端oplog最老的时间是否
	// 大于给定的时间，如果是则会直接报错退出。
	StartPosition string `json:"start_position,omitempty"`
}

// TunnelType 通道模式。
// tunnel pipeline type. now we support rpc,file,kafka,mock,direct
type TunnelType string

const (
	TunnelTypeDirect TunnelType = "direct"
	TunnelTypeFile   TunnelType = "file"
	TunnelTypeKafka  TunnelType = "kafka"
	TunnelTypeMock   TunnelType = "mock"
	TunnelTypeRpc    TunnelType = "rpc"
)

// MongoConnectMode 连接模式
// connect mode:
// primary: fetch data from primary.
// secondaryPreferred: fetch data from secondary if has, otherwise primary.(default)
// standalone: fetch data from given 1 node, no matter primary, secondary or hidden. This is only
// support when tunnel type is direct.
// primary表示从主上拉取，secondaryPreferred表示优先从secondary拉取（默认建议值），
// standalone表示从任意单个结点拉取。
type MongoConnectMode string

const (
	MongoConnectModePrimary            MongoConnectMode = "primary"
	MongoConnectModeSecondaryPreferred MongoConnectMode = "secondaryPreferred"
	MongoConnectModeStandalone         MongoConnectMode = "standalone"
)

type TunnelMessageFormat string

const (
	TunnelMessageFormatRaw  TunnelMessageFormat = "raw"
	TunnelMessageFormatJson TunnelMessageFormat = "json"
	TunnelMessageFormatBson TunnelMessageFormat = "bson"
)

// TunnelSpec 通道信息配置
type TunnelSpec struct {
	Type                 TunnelType          `json:"type,omitempty"`
	Address              string              `json:"address,omitempty"`
	Message              TunnelMessageFormat `json:"message,omitempty"`
	KafkaPartitionNumber int32               `json:"kafka_partition_number,omitempty"`
	JsonFormat           string              `json:"json_format,omitempty"`
	MongoSslRootCaFile   string              `json:"mongo_ssl_root_ca_file,omitempty"`
}

type FilterSpec struct {
	// 黑白名单过滤，目前不支持正则，白名单表示通过的namespace，黑名单表示过滤的namespace，
	// 不能同时指定。分号分割不同namespace，每个namespace可以是db，也可以是db.collection。
	NamespaceBlack string `json:"namespace_black,omitempty"`
	NamespaceWhite string `json:"namespace_white,omitempty"`
	// 正常情况下，不建议配置该参数，但对于有些非常特殊的场景，用户可以启用admin，mongoshake等库的同步，
	// 以分号分割，例如：admin;mongoshake。
	PassSpecialDb string `json:"pass_special_db,omitempty"`
	// 是否需要开启DDL同步，true表示开启，源是sharding暂时不支持开启。
	// 如果目的端是sharding，暂时不支持applyOps命令，包括事务。
	DdlEnable bool `json:"ddl_enable,omitempty"`
	// filter oplog gid if enabled.
	// 如果MongoDB启用了gid，但是目的端MongoDB不支持gid导致同步会失败，可以启用gid过滤，将会去掉gid字段。
	// 谨慎建议开启，shake本身性能受损很大。
	OplogGids bool `json:"oplog_gids,omitempty"`
}

// MongoSpec defines  connect source mongodb of MongoShake;
// MongoSpec 源MongoDB连接信息
type MongoSpec struct {
	// connect source mongodb, set username and password if enable authority. Please note: password shouldn't contain '@'.
	// split by comma(,) if use multiple instance in one replica-set. E.g., mongodb://username1:password1@primaryA,secondaryB,secondaryC
	// split by semicolon(;) if sharding enable. E.g., mongodb://username1:password1@primaryA,secondaryB,secondaryC;mongodb://username2:password2@primaryX,secondaryY,secondaryZ
	// 源MongoDB连接串信息，逗号分隔同一个副本集内的结点，分号分隔分片sharding实例，免密模式
	// 可以忽略“username:password@”，注意，密码里面不能含有'@'符号。
	// 举例：
	// 副本集：mongodb://username1:password1@primaryA,secondaryB,secondaryC
	// 分片集：mongodb://username1:password1@primaryA,secondaryB,secondaryC;mongodb://username2:password2@primaryX,secondaryY,secondaryZ
	MongoUrls string `json:"mongo_urls,omitempty"`
	// please fill the source config server url if source mongodb is sharding.
	MongoCsUrl string `json:"mongo_cs_url,omitempty"`
	// please give at least one mongos address if source is sharding.
	MongoSUrl string `json:"mongo_s_url,omitempty"`
	// enable source ssl
	MongoSslRootCaFile string `json:"mongo_ssl_root_ca_file"`
}

// SyncMode 同步模式，all表示全量+增量同步，full表示全量同步，incr表示增量同步。
type SyncMode string

// sync mode: all/full/incr. default is incr.
// all means full synchronization + incremental synchronization.
// full means full synchronization only.
// incr means incremental synchronization only.
const (
	SyncModeAll   SyncMode = "all"
	SyncModeFull  SyncMode = "full"
	SyncModelIncr SyncMode = "incr"
)

type LogLevel string

const (
	LogLevelInfo    LogLevel = "info"
	LogLevelDebug   LogLevel = "debug"
	LogLevelWarning LogLevel = "warning"
	LogLevelError   LogLevel = "error"
)

type LogSpec struct {
	// global log level: debug, info, warning, error. lower level message will be filter
	Level LogLevel `json:"level,omitempty"`
	// log directory. log and pid file will be stored into this file.
	// if not set, default is "./logs/"
	// log和pid文件的目录，如果不设置默认打到当前路径的logs目录。
	Dir string `json:"dir,omitempty"`
	// log file name.
	// log文件名。
	File string `json:"file,omitempty"`
	// log flush enable. If set false, logs may not be print when exit. If
	// set true, performance will be decreased extremely
	// 设置log刷新，false表示包含缓存，如果true那么每条log都会直接刷屏，但对性能有影响；
	// 反之，退出不一定能打印所有的log，调试时建议配置true。
	Flush bool `json:"flush,omitempty"`
}

type VolumeSpec struct {
	// EmptyDir represents a temporary directory that shares a pod's lifetime.
	EmptyDir *corev1.EmptyDirVolumeSource `json:"emptyDir,omitempty"`

	// HostPath represents a pre-existing file or directory on the host machine
	// that is directly exposed to the container.
	HostPath *corev1.HostPathVolumeSource `json:"hostPath,omitempty"`

	// PersistentVolumeClaim represents a reference to a PersistentVolumeClaim.
	// It has the highest level of precedence, followed by HostPath and
	// EmptyDir. And represents the PVC specification.
	PersistentVolumeClaim *corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaim,omitempty"`
}

type ResourceSpecRequirements struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type ResourcesSpec struct {
	Limits   *ResourceSpecRequirements `json:"limits,omitempty"`
	Requests *ResourceSpecRequirements `json:"requests,omitempty"`
}
type LivenessProbeExtended struct {
	corev1.Probe        `json:",inline"`
	StartupDelaySeconds int `json:"startupDelaySeconds,omitempty"`
}

type AppState string

const (
	AppStateInit     AppState = "initializing"
	AppStateRunning  AppState = "running"
	AppStateStopping AppState = "stopping"
	AppStatePaused   AppState = "paused"
	AppStateComplete AppState = "complete"
	AppStateError    AppState = "error"
)

func (l LivenessProbeExtended) CommandHas(flag string) bool {
	if l.Handler.Exec == nil {
		return false
	}

	for _, v := range l.Handler.Exec.Command {
		if v == flag {
			return true
		}
	}

	return false
}

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

type ClusterCondition struct {
	Status             ConditionStatus `json:"status"`
	Type               AppState        `json:"type"`
	LastTransitionTime metav1.Time     `json:"lastTransitionTime,omitempty"`
	Reason             string          `json:"reason,omitempty"`
	Message            string          `json:"message,omitempty"`
}

type HealthyCheckType string

const (
	HealthyCheckTypeOk   HealthyCheckType = "success"
	HealthyCheckTypeFail HealthyCheckType = "failed"
)

const maxStatusesQuantity = 20

func (s *MongoShakeStatus) AddCondition(c ClusterCondition) {
	if len(s.Conditions) == 0 {
		s.Conditions = append(s.Conditions, c)
		return
	}

	if s.Conditions[len(s.Conditions)-1].Type != c.Type {
		s.Conditions = append(s.Conditions, c)
	}

	if len(s.Conditions) > maxStatusesQuantity {
		s.Conditions = s.Conditions[len(s.Conditions)-maxStatusesQuantity:]
	}
}

func init() {
	SchemeBuilder.Register(&MongoShake{}, &MongoShakeList{})
}
