package v1

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

var (
	defaultCollectorId                   = "mongoshake"
	defaultFullSyncHttpPort        int32 = 9101
	defaultIncrSyncHttpPort        int32 = 9100
	defaultSystemProfilePort       int32 = 9200
	defaultLogFileName                   = "collector.log"
	defaultStorageDbName                 = "mongoshake"
	defaultStorageCollectionName         = "ckpt_default"
	defaultCheckPointStartPosition       = "1970-01-01T00:00:00Z"
	defaultImagePullPolicy               = corev1.PullAlways
)

func (cr *MongoShake) CheckNSetDefaults(logger logr.Logger) error {
	// 如果Configuration不为空，则直接使用该配置信息，不使用CollectorSpec传参信息
	if cr.Spec.Configuration != "" {
		return nil
	}
	if cr.Spec.Collector == nil {
		return errors.New("Required value for spec.collector or spec.configuration")
	}
	if cr.Spec.Collector.Mode == "" {
		return errors.New("Required value for spec.collector.mode")
	}
	if cr.Spec.Collector.Mongo == nil {
		return errors.New("Required value for spec.collector.mongo")
	}
	if cr.Spec.Collector.Tunnel == nil {
		return errors.New("Required value for spec.collector.tunnel")
	}
	if cr.Spec.Image == "" {
		return errors.New("Required value for spec.image")
	}
	if cr.Spec.ImagePullPolicy == "" {
		cr.Spec.ImagePullPolicy = defaultImagePullPolicy
	}

	if cr.Spec.Collector.ConfVersion == 0 {
		cr.Spec.Collector.ConfVersion = 10
	}
	if cr.Spec.Collector.Id == "" {
		cr.Spec.Collector.Id = defaultCollectorId
	}
	if cr.Spec.Collector.FullSyncHttpPort == 0 {
		cr.Spec.Collector.FullSyncHttpPort = defaultFullSyncHttpPort
	}
	if cr.Spec.Collector.IncrSyncHttpPort == 0 {
		cr.Spec.Collector.IncrSyncHttpPort = defaultIncrSyncHttpPort
	}
	if cr.Spec.Collector.SystemProfilePort == 0 {
		cr.Spec.Collector.SystemProfilePort = defaultSystemProfilePort
	}
	// 日志配置检查
	if cr.Spec.Collector.Log == nil {
		cr.Spec.Collector.Log = &LogSpec{}
	}
	if cr.Spec.Collector.Log.Level == "" {
		cr.Spec.Collector.Log.Level = LogLevelInfo
	}
	if cr.Spec.Collector.Log.File == "" {
		cr.Spec.Collector.Log.File = defaultLogFileName
	}
	if cr.Spec.Collector.Tunnel.Type == "" {
		cr.Spec.Collector.Tunnel.Type = TunnelTypeDirect
	}
	if cr.Spec.Collector.MongoConnect == "" {
		cr.Spec.Collector.MongoConnect = MongoConnectModeSecondaryPreferred
	}

	if cr.Spec.Collector.Filter == nil {
		cr.Spec.Collector.Filter = &FilterSpec{}
	}

	// checkpoint断点续传配置检查
	if cr.Spec.Collector.CheckPoint == nil {
		cr.Spec.Collector.CheckPoint = &CheckPointSpec{}
	}
	if cr.Spec.Collector.CheckPoint.StorageDb == "" {
		cr.Spec.Collector.CheckPoint.StorageDb = defaultStorageDbName
	}
	if cr.Spec.Collector.CheckPoint.StorageCollection == "" {
		cr.Spec.Collector.CheckPoint.StorageCollection = defaultStorageCollectionName
	}
	if cr.Spec.Collector.CheckPoint.StartPosition == "" {
		cr.Spec.Collector.CheckPoint.StartPosition = defaultCheckPointStartPosition
	}

	// 全量同步配置检查
	if cr.Spec.Collector.FullSync == nil {
		cr.Spec.Collector.FullSync = &FullSyncSpec{}
	}
	if cr.Spec.Collector.FullSync.ReaderCollectionParallel == 0 {
		cr.Spec.Collector.FullSync.ReaderCollectionParallel = 6
	}
	if cr.Spec.Collector.FullSync.ReaderWriteDocumentParallel == 0 {
		cr.Spec.Collector.FullSync.ReaderWriteDocumentParallel = 8
	}
	if cr.Spec.Collector.FullSync.ReaderDocumentBatchSize == 0 {
		cr.Spec.Collector.FullSync.ReaderDocumentBatchSize = 128
	}
	if cr.Spec.Collector.FullSync.ReaderParallelThread == 0 {
		cr.Spec.Collector.FullSync.ReaderParallelThread = 1
	}
	if cr.Spec.Collector.FullSync.ReaderParallelIndex == "" {
		cr.Spec.Collector.FullSync.ReaderParallelIndex = "_id"
	}
	if cr.Spec.Collector.FullSync.CreateIndex == "" {
		cr.Spec.Collector.FullSync.CreateIndex = CreateIndexOptionNone
	}

	if cr.Spec.Collector.IncrSync == nil {
		cr.Spec.Collector.IncrSync = &IncrSyncSpec{}
	}
	if cr.Spec.Collector.IncrSync.MongoFetchMethod == "" {
		cr.Spec.Collector.IncrSync.MongoFetchMethod = FetchMethodOplog
	}
	if cr.Spec.Collector.IncrSync.ShardKey == "" {
		cr.Spec.Collector.IncrSync.ShardKey = ShardKeyCollection
	}
	if cr.Spec.Collector.IncrSync.Worker == 0 {
		cr.Spec.Collector.IncrSync.Worker = 8
	}

	if cr.Spec.Collector.Tunnel.Type != TunnelTypeDirect && cr.Spec.Collector.IncrSync.TunnelWriteThread == 0 {
		cr.Spec.Collector.IncrSync.TunnelWriteThread = cr.Spec.Collector.IncrSync.Worker * 1
	}

	if cr.Spec.Collector.IncrSync.WorkerBatchQueueSize == 0 {
		cr.Spec.Collector.IncrSync.WorkerBatchQueueSize = 64
	}
	if cr.Spec.Collector.IncrSync.AdaptiveBatchingMaxSize == 0 {
		cr.Spec.Collector.IncrSync.WorkerBatchQueueSize = 1024
	}
	if cr.Spec.Collector.IncrSync.FetcherBufferCapacity == 0 {
		cr.Spec.Collector.IncrSync.FetcherBufferCapacity = 256
	}
	if cr.Spec.Collector.IncrSync.ConflictWriteTo == "" {
		cr.Spec.Collector.IncrSync.ConflictWriteTo = "none"
	}
	if cr.Spec.Configuration == "" {
		cr.Spec.Configuration = completedConfig(cr.Spec.Collector)
	}
	return nil
}

