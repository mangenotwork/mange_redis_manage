package install

import (
	_ "database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_redis_manage/common/cache"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/common/sqlitedb"
	_ "github.com/mattn/go-sqlite3"
)

//检查安装中间件
//如果已经安装，无法访问install路由
//如果是未安装才能访问install路由
func CheckInstall() gin.HandlerFunc {
	manlog.Debug("CheckInstall")
	return func(c *gin.Context) {
		if IsInstall() {
			if isInstallUrl(c.Request.URL) {
				c.Abort()
				c.Redirect(http.StatusMovedPermanently, "/")
				return
			}
			c.Next()
			return
		} else {
			if isInstallUrl(c.Request.URL) {
				c.Next()
				return
			}
			c.Abort()
			c.Redirect(http.StatusMovedPermanently, "/install/index")
			return
		}
	}
}

func IsInstall() bool {
	//1.检查缓存的安装标识
	v, isv := cache.Get("install")
	if isv && v == "yes" {
		return true
	}
	//2.如果缓存不存在，则检查数据安装标识
	rows, err := sqlitedb.Query(fmt.Sprintf("SELECT * FROM install limit 1;"))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var isInstall string
			err := rows.Scan(&isInstall)
			if err == nil && isInstall == "yes" {
				cache.SetAlways("install", "yes")
				return true
			}
		}
	}
	return false
}

func isInstallUrl(urls *url.URL) bool {
	urls_list := strings.Split(urls.String(), "/")
	if len(urls_list) > 2 && urls_list[1] == "install" {
		return true
	}
	return false
}

//安装页面
func InstallPG(c *gin.Context) {
	c.HTML(http.StatusOK, "install_index.html", gin.H{
		"title":     "ManGe Redis管理工具v0.1",
		"welcome":   "欢迎使用ManGe Redis管理工具v0.1",
		"thank":     "感谢圆梦时刻提供技术支持!",
		"thank_url": "https://www.ymzy.cn",
		"author":    "ManGe (2912882908@qq.com)",
	})
}

//安装步骤1
func Run(c *gin.Context) {
	admin_name := c.DefaultQuery("admin_name", "")
	admin_password := c.DefaultQuery("admin_password", "")
	secret := c.DefaultQuery("secret", "mangeredismanage2020")
	if secret == "" {
		secret = "mangeredismanage2020"
	}

	if admin_name == "" || admin_password == "" {
		installErrPG(c, "未设置超级管理员")
		return
	}

	//1.创建db目录与文件
	pwd, _ := os.Getwd()

	for _, v := range []string{"db", "cache", "log", "report"} {
		path := pwd + "/" + v
		manlog.Debug(path)
		isdir_db, err := isFileExist(path)
		if err != nil {
			installErrPG(c, fmt.Sprintln("创建%s目录失败,err=%v", path, err))
			return
		}
		if !isdir_db {
			os.Mkdir(path, os.ModePerm)
		}
	}

	for _, f := range []string{"/db/mange_redis_manage.db", "/cache/init.temp"} {
		f_path := pwd + f
		_, err := os.OpenFile(f_path, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			manlog.Error("error:", err)
			installErrPG(c, fmt.Sprintln("创建%s文件失败,err=%v", f_path, err))
			return
		}
	}

	//2.创建表
	all_table := initTable()
	for _, table := range all_table {
		//sqlitedb.Prepare(table)
		sqlitedb.Exec(table)
	}

	//3.填入基础数据
	sqlitedb.Prepare("INSERT INTO install(init) values(?)", "yes")
	sqlitedb.Prepare("INSERT INTO jwt_secret(init) values(?)", secret)
	sqlitedb.Prepare("INSERT INTO table_user(uname,upassword,ugroup) values(?,?,?)", admin_name, admin_password, 0)

	c.Redirect(http.StatusMovedPermanently, "/")
	return
}

//安装错误页面
func installErrPG(c *gin.Context, errinfo string) {
	manlog.Error(errinfo)
	c.HTML(http.StatusOK, "install_err.html", gin.H{
		"title":   "ManGe Redis管理工具v0.1",
		"errinfo": errinfo,
	})
}

