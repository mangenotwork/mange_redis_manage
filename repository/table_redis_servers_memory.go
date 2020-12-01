package models

/*
CREATE TABLE table_redis_servers_memory(
    id INTEGER PRIMARY KEY,
    hid TEXT NOT NULL,
    get_time BIGINT NOT NULL,
   	used_memory BIGINT NOT NULL,
	used_memory_human TEXT NOT NULL,
	used_memory_rss BIGINT NOT NULL,
	used_memory_rss_human TEXT NOT NULL,
	used_memory_peak BIGINT NOT NULL,
	used_memory_peak_human TEXT NOT NULL,
	used_memory_peak_perc TEXT NOT NULL,
	used_memory_overhead BIGINT NOT NULL,
	used_memory_startup BIGINT NOT NULL,
	used_memory_dataset BIGINT NOT NULL,
	used_memory_dataset_perc TEXT NOT NULL,
	allocator_allocated BIGINT NOT NULL,
	allocator_active BIGINT NOT NULL,
	allocator_resident BIGINT NOT NULL,
	total_system_memory BIGINT NOT NULL,
	total_system_memory_human TEXT NOT NULL,
	used_memory_lua BIGINT NOT NULL,
	used_memory_lua_human TEXT NOT NULL,
	used_memory_scripts BIGINT NOT NULL,
	used_memory_scripts_human TEXT NOT NULL,
	number_of_cached_scripts BIGINT NOT NULL,
	maxmemory BIGINT NOT NULL,
	maxmemory_human TEXT NOT NULL,
	maxmemory_policy TEXT NOT NULL,
	allocator_frag_ratio DOUBLE NOT NULL,
	allocator_frag_bytes BIGINT NOT NULL,
	allocator_rss_ratio DOUBLE NOT NULL
	allocator_rss_bytes BIGINT NOT NULL,
	rss_overhead_ratio DOUBLE NOT NULL,
	rss_overhead_bytes BIGINT NOT NULL,
	mem_fragmentation_ratio DOUBLE NOT NULL,
	mem_fragmentation_bytes BIGINT NOT NULL,
	mem_not_counted_for_evict BIGINT NOT NULL,
	mem_replication_backlog BIGINT NOT NULL,
	mem_clients_slaves BIGINT NOT NULL,
	mem_clients_normal BIGINT NOT NULL,
	mem_aof_buffer BIGINT NOT NULL,
	mem_allocator TEXT NOT NULL,
	active_defrag_running BIGINT NOT NULL,
	lazyfree_pending_objects BIGINT NOT NULL
);
*/