func completedConfig(collector *CollectorSpec) string {
	var config strings.Builder
	config.WriteString(fmt.Sprintf("conf.version = %d\n", collector.ConfVersion))
	config.WriteString(fmt.Sprintf("id = %s\n", collector.Id))
	config.WriteString(fmt.Sprintf("master_quorum = %t\n", collector.MasterQuorum))
	config.WriteString(fmt.Sprintf("full_sync.http_port = %d\n", collector.FullSyncHttpPort))
	config.WriteString(fmt.Sprintf("incr_sync.http_port = %d\n", collector.IncrSyncHttpPort))
	config.WriteString(fmt.Sprintf("system_profile_port = %d\n", collector.SystemProfilePort))
	config.WriteString(fmt.Sprintf("log.level = %s\n", collector.Log.Level))
	config.WriteString(fmt.Sprintf("log.dir = %s\n", collector.Log.Dir))
	config.WriteString(fmt.Sprintf("log.file = %s\n", collector.Log.File))
	config.WriteString(fmt.Sprintf("log.flush = %t\n", collector.Log.Flush))
	config.WriteString(fmt.Sprintf("sync_mode = %s\n", collector.Mode))
	config.WriteString(fmt.Sprintf("mongo_urls = %s\n", collector.Mongo.MongoUrls))
	config.WriteString(fmt.Sprintf("mongo_cs_url = %s\n", collector.Mongo.MongoCsUrl))
	config.WriteString(fmt.Sprintf("mongo_s_url = %s\n", collector.Mongo.MongoSUrl))
	config.WriteString(fmt.Sprintf("mongo_ssl_root_ca_file = %s\n", collector.Mongo.MongoSslRootCaFile))
	config.WriteString(fmt.Sprintf("tunnel = %s\n", collector.Tunnel.Type))
	config.WriteString(fmt.Sprintf("tunnel.address = %s\n", collector.Tunnel.Address))
	config.WriteString(fmt.Sprintf("tunnel.message = %s\n", collector.Tunnel.Message))
	config.WriteString(fmt.Sprintf("tunnel.kafka.partition_number = %d\n", collector.Tunnel.KafkaPartitionNumber))
	config.WriteString(fmt.Sprintf("tunnel.json.format = %s\n", collector.Tunnel.JsonFormat))
	config.WriteString(fmt.Sprintf("tunnel.mongo_ssl_root_ca_file = %s\n", collector.Tunnel.MongoSslRootCaFile))
	config.WriteString(fmt.Sprintf("mongo_connect_mode = %d\n", collector.ConfVersion))
	config.WriteString(fmt.Sprintf("filter.namespace.black = %s\n", collector.Filter.NamespaceBlack))
	config.WriteString(fmt.Sprintf("filter.namespace.white = %s\n", collector.Filter.NamespaceWhite))
	config.WriteString(fmt.Sprintf("filter.pass.special.db = %s\n", collector.Filter.PassSpecialDb))
	config.WriteString(fmt.Sprintf("filter.ddl_enable = %t\n", collector.Filter.DdlEnable))
	config.WriteString(fmt.Sprintf("filter.oplog.gids = %t\n", collector.Filter.OplogGids))
	config.WriteString(fmt.Sprintf("checkpoint.storage.url = %s\n", collector.CheckPoint.StorageUrl))
	config.WriteString(fmt.Sprintf("checkpoint.storage.db = %s\n", collector.CheckPoint.StorageDb))
	config.WriteString(fmt.Sprintf("checkpoint.storage.collection = %s\n", collector.CheckPoint.StorageCollection))
	config.WriteString(fmt.Sprintf("checkpoint.start_position = %s\n", collector.CheckPoint.StartPosition))
	config.WriteString(fmt.Sprintf("transform.namespace = %s\n", collector.TransformNamespace))
	config.WriteString(fmt.Sprintf("full_sync.reader.collection_parallel = %d\n", collector.FullSync.ReaderCollectionParallel))
	config.WriteString(fmt.Sprintf("full_sync.reader.write_document_parallel = %d\n", collector.FullSync.ReaderWriteDocumentParallel))
	config.WriteString(fmt.Sprintf("full_sync.reader.document_batch_size = %d\n", collector.FullSync.ReaderDocumentBatchSize))
	config.WriteString(fmt.Sprintf("full_sync.reader.parallel_thread = %d\n", collector.FullSync.ReaderParallelThread))
	config.WriteString(fmt.Sprintf("full_sync.reader.parallel_index = %s\n", collector.FullSync.ReaderParallelIndex))
	config.WriteString(fmt.Sprintf("full_sync.collection_exist_drop = %t\n", collector.FullSync.CollectionExistDrop))
	config.WriteString(fmt.Sprintf("full_sync.create_index = %s\n", collector.FullSync.CreateIndex))
	config.WriteString(fmt.Sprintf("full_sync.executor.insert_on_dup_update = %t\n", collector.FullSync.ExecutorInsertOnDupUpdate))
	config.WriteString(fmt.Sprintf("full_sync.executor.filter.orphan_document = %t\n", collector.FullSync.ExecutorFilterOrphanDocument))
	config.WriteString(fmt.Sprintf("full_sync.executor.majority_enable = %t\n", collector.FullSync.ExecutorMajorityEnable))
	config.WriteString(fmt.Sprintf("incr_sync.mongo_fetch_method = %s\n", collector.IncrSync.MongoFetchMethod))
	config.WriteString(fmt.Sprintf("incr_sync.change_stream.watch_full_document = %t\n", collector.IncrSync.ChangeStreamWatchFullDocument))
	config.WriteString(fmt.Sprintf("incr_sync.oplog.gids = %s\n", collector.IncrSync.OplogGids))
	config.WriteString(fmt.Sprintf("incr_sync.shard_key = %s\n", collector.IncrSync.ShardKey))
	config.WriteString(fmt.Sprintf("incr_sync.shard_by_object_id_whitelist = %s\n", collector.IncrSync.ShardByObjectIdWhitelist))
	config.WriteString(fmt.Sprintf("incr_sync.worker = %d\n", collector.IncrSync.Worker))
	config.WriteString(fmt.Sprintf("incr_sync.tunnel.write_thread = %d\n", collector.IncrSync.TunnelWriteThread))
	config.WriteString(fmt.Sprintf("incr_sync.target_delay = %d\n", collector.IncrSync.TargetDelay))
	config.WriteString(fmt.Sprintf("incr_sync.worker.batch_queue_size = %d\n", collector.IncrSync.WorkerBatchQueueSize))
	config.WriteString(fmt.Sprintf("incr_sync.fetcher.buffer_capacity = %d\n", collector.IncrSync.FetcherBufferCapacity))
	config.WriteString(fmt.Sprintf("incr_sync.executor.upsert = %t\n", collector.IncrSync.ExecutorUpsert))
	config.WriteString(fmt.Sprintf("incr_sync.executor.insert_on_dup_update = %t\n", collector.IncrSync.ExecutorInsertOnDupUpdate))
	config.WriteString(fmt.Sprintf("incr_sync.conflict_write_to = %s\n", collector.IncrSync.ConflictWriteTo))
	config.WriteString(fmt.Sprintf("incr_sync.executor.majority_enable = %t\n", collector.IncrSync.ExecutorMajorityEnable))
	config.WriteString(fmt.Sprintf("special.source.db.flag = %s\n", collector.SpecialSourceDbFlag))

	return config.String()
}