//判断文件文件夹是否存在
func isFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func initTable() map[string]string {
	return map[string]string{
		//安装确认表
		"install": `
			CREATE TABLE install(
				init TEXT NOT NULL
			);`,
		//redis连接信息表
		"table_redis_info": `
			CREATE TABLE table_redis_info(
			   id INTEGER PRIMARY KEY,
			   uid INT NOT NULL,
			   conn_name TEXT NOT NULL,
			   conn_host TEXT NOT NULL,
			   conn_port INT NOT NULL,
			   conn_password TEXT NOT NULL,
			   is_ssh BOOLEAN NOT NULL,
			   ssh_url TEXT NOT NULL,
			   ssh_user TEXT NOT NULL,
			   ssh_password TEXT NOT NULL,
			   conn_create BIGINT NOT NULL
			);`,
		//用户表
		"table_user": `
			CREATE TABLE table_user(
			   uid INTEGER PRIMARY KEY,
			   uname TEXT NOT NULL,
			   upassword TEXT NOT NULL,
			   ugroup INT NOT NULL
			);`,
		//jwt secret 表
		"jwt_secret": `
			CREATE TABLE jwt_secret(
				init TEXT NOT NULL
			);`,
		//redis 服务客户端连接信息表
		`table_redis_servers_clients`: `
			CREATE TABLE table_redis_servers_clients(
			   id INTEGER PRIMARY KEY,
			   hid TEXT NOT NULL,
			   get_time BIGINT NOT NULL,
			   connected_clients BIGINT NOT NULL,
			   client_recent_max_input_buffer BIGINT NOT NULL,
			   client_recent_max_output_buffer BIGINT NOT NULL,
			   clocked_clients BIGINT NOT NULL
			);`,
		//redis 服务集群信息表
		`table_redis_servers_cluster`: `
			CREATE TABLE table_redis_servers_cluster(
			   id INTEGER PRIMARY KEY,
			   hid TEXT NOT NULL,
			   get_time BIGINT NOT NULL,
			   cluster_enabled TEXT NOT NULL
			);`,
		//redis 服务cpu信息表
		`table_redis_servers_cpu`: `
			CREATE TABLE table_redis_servers_cpu(
			   id INTEGER PRIMARY KEY,
			   hid TEXT NOT NULL,
			   get_time BIGINT NOT NULL,
			   used_cpu_sys DOUBLE NOT NULL,
			   used_cpu_user DOUBLE NOT NULL,
			   used_cpu_sys_children DOUBLE NOT NULL,
			   used_cpu_user_children DOUBLE NOT NULL
			);`,
		//redis 服务信息表
		`table_redis_servers_infos`: `
			CREATE TABLE table_redis_servers_infos(
			    id INTEGER PRIMARY KEY,
			    hid TEXT NOT NULL,
			    get_time BIGINT NOT NULL,
			    redis_version TEXT NOT NULL,
				redis_git_sha1 TEXT NOT NULL,
				redis_git_dirty TEXT NOT NULL,
				redis_build_id TEXT NOT NULL,
				redis_mode TEXT NOT NULL,
				os TEXT NOT NULL,
				arch_bits TEXT NOT NULL,
				multiplexing_api TEXT NOT NULL,
				atomicvar_api TEXT NOT NULL,
				gcc_version TEXT NOT NULL,
				process_id TEXT NOT NULL,
				run_id TEXT NOT NULL,
				tcp_port BIGINT NOT NULL,
				uptime_in_seconds BIGINT NOT NULL,
				uptime_in_days BIGINT NOT NULL,
				hz BIGINT NOT NULL,
				configured_hz BIGINT NOT NULL,
				lru_clock BIGINT NOT NULL,
				executable TEXT NOT NULL,
				config_file TEXT NOT NULL
			);`,
		//redis 服务db key信息表
		`table_redis_servers_keyspace`: `
			CREATE TABLE table_redis_servers_keyspace(
			    id INTEGER PRIMARY KEY,
			    hid TEXT NOT NULL,
			    get_time BIGINT NOT NULL,
			   	db_id BIGINT NOT NULL,
				keys_count BIGINT NOT NULL,
				expires BIGINT NOT NULL,
				avgttl BIGINT NOT NULL
			);`,
		//redis 服务内存信息表
		`table_redis_servers_memory`: `
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
				allocator_rss_ratio DOUBLE NOT NULL,
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
			);`,
		//redis 服务persistence 信息表
		`table_redis_servers_persistence`: `
			CREATE TABLE table_redis_servers_persistence(
			    id INTEGER PRIMARY KEY,
			    hid TEXT NOT NULL,
			    get_time BIGINT NOT NULL,
			   	loading BIGINT NOT NULL,
				rdb_changes_since_last_save BIGINT NOT NULL,
				rdb_bgsave_in_progress BIGINT NOT NULL,
				rdb_last_save_time BIGINT NOT NULL,
				rdb_last_bgsave_status TEXT NOT NULL,
				rdb_last_bgsave_time_sec BIGINT NOT NULL,
				rdb_current_bgsave_time_sec BIGINT NOT NULL,
				rdb_last_cow_size BIGINT NOT NULL,
				aof_enabled BIGINT NOT NULL,
				aof_rewrite_in_progress BIGINT NOT NULL,
				aof_rewrite_scheduled BIGINT NOT NULL,
				aof_last_rewrite_time_sec BIGINT NOT NULL,
				aof_current_rewrite_time_sec BIGINT NOT NULL,
				aof_last_bgrewrite_status TEXT NOT NULL,
				aof_last_write_status TEXT NOT NULL,
				aof_last_cow_size BIGINT NOT NULL,
				aof_current_size BIGINT NOT NULL,
				aof_base_size BIGINT NOT NULL,
				aof_pending_rewrite BIGINT NOT NULL,
				aof_buffer_length BIGINT NOT NULL,
				aof_rewrite_buffer_length BIGINT NOT NULL,
				aof_pending_bio_fsync BIGINT NOT NULL,
				aof_delayed_fsync BIGINT NOT NULL
			);`,
		//redis 服务replication 信息表
		`table_redis_servers_replication`: `
			CREATE TABLE table_redis_servers_replication(
			    id INTEGER PRIMARY KEY,
			    hid TEXT NOT NULL,
			    get_time BIGINT NOT NULL,
			    role_value TEXT NOT NULL,
			    connected_slaves TEXT NOT NULL,
			    master_replid TEXT NOT NULL,
			    master_replid2 TEXT NOT NULL,
			    master_repl_offset TEXT NOT NULL,
			    second_repl_offset TEXT NOT NULL,
			    repl_backlog_active TEXT NOT NULL,
			    repl_backlog_size BIGINT NOT NULL,
			    repl_backlog_first_byteoffset BIGINT NOT NULL,
			    repl_backlog_histlen BIGINT NOT NULL
			);`,
		//redis 服务状态信息表
		`table_redis_servers_stats`: `
			CREATE TABLE table_redis_servers_stats(
			    id INTEGER PRIMARY KEY,
			    hid TEXT NOT NULL,
			    get_time BIGINT NOT NULL,
				total_connections_received BIGINT NOT NULL,
				total_commands_processed BIGINT NOT NULL,
				instantaneous_ops_per_sec BIGINT NOT NULL,
				total_net_input_bytes BIGINT NOT NULL,
				total_net_output_bytes BIGINT NOT NULL,
				instantaneous_input_kbps DOUBLE NOT NULL,
				instantaneous_output_kbps DOUBLE NOT NULL,
				rejected_connections BIGINT NOT NULL,
				sync_full BIGINT NOT NULL,
				sync_partial_ok BIGINT NOT NULL,
				sync_partial_err BIGINT NOT NULL,
				expired_keys BIGINT NOT NULL,
				expired_stale_perc TEXT NOT NULL,
				expired_time_cap_reached_count BIGINT NOT NULL,
				evicted_keys BIGINT NOT NULL,
				keyspace_hits BIGINT NOT NULL,
				keyspace_misses BIGINT NOT NULL,
				pubsub_channels BIGINT NOT NULL,
				pubsub_patterns BIGINT NOT NULL,
				latest_fork_usec BIGINT NOT NULL,
				migrate_cached_sockets BIGINT NOT NULL,
				slave_expires_tracked_keys BIGINT NOT NULL,
				active_defrag_hits BIGINT NOT NULL,
				active_defrag_misses BIGINT NOT NULL,
				active_defrag_key_hits BIGINT NOT NULL,
				active_defrag_key_misses BIGINT NOT NULL
			);`,
		//redis 服务慢日志
		`table_redis_slowlog`: `
			CREATE TABLE table_redis_slowlog(
			   id INTEGER PRIMARY KEY,
			   hid TEXT NOT NULL,
			   get_time BIGINT NOT NULL,
			   only_id BIGINT NOT NULL,
			   time BIGINT NOT NULL,
			   duration BIGINT NOT NULL,
			   cmd TEXT NOT NULL,
			   client TEXT NOT NULL,
			);`,
	}
}