type RedisMemoryDB struct {
	ID                     int64   `gorm:"primary_key;column:id" json:"-"`
	Hid                    string  `gorm:"column:hid" json:"host_id"`
	GetTime                int64   `gorm:"column:get_time" json:"get_time"`
	UsedMemory             int64   `gorm:"column:used_memory" json:"used_memory"`                             //由redis分配器分配的内存总量，单位字节
	UsedMemoryHuman        string  `gorm:"column:used_memory_human" json:"used_memory_human"`                 //
	UsedMemoryRss          int64   `gorm:"column:used_memory_rss" json:"used_memory_rss"`                     //从操作系统角度，返回redis已分配内存总量
	UsedMemoryRssHuman     string  `gorm:"column:used_memory_rss_human" json:"used_memory_rss_human"`         //
	UsedMemoryPeak         int64   `gorm:"column:used_memory_peak" json:"used_memory_peak"`                   //redis的内存消耗峰值（以字节为单位）
	UsedMemoryPeakHuman    string  `gorm:"column:used_memory_peak_human" json:"used_memory_peak_human"`       //
	UsedMemoryPeakPerc     string  `gorm:"column:used_memory_peak_perc" json:"used_memory_peak_perc"`         //
	UsedMemoryOverhead     int64   `gorm:"column:used_memory_overhead" json:"used_memory_overhead"`           //
	UsedMemoryStartup      int64   `gorm:"column:used_memory_startup" json:"used_memory_startup"`             //
	UsedMemoryDataset      int64   `gorm:"column:used_memory_dataset" json:"used_memory_dataset"`             //
	UsedMemoryDatasetPerc  string  `gorm:"column:used_memory_dataset_perc" json:"used_memory_dataset_perc"`   //
	AllocatorAllocated     int64   `gorm:"column:allocator_allocated" json:"allocator_allocated"`             //
	AllocatorActive        int64   `gorm:"column:allocator_active" json:"allocator_active"`                   //
	AllocatorResident      int64   `gorm:"column:allocator_resident" json:"allocator_resident"`               //
	TotalSystemMemory      int64   `gorm:"column:total_system_memory" json:"total_system_memory"`             //
	TotalSystemMemoryHuman string  `gorm:"column:total_system_memory_human" json:"total_system_memory_human"` //
	UsedMemoryLua          int64   `gorm:"column:used_memory_lua" json:"used_memory_lua"`                     //lua引擎所使用的内存大小（单位字节）
	UsedMemoryLuaHuman     string  `gorm:"column:used_memory_lua_human" json:"used_memory_lua_human"`         //
	UsedMemoryScripts      int64   `gorm:"column:used_memory_scripts" json:"used_memory_scripts"`             //
	UsedMemoryScriptsHuman string  `gorm:"column:used_memory_scripts_human" json:"used_memory_scripts_human"` //
	NumberOfCachedScripts  int64   `gorm:"column:number_of_cached_scripts" json:"number_of_cached_scripts"`   //
	Maxmemory              int64   `gorm:"column:maxmemory" json:"maxmemory"`                                 //
	MaxmemoryHuman         string  `gorm:"column:maxmemory_human" json:"maxmemory_human"`                     //
	MaxmemoryPolicy        string  `gorm:"column:maxmemory_policy" json:"maxmemory_policy"`                   //
	AllocatorFragRatio     float64 `gorm:"column:allocator_frag_ratio" json:"allocator_frag_ratio"`           //
	AllocatorFragBytes     int64   `gorm:"column:allocator_frag_bytes" json:"allocator_frag_bytes"`           //
	AllocatorRssRatio      float64 `gorm:"column:allocator_rss_ratio" json:"allocator_rss_ratio"`             //
	AllocatorRssBytes      int64   `gorm:"column:allocator_rss_bytes" json:"allocator_rss_bytes"`             //
	RssOverheadRatio       float64 `gorm:"column:rss_overhead_ratio" json:"rss_overhead_ratio"`               //
	RssOverheadBytes       int64   `gorm:"column:rss_overhead_bytes" json:"rss_overhead_bytes"`               //
	MemFragmentationRatio  float64 `gorm:"column:mem_fragmentation_ratio" json:"mem_fragmentation_ratio"`     //used_memory_rss 和 used_memory 之间的比率
	MemFragmentationBytes  int64   `gorm:"column:mem_fragmentation_bytes" json:"mem_fragmentation_bytes"`     //
	MemNotCountedForEvict  int64   `gorm:"column:mem_not_counted_for_evict" json:"mem_not_counted_for_evict"` //
	MemReplicationBacklog  int64   `gorm:"column:mem_replication_backlog" json:"mem_replication_backlog"`     //
	MemClientsSlaves       int64   `gorm:"column:mem_clients_slaves" json:"mem_clients_slaves"`               //
	MemClientsNormal       int64   `gorm:"column:mem_clients_normal" json:"mem_clients_normal"`               //
	MemAofBuffer           int64   `gorm:"column:mem_aof_buffer" json:"mem_aof_buffer"`                       //
	MemAllocator           string  `gorm:"column:mem_allocator" json:"mem_allocator"`                         //编译时指定的redis的内存分配器。越好的分配器内存碎片化率越低，低版本建议升级
	ActiveDefragRunning    int64   `gorm:"column:active_defrag_running" json:"active_defrag_running"`         //
	LazyfreePendingObjects int64   `gorm:"column:lazyfree_pending_objects" json:"lazyfree_pending_objects"`   //
}
