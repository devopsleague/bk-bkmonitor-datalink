// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package http

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/bkapi"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/consul"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/featureFlag"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/influxdb"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/influxdb/decoder"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/internal/function"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/internal/json"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/log"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/metadata"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/mock"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/query/promql"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/query/structured"
	redisUtil "github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/redis"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/tsdb/elasticsearch"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/tsdb/redis"
)

func TestQueryTsWithDoris(t *testing.T) {
	ctx := metadata.InitHashID(context.Background())

	//viper.Set(bkapi.BkAPIAddressConfigPath, mock.BkSQLUrlDomain)

	spaceUid := influxdb.SpaceUid
	tableID := influxdb.ResultTableDoris

	mock.Init()
	promql.MockEngine()
	influxdb.MockSpaceRouter(ctx)

	defaultStart := time.UnixMilli(1744662513046)
	defaultEnd := time.UnixMilli(1744684113046)

	mock.BkSQL.Set(map[string]any{
		"SHOW CREATE TABLE `2_bklog_bkunify_query_doris`.doris": `{"result":true,"message":"成功","code":"00","data":{"result_table_scan_range":{},"cluster":"doris-test","totalRecords":18,"external_api_call_time_mills":{"bkbase_auth_api":43,"bkbase_meta_api":0,"bkbase_apigw_api":33},"resource_use_summary":{"cpu_time_mills":0,"memory_bytes":0,"processed_bytes":0,"processed_rows":0},"source":"","list":[{"Field":"thedate","Type":"int","Null":"NO","Key":"YES","Default":null,"Extra":""},{"Field":"dteventtimestamp","Type":"bigint","Null":"NO","Key":"YES","Default":null,"Extra":""},{"Field":"dteventtime","Type":"varchar(32)","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"localtime","Type":"varchar(32)","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"__shard_key__","Type":"bigint","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"__ext","Type":"variant","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"cloudid","Type":"double","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"file","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"gseindex","Type":"double","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"iterationindex","Type":"double","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"level","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"log","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"message","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"path","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"report_time","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"serverip","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"time","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"trace_id","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"}],"stage_elapsed_time_mills":{"check_query_syntax":0,"query_db":5,"get_query_driver":0,"match_query_forbidden_config":0,"convert_query_statement":2,"connect_db":45,"match_query_routing_rule":0,"check_permission":43,"check_query_semantic":0,"pick_valid_storage":1},"select_fields_order":["Field","Type","Null","Key","Default","Extra"],"sql":"SHOW COLUMNS FROM mapleleaf_2.bklog_bkunify_query_doris_2","total_record_size":11776,"timetaken":0.096,"result_schema":[{"field_type":"string","field_name":"Field","field_alias":"Field","field_index":0},{"field_type":"string","field_name":"Type","field_alias":"Type","field_index":1},{"field_type":"string","field_name":"Null","field_alias":"Null","field_index":2},{"field_type":"string","field_name":"Key","field_alias":"Key","field_index":3},{"field_type":"string","field_name":"Default","field_alias":"Default","field_index":4},{"field_type":"string","field_name":"Extra","field_alias":"Extra","field_index":5}],"bksql_call_elapsed_time":0,"device":"doris","result_table_ids":["2_bklog_bkunify_query_doris"]},"errors":null,"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000"}`,

		// 查询 1 条原始数据，按照字段正向排序
		"SELECT *, `gseIndex` AS `_value_`, `dtEventTimeStamp` AS `_timestamp_` FROM `2_bklog_bkunify_query_doris`.doris WHERE `dtEventTimeStamp` >= 1744662180000 AND `dtEventTimeStamp` <= 1744684113000 AND `thedate` = '20250415' ORDER BY `_value_` ASC LIMIT 1": `{"result":true,"message":"成功","code":"00","data":{"result_table_scan_range":{"2_bklog_bkunify_query_doris":{"start":"2025041500","end":"2025041523"}},"cluster":"doris-test","totalRecords":1,"external_api_call_time_mills":{"bkbase_auth_api":12,"bkbase_meta_api":0,"bkbase_apigw_api":0},"resource_use_summary":{"cpu_time_mills":0,"memory_bytes":0,"processed_bytes":0,"processed_rows":0},"source":"","list":[{"thedate":20250415,"dteventtimestamp":1744662180000,"dteventtime":null,"localtime":null,"__shard_key__":29077703953,"__ext":"{\"container_id\":\"375597ee636fd5d53cb7b0958823d9ba6534bd24cd698e485c41ca2f01b78ed2\",\"container_image\":\"sha256:3a0506f06f1467e93c3a582203aac1a7501e77091572ec9612ddeee4a4dbbdb8\",\"container_name\":\"unify-query\",\"io_kubernetes_pod\":\"bk-datalink-unify-query-6df8bcc4c9-rk4sc\",\"io_kubernetes_pod_ip\":\"127.0.0.1\",\"io_kubernetes_pod_namespace\":\"blueking\",\"io_kubernetes_pod_uid\":\"558c5b17-b221-47e1-aa66-036cc9b43e2a\",\"io_kubernetes_workload_name\":\"bk-datalink-unify-query-6df8bcc4c9\",\"io_kubernetes_workload_type\":\"ReplicaSet\"}","cloudid":0.0,"file":"http/handler.go:320","gseindex":2450131.0,"iterationindex":19.0,"level":"info","log":"2025-04-14T20:22:59.982Z\tinfo\thttp/handler.go:320\t[5108397435e997364f8dc1251533e65e] header: map[Accept:[*/*] Accept-Encoding:[gzip, deflate] Bk-Query-Source:[strategy:9155] Connection:[keep-alive] Content-Length:[863] Content-Type:[application/json] Traceparent:[00-5108397435e997364f8dc1251533e65e-ca18e72c0f0eafd4-00] User-Agent:[python-requests/2.31.0] X-Bk-Scope-Space-Uid:[bkcc__2]], body: {\"space_uid\":\"bkcc__2\",\"query_list\":[{\"field_name\":\"bscp_config_consume_total_file_change_count\",\"is_regexp\":false,\"function\":[{\"method\":\"mean\",\"without\":false,\"dimensions\":[\"app\",\"biz\",\"clientType\"]}],\"time_aggregation\":{\"function\":\"increase\",\"window\":\"1m\"},\"is_dom_sampled\":false,\"reference_name\":\"a\",\"dimensions\":[\"app\",\"biz\",\"clientType\"],\"conditions\":{\"field_list\":[{\"field_name\":\"releaseChangeStatus\",\"value\":[\"Failed\"],\"op\":\"contains\"},{\"field_name\":\"bcs_cluster_id\",\"value\":[\"BCS-K8S-00000\"],\"op\":\"contains\"}],\"condition_list\":[\"and\"]},\"keep_columns\":[\"_time\",\"a\",\"app\",\"biz\",\"clientType\"],\"query_string\":\"\"}],\"metric_merge\":\"a\",\"start_time\":\"1744660260\",\"end_time\":\"1744662120\",\"step\":\"60s\",\"timezone\":\"Asia/Shanghai\",\"instant\":false}","message":" header: map[Accept:[*/*] Accept-Encoding:[gzip, deflate] Bk-Query-Source:[strategy:9155] Connection:[keep-alive] Content-Length:[863] Content-Type:[application/json] Traceparent:[00-5108397435e997364f8dc1251533e65e-ca18e72c0f0eafd4-00] User-Agent:[python-requests/2.31.0] X-Bk-Scope-Space-Uid:[bkcc__2]], body: {\"space_uid\":\"bkcc__2\",\"query_list\":[{\"field_name\":\"bscp_config_consume_total_file_change_count\",\"is_regexp\":false,\"function\":[{\"method\":\"mean\",\"without\":false,\"dimensions\":[\"app\",\"biz\",\"clientType\"]}],\"time_aggregation\":{\"function\":\"increase\",\"window\":\"1m\"},\"is_dom_sampled\":false,\"reference_name\":\"a\",\"dimensions\":[\"app\",\"biz\",\"clientType\"],\"conditions\":{\"field_list\":[{\"field_name\":\"releaseChangeStatus\",\"value\":[\"Failed\"],\"op\":\"contains\"},{\"field_name\":\"bcs_cluster_id\",\"value\":[\"BCS-K8S-00000\"],\"op\":\"contains\"}],\"condition_list\":[\"and\"]},\"keep_columns\":[\"_time\",\"a\",\"app\",\"biz\",\"clientType\"],\"query_string\":\"\"}],\"metric_merge\":\"a\",\"start_time\":\"1744660260\",\"end_time\":\"1744662120\",\"step\":\"60s\",\"timezone\":\"Asia/Shanghai\",\"instant\":false}","path":"/var/host/data/bcs/lib/docker/containers/375597ee636fd5d53cb7b0958823d9ba6534bd24cd698e485c41ca2f01b78ed2/375597ee636fd5d53cb7b0958823d9ba6534bd24cd698e485c41ca2f01b78ed2-json.log","report_time":"2025-04-14T20:22:59.982Z","serverip":"127.0.0.1","time":"1744662180000","trace_id":"5108397435e997364f8dc1251533e65e","_value_":2450131.0,"_timestamp_":1744662180000}],"stage_elapsed_time_mills":{"check_query_syntax":1,"query_db":182,"get_query_driver":0,"match_query_forbidden_config":0,"convert_query_statement":2,"connect_db":56,"match_query_routing_rule":0,"check_permission":13,"check_query_semantic":0,"pick_valid_storage":1},"select_fields_order":["thedate","dteventtimestamp","dteventtime","localtime","__shard_key__","__ext","cloudid","file","gseindex","iterationindex","level","log","message","path","report_time","serverip","time","trace_id","_value_","_timestamp_"],"total_record_size":8856,"timetaken":0.255,"result_schema":[{"field_type":"int","field_name":"__c0","field_alias":"thedate","field_index":0},{"field_type":"long","field_name":"__c1","field_alias":"dteventtimestamp","field_index":1},{"field_type":"string","field_name":"__c2","field_alias":"dteventtime","field_index":2},{"field_type":"string","field_name":"__c3","field_alias":"localtime","field_index":3},{"field_type":"long","field_name":"__c4","field_alias":"__shard_key__","field_index":4},{"field_type":"string","field_name":"__c5","field_alias":"__ext","field_index":5},{"field_type":"double","field_name":"__c6","field_alias":"cloudid","field_index":6},{"field_type":"string","field_name":"__c7","field_alias":"file","field_index":7},{"field_type":"double","field_name":"__c8","field_alias":"gseindex","field_index":8},{"field_type":"double","field_name":"__c9","field_alias":"iterationindex","field_index":9},{"field_type":"string","field_name":"__c10","field_alias":"level","field_index":10},{"field_type":"string","field_name":"__c11","field_alias":"log","field_index":11},{"field_type":"string","field_name":"__c12","field_alias":"message","field_index":12},{"field_type":"string","field_name":"__c13","field_alias":"path","field_index":13},{"field_type":"string","field_name":"__c14","field_alias":"report_time","field_index":14},{"field_type":"string","field_name":"__c15","field_alias":"serverip","field_index":15},{"field_type":"string","field_name":"__c16","field_alias":"time","field_index":16},{"field_type":"string","field_name":"__c17","field_alias":"trace_id","field_index":17},{"field_type":"double","field_name":"__c18","field_alias":"_value_","field_index":18},{"field_type":"long","field_name":"__c19","field_alias":"_timestamp_","field_index":19}],"bksql_call_elapsed_time":0,"device":"doris","result_table_ids":["2_bklog_bkunify_query_doris"]},"errors":null,"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000"}`,

		// 根据维度 __ext.container_name 进行 count 聚合，同时用值正向排序
		"SELECT CAST(__ext[\"container_name\"] AS STRING) AS `__ext__bk_46__container_name`, COUNT(`gseIndex`) AS `_value_`, CAST(FLOOR(dtEventTimeStamp / 30000) AS INT) * 30000  AS `_timestamp_` FROM `2_bklog_bkunify_query_doris`.doris WHERE `dtEventTimeStamp` >= 1744662509999 AND `dtEventTimeStamp` <= 1744684142999 AND `thedate` = '20250415' GROUP BY __ext__bk_46__container_name, _timestamp_ ORDER BY `_timestamp_` ASC, `_value_` ASC LIMIT 2000005": `{"result":true,"message":"成功","code":"00","data":{"result_table_scan_range":{"2_bklog_bkunify_query_doris":{"start":"2025041500","end":"2025041523"}},"cluster":"doris-test","totalRecords":722,"external_api_call_time_mills":{"bkbase_auth_api":72,"bkbase_meta_api":6,"bkbase_apigw_api":28},"resource_use_summary":{"cpu_time_mills":0,"memory_bytes":0,"processed_bytes":0,"processed_rows":0},"source":"","list":[{"__ext__bk_46__container_name":"unify-query","_value_":3684,"_timestamp_":1744662510000},{"__ext__bk_46__container_name":"unify-query","_value_":4012,"_timestamp_":1744662540000},{"__ext__bk_46__container_name":"unify-query","_value_":3671,"_timestamp_":1744662570000},{"__ext__bk_46__container_name":"unify-query","_value_":17092,"_timestamp_":1744662600000},{"__ext__bk_46__container_name":"unify-query","_value_":12881,"_timestamp_":1744662630000},{"__ext__bk_46__container_name":"unify-query","_value_":5902,"_timestamp_":1744662660000},{"__ext__bk_46__container_name":"unify-query","_value_":10443,"_timestamp_":1744662690000},{"__ext__bk_46__container_name":"unify-query","_value_":4388,"_timestamp_":1744662720000},{"__ext__bk_46__container_name":"unify-query","_value_":3357,"_timestamp_":1744662750000},{"__ext__bk_46__container_name":"unify-query","_value_":4381,"_timestamp_":1744662780000},{"__ext__bk_46__container_name":"unify-query","_value_":3683,"_timestamp_":1744662810000},{"__ext__bk_46__container_name":"unify-query","_value_":4353,"_timestamp_":1744662840000},{"__ext__bk_46__container_name":"unify-query","_value_":3441,"_timestamp_":1744662870000},{"__ext__bk_46__container_name":"unify-query","_value_":4251,"_timestamp_":1744662900000},{"__ext__bk_46__container_name":"unify-query","_value_":3476,"_timestamp_":1744662930000},{"__ext__bk_46__container_name":"unify-query","_value_":4036,"_timestamp_":1744662960000},{"__ext__bk_46__container_name":"unify-query","_value_":3549,"_timestamp_":1744662990000},{"__ext__bk_46__container_name":"unify-query","_value_":4351,"_timestamp_":1744663020000},{"__ext__bk_46__container_name":"unify-query","_value_":3651,"_timestamp_":1744663050000},{"__ext__bk_46__container_name":"unify-query","_value_":4096,"_timestamp_":1744663080000},{"__ext__bk_46__container_name":"unify-query","_value_":3618,"_timestamp_":1744663110000},{"__ext__bk_46__container_name":"unify-query","_value_":4100,"_timestamp_":1744663140000},{"__ext__bk_46__container_name":"unify-query","_value_":3622,"_timestamp_":1744663170000},{"__ext__bk_46__container_name":"unify-query","_value_":6044,"_timestamp_":1744663200000},{"__ext__bk_46__container_name":"unify-query","_value_":3766,"_timestamp_":1744663230000},{"__ext__bk_46__container_name":"unify-query","_value_":4461,"_timestamp_":1744663260000},{"__ext__bk_46__container_name":"unify-query","_value_":3783,"_timestamp_":1744663290000},{"__ext__bk_46__container_name":"unify-query","_value_":4559,"_timestamp_":1744663320000},{"__ext__bk_46__container_name":"unify-query","_value_":3634,"_timestamp_":1744663350000},{"__ext__bk_46__container_name":"unify-query","_value_":3869,"_timestamp_":1744663380000},{"__ext__bk_46__container_name":"unify-query","_value_":3249,"_timestamp_":1744663410000},{"__ext__bk_46__container_name":"unify-query","_value_":4473,"_timestamp_":1744663440000},{"__ext__bk_46__container_name":"unify-query","_value_":3514,"_timestamp_":1744663470000},{"__ext__bk_46__container_name":"unify-query","_value_":4923,"_timestamp_":1744663500000},{"__ext__bk_46__container_name":"unify-query","_value_":3379,"_timestamp_":1744663530000},{"__ext__bk_46__container_name":"unify-query","_value_":4489,"_timestamp_":1744663560000},{"__ext__bk_46__container_name":"unify-query","_value_":3411,"_timestamp_":1744663590000},{"__ext__bk_46__container_name":"unify-query","_value_":4374,"_timestamp_":1744663620000},{"__ext__bk_46__container_name":"unify-query","_value_":3370,"_timestamp_":1744663650000},{"__ext__bk_46__container_name":"unify-query","_value_":4310,"_timestamp_":1744663680000},{"__ext__bk_46__container_name":"unify-query","_value_":3609,"_timestamp_":1744663710000},{"__ext__bk_46__container_name":"unify-query","_value_":4318,"_timestamp_":1744663740000},{"__ext__bk_46__container_name":"unify-query","_value_":3570,"_timestamp_":1744663770000},{"__ext__bk_46__container_name":"unify-query","_value_":4334,"_timestamp_":1744663800000},{"__ext__bk_46__container_name":"unify-query","_value_":3767,"_timestamp_":1744663830000},{"__ext__bk_46__container_name":"unify-query","_value_":4455,"_timestamp_":1744663860000},{"__ext__bk_46__container_name":"unify-query","_value_":3703,"_timestamp_":1744663890000},{"__ext__bk_46__container_name":"unify-query","_value_":4511,"_timestamp_":1744663920000},{"__ext__bk_46__container_name":"unify-query","_value_":3667,"_timestamp_":1744663950000},{"__ext__bk_46__container_name":"unify-query","_value_":3998,"_timestamp_":1744663980000},{"__ext__bk_46__container_name":"unify-query","_value_":3579,"_timestamp_":1744664010000},{"__ext__bk_46__container_name":"unify-query","_value_":4156,"_timestamp_":1744664040000},{"__ext__bk_46__container_name":"unify-query","_value_":3340,"_timestamp_":1744664070000},{"__ext__bk_46__container_name":"unify-query","_value_":4344,"_timestamp_":1744664100000},{"__ext__bk_46__container_name":"unify-query","_value_":3590,"_timestamp_":1744664130000},{"__ext__bk_46__container_name":"unify-query","_value_":4161,"_timestamp_":1744664160000},{"__ext__bk_46__container_name":"unify-query","_value_":3484,"_timestamp_":1744664190000},{"__ext__bk_46__container_name":"unify-query","_value_":4273,"_timestamp_":1744664220000},{"__ext__bk_46__container_name":"unify-query","_value_":3494,"_timestamp_":1744664250000},{"__ext__bk_46__container_name":"unify-query","_value_":4230,"_timestamp_":1744664280000},{"__ext__bk_46__container_name":"unify-query","_value_":3619,"_timestamp_":1744664310000},{"__ext__bk_46__container_name":"unify-query","_value_":4013,"_timestamp_":1744664340000},{"__ext__bk_46__container_name":"unify-query","_value_":3565,"_timestamp_":1744664370000},{"__ext__bk_46__container_name":"unify-query","_value_":18144,"_timestamp_":1744664400000},{"__ext__bk_46__container_name":"unify-query","_value_":13615,"_timestamp_":1744664430000},{"__ext__bk_46__container_name":"unify-query","_value_":3178,"_timestamp_":1744664460000},{"__ext__bk_46__container_name":"unify-query","_value_":13044,"_timestamp_":1744664490000},{"__ext__bk_46__container_name":"unify-query","_value_":4767,"_timestamp_":1744664520000},{"__ext__bk_46__container_name":"unify-query","_value_":3528,"_timestamp_":1744664550000},{"__ext__bk_46__container_name":"unify-query","_value_":4316,"_timestamp_":1744664580000},{"__ext__bk_46__container_name":"unify-query","_value_":3317,"_timestamp_":1744664610000},{"__ext__bk_46__container_name":"unify-query","_value_":4395,"_timestamp_":1744664640000},{"__ext__bk_46__container_name":"unify-query","_value_":3599,"_timestamp_":1744664670000},{"__ext__bk_46__container_name":"unify-query","_value_":4149,"_timestamp_":1744664700000},{"__ext__bk_46__container_name":"unify-query","_value_":3474,"_timestamp_":1744664730000},{"__ext__bk_46__container_name":"unify-query","_value_":4201,"_timestamp_":1744664760000},{"__ext__bk_46__container_name":"unify-query","_value_":3384,"_timestamp_":1744664790000},{"__ext__bk_46__container_name":"unify-query","_value_":4442,"_timestamp_":1744664820000},{"__ext__bk_46__container_name":"unify-query","_value_":3559,"_timestamp_":1744664850000},{"__ext__bk_46__container_name":"unify-query","_value_":4166,"_timestamp_":1744664880000},{"__ext__bk_46__container_name":"unify-query","_value_":3438,"_timestamp_":1744664910000},{"__ext__bk_46__container_name":"unify-query","_value_":4244,"_timestamp_":1744664940000},{"__ext__bk_46__container_name":"unify-query","_value_":3640,"_timestamp_":1744664970000},{"__ext__bk_46__container_name":"unify-query","_value_":4305,"_timestamp_":1744665000000},{"__ext__bk_46__container_name":"unify-query","_value_":3771,"_timestamp_":1744665030000},{"__ext__bk_46__container_name":"unify-query","_value_":4485,"_timestamp_":1744665060000},{"__ext__bk_46__container_name":"unify-query","_value_":3842,"_timestamp_":1744665090000},{"__ext__bk_46__container_name":"unify-query","_value_":4423,"_timestamp_":1744665120000},{"__ext__bk_46__container_name":"unify-query","_value_":3610,"_timestamp_":1744665150000},{"__ext__bk_46__container_name":"unify-query","_value_":4125,"_timestamp_":1744665180000},{"__ext__bk_46__container_name":"unify-query","_value_":3500,"_timestamp_":1744665210000},{"__ext__bk_46__container_name":"unify-query","_value_":4252,"_timestamp_":1744665240000},{"__ext__bk_46__container_name":"unify-query","_value_":3427,"_timestamp_":1744665270000},{"__ext__bk_46__container_name":"unify-query","_value_":5089,"_timestamp_":1744665300000},{"__ext__bk_46__container_name":"unify-query","_value_":3450,"_timestamp_":1744665330000},{"__ext__bk_46__container_name":"unify-query","_value_":4349,"_timestamp_":1744665360000},{"__ext__bk_46__container_name":"unify-query","_value_":3188,"_timestamp_":1744665390000},{"__ext__bk_46__container_name":"unify-query","_value_":4556,"_timestamp_":1744665420000},{"__ext__bk_46__container_name":"unify-query","_value_":3372,"_timestamp_":1744665450000},{"__ext__bk_46__container_name":"unify-query","_value_":4408,"_timestamp_":1744665480000},{"__ext__bk_46__container_name":"unify-query","_value_":3445,"_timestamp_":1744665510000},{"__ext__bk_46__container_name":"unify-query","_value_":4213,"_timestamp_":1744665540000},{"__ext__bk_46__container_name":"unify-query","_value_":3408,"_timestamp_":1744665570000},{"__ext__bk_46__container_name":"unify-query","_value_":6235,"_timestamp_":1744665600000},{"__ext__bk_46__container_name":"unify-query","_value_":3641,"_timestamp_":1744665630000},{"__ext__bk_46__container_name":"unify-query","_value_":4577,"_timestamp_":1744665660000},{"__ext__bk_46__container_name":"unify-query","_value_":3719,"_timestamp_":1744665690000},{"__ext__bk_46__container_name":"unify-query","_value_":4548,"_timestamp_":1744665720000},{"__ext__bk_46__container_name":"unify-query","_value_":3420,"_timestamp_":1744665750000},{"__ext__bk_46__container_name":"unify-query","_value_":4246,"_timestamp_":1744665780000},{"__ext__bk_46__container_name":"unify-query","_value_":3359,"_timestamp_":1744665810000},{"__ext__bk_46__container_name":"unify-query","_value_":4332,"_timestamp_":1744665840000},{"__ext__bk_46__container_name":"unify-query","_value_":3422,"_timestamp_":1744665870000},{"__ext__bk_46__container_name":"unify-query","_value_":4229,"_timestamp_":1744665900000},{"__ext__bk_46__container_name":"unify-query","_value_":3610,"_timestamp_":1744665930000},{"__ext__bk_46__container_name":"unify-query","_value_":4119,"_timestamp_":1744665960000},{"__ext__bk_46__container_name":"unify-query","_value_":3570,"_timestamp_":1744665990000},{"__ext__bk_46__container_name":"unify-query","_value_":4144,"_timestamp_":1744666020000},{"__ext__bk_46__container_name":"unify-query","_value_":3302,"_timestamp_":1744666050000},{"__ext__bk_46__container_name":"unify-query","_value_":4398,"_timestamp_":1744666080000},{"__ext__bk_46__container_name":"unify-query","_value_":3559,"_timestamp_":1744666110000},{"__ext__bk_46__container_name":"unify-query","_value_":4097,"_timestamp_":1744666140000},{"__ext__bk_46__container_name":"unify-query","_value_":3315,"_timestamp_":1744666170000},{"__ext__bk_46__container_name":"unify-query","_value_":16721,"_timestamp_":1744666200000},{"__ext__bk_46__container_name":"unify-query","_value_":13631,"_timestamp_":1744666230000},{"__ext__bk_46__container_name":"unify-query","_value_":2982,"_timestamp_":1744666260000},{"__ext__bk_46__container_name":"unify-query","_value_":11858,"_timestamp_":1744666290000},{"__ext__bk_46__container_name":"unify-query","_value_":5515,"_timestamp_":1744666320000},{"__ext__bk_46__container_name":"unify-query","_value_":2869,"_timestamp_":1744666350000},{"__ext__bk_46__container_name":"unify-query","_value_":4795,"_timestamp_":1744666380000},{"__ext__bk_46__container_name":"unify-query","_value_":3603,"_timestamp_":1744666410000},{"__ext__bk_46__container_name":"unify-query","_value_":4204,"_timestamp_":1744666440000},{"__ext__bk_46__container_name":"unify-query","_value_":3264,"_timestamp_":1744666470000},{"__ext__bk_46__container_name":"unify-query","_value_":4377,"_timestamp_":1744666500000},{"__ext__bk_46__container_name":"unify-query","_value_":3443,"_timestamp_":1744666530000},{"__ext__bk_46__container_name":"unify-query","_value_":4307,"_timestamp_":1744666560000},{"__ext__bk_46__container_name":"unify-query","_value_":3459,"_timestamp_":1744666590000},{"__ext__bk_46__container_name":"unify-query","_value_":4342,"_timestamp_":1744666620000},{"__ext__bk_46__container_name":"unify-query","_value_":3598,"_timestamp_":1744666650000},{"__ext__bk_46__container_name":"unify-query","_value_":4052,"_timestamp_":1744666680000},{"__ext__bk_46__container_name":"unify-query","_value_":3577,"_timestamp_":1744666710000},{"__ext__bk_46__container_name":"unify-query","_value_":4128,"_timestamp_":1744666740000},{"__ext__bk_46__container_name":"unify-query","_value_":3499,"_timestamp_":1744666770000},{"__ext__bk_46__container_name":"unify-query","_value_":6209,"_timestamp_":1744666800000},{"__ext__bk_46__container_name":"unify-query","_value_":3575,"_timestamp_":1744666830000},{"__ext__bk_46__container_name":"unify-query","_value_":4543,"_timestamp_":1744666860000},{"__ext__bk_46__container_name":"unify-query","_value_":3604,"_timestamp_":1744666890000},{"__ext__bk_46__container_name":"unify-query","_value_":4579,"_timestamp_":1744666920000},{"__ext__bk_46__container_name":"unify-query","_value_":3531,"_timestamp_":1744666950000},{"__ext__bk_46__container_name":"unify-query","_value_":4314,"_timestamp_":1744666980000},{"__ext__bk_46__container_name":"unify-query","_value_":3416,"_timestamp_":1744667010000},{"__ext__bk_46__container_name":"unify-query","_value_":4320,"_timestamp_":1744667040000},{"__ext__bk_46__container_name":"unify-query","_value_":3488,"_timestamp_":1744667070000},{"__ext__bk_46__container_name":"unify-query","_value_":5054,"_timestamp_":1744667100000},{"__ext__bk_46__container_name":"unify-query","_value_":3525,"_timestamp_":1744667130000},{"__ext__bk_46__container_name":"unify-query","_value_":4313,"_timestamp_":1744667160000},{"__ext__bk_46__container_name":"unify-query","_value_":3607,"_timestamp_":1744667190000},{"__ext__bk_46__container_name":"unify-query","_value_":4118,"_timestamp_":1744667220000},{"__ext__bk_46__container_name":"unify-query","_value_":3350,"_timestamp_":1744667250000},{"__ext__bk_46__container_name":"unify-query","_value_":4280,"_timestamp_":1744667280000},{"__ext__bk_46__container_name":"unify-query","_value_":3634,"_timestamp_":1744667310000},{"__ext__bk_46__container_name":"unify-query","_value_":4174,"_timestamp_":1744667340000},{"__ext__bk_46__container_name":"unify-query","_value_":3807,"_timestamp_":1744667370000},{"__ext__bk_46__container_name":"unify-query","_value_":4358,"_timestamp_":1744667400000},{"__ext__bk_46__container_name":"unify-query","_value_":3595,"_timestamp_":1744667430000},{"__ext__bk_46__container_name":"unify-query","_value_":4630,"_timestamp_":1744667460000},{"__ext__bk_46__container_name":"unify-query","_value_":3845,"_timestamp_":1744667490000},{"__ext__bk_46__container_name":"unify-query","_value_":4361,"_timestamp_":1744667520000},{"__ext__bk_46__container_name":"unify-query","_value_":3572,"_timestamp_":1744667550000},{"__ext__bk_46__container_name":"unify-query","_value_":4095,"_timestamp_":1744667580000},{"__ext__bk_46__container_name":"unify-query","_value_":3535,"_timestamp_":1744667610000},{"__ext__bk_46__container_name":"unify-query","_value_":4200,"_timestamp_":1744667640000},{"__ext__bk_46__container_name":"unify-query","_value_":3390,"_timestamp_":1744667670000},{"__ext__bk_46__container_name":"unify-query","_value_":4262,"_timestamp_":1744667700000},{"__ext__bk_46__container_name":"unify-query","_value_":3398,"_timestamp_":1744667730000},{"__ext__bk_46__container_name":"unify-query","_value_":4320,"_timestamp_":1744667760000},{"__ext__bk_46__container_name":"unify-query","_value_":3429,"_timestamp_":1744667790000},{"__ext__bk_46__container_name":"unify-query","_value_":4288,"_timestamp_":1744667820000},{"__ext__bk_46__container_name":"unify-query","_value_":3482,"_timestamp_":1744667850000},{"__ext__bk_46__container_name":"unify-query","_value_":4166,"_timestamp_":1744667880000},{"__ext__bk_46__container_name":"unify-query","_value_":3612,"_timestamp_":1744667910000},{"__ext__bk_46__container_name":"unify-query","_value_":4194,"_timestamp_":1744667940000},{"__ext__bk_46__container_name":"unify-query","_value_":3423,"_timestamp_":1744667970000},{"__ext__bk_46__container_name":"unify-query","_value_":18203,"_timestamp_":1744668000000},{"__ext__bk_46__container_name":"unify-query","_value_":13685,"_timestamp_":1744668030000},{"__ext__bk_46__container_name":"unify-query","_value_":3281,"_timestamp_":1744668060000},{"__ext__bk_46__container_name":"unify-query","_value_":12556,"_timestamp_":1744668090000},{"__ext__bk_46__container_name":"unify-query","_value_":4893,"_timestamp_":1744668120000},{"__ext__bk_46__container_name":"unify-query","_value_":3607,"_timestamp_":1744668150000},{"__ext__bk_46__container_name":"unify-query","_value_":4336,"_timestamp_":1744668180000},{"__ext__bk_46__container_name":"unify-query","_value_":3609,"_timestamp_":1744668210000},{"__ext__bk_46__container_name":"unify-query","_value_":4097,"_timestamp_":1744668240000},{"__ext__bk_46__container_name":"unify-query","_value_":3669,"_timestamp_":1744668270000},{"__ext__bk_46__container_name":"unify-query","_value_":3997,"_timestamp_":1744668300000},{"__ext__bk_46__container_name":"unify-query","_value_":3494,"_timestamp_":1744668330000},{"__ext__bk_46__container_name":"unify-query","_value_":4172,"_timestamp_":1744668360000},{"__ext__bk_46__container_name":"unify-query","_value_":3523,"_timestamp_":1744668390000},{"__ext__bk_46__container_name":"unify-query","_value_":3877,"_timestamp_":1744668420000},{"__ext__bk_46__container_name":"unify-query","_value_":3565,"_timestamp_":1744668450000},{"__ext__bk_46__container_name":"unify-query","_value_":4230,"_timestamp_":1744668480000},{"__ext__bk_46__container_name":"unify-query","_value_":3469,"_timestamp_":1744668510000},{"__ext__bk_46__container_name":"unify-query","_value_":4243,"_timestamp_":1744668540000},{"__ext__bk_46__container_name":"unify-query","_value_":3304,"_timestamp_":1744668570000},{"__ext__bk_46__container_name":"unify-query","_value_":4690,"_timestamp_":1744668600000},{"__ext__bk_46__container_name":"unify-query","_value_":3717,"_timestamp_":1744668630000},{"__ext__bk_46__container_name":"unify-query","_value_":4618,"_timestamp_":1744668660000},{"__ext__bk_46__container_name":"unify-query","_value_":3732,"_timestamp_":1744668690000},{"__ext__bk_46__container_name":"unify-query","_value_":4477,"_timestamp_":1744668720000},{"__ext__bk_46__container_name":"unify-query","_value_":3615,"_timestamp_":1744668750000},{"__ext__bk_46__container_name":"unify-query","_value_":4154,"_timestamp_":1744668780000},{"__ext__bk_46__container_name":"unify-query","_value_":3367,"_timestamp_":1744668810000},{"__ext__bk_46__container_name":"unify-query","_value_":4193,"_timestamp_":1744668840000},{"__ext__bk_46__container_name":"unify-query","_value_":3592,"_timestamp_":1744668870000},{"__ext__bk_46__container_name":"unify-query","_value_":4971,"_timestamp_":1744668900000},{"__ext__bk_46__container_name":"unify-query","_value_":3359,"_timestamp_":1744668930000},{"__ext__bk_46__container_name":"unify-query","_value_":4540,"_timestamp_":1744668960000},{"__ext__bk_46__container_name":"unify-query","_value_":3406,"_timestamp_":1744668990000},{"__ext__bk_46__container_name":"unify-query","_value_":4375,"_timestamp_":1744669020000},{"__ext__bk_46__container_name":"unify-query","_value_":3386,"_timestamp_":1744669050000},{"__ext__bk_46__container_name":"unify-query","_value_":4281,"_timestamp_":1744669080000},{"__ext__bk_46__container_name":"unify-query","_value_":3410,"_timestamp_":1744669110000},{"__ext__bk_46__container_name":"unify-query","_value_":4545,"_timestamp_":1744669140000},{"__ext__bk_46__container_name":"unify-query","_value_":3724,"_timestamp_":1744669170000},{"__ext__bk_46__container_name":"unify-query","_value_":5903,"_timestamp_":1744669200000},{"__ext__bk_46__container_name":"unify-query","_value_":3672,"_timestamp_":1744669230000},{"__ext__bk_46__container_name":"unify-query","_value_":4413,"_timestamp_":1744669260000},{"__ext__bk_46__container_name":"unify-query","_value_":3792,"_timestamp_":1744669290000},{"__ext__bk_46__container_name":"unify-query","_value_":4422,"_timestamp_":1744669320000},{"__ext__bk_46__container_name":"unify-query","_value_":3718,"_timestamp_":1744669350000},{"__ext__bk_46__container_name":"unify-query","_value_":4213,"_timestamp_":1744669380000},{"__ext__bk_46__container_name":"unify-query","_value_":3622,"_timestamp_":1744669410000},{"__ext__bk_46__container_name":"unify-query","_value_":4043,"_timestamp_":1744669440000},{"__ext__bk_46__container_name":"unify-query","_value_":3542,"_timestamp_":1744669470000},{"__ext__bk_46__container_name":"unify-query","_value_":4179,"_timestamp_":1744669500000},{"__ext__bk_46__container_name":"unify-query","_value_":3368,"_timestamp_":1744669530000},{"__ext__bk_46__container_name":"unify-query","_value_":4354,"_timestamp_":1744669560000},{"__ext__bk_46__container_name":"unify-query","_value_":3368,"_timestamp_":1744669590000},{"__ext__bk_46__container_name":"unify-query","_value_":4229,"_timestamp_":1744669620000},{"__ext__bk_46__container_name":"unify-query","_value_":3458,"_timestamp_":1744669650000},{"__ext__bk_46__container_name":"unify-query","_value_":4310,"_timestamp_":1744669680000},{"__ext__bk_46__container_name":"unify-query","_value_":3512,"_timestamp_":1744669710000},{"__ext__bk_46__container_name":"unify-query","_value_":4188,"_timestamp_":1744669740000},{"__ext__bk_46__container_name":"unify-query","_value_":3436,"_timestamp_":1744669770000},{"__ext__bk_46__container_name":"unify-query","_value_":12171,"_timestamp_":1744669800000},{"__ext__bk_46__container_name":"unify-query","_value_":18129,"_timestamp_":1744669830000},{"__ext__bk_46__container_name":"unify-query","_value_":7142,"_timestamp_":1744669860000},{"__ext__bk_46__container_name":"unify-query","_value_":9153,"_timestamp_":1744669890000},{"__ext__bk_46__container_name":"unify-query","_value_":4566,"_timestamp_":1744669920000},{"__ext__bk_46__container_name":"unify-query","_value_":3225,"_timestamp_":1744669950000},{"__ext__bk_46__container_name":"unify-query","_value_":4378,"_timestamp_":1744669980000},{"__ext__bk_46__container_name":"unify-query","_value_":3623,"_timestamp_":1744670010000},{"__ext__bk_46__container_name":"unify-query","_value_":4266,"_timestamp_":1744670040000},{"__ext__bk_46__container_name":"unify-query","_value_":3645,"_timestamp_":1744670070000},{"__ext__bk_46__container_name":"unify-query","_value_":4043,"_timestamp_":1744670100000},{"__ext__bk_46__container_name":"unify-query","_value_":3350,"_timestamp_":1744670130000},{"__ext__bk_46__container_name":"unify-query","_value_":4333,"_timestamp_":1744670160000},{"__ext__bk_46__container_name":"unify-query","_value_":3489,"_timestamp_":1744670190000},{"__ext__bk_46__container_name":"unify-query","_value_":4303,"_timestamp_":1744670220000},{"__ext__bk_46__container_name":"unify-query","_value_":3560,"_timestamp_":1744670250000},{"__ext__bk_46__container_name":"unify-query","_value_":4121,"_timestamp_":1744670280000},{"__ext__bk_46__container_name":"unify-query","_value_":3374,"_timestamp_":1744670310000},{"__ext__bk_46__container_name":"unify-query","_value_":4362,"_timestamp_":1744670340000},{"__ext__bk_46__container_name":"unify-query","_value_":3242,"_timestamp_":1744670370000},{"__ext__bk_46__container_name":"unify-query","_value_":6416,"_timestamp_":1744670400000},{"__ext__bk_46__container_name":"unify-query","_value_":3697,"_timestamp_":1744670430000},{"__ext__bk_46__container_name":"unify-query","_value_":4506,"_timestamp_":1744670460000},{"__ext__bk_46__container_name":"unify-query","_value_":3749,"_timestamp_":1744670490000},{"__ext__bk_46__container_name":"unify-query","_value_":4587,"_timestamp_":1744670520000},{"__ext__bk_46__container_name":"unify-query","_value_":3538,"_timestamp_":1744670550000},{"__ext__bk_46__container_name":"unify-query","_value_":4221,"_timestamp_":1744670580000},{"__ext__bk_46__container_name":"unify-query","_value_":3476,"_timestamp_":1744670610000},{"__ext__bk_46__container_name":"unify-query","_value_":4227,"_timestamp_":1744670640000},{"__ext__bk_46__container_name":"unify-query","_value_":3587,"_timestamp_":1744670670000},{"__ext__bk_46__container_name":"unify-query","_value_":4848,"_timestamp_":1744670700000},{"__ext__bk_46__container_name":"unify-query","_value_":3551,"_timestamp_":1744670730000},{"__ext__bk_46__container_name":"unify-query","_value_":4068,"_timestamp_":1744670760000},{"__ext__bk_46__container_name":"unify-query","_value_":3387,"_timestamp_":1744670790000},{"__ext__bk_46__container_name":"unify-query","_value_":4366,"_timestamp_":1744670820000},{"__ext__bk_46__container_name":"unify-query","_value_":3635,"_timestamp_":1744670850000},{"__ext__bk_46__container_name":"unify-query","_value_":4256,"_timestamp_":1744670880000},{"__ext__bk_46__container_name":"unify-query","_value_":3690,"_timestamp_":1744670910000},{"__ext__bk_46__container_name":"unify-query","_value_":4155,"_timestamp_":1744670940000},{"__ext__bk_46__container_name":"unify-query","_value_":3318,"_timestamp_":1744670970000},{"__ext__bk_46__container_name":"unify-query","_value_":4661,"_timestamp_":1744671000000},{"__ext__bk_46__container_name":"unify-query","_value_":3494,"_timestamp_":1744671030000},{"__ext__bk_46__container_name":"unify-query","_value_":4442,"_timestamp_":1744671060000},{"__ext__bk_46__container_name":"unify-query","_value_":3643,"_timestamp_":1744671090000},{"__ext__bk_46__container_name":"unify-query","_value_":4755,"_timestamp_":1744671120000},{"__ext__bk_46__container_name":"unify-query","_value_":3607,"_timestamp_":1744671150000},{"__ext__bk_46__container_name":"unify-query","_value_":4284,"_timestamp_":1744671180000},{"__ext__bk_46__container_name":"unify-query","_value_":3258,"_timestamp_":1744671210000},{"__ext__bk_46__container_name":"unify-query","_value_":4453,"_timestamp_":1744671240000},{"__ext__bk_46__container_name":"unify-query","_value_":3431,"_timestamp_":1744671270000},{"__ext__bk_46__container_name":"unify-query","_value_":4231,"_timestamp_":1744671300000},{"__ext__bk_46__container_name":"unify-query","_value_":3623,"_timestamp_":1744671330000},{"__ext__bk_46__container_name":"unify-query","_value_":3907,"_timestamp_":1744671360000},{"__ext__bk_46__container_name":"unify-query","_value_":3524,"_timestamp_":1744671390000},{"__ext__bk_46__container_name":"unify-query","_value_":4438,"_timestamp_":1744671420000},{"__ext__bk_46__container_name":"unify-query","_value_":3547,"_timestamp_":1744671450000},{"__ext__bk_46__container_name":"unify-query","_value_":4033,"_timestamp_":1744671480000},{"__ext__bk_46__container_name":"unify-query","_value_":3632,"_timestamp_":1744671510000},{"__ext__bk_46__container_name":"unify-query","_value_":4162,"_timestamp_":1744671540000},{"__ext__bk_46__container_name":"unify-query","_value_":3588,"_timestamp_":1744671570000},{"__ext__bk_46__container_name":"unify-query","_value_":16444,"_timestamp_":1744671600000},{"__ext__bk_46__container_name":"unify-query","_value_":15396,"_timestamp_":1744671630000},{"__ext__bk_46__container_name":"unify-query","_value_":3024,"_timestamp_":1744671660000},{"__ext__bk_46__container_name":"unify-query","_value_":12656,"_timestamp_":1744671690000},{"__ext__bk_46__container_name":"unify-query","_value_":4733,"_timestamp_":1744671720000},{"__ext__bk_46__container_name":"unify-query","_value_":3766,"_timestamp_":1744671750000},{"__ext__bk_46__container_name":"unify-query","_value_":4388,"_timestamp_":1744671780000},{"__ext__bk_46__container_name":"unify-query","_value_":3340,"_timestamp_":1744671810000},{"__ext__bk_46__container_name":"unify-query","_value_":4487,"_timestamp_":1744671840000},{"__ext__bk_46__container_name":"unify-query","_value_":3549,"_timestamp_":1744671870000},{"__ext__bk_46__container_name":"unify-query","_value_":4154,"_timestamp_":1744671900000},{"__ext__bk_46__container_name":"unify-query","_value_":3406,"_timestamp_":1744671930000},{"__ext__bk_46__container_name":"unify-query","_value_":4314,"_timestamp_":1744671960000},{"__ext__bk_46__container_name":"unify-query","_value_":3472,"_timestamp_":1744671990000},{"__ext__bk_46__container_name":"unify-query","_value_":4309,"_timestamp_":1744672020000},{"__ext__bk_46__container_name":"unify-query","_value_":3458,"_timestamp_":1744672050000},{"__ext__bk_46__container_name":"unify-query","_value_":4191,"_timestamp_":1744672080000},{"__ext__bk_46__container_name":"unify-query","_value_":3475,"_timestamp_":1744672110000},{"__ext__bk_46__container_name":"unify-query","_value_":4194,"_timestamp_":1744672140000},{"__ext__bk_46__container_name":"unify-query","_value_":3525,"_timestamp_":1744672170000},{"__ext__bk_46__container_name":"unify-query","_value_":4445,"_timestamp_":1744672200000},{"__ext__bk_46__container_name":"unify-query","_value_":3822,"_timestamp_":1744672230000},{"__ext__bk_46__container_name":"unify-query","_value_":4346,"_timestamp_":1744672260000},{"__ext__bk_46__container_name":"unify-query","_value_":3700,"_timestamp_":1744672290000},{"__ext__bk_46__container_name":"unify-query","_value_":4615,"_timestamp_":1744672320000},{"__ext__bk_46__container_name":"unify-query","_value_":3591,"_timestamp_":1744672350000},{"__ext__bk_46__container_name":"unify-query","_value_":4056,"_timestamp_":1744672380000},{"__ext__bk_46__container_name":"unify-query","_value_":3544,"_timestamp_":1744672410000},{"__ext__bk_46__container_name":"unify-query","_value_":4188,"_timestamp_":1744672440000},{"__ext__bk_46__container_name":"unify-query","_value_":3647,"_timestamp_":1744672470000},{"__ext__bk_46__container_name":"unify-query","_value_":4887,"_timestamp_":1744672500000},{"__ext__bk_46__container_name":"unify-query","_value_":3450,"_timestamp_":1744672530000},{"__ext__bk_46__container_name":"unify-query","_value_":4302,"_timestamp_":1744672560000},{"__ext__bk_46__container_name":"unify-query","_value_":3425,"_timestamp_":1744672590000},{"__ext__bk_46__container_name":"unify-query","_value_":4320,"_timestamp_":1744672620000},{"__ext__bk_46__container_name":"unify-query","_value_":3532,"_timestamp_":1744672650000},{"__ext__bk_46__container_name":"unify-query","_value_":4282,"_timestamp_":1744672680000},{"__ext__bk_46__container_name":"unify-query","_value_":3571,"_timestamp_":1744672710000},{"__ext__bk_46__container_name":"unify-query","_value_":4182,"_timestamp_":1744672740000},{"__ext__bk_46__container_name":"unify-query","_value_":3210,"_timestamp_":1744672770000},{"__ext__bk_46__container_name":"unify-query","_value_":6383,"_timestamp_":1744672800000},{"__ext__bk_46__container_name":"unify-query","_value_":3622,"_timestamp_":1744672830000},{"__ext__bk_46__container_name":"unify-query","_value_":4408,"_timestamp_":1744672860000},{"__ext__bk_46__container_name":"unify-query","_value_":3611,"_timestamp_":1744672890000},{"__ext__bk_46__container_name":"unify-query","_value_":4795,"_timestamp_":1744672920000},{"__ext__bk_46__container_name":"unify-query","_value_":3632,"_timestamp_":1744672950000},{"__ext__bk_46__container_name":"unify-query","_value_":4102,"_timestamp_":1744672980000},{"__ext__bk_46__container_name":"unify-query","_value_":3534,"_timestamp_":1744673010000},{"__ext__bk_46__container_name":"unify-query","_value_":4212,"_timestamp_":1744673040000},{"__ext__bk_46__container_name":"unify-query","_value_":3380,"_timestamp_":1744673070000},{"__ext__bk_46__container_name":"unify-query","_value_":4289,"_timestamp_":1744673100000},{"__ext__bk_46__container_name":"unify-query","_value_":3565,"_timestamp_":1744673130000},{"__ext__bk_46__container_name":"unify-query","_value_":4120,"_timestamp_":1744673160000},{"__ext__bk_46__container_name":"unify-query","_value_":3526,"_timestamp_":1744673190000},{"__ext__bk_46__container_name":"unify-query","_value_":4200,"_timestamp_":1744673220000},{"__ext__bk_46__container_name":"unify-query","_value_":3302,"_timestamp_":1744673250000},{"__ext__bk_46__container_name":"unify-query","_value_":4370,"_timestamp_":1744673280000},{"__ext__bk_46__container_name":"unify-query","_value_":3462,"_timestamp_":1744673310000},{"__ext__bk_46__container_name":"unify-query","_value_":4223,"_timestamp_":1744673340000},{"__ext__bk_46__container_name":"unify-query","_value_":3564,"_timestamp_":1744673370000},{"__ext__bk_46__container_name":"unify-query","_value_":12072,"_timestamp_":1744673400000},{"__ext__bk_46__container_name":"unify-query","_value_":17986,"_timestamp_":1744673430000},{"__ext__bk_46__container_name":"unify-query","_value_":4089,"_timestamp_":1744673460000},{"__ext__bk_46__container_name":"unify-query","_value_":12000,"_timestamp_":1744673490000},{"__ext__bk_46__container_name":"unify-query","_value_":4790,"_timestamp_":1744673520000},{"__ext__bk_46__container_name":"unify-query","_value_":3637,"_timestamp_":1744673550000},{"__ext__bk_46__container_name":"unify-query","_value_":4177,"_timestamp_":1744673580000},{"__ext__bk_46__container_name":"unify-query","_value_":3438,"_timestamp_":1744673610000},{"__ext__bk_46__container_name":"unify-query","_value_":4465,"_timestamp_":1744673640000},{"__ext__bk_46__container_name":"unify-query","_value_":3627,"_timestamp_":1744673670000},{"__ext__bk_46__container_name":"unify-query","_value_":4131,"_timestamp_":1744673700000},{"__ext__bk_46__container_name":"unify-query","_value_":3396,"_timestamp_":1744673730000},{"__ext__bk_46__container_name":"unify-query","_value_":4395,"_timestamp_":1744673760000},{"__ext__bk_46__container_name":"unify-query","_value_":3638,"_timestamp_":1744673790000},{"__ext__bk_46__container_name":"unify-query","_value_":4093,"_timestamp_":1744673820000},{"__ext__bk_46__container_name":"unify-query","_value_":3584,"_timestamp_":1744673850000},{"__ext__bk_46__container_name":"unify-query","_value_":4082,"_timestamp_":1744673880000},{"__ext__bk_46__container_name":"unify-query","_value_":3475,"_timestamp_":1744673910000},{"__ext__bk_46__container_name":"unify-query","_value_":4051,"_timestamp_":1744673940000},{"__ext__bk_46__container_name":"unify-query","_value_":3354,"_timestamp_":1744673970000},{"__ext__bk_46__container_name":"unify-query","_value_":6296,"_timestamp_":1744674000000},{"__ext__bk_46__container_name":"unify-query","_value_":3473,"_timestamp_":1744674030000},{"__ext__bk_46__container_name":"unify-query","_value_":4412,"_timestamp_":1744674060000},{"__ext__bk_46__container_name":"unify-query","_value_":3793,"_timestamp_":1744674090000},{"__ext__bk_46__container_name":"unify-query","_value_":4391,"_timestamp_":1744674120000},{"__ext__bk_46__container_name":"unify-query","_value_":3836,"_timestamp_":1744674150000},{"__ext__bk_46__container_name":"unify-query","_value_":4190,"_timestamp_":1744674180000},{"__ext__bk_46__container_name":"unify-query","_value_":3478,"_timestamp_":1744674210000},{"__ext__bk_46__container_name":"unify-query","_value_":4230,"_timestamp_":1744674240000},{"__ext__bk_46__container_name":"unify-query","_value_":3488,"_timestamp_":1744674270000},{"__ext__bk_46__container_name":"unify-query","_value_":4964,"_timestamp_":1744674300000},{"__ext__bk_46__container_name":"unify-query","_value_":3455,"_timestamp_":1744674330000},{"__ext__bk_46__container_name":"unify-query","_value_":4116,"_timestamp_":1744674360000},{"__ext__bk_46__container_name":"unify-query","_value_":3250,"_timestamp_":1744674390000},{"__ext__bk_46__container_name":"unify-query","_value_":4494,"_timestamp_":1744674420000},{"__ext__bk_46__container_name":"unify-query","_value_":3326,"_timestamp_":1744674450000},{"__ext__bk_46__container_name":"unify-query","_value_":4590,"_timestamp_":1744674480000},{"__ext__bk_46__container_name":"unify-query","_value_":3580,"_timestamp_":1744674510000},{"__ext__bk_46__container_name":"unify-query","_value_":4368,"_timestamp_":1744674540000},{"__ext__bk_46__container_name":"unify-query","_value_":3685,"_timestamp_":1744674570000},{"__ext__bk_46__container_name":"unify-query","_value_":4381,"_timestamp_":1744674600000},{"__ext__bk_46__container_name":"unify-query","_value_":3699,"_timestamp_":1744674630000},{"__ext__bk_46__container_name":"unify-query","_value_":4513,"_timestamp_":1744674660000},{"__ext__bk_46__container_name":"unify-query","_value_":3729,"_timestamp_":1744674690000},{"__ext__bk_46__container_name":"unify-query","_value_":4500,"_timestamp_":1744674720000},{"__ext__bk_46__container_name":"unify-query","_value_":3639,"_timestamp_":1744674750000},{"__ext__bk_46__container_name":"unify-query","_value_":4018,"_timestamp_":1744674780000},{"__ext__bk_46__container_name":"unify-query","_value_":3587,"_timestamp_":1744674810000},{"__ext__bk_46__container_name":"unify-query","_value_":4168,"_timestamp_":1744674840000},{"__ext__bk_46__container_name":"unify-query","_value_":3389,"_timestamp_":1744674870000},{"__ext__bk_46__container_name":"unify-query","_value_":4289,"_timestamp_":1744674900000},{"__ext__bk_46__container_name":"unify-query","_value_":3540,"_timestamp_":1744674930000},{"__ext__bk_46__container_name":"unify-query","_value_":4106,"_timestamp_":1744674960000},{"__ext__bk_46__container_name":"unify-query","_value_":3478,"_timestamp_":1744674990000},{"__ext__bk_46__container_name":"unify-query","_value_":4268,"_timestamp_":1744675020000},{"__ext__bk_46__container_name":"unify-query","_value_":3577,"_timestamp_":1744675050000},{"__ext__bk_46__container_name":"unify-query","_value_":4087,"_timestamp_":1744675080000},{"__ext__bk_46__container_name":"unify-query","_value_":3511,"_timestamp_":1744675110000},{"__ext__bk_46__container_name":"unify-query","_value_":4174,"_timestamp_":1744675140000},{"__ext__bk_46__container_name":"unify-query","_value_":3573,"_timestamp_":1744675170000},{"__ext__bk_46__container_name":"unify-query","_value_":17095,"_timestamp_":1744675200000},{"__ext__bk_46__container_name":"unify-query","_value_":14907,"_timestamp_":1744675230000},{"__ext__bk_46__container_name":"unify-query","_value_":6455,"_timestamp_":1744675260000},{"__ext__bk_46__container_name":"unify-query","_value_":9818,"_timestamp_":1744675290000},{"__ext__bk_46__container_name":"unify-query","_value_":5253,"_timestamp_":1744675320000},{"__ext__bk_46__container_name":"unify-query","_value_":3567,"_timestamp_":1744675350000},{"__ext__bk_46__container_name":"unify-query","_value_":4047,"_timestamp_":1744675380000},{"__ext__bk_46__container_name":"unify-query","_value_":3342,"_timestamp_":1744675410000},{"__ext__bk_46__container_name":"unify-query","_value_":4605,"_timestamp_":1744675440000},{"__ext__bk_46__container_name":"unify-query","_value_":3394,"_timestamp_":1744675470000},{"__ext__bk_46__container_name":"unify-query","_value_":4260,"_timestamp_":1744675500000},{"__ext__bk_46__container_name":"unify-query","_value_":3373,"_timestamp_":1744675530000},{"__ext__bk_46__container_name":"unify-query","_value_":4341,"_timestamp_":1744675560000},{"__ext__bk_46__container_name":"unify-query","_value_":3559,"_timestamp_":1744675590000},{"__ext__bk_46__container_name":"unify-query","_value_":4188,"_timestamp_":1744675620000},{"__ext__bk_46__container_name":"unify-query","_value_":3519,"_timestamp_":1744675650000},{"__ext__bk_46__container_name":"unify-query","_value_":4143,"_timestamp_":1744675680000},{"__ext__bk_46__container_name":"unify-query","_value_":3630,"_timestamp_":1744675710000},{"__ext__bk_46__container_name":"unify-query","_value_":4042,"_timestamp_":1744675740000},{"__ext__bk_46__container_name":"unify-query","_value_":3653,"_timestamp_":1744675770000},{"__ext__bk_46__container_name":"unify-query","_value_":4358,"_timestamp_":1744675800000},{"__ext__bk_46__container_name":"unify-query","_value_":3688,"_timestamp_":1744675830000},{"__ext__bk_46__container_name":"unify-query","_value_":4450,"_timestamp_":1744675860000},{"__ext__bk_46__container_name":"unify-query","_value_":3387,"_timestamp_":1744675890000},{"__ext__bk_46__container_name":"unify-query","_value_":4864,"_timestamp_":1744675920000},{"__ext__bk_46__container_name":"unify-query","_value_":3629,"_timestamp_":1744675950000},{"__ext__bk_46__container_name":"unify-query","_value_":4127,"_timestamp_":1744675980000},{"__ext__bk_46__container_name":"unify-query","_value_":3424,"_timestamp_":1744676010000},{"__ext__bk_46__container_name":"unify-query","_value_":4267,"_timestamp_":1744676040000},{"__ext__bk_46__container_name":"unify-query","_value_":3328,"_timestamp_":1744676070000},{"__ext__bk_46__container_name":"unify-query","_value_":5128,"_timestamp_":1744676100000},{"__ext__bk_46__container_name":"unify-query","_value_":3657,"_timestamp_":1744676130000},{"__ext__bk_46__container_name":"unify-query","_value_":4185,"_timestamp_":1744676160000},{"__ext__bk_46__container_name":"unify-query","_value_":3336,"_timestamp_":1744676190000},{"__ext__bk_46__container_name":"unify-query","_value_":4532,"_timestamp_":1744676220000},{"__ext__bk_46__container_name":"unify-query","_value_":3700,"_timestamp_":1744676250000},{"__ext__bk_46__container_name":"unify-query","_value_":4174,"_timestamp_":1744676280000},{"__ext__bk_46__container_name":"unify-query","_value_":3318,"_timestamp_":1744676310000},{"__ext__bk_46__container_name":"unify-query","_value_":4463,"_timestamp_":1744676340000},{"__ext__bk_46__container_name":"unify-query","_value_":3502,"_timestamp_":1744676370000},{"__ext__bk_46__container_name":"unify-query","_value_":6064,"_timestamp_":1744676400000},{"__ext__bk_46__container_name":"unify-query","_value_":3292,"_timestamp_":1744676430000},{"__ext__bk_46__container_name":"unify-query","_value_":4858,"_timestamp_":1744676460000},{"__ext__bk_46__container_name":"unify-query","_value_":3543,"_timestamp_":1744676490000},{"__ext__bk_46__container_name":"unify-query","_value_":4620,"_timestamp_":1744676520000},{"__ext__bk_46__container_name":"unify-query","_value_":3750,"_timestamp_":1744676550000},{"__ext__bk_46__container_name":"unify-query","_value_":4043,"_timestamp_":1744676580000},{"__ext__bk_46__container_name":"unify-query","_value_":3595,"_timestamp_":1744676610000},{"__ext__bk_46__container_name":"unify-query","_value_":4152,"_timestamp_":1744676640000},{"__ext__bk_46__container_name":"unify-query","_value_":3550,"_timestamp_":1744676670000},{"__ext__bk_46__container_name":"unify-query","_value_":4011,"_timestamp_":1744676700000},{"__ext__bk_46__container_name":"unify-query","_value_":3502,"_timestamp_":1744676730000},{"__ext__bk_46__container_name":"unify-query","_value_":4050,"_timestamp_":1744676760000},{"__ext__bk_46__container_name":"unify-query","_value_":3118,"_timestamp_":1744676790000},{"__ext__bk_46__container_name":"unify-query","_value_":4628,"_timestamp_":1744676820000},{"__ext__bk_46__container_name":"unify-query","_value_":3441,"_timestamp_":1744676850000},{"__ext__bk_46__container_name":"unify-query","_value_":4366,"_timestamp_":1744676880000},{"__ext__bk_46__container_name":"unify-query","_value_":3500,"_timestamp_":1744676910000},{"__ext__bk_46__container_name":"unify-query","_value_":4160,"_timestamp_":1744676940000},{"__ext__bk_46__container_name":"unify-query","_value_":3662,"_timestamp_":1744676970000},{"__ext__bk_46__container_name":"unify-query","_value_":11392,"_timestamp_":1744677000000},{"__ext__bk_46__container_name":"unify-query","_value_":18649,"_timestamp_":1744677030000},{"__ext__bk_46__container_name":"unify-query","_value_":7107,"_timestamp_":1744677060000},{"__ext__bk_46__container_name":"unify-query","_value_":9213,"_timestamp_":1744677090000},{"__ext__bk_46__container_name":"unify-query","_value_":4235,"_timestamp_":1744677120000},{"__ext__bk_46__container_name":"unify-query","_value_":3623,"_timestamp_":1744677150000},{"__ext__bk_46__container_name":"unify-query","_value_":4412,"_timestamp_":1744677180000},{"__ext__bk_46__container_name":"unify-query","_value_":3436,"_timestamp_":1744677210000},{"__ext__bk_46__container_name":"unify-query","_value_":4233,"_timestamp_":1744677240000},{"__ext__bk_46__container_name":"unify-query","_value_":3440,"_timestamp_":1744677270000},{"__ext__bk_46__container_name":"unify-query","_value_":4383,"_timestamp_":1744677300000},{"__ext__bk_46__container_name":"unify-query","_value_":3507,"_timestamp_":1744677330000},{"__ext__bk_46__container_name":"unify-query","_value_":4288,"_timestamp_":1744677360000},{"__ext__bk_46__container_name":"unify-query","_value_":3197,"_timestamp_":1744677390000},{"__ext__bk_46__container_name":"unify-query","_value_":4605,"_timestamp_":1744677420000},{"__ext__bk_46__container_name":"unify-query","_value_":3249,"_timestamp_":1744677450000},{"__ext__bk_46__container_name":"unify-query","_value_":4421,"_timestamp_":1744677480000},{"__ext__bk_46__container_name":"unify-query","_value_":2998,"_timestamp_":1744677510000},{"__ext__bk_46__container_name":"unify-query","_value_":4700,"_timestamp_":1744677540000},{"__ext__bk_46__container_name":"unify-query","_value_":3598,"_timestamp_":1744677570000},{"__ext__bk_46__container_name":"unify-query","_value_":5781,"_timestamp_":1744677600000},{"__ext__bk_46__container_name":"unify-query","_value_":3734,"_timestamp_":1744677630000},{"__ext__bk_46__container_name":"unify-query","_value_":4510,"_timestamp_":1744677660000},{"__ext__bk_46__container_name":"unify-query","_value_":3752,"_timestamp_":1744677690000},{"__ext__bk_46__container_name":"unify-query","_value_":4447,"_timestamp_":1744677720000},{"__ext__bk_46__container_name":"unify-query","_value_":3523,"_timestamp_":1744677750000},{"__ext__bk_46__container_name":"unify-query","_value_":4187,"_timestamp_":1744677780000},{"__ext__bk_46__container_name":"unify-query","_value_":3640,"_timestamp_":1744677810000},{"__ext__bk_46__container_name":"unify-query","_value_":3900,"_timestamp_":1744677840000},{"__ext__bk_46__container_name":"unify-query","_value_":3514,"_timestamp_":1744677870000},{"__ext__bk_46__container_name":"unify-query","_value_":4863,"_timestamp_":1744677900000},{"__ext__bk_46__container_name":"unify-query","_value_":3565,"_timestamp_":1744677930000},{"__ext__bk_46__container_name":"unify-query","_value_":4335,"_timestamp_":1744677960000},{"__ext__bk_46__container_name":"unify-query","_value_":3533,"_timestamp_":1744677990000},{"__ext__bk_46__container_name":"unify-query","_value_":4307,"_timestamp_":1744678020000},{"__ext__bk_46__container_name":"unify-query","_value_":3556,"_timestamp_":1744678050000},{"__ext__bk_46__container_name":"unify-query","_value_":4179,"_timestamp_":1744678080000},{"__ext__bk_46__container_name":"unify-query","_value_":3664,"_timestamp_":1744678110000},{"__ext__bk_46__container_name":"unify-query","_value_":4362,"_timestamp_":1744678140000},{"__ext__bk_46__container_name":"unify-query","_value_":3222,"_timestamp_":1744678170000},{"__ext__bk_46__container_name":"unify-query","_value_":4750,"_timestamp_":1744678200000},{"__ext__bk_46__container_name":"unify-query","_value_":3546,"_timestamp_":1744678230000},{"__ext__bk_46__container_name":"unify-query","_value_":4601,"_timestamp_":1744678260000},{"__ext__bk_46__container_name":"unify-query","_value_":3702,"_timestamp_":1744678290000},{"__ext__bk_46__container_name":"unify-query","_value_":4564,"_timestamp_":1744678320000},{"__ext__bk_46__container_name":"unify-query","_value_":3610,"_timestamp_":1744678350000},{"__ext__bk_46__container_name":"unify-query","_value_":4130,"_timestamp_":1744678380000},{"__ext__bk_46__container_name":"unify-query","_value_":3412,"_timestamp_":1744678410000},{"__ext__bk_46__container_name":"unify-query","_value_":4614,"_timestamp_":1744678440000},{"__ext__bk_46__container_name":"unify-query","_value_":3522,"_timestamp_":1744678470000},{"__ext__bk_46__container_name":"unify-query","_value_":4148,"_timestamp_":1744678500000},{"__ext__bk_46__container_name":"unify-query","_value_":3408,"_timestamp_":1744678530000},{"__ext__bk_46__container_name":"unify-query","_value_":4261,"_timestamp_":1744678560000},{"__ext__bk_46__container_name":"unify-query","_value_":3607,"_timestamp_":1744678590000},{"__ext__bk_46__container_name":"unify-query","_value_":4172,"_timestamp_":1744678620000},{"__ext__bk_46__container_name":"unify-query","_value_":3529,"_timestamp_":1744678650000},{"__ext__bk_46__container_name":"unify-query","_value_":4227,"_timestamp_":1744678680000},{"__ext__bk_46__container_name":"unify-query","_value_":3487,"_timestamp_":1744678710000},{"__ext__bk_46__container_name":"unify-query","_value_":4298,"_timestamp_":1744678740000},{"__ext__bk_46__container_name":"unify-query","_value_":3609,"_timestamp_":1744678770000},{"__ext__bk_46__container_name":"unify-query","_value_":7230,"_timestamp_":1744678800000},{"__ext__bk_46__container_name":"unify-query","_value_":3818,"_timestamp_":1744678830000},{"__ext__bk_46__container_name":"unify-query","_value_":11924,"_timestamp_":1744678860000},{"__ext__bk_46__container_name":"unify-query","_value_":27269,"_timestamp_":1744678890000},{"__ext__bk_46__container_name":"unify-query","_value_":5073,"_timestamp_":1744678920000},{"__ext__bk_46__container_name":"unify-query","_value_":3474,"_timestamp_":1744678950000},{"__ext__bk_46__container_name":"unify-query","_value_":4474,"_timestamp_":1744678980000},{"__ext__bk_46__container_name":"unify-query","_value_":3536,"_timestamp_":1744679010000},{"__ext__bk_46__container_name":"unify-query","_value_":4525,"_timestamp_":1744679040000},{"__ext__bk_46__container_name":"unify-query","_value_":3503,"_timestamp_":1744679070000},{"__ext__bk_46__container_name":"unify-query","_value_":4194,"_timestamp_":1744679100000},{"__ext__bk_46__container_name":"unify-query","_value_":3557,"_timestamp_":1744679130000},{"__ext__bk_46__container_name":"unify-query","_value_":4259,"_timestamp_":1744679160000},{"__ext__bk_46__container_name":"unify-query","_value_":3611,"_timestamp_":1744679190000},{"__ext__bk_46__container_name":"unify-query","_value_":4218,"_timestamp_":1744679220000},{"__ext__bk_46__container_name":"unify-query","_value_":3622,"_timestamp_":1744679250000},{"__ext__bk_46__container_name":"unify-query","_value_":4417,"_timestamp_":1744679280000},{"__ext__bk_46__container_name":"unify-query","_value_":3730,"_timestamp_":1744679310000},{"__ext__bk_46__container_name":"unify-query","_value_":4204,"_timestamp_":1744679340000},{"__ext__bk_46__container_name":"unify-query","_value_":3641,"_timestamp_":1744679370000},{"__ext__bk_46__container_name":"unify-query","_value_":4849,"_timestamp_":1744679400000},{"__ext__bk_46__container_name":"unify-query","_value_":3803,"_timestamp_":1744679430000},{"__ext__bk_46__container_name":"unify-query","_value_":4398,"_timestamp_":1744679460000},{"__ext__bk_46__container_name":"unify-query","_value_":3674,"_timestamp_":1744679490000},{"__ext__bk_46__container_name":"unify-query","_value_":4727,"_timestamp_":1744679520000},{"__ext__bk_46__container_name":"unify-query","_value_":3926,"_timestamp_":1744679550000},{"__ext__bk_46__container_name":"unify-query","_value_":4173,"_timestamp_":1744679580000},{"__ext__bk_46__container_name":"unify-query","_value_":3531,"_timestamp_":1744679610000},{"__ext__bk_46__container_name":"unify-query","_value_":4968,"_timestamp_":1744679640000},{"__ext__bk_46__container_name":"unify-query","_value_":3432,"_timestamp_":1744679670000},{"__ext__bk_46__container_name":"unify-query","_value_":5059,"_timestamp_":1744679700000},{"__ext__bk_46__container_name":"unify-query","_value_":3560,"_timestamp_":1744679730000},{"__ext__bk_46__container_name":"unify-query","_value_":4087,"_timestamp_":1744679760000},{"__ext__bk_46__container_name":"unify-query","_value_":3590,"_timestamp_":1744679790000},{"__ext__bk_46__container_name":"unify-query","_value_":4436,"_timestamp_":1744679820000},{"__ext__bk_46__container_name":"unify-query","_value_":5299,"_timestamp_":1744679850000},{"__ext__bk_46__container_name":"unify-query","_value_":4320,"_timestamp_":1744679880000},{"__ext__bk_46__container_name":"unify-query","_value_":3861,"_timestamp_":1744679910000},{"__ext__bk_46__container_name":"unify-query","_value_":4511,"_timestamp_":1744679940000},{"__ext__bk_46__container_name":"unify-query","_value_":3711,"_timestamp_":1744679970000},{"__ext__bk_46__container_name":"unify-query","_value_":6021,"_timestamp_":1744680000000},{"__ext__bk_46__container_name":"unify-query","_value_":3942,"_timestamp_":1744680030000},{"__ext__bk_46__container_name":"unify-query","_value_":4800,"_timestamp_":1744680060000},{"__ext__bk_46__container_name":"unify-query","_value_":3681,"_timestamp_":1744680090000},{"__ext__bk_46__container_name":"unify-query","_value_":4592,"_timestamp_":1744680120000},{"__ext__bk_46__container_name":"unify-query","_value_":3560,"_timestamp_":1744680150000},{"__ext__bk_46__container_name":"unify-query","_value_":4194,"_timestamp_":1744680180000},{"__ext__bk_46__container_name":"unify-query","_value_":3490,"_timestamp_":1744680210000},{"__ext__bk_46__container_name":"unify-query","_value_":4971,"_timestamp_":1744680240000},{"__ext__bk_46__container_name":"unify-query","_value_":4009,"_timestamp_":1744680270000},{"__ext__bk_46__container_name":"unify-query","_value_":4837,"_timestamp_":1744680300000},{"__ext__bk_46__container_name":"unify-query","_value_":3227,"_timestamp_":1744680330000},{"__ext__bk_46__container_name":"unify-query","_value_":4531,"_timestamp_":1744680360000},{"__ext__bk_46__container_name":"unify-query","_value_":2888,"_timestamp_":1744680390000},{"__ext__bk_46__container_name":"unify-query","_value_":5083,"_timestamp_":1744680420000},{"__ext__bk_46__container_name":"unify-query","_value_":3557,"_timestamp_":1744680450000},{"__ext__bk_46__container_name":"unify-query","_value_":4207,"_timestamp_":1744680480000},{"__ext__bk_46__container_name":"unify-query","_value_":3373,"_timestamp_":1744680510000},{"__ext__bk_46__container_name":"unify-query","_value_":4482,"_timestamp_":1744680540000},{"__ext__bk_46__container_name":"unify-query","_value_":3110,"_timestamp_":1744680570000},{"__ext__bk_46__container_name":"unify-query","_value_":13551,"_timestamp_":1744680600000},{"__ext__bk_46__container_name":"unify-query","_value_":17159,"_timestamp_":1744680630000},{"__ext__bk_46__container_name":"unify-query","_value_":6284,"_timestamp_":1744680660000},{"__ext__bk_46__container_name":"unify-query","_value_":9924,"_timestamp_":1744680690000},{"__ext__bk_46__container_name":"unify-query","_value_":4547,"_timestamp_":1744680720000},{"__ext__bk_46__container_name":"unify-query","_value_":3474,"_timestamp_":1744680750000},{"__ext__bk_46__container_name":"unify-query","_value_":4312,"_timestamp_":1744680780000},{"__ext__bk_46__container_name":"unify-query","_value_":3689,"_timestamp_":1744680810000},{"__ext__bk_46__container_name":"unify-query","_value_":4680,"_timestamp_":1744680840000},{"__ext__bk_46__container_name":"unify-query","_value_":3609,"_timestamp_":1744680870000},{"__ext__bk_46__container_name":"unify-query","_value_":4886,"_timestamp_":1744680900000},{"__ext__bk_46__container_name":"unify-query","_value_":3842,"_timestamp_":1744680930000},{"__ext__bk_46__container_name":"unify-query","_value_":4810,"_timestamp_":1744680960000},{"__ext__bk_46__container_name":"unify-query","_value_":4102,"_timestamp_":1744680990000},{"__ext__bk_46__container_name":"unify-query","_value_":4594,"_timestamp_":1744681020000},{"__ext__bk_46__container_name":"unify-query","_value_":4168,"_timestamp_":1744681050000},{"__ext__bk_46__container_name":"unify-query","_value_":4562,"_timestamp_":1744681080000},{"__ext__bk_46__container_name":"unify-query","_value_":4506,"_timestamp_":1744681110000},{"__ext__bk_46__container_name":"unify-query","_value_":5243,"_timestamp_":1744681140000},{"__ext__bk_46__container_name":"unify-query","_value_":5135,"_timestamp_":1744681170000},{"__ext__bk_46__container_name":"unify-query","_value_":6671,"_timestamp_":1744681200000},{"__ext__bk_46__container_name":"unify-query","_value_":3806,"_timestamp_":1744681230000},{"__ext__bk_46__container_name":"unify-query","_value_":4535,"_timestamp_":1744681260000},{"__ext__bk_46__container_name":"unify-query","_value_":3721,"_timestamp_":1744681290000},{"__ext__bk_46__container_name":"unify-query","_value_":4799,"_timestamp_":1744681320000},{"__ext__bk_46__container_name":"unify-query","_value_":3909,"_timestamp_":1744681350000},{"__ext__bk_46__container_name":"unify-query","_value_":4261,"_timestamp_":1744681380000},{"__ext__bk_46__container_name":"unify-query","_value_":3671,"_timestamp_":1744681410000},{"__ext__bk_46__container_name":"unify-query","_value_":4359,"_timestamp_":1744681440000},{"__ext__bk_46__container_name":"unify-query","_value_":4063,"_timestamp_":1744681470000},{"__ext__bk_46__container_name":"unify-query","_value_":5231,"_timestamp_":1744681500000},{"__ext__bk_46__container_name":"unify-query","_value_":3778,"_timestamp_":1744681530000},{"__ext__bk_46__container_name":"unify-query","_value_":4684,"_timestamp_":1744681560000},{"__ext__bk_46__container_name":"unify-query","_value_":4072,"_timestamp_":1744681590000},{"__ext__bk_46__container_name":"unify-query","_value_":5029,"_timestamp_":1744681620000},{"__ext__bk_46__container_name":"unify-query","_value_":3700,"_timestamp_":1744681650000},{"__ext__bk_46__container_name":"unify-query","_value_":4670,"_timestamp_":1744681680000},{"__ext__bk_46__container_name":"unify-query","_value_":3557,"_timestamp_":1744681710000},{"__ext__bk_46__container_name":"unify-query","_value_":4590,"_timestamp_":1744681740000},{"__ext__bk_46__container_name":"unify-query","_value_":3041,"_timestamp_":1744681770000},{"__ext__bk_46__container_name":"unify-query","_value_":5043,"_timestamp_":1744681800000},{"__ext__bk_46__container_name":"unify-query","_value_":3530,"_timestamp_":1744681830000},{"__ext__bk_46__container_name":"unify-query","_value_":6807,"_timestamp_":1744681860000},{"__ext__bk_46__container_name":"unify-query","_value_":4455,"_timestamp_":1744681890000},{"__ext__bk_46__container_name":"unify-query","_value_":6841,"_timestamp_":1744681920000},{"__ext__bk_46__container_name":"unify-query","_value_":4519,"_timestamp_":1744681950000},{"__ext__bk_46__container_name":"unify-query","_value_":6617,"_timestamp_":1744681980000},{"__ext__bk_46__container_name":"unify-query","_value_":4633,"_timestamp_":1744682010000},{"__ext__bk_46__container_name":"unify-query","_value_":5997,"_timestamp_":1744682040000},{"__ext__bk_46__container_name":"unify-query","_value_":4446,"_timestamp_":1744682070000},{"__ext__bk_46__container_name":"unify-query","_value_":5569,"_timestamp_":1744682100000},{"__ext__bk_46__container_name":"unify-query","_value_":4324,"_timestamp_":1744682130000},{"__ext__bk_46__container_name":"unify-query","_value_":5354,"_timestamp_":1744682160000},{"__ext__bk_46__container_name":"unify-query","_value_":7245,"_timestamp_":1744682190000},{"__ext__bk_46__container_name":"unify-query","_value_":5258,"_timestamp_":1744682220000},{"__ext__bk_46__container_name":"unify-query","_value_":4296,"_timestamp_":1744682250000},{"__ext__bk_46__container_name":"unify-query","_value_":5349,"_timestamp_":1744682280000},{"__ext__bk_46__container_name":"unify-query","_value_":4479,"_timestamp_":1744682310000},{"__ext__bk_46__container_name":"unify-query","_value_":5127,"_timestamp_":1744682340000},{"__ext__bk_46__container_name":"unify-query","_value_":4006,"_timestamp_":1744682370000},{"__ext__bk_46__container_name":"unify-query","_value_":19058,"_timestamp_":1744682400000},{"__ext__bk_46__container_name":"unify-query","_value_":14501,"_timestamp_":1744682430000},{"__ext__bk_46__container_name":"unify-query","_value_":3810,"_timestamp_":1744682460000},{"__ext__bk_46__container_name":"unify-query","_value_":12368,"_timestamp_":1744682490000},{"__ext__bk_46__container_name":"unify-query","_value_":6976,"_timestamp_":1744682520000},{"__ext__bk_46__container_name":"unify-query","_value_":4399,"_timestamp_":1744682550000},{"__ext__bk_46__container_name":"unify-query","_value_":5482,"_timestamp_":1744682580000},{"__ext__bk_46__container_name":"unify-query","_value_":4524,"_timestamp_":1744682610000},{"__ext__bk_46__container_name":"unify-query","_value_":5478,"_timestamp_":1744682640000},{"__ext__bk_46__container_name":"unify-query","_value_":4920,"_timestamp_":1744682670000},{"__ext__bk_46__container_name":"unify-query","_value_":5347,"_timestamp_":1744682700000},{"__ext__bk_46__container_name":"unify-query","_value_":4427,"_timestamp_":1744682730000},{"__ext__bk_46__container_name":"unify-query","_value_":5102,"_timestamp_":1744682760000},{"__ext__bk_46__container_name":"unify-query","_value_":4441,"_timestamp_":1744682790000},{"__ext__bk_46__container_name":"unify-query","_value_":5596,"_timestamp_":1744682820000},{"__ext__bk_46__container_name":"unify-query","_value_":4888,"_timestamp_":1744682850000},{"__ext__bk_46__container_name":"unify-query","_value_":5306,"_timestamp_":1744682880000},{"__ext__bk_46__container_name":"unify-query","_value_":4825,"_timestamp_":1744682910000},{"__ext__bk_46__container_name":"unify-query","_value_":5897,"_timestamp_":1744682940000},{"__ext__bk_46__container_name":"unify-query","_value_":4481,"_timestamp_":1744682970000},{"__ext__bk_46__container_name":"unify-query","_value_":6086,"_timestamp_":1744683000000},{"__ext__bk_46__container_name":"unify-query","_value_":4910,"_timestamp_":1744683030000},{"__ext__bk_46__container_name":"unify-query","_value_":5676,"_timestamp_":1744683060000},{"__ext__bk_46__container_name":"unify-query","_value_":3626,"_timestamp_":1744683090000},{"__ext__bk_46__container_name":"unify-query","_value_":6929,"_timestamp_":1744683120000},{"__ext__bk_46__container_name":"unify-query","_value_":4601,"_timestamp_":1744683150000},{"__ext__bk_46__container_name":"unify-query","_value_":5525,"_timestamp_":1744683180000},{"__ext__bk_46__container_name":"unify-query","_value_":4500,"_timestamp_":1744683210000},{"__ext__bk_46__container_name":"unify-query","_value_":5617,"_timestamp_":1744683240000},{"__ext__bk_46__container_name":"unify-query","_value_":4503,"_timestamp_":1744683270000},{"__ext__bk_46__container_name":"unify-query","_value_":6328,"_timestamp_":1744683300000},{"__ext__bk_46__container_name":"unify-query","_value_":4557,"_timestamp_":1744683330000},{"__ext__bk_46__container_name":"unify-query","_value_":5356,"_timestamp_":1744683360000},{"__ext__bk_46__container_name":"unify-query","_value_":4413,"_timestamp_":1744683390000},{"__ext__bk_46__container_name":"unify-query","_value_":5335,"_timestamp_":1744683420000},{"__ext__bk_46__container_name":"unify-query","_value_":4640,"_timestamp_":1744683450000},{"__ext__bk_46__container_name":"unify-query","_value_":5399,"_timestamp_":1744683480000},{"__ext__bk_46__container_name":"unify-query","_value_":4298,"_timestamp_":1744683510000},{"__ext__bk_46__container_name":"unify-query","_value_":5415,"_timestamp_":1744683540000},{"__ext__bk_46__container_name":"unify-query","_value_":4540,"_timestamp_":1744683570000},{"__ext__bk_46__container_name":"unify-query","_value_":6949,"_timestamp_":1744683600000},{"__ext__bk_46__container_name":"unify-query","_value_":4574,"_timestamp_":1744683630000},{"__ext__bk_46__container_name":"unify-query","_value_":5757,"_timestamp_":1744683660000},{"__ext__bk_46__container_name":"unify-query","_value_":4669,"_timestamp_":1744683690000},{"__ext__bk_46__container_name":"unify-query","_value_":5706,"_timestamp_":1744683720000},{"__ext__bk_46__container_name":"unify-query","_value_":4472,"_timestamp_":1744683750000},{"__ext__bk_46__container_name":"unify-query","_value_":5386,"_timestamp_":1744683780000},{"__ext__bk_46__container_name":"unify-query","_value_":4490,"_timestamp_":1744683810000},{"__ext__bk_46__container_name":"unify-query","_value_":5104,"_timestamp_":1744683840000},{"__ext__bk_46__container_name":"unify-query","_value_":4201,"_timestamp_":1744683870000},{"__ext__bk_46__container_name":"unify-query","_value_":5979,"_timestamp_":1744683900000},{"__ext__bk_46__container_name":"unify-query","_value_":4853,"_timestamp_":1744683930000},{"__ext__bk_46__container_name":"unify-query","_value_":6691,"_timestamp_":1744683960000},{"__ext__bk_46__container_name":"unify-query","_value_":4572,"_timestamp_":1744683990000},{"__ext__bk_46__container_name":"unify-query","_value_":5554,"_timestamp_":1744684020000},{"__ext__bk_46__container_name":"unify-query","_value_":5244,"_timestamp_":1744684050000},{"__ext__bk_46__container_name":"unify-query","_value_":5392,"_timestamp_":1744684080000},{"__ext__bk_46__container_name":"unify-query","_value_":4550,"_timestamp_":1744684110000},{"__ext__bk_46__container_name":"unify-query","_value_":520,"_timestamp_":1744684140000}],"stage_elapsed_time_mills":{"check_query_syntax":2,"query_db":52,"get_query_driver":0,"match_query_forbidden_config":0,"convert_query_statement":8,"connect_db":55,"match_query_routing_rule":0,"check_permission":73,"check_query_semantic":0,"pick_valid_storage":1},"total_record_size":269248,"timetaken":0.191,"result_schema":[{"field_type":"string","field_name":"__c0","field_alias":"__ext__bk_46__container_name","field_index":0},{"field_type":"long","field_name":"__c1","field_alias":"_value_","field_index":1},{"field_type":"long","field_name":"__c2","field_alias":"_timestamp_","field_index":2}],"bksql_call_elapsed_time":0,"device":"doris","result_table_ids":["2_bklog_bkunify_query_doris"]},"errors":null,"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000"}`,
	})

	for i, c := range map[string]struct {
		queryTs *structured.QueryTs
		result  string
	}{
		"查询 1 条原始数据，按照字段正向排序": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "gseIndex",
						Limit:         1,
						From:          0,
						ReferenceName: "a",
					},
				},
				OrderBy: structured.OrderBy{
					"_value",
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     false,
				SpaceUid:    spaceUid,
			},
			result: `{
  "series" : [ {
    "name" : "_result0",
    "metric_name" : "",
    "columns" : [ "_time", "_value" ],
    "types" : [ "float", "float" ],
    "group_keys" : [ "__ext.container_id", "__ext.container_image", "__ext.container_name", "__ext.io_kubernetes_pod", "__ext.io_kubernetes_pod_ip", "__ext.io_kubernetes_pod_namespace", "__ext.io_kubernetes_pod_uid", "__ext.io_kubernetes_workload_name", "__ext.io_kubernetes_workload_type", "__name__", "cloudid", "file", "gseindex", "iterationindex", "level", "log", "message", "path", "report_time", "serverip", "time", "trace_id" ],
    "group_values" : [ "375597ee636fd5d53cb7b0958823d9ba6534bd24cd698e485c41ca2f01b78ed2", "sha256:3a0506f06f1467e93c3a582203aac1a7501e77091572ec9612ddeee4a4dbbdb8", "unify-query", "bk-datalink-unify-query-6df8bcc4c9-rk4sc", "127.0.0.1", "blueking", "558c5b17-b221-47e1-aa66-036cc9b43e2a", "bk-datalink-unify-query-6df8bcc4c9", "ReplicaSet", "bklog:result_table:doris:gseIndex", "0", "http/handler.go:320", "2450131", "19", "info", "2025-04-14T20:22:59.982Z\tinfo\thttp/handler.go:320\t[5108397435e997364f8dc1251533e65e] header: map[Accept:[*/*] Accept-Encoding:[gzip, deflate] Bk-Query-Source:[strategy:9155] Connection:[keep-alive] Content-Length:[863] Content-Type:[application/json] Traceparent:[00-5108397435e997364f8dc1251533e65e-ca18e72c0f0eafd4-00] User-Agent:[python-requests/2.31.0] X-Bk-Scope-Space-Uid:[bkcc__2]], body: {\"space_uid\":\"bkcc__2\",\"query_list\":[{\"field_name\":\"bscp_config_consume_total_file_change_count\",\"is_regexp\":false,\"function\":[{\"method\":\"mean\",\"without\":false,\"dimensions\":[\"app\",\"biz\",\"clientType\"]}],\"time_aggregation\":{\"function\":\"increase\",\"window\":\"1m\"},\"is_dom_sampled\":false,\"reference_name\":\"a\",\"dimensions\":[\"app\",\"biz\",\"clientType\"],\"conditions\":{\"field_list\":[{\"field_name\":\"releaseChangeStatus\",\"value\":[\"Failed\"],\"op\":\"contains\"},{\"field_name\":\"bcs_cluster_id\",\"value\":[\"BCS-K8S-00000\"],\"op\":\"contains\"}],\"condition_list\":[\"and\"]},\"keep_columns\":[\"_time\",\"a\",\"app\",\"biz\",\"clientType\"],\"query_string\":\"\"}],\"metric_merge\":\"a\",\"start_time\":\"1744660260\",\"end_time\":\"1744662120\",\"step\":\"60s\",\"timezone\":\"Asia/Shanghai\",\"instant\":false}", " header: map[Accept:[*/*] Accept-Encoding:[gzip, deflate] Bk-Query-Source:[strategy:9155] Connection:[keep-alive] Content-Length:[863] Content-Type:[application/json] Traceparent:[00-5108397435e997364f8dc1251533e65e-ca18e72c0f0eafd4-00] User-Agent:[python-requests/2.31.0] X-Bk-Scope-Space-Uid:[bkcc__2]], body: {\"space_uid\":\"bkcc__2\",\"query_list\":[{\"field_name\":\"bscp_config_consume_total_file_change_count\",\"is_regexp\":false,\"function\":[{\"method\":\"mean\",\"without\":false,\"dimensions\":[\"app\",\"biz\",\"clientType\"]}],\"time_aggregation\":{\"function\":\"increase\",\"window\":\"1m\"},\"is_dom_sampled\":false,\"reference_name\":\"a\",\"dimensions\":[\"app\",\"biz\",\"clientType\"],\"conditions\":{\"field_list\":[{\"field_name\":\"releaseChangeStatus\",\"value\":[\"Failed\"],\"op\":\"contains\"},{\"field_name\":\"bcs_cluster_id\",\"value\":[\"BCS-K8S-00000\"],\"op\":\"contains\"}],\"condition_list\":[\"and\"]},\"keep_columns\":[\"_time\",\"a\",\"app\",\"biz\",\"clientType\"],\"query_string\":\"\"}],\"metric_merge\":\"a\",\"start_time\":\"1744660260\",\"end_time\":\"1744662120\",\"step\":\"60s\",\"timezone\":\"Asia/Shanghai\",\"instant\":false}", "/var/host/data/bcs/lib/docker/containers/375597ee636fd5d53cb7b0958823d9ba6534bd24cd698e485c41ca2f01b78ed2/375597ee636fd5d53cb7b0958823d9ba6534bd24cd698e485c41ca2f01b78ed2-json.log", "2025-04-14T20:22:59.982Z", "127.0.0.1", "1744662180000", "5108397435e997364f8dc1251533e65e" ],
    "values" : [ [ 1744662480000, 2450131 ] ]
  } ]
}`,
		},
		"根据维度 __ext.container_name 进行 count 聚合，同时用值正向排序": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "gseIndex",
						ReferenceName: "a",
						TimeAggregation: structured.TimeAggregation{
							Function: "count_over_time",
							Window:   "30s",
						},
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method:     "sum",
								Dimensions: []string{"__ext.container_name"},
							},
							{
								Method: "topk",
								VArgsList: []interface{}{
									5,
								},
							},
						},
					},
				},
				OrderBy: structured.OrderBy{
					"_value",
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     false,
				SpaceUid:    spaceUid,
				Step:        "30s",
			},
			result: `{
  "series" : [ {
    "name" : "_result0",
    "metric_name" : "",
    "columns" : [ "_time", "_value" ],
    "types" : [ "float", "float" ],
    "group_keys" : [ "__ext.container_name" ],
    "group_values" : [ "unify-query" ],
    "values" : [ [ 1744662510000, 3684 ], [ 1744662540000, 4012 ], [ 1744662570000, 3671 ], [ 1744662600000, 17092 ], [ 1744662630000, 12881 ], [ 1744662660000, 5902 ], [ 1744662690000, 10443 ], [ 1744662720000, 4388 ], [ 1744662750000, 3357 ], [ 1744662780000, 4381 ], [ 1744662810000, 3683 ], [ 1744662840000, 4353 ], [ 1744662870000, 3441 ], [ 1744662900000, 4251 ], [ 1744662930000, 3476 ], [ 1744662960000, 4036 ], [ 1744662990000, 3549 ], [ 1744663020000, 4351 ], [ 1744663050000, 3651 ], [ 1744663080000, 4096 ], [ 1744663110000, 3618 ], [ 1744663140000, 4100 ], [ 1744663170000, 3622 ], [ 1744663200000, 6044 ], [ 1744663230000, 3766 ], [ 1744663260000, 4461 ], [ 1744663290000, 3783 ], [ 1744663320000, 4559 ], [ 1744663350000, 3634 ], [ 1744663380000, 3869 ], [ 1744663410000, 3249 ], [ 1744663440000, 4473 ], [ 1744663470000, 3514 ], [ 1744663500000, 4923 ], [ 1744663530000, 3379 ], [ 1744663560000, 4489 ], [ 1744663590000, 3411 ], [ 1744663620000, 4374 ], [ 1744663650000, 3370 ], [ 1744663680000, 4310 ], [ 1744663710000, 3609 ], [ 1744663740000, 4318 ], [ 1744663770000, 3570 ], [ 1744663800000, 4334 ], [ 1744663830000, 3767 ], [ 1744663860000, 4455 ], [ 1744663890000, 3703 ], [ 1744663920000, 4511 ], [ 1744663950000, 3667 ], [ 1744663980000, 3998 ], [ 1744664010000, 3579 ], [ 1744664040000, 4156 ], [ 1744664070000, 3340 ], [ 1744664100000, 4344 ], [ 1744664130000, 3590 ], [ 1744664160000, 4161 ], [ 1744664190000, 3484 ], [ 1744664220000, 4273 ], [ 1744664250000, 3494 ], [ 1744664280000, 4230 ], [ 1744664310000, 3619 ], [ 1744664340000, 4013 ], [ 1744664370000, 3565 ], [ 1744664400000, 18144 ], [ 1744664430000, 13615 ], [ 1744664460000, 3178 ], [ 1744664490000, 13044 ], [ 1744664520000, 4767 ], [ 1744664550000, 3528 ], [ 1744664580000, 4316 ], [ 1744664610000, 3317 ], [ 1744664640000, 4395 ], [ 1744664670000, 3599 ], [ 1744664700000, 4149 ], [ 1744664730000, 3474 ], [ 1744664760000, 4201 ], [ 1744664790000, 3384 ], [ 1744664820000, 4442 ], [ 1744664850000, 3559 ], [ 1744664880000, 4166 ], [ 1744664910000, 3438 ], [ 1744664940000, 4244 ], [ 1744664970000, 3640 ], [ 1744665000000, 4305 ], [ 1744665030000, 3771 ], [ 1744665060000, 4485 ], [ 1744665090000, 3842 ], [ 1744665120000, 4423 ], [ 1744665150000, 3610 ], [ 1744665180000, 4125 ], [ 1744665210000, 3500 ], [ 1744665240000, 4252 ], [ 1744665270000, 3427 ], [ 1744665300000, 5089 ], [ 1744665330000, 3450 ], [ 1744665360000, 4349 ], [ 1744665390000, 3188 ], [ 1744665420000, 4556 ], [ 1744665450000, 3372 ], [ 1744665480000, 4408 ], [ 1744665510000, 3445 ], [ 1744665540000, 4213 ], [ 1744665570000, 3408 ], [ 1744665600000, 6235 ], [ 1744665630000, 3641 ], [ 1744665660000, 4577 ], [ 1744665690000, 3719 ], [ 1744665720000, 4548 ], [ 1744665750000, 3420 ], [ 1744665780000, 4246 ], [ 1744665810000, 3359 ], [ 1744665840000, 4332 ], [ 1744665870000, 3422 ], [ 1744665900000, 4229 ], [ 1744665930000, 3610 ], [ 1744665960000, 4119 ], [ 1744665990000, 3570 ], [ 1744666020000, 4144 ], [ 1744666050000, 3302 ], [ 1744666080000, 4398 ], [ 1744666110000, 3559 ], [ 1744666140000, 4097 ], [ 1744666170000, 3315 ], [ 1744666200000, 16721 ], [ 1744666230000, 13631 ], [ 1744666260000, 2982 ], [ 1744666290000, 11858 ], [ 1744666320000, 5515 ], [ 1744666350000, 2869 ], [ 1744666380000, 4795 ], [ 1744666410000, 3603 ], [ 1744666440000, 4204 ], [ 1744666470000, 3264 ], [ 1744666500000, 4377 ], [ 1744666530000, 3443 ], [ 1744666560000, 4307 ], [ 1744666590000, 3459 ], [ 1744666620000, 4342 ], [ 1744666650000, 3598 ], [ 1744666680000, 4052 ], [ 1744666710000, 3577 ], [ 1744666740000, 4128 ], [ 1744666770000, 3499 ], [ 1744666800000, 6209 ], [ 1744666830000, 3575 ], [ 1744666860000, 4543 ], [ 1744666890000, 3604 ], [ 1744666920000, 4579 ], [ 1744666950000, 3531 ], [ 1744666980000, 4314 ], [ 1744667010000, 3416 ], [ 1744667040000, 4320 ], [ 1744667070000, 3488 ], [ 1744667100000, 5054 ], [ 1744667130000, 3525 ], [ 1744667160000, 4313 ], [ 1744667190000, 3607 ], [ 1744667220000, 4118 ], [ 1744667250000, 3350 ], [ 1744667280000, 4280 ], [ 1744667310000, 3634 ], [ 1744667340000, 4174 ], [ 1744667370000, 3807 ], [ 1744667400000, 4358 ], [ 1744667430000, 3595 ], [ 1744667460000, 4630 ], [ 1744667490000, 3845 ], [ 1744667520000, 4361 ], [ 1744667550000, 3572 ], [ 1744667580000, 4095 ], [ 1744667610000, 3535 ], [ 1744667640000, 4200 ], [ 1744667670000, 3390 ], [ 1744667700000, 4262 ], [ 1744667730000, 3398 ], [ 1744667760000, 4320 ], [ 1744667790000, 3429 ], [ 1744667820000, 4288 ], [ 1744667850000, 3482 ], [ 1744667880000, 4166 ], [ 1744667910000, 3612 ], [ 1744667940000, 4194 ], [ 1744667970000, 3423 ], [ 1744668000000, 18203 ], [ 1744668030000, 13685 ], [ 1744668060000, 3281 ], [ 1744668090000, 12556 ], [ 1744668120000, 4893 ], [ 1744668150000, 3607 ], [ 1744668180000, 4336 ], [ 1744668210000, 3609 ], [ 1744668240000, 4097 ], [ 1744668270000, 3669 ], [ 1744668300000, 3997 ], [ 1744668330000, 3494 ], [ 1744668360000, 4172 ], [ 1744668390000, 3523 ], [ 1744668420000, 3877 ], [ 1744668450000, 3565 ], [ 1744668480000, 4230 ], [ 1744668510000, 3469 ], [ 1744668540000, 4243 ], [ 1744668570000, 3304 ], [ 1744668600000, 4690 ], [ 1744668630000, 3717 ], [ 1744668660000, 4618 ], [ 1744668690000, 3732 ], [ 1744668720000, 4477 ], [ 1744668750000, 3615 ], [ 1744668780000, 4154 ], [ 1744668810000, 3367 ], [ 1744668840000, 4193 ], [ 1744668870000, 3592 ], [ 1744668900000, 4971 ], [ 1744668930000, 3359 ], [ 1744668960000, 4540 ], [ 1744668990000, 3406 ], [ 1744669020000, 4375 ], [ 1744669050000, 3386 ], [ 1744669080000, 4281 ], [ 1744669110000, 3410 ], [ 1744669140000, 4545 ], [ 1744669170000, 3724 ], [ 1744669200000, 5903 ], [ 1744669230000, 3672 ], [ 1744669260000, 4413 ], [ 1744669290000, 3792 ], [ 1744669320000, 4422 ], [ 1744669350000, 3718 ], [ 1744669380000, 4213 ], [ 1744669410000, 3622 ], [ 1744669440000, 4043 ], [ 1744669470000, 3542 ], [ 1744669500000, 4179 ], [ 1744669530000, 3368 ], [ 1744669560000, 4354 ], [ 1744669590000, 3368 ], [ 1744669620000, 4229 ], [ 1744669650000, 3458 ], [ 1744669680000, 4310 ], [ 1744669710000, 3512 ], [ 1744669740000, 4188 ], [ 1744669770000, 3436 ], [ 1744669800000, 12171 ], [ 1744669830000, 18129 ], [ 1744669860000, 7142 ], [ 1744669890000, 9153 ], [ 1744669920000, 4566 ], [ 1744669950000, 3225 ], [ 1744669980000, 4378 ], [ 1744670010000, 3623 ], [ 1744670040000, 4266 ], [ 1744670070000, 3645 ], [ 1744670100000, 4043 ], [ 1744670130000, 3350 ], [ 1744670160000, 4333 ], [ 1744670190000, 3489 ], [ 1744670220000, 4303 ], [ 1744670250000, 3560 ], [ 1744670280000, 4121 ], [ 1744670310000, 3374 ], [ 1744670340000, 4362 ], [ 1744670370000, 3242 ], [ 1744670400000, 6416 ], [ 1744670430000, 3697 ], [ 1744670460000, 4506 ], [ 1744670490000, 3749 ], [ 1744670520000, 4587 ], [ 1744670550000, 3538 ], [ 1744670580000, 4221 ], [ 1744670610000, 3476 ], [ 1744670640000, 4227 ], [ 1744670670000, 3587 ], [ 1744670700000, 4848 ], [ 1744670730000, 3551 ], [ 1744670760000, 4068 ], [ 1744670790000, 3387 ], [ 1744670820000, 4366 ], [ 1744670850000, 3635 ], [ 1744670880000, 4256 ], [ 1744670910000, 3690 ], [ 1744670940000, 4155 ], [ 1744670970000, 3318 ], [ 1744671000000, 4661 ], [ 1744671030000, 3494 ], [ 1744671060000, 4442 ], [ 1744671090000, 3643 ], [ 1744671120000, 4755 ], [ 1744671150000, 3607 ], [ 1744671180000, 4284 ], [ 1744671210000, 3258 ], [ 1744671240000, 4453 ], [ 1744671270000, 3431 ], [ 1744671300000, 4231 ], [ 1744671330000, 3623 ], [ 1744671360000, 3907 ], [ 1744671390000, 3524 ], [ 1744671420000, 4438 ], [ 1744671450000, 3547 ], [ 1744671480000, 4033 ], [ 1744671510000, 3632 ], [ 1744671540000, 4162 ], [ 1744671570000, 3588 ], [ 1744671600000, 16444 ], [ 1744671630000, 15396 ], [ 1744671660000, 3024 ], [ 1744671690000, 12656 ], [ 1744671720000, 4733 ], [ 1744671750000, 3766 ], [ 1744671780000, 4388 ], [ 1744671810000, 3340 ], [ 1744671840000, 4487 ], [ 1744671870000, 3549 ], [ 1744671900000, 4154 ], [ 1744671930000, 3406 ], [ 1744671960000, 4314 ], [ 1744671990000, 3472 ], [ 1744672020000, 4309 ], [ 1744672050000, 3458 ], [ 1744672080000, 4191 ], [ 1744672110000, 3475 ], [ 1744672140000, 4194 ], [ 1744672170000, 3525 ], [ 1744672200000, 4445 ], [ 1744672230000, 3822 ], [ 1744672260000, 4346 ], [ 1744672290000, 3700 ], [ 1744672320000, 4615 ], [ 1744672350000, 3591 ], [ 1744672380000, 4056 ], [ 1744672410000, 3544 ], [ 1744672440000, 4188 ], [ 1744672470000, 3647 ], [ 1744672500000, 4887 ], [ 1744672530000, 3450 ], [ 1744672560000, 4302 ], [ 1744672590000, 3425 ], [ 1744672620000, 4320 ], [ 1744672650000, 3532 ], [ 1744672680000, 4282 ], [ 1744672710000, 3571 ], [ 1744672740000, 4182 ], [ 1744672770000, 3210 ], [ 1744672800000, 6383 ], [ 1744672830000, 3622 ], [ 1744672860000, 4408 ], [ 1744672890000, 3611 ], [ 1744672920000, 4795 ], [ 1744672950000, 3632 ], [ 1744672980000, 4102 ], [ 1744673010000, 3534 ], [ 1744673040000, 4212 ], [ 1744673070000, 3380 ], [ 1744673100000, 4289 ], [ 1744673130000, 3565 ], [ 1744673160000, 4120 ], [ 1744673190000, 3526 ], [ 1744673220000, 4200 ], [ 1744673250000, 3302 ], [ 1744673280000, 4370 ], [ 1744673310000, 3462 ], [ 1744673340000, 4223 ], [ 1744673370000, 3564 ], [ 1744673400000, 12072 ], [ 1744673430000, 17986 ], [ 1744673460000, 4089 ], [ 1744673490000, 12000 ], [ 1744673520000, 4790 ], [ 1744673550000, 3637 ], [ 1744673580000, 4177 ], [ 1744673610000, 3438 ], [ 1744673640000, 4465 ], [ 1744673670000, 3627 ], [ 1744673700000, 4131 ], [ 1744673730000, 3396 ], [ 1744673760000, 4395 ], [ 1744673790000, 3638 ], [ 1744673820000, 4093 ], [ 1744673850000, 3584 ], [ 1744673880000, 4082 ], [ 1744673910000, 3475 ], [ 1744673940000, 4051 ], [ 1744673970000, 3354 ], [ 1744674000000, 6296 ], [ 1744674030000, 3473 ], [ 1744674060000, 4412 ], [ 1744674090000, 3793 ], [ 1744674120000, 4391 ], [ 1744674150000, 3836 ], [ 1744674180000, 4190 ], [ 1744674210000, 3478 ], [ 1744674240000, 4230 ], [ 1744674270000, 3488 ], [ 1744674300000, 4964 ], [ 1744674330000, 3455 ], [ 1744674360000, 4116 ], [ 1744674390000, 3250 ], [ 1744674420000, 4494 ], [ 1744674450000, 3326 ], [ 1744674480000, 4590 ], [ 1744674510000, 3580 ], [ 1744674540000, 4368 ], [ 1744674570000, 3685 ], [ 1744674600000, 4381 ], [ 1744674630000, 3699 ], [ 1744674660000, 4513 ], [ 1744674690000, 3729 ], [ 1744674720000, 4500 ], [ 1744674750000, 3639 ], [ 1744674780000, 4018 ], [ 1744674810000, 3587 ], [ 1744674840000, 4168 ], [ 1744674870000, 3389 ], [ 1744674900000, 4289 ], [ 1744674930000, 3540 ], [ 1744674960000, 4106 ], [ 1744674990000, 3478 ], [ 1744675020000, 4268 ], [ 1744675050000, 3577 ], [ 1744675080000, 4087 ], [ 1744675110000, 3511 ], [ 1744675140000, 4174 ], [ 1744675170000, 3573 ], [ 1744675200000, 17095 ], [ 1744675230000, 14907 ], [ 1744675260000, 6455 ], [ 1744675290000, 9818 ], [ 1744675320000, 5253 ], [ 1744675350000, 3567 ], [ 1744675380000, 4047 ], [ 1744675410000, 3342 ], [ 1744675440000, 4605 ], [ 1744675470000, 3394 ], [ 1744675500000, 4260 ], [ 1744675530000, 3373 ], [ 1744675560000, 4341 ], [ 1744675590000, 3559 ], [ 1744675620000, 4188 ], [ 1744675650000, 3519 ], [ 1744675680000, 4143 ], [ 1744675710000, 3630 ], [ 1744675740000, 4042 ], [ 1744675770000, 3653 ], [ 1744675800000, 4358 ], [ 1744675830000, 3688 ], [ 1744675860000, 4450 ], [ 1744675890000, 3387 ], [ 1744675920000, 4864 ], [ 1744675950000, 3629 ], [ 1744675980000, 4127 ], [ 1744676010000, 3424 ], [ 1744676040000, 4267 ], [ 1744676070000, 3328 ], [ 1744676100000, 5128 ], [ 1744676130000, 3657 ], [ 1744676160000, 4185 ], [ 1744676190000, 3336 ], [ 1744676220000, 4532 ], [ 1744676250000, 3700 ], [ 1744676280000, 4174 ], [ 1744676310000, 3318 ], [ 1744676340000, 4463 ], [ 1744676370000, 3502 ], [ 1744676400000, 6064 ], [ 1744676430000, 3292 ], [ 1744676460000, 4858 ], [ 1744676490000, 3543 ], [ 1744676520000, 4620 ], [ 1744676550000, 3750 ], [ 1744676580000, 4043 ], [ 1744676610000, 3595 ], [ 1744676640000, 4152 ], [ 1744676670000, 3550 ], [ 1744676700000, 4011 ], [ 1744676730000, 3502 ], [ 1744676760000, 4050 ], [ 1744676790000, 3118 ], [ 1744676820000, 4628 ], [ 1744676850000, 3441 ], [ 1744676880000, 4366 ], [ 1744676910000, 3500 ], [ 1744676940000, 4160 ], [ 1744676970000, 3662 ], [ 1744677000000, 11392 ], [ 1744677030000, 18649 ], [ 1744677060000, 7107 ], [ 1744677090000, 9213 ], [ 1744677120000, 4235 ], [ 1744677150000, 3623 ], [ 1744677180000, 4412 ], [ 1744677210000, 3436 ], [ 1744677240000, 4233 ], [ 1744677270000, 3440 ], [ 1744677300000, 4383 ], [ 1744677330000, 3507 ], [ 1744677360000, 4288 ], [ 1744677390000, 3197 ], [ 1744677420000, 4605 ], [ 1744677450000, 3249 ], [ 1744677480000, 4421 ], [ 1744677510000, 2998 ], [ 1744677540000, 4700 ], [ 1744677570000, 3598 ], [ 1744677600000, 5781 ], [ 1744677630000, 3734 ], [ 1744677660000, 4510 ], [ 1744677690000, 3752 ], [ 1744677720000, 4447 ], [ 1744677750000, 3523 ], [ 1744677780000, 4187 ], [ 1744677810000, 3640 ], [ 1744677840000, 3900 ], [ 1744677870000, 3514 ], [ 1744677900000, 4863 ], [ 1744677930000, 3565 ], [ 1744677960000, 4335 ], [ 1744677990000, 3533 ], [ 1744678020000, 4307 ], [ 1744678050000, 3556 ], [ 1744678080000, 4179 ], [ 1744678110000, 3664 ], [ 1744678140000, 4362 ], [ 1744678170000, 3222 ], [ 1744678200000, 4750 ], [ 1744678230000, 3546 ], [ 1744678260000, 4601 ], [ 1744678290000, 3702 ], [ 1744678320000, 4564 ], [ 1744678350000, 3610 ], [ 1744678380000, 4130 ], [ 1744678410000, 3412 ], [ 1744678440000, 4614 ], [ 1744678470000, 3522 ], [ 1744678500000, 4148 ], [ 1744678530000, 3408 ], [ 1744678560000, 4261 ], [ 1744678590000, 3607 ], [ 1744678620000, 4172 ], [ 1744678650000, 3529 ], [ 1744678680000, 4227 ], [ 1744678710000, 3487 ], [ 1744678740000, 4298 ], [ 1744678770000, 3609 ], [ 1744678800000, 7230 ], [ 1744678830000, 3818 ], [ 1744678860000, 11924 ], [ 1744678890000, 27269 ], [ 1744678920000, 5073 ], [ 1744678950000, 3474 ], [ 1744678980000, 4474 ], [ 1744679010000, 3536 ], [ 1744679040000, 4525 ], [ 1744679070000, 3503 ], [ 1744679100000, 4194 ], [ 1744679130000, 3557 ], [ 1744679160000, 4259 ], [ 1744679190000, 3611 ], [ 1744679220000, 4218 ], [ 1744679250000, 3622 ], [ 1744679280000, 4417 ], [ 1744679310000, 3730 ], [ 1744679340000, 4204 ], [ 1744679370000, 3641 ], [ 1744679400000, 4849 ], [ 1744679430000, 3803 ], [ 1744679460000, 4398 ], [ 1744679490000, 3674 ], [ 1744679520000, 4727 ], [ 1744679550000, 3926 ], [ 1744679580000, 4173 ], [ 1744679610000, 3531 ], [ 1744679640000, 4968 ], [ 1744679670000, 3432 ], [ 1744679700000, 5059 ], [ 1744679730000, 3560 ], [ 1744679760000, 4087 ], [ 1744679790000, 3590 ], [ 1744679820000, 4436 ], [ 1744679850000, 5299 ], [ 1744679880000, 4320 ], [ 1744679910000, 3861 ], [ 1744679940000, 4511 ], [ 1744679970000, 3711 ], [ 1744680000000, 6021 ], [ 1744680030000, 3942 ], [ 1744680060000, 4800 ], [ 1744680090000, 3681 ], [ 1744680120000, 4592 ], [ 1744680150000, 3560 ], [ 1744680180000, 4194 ], [ 1744680210000, 3490 ], [ 1744680240000, 4971 ], [ 1744680270000, 4009 ], [ 1744680300000, 4837 ], [ 1744680330000, 3227 ], [ 1744680360000, 4531 ], [ 1744680390000, 2888 ], [ 1744680420000, 5083 ], [ 1744680450000, 3557 ], [ 1744680480000, 4207 ], [ 1744680510000, 3373 ], [ 1744680540000, 4482 ], [ 1744680570000, 3110 ], [ 1744680600000, 13551 ], [ 1744680630000, 17159 ], [ 1744680660000, 6284 ], [ 1744680690000, 9924 ], [ 1744680720000, 4547 ], [ 1744680750000, 3474 ], [ 1744680780000, 4312 ], [ 1744680810000, 3689 ], [ 1744680840000, 4680 ], [ 1744680870000, 3609 ], [ 1744680900000, 4886 ], [ 1744680930000, 3842 ], [ 1744680960000, 4810 ], [ 1744680990000, 4102 ], [ 1744681020000, 4594 ], [ 1744681050000, 4168 ], [ 1744681080000, 4562 ], [ 1744681110000, 4506 ], [ 1744681140000, 5243 ], [ 1744681170000, 5135 ], [ 1744681200000, 6671 ], [ 1744681230000, 3806 ], [ 1744681260000, 4535 ], [ 1744681290000, 3721 ], [ 1744681320000, 4799 ], [ 1744681350000, 3909 ], [ 1744681380000, 4261 ], [ 1744681410000, 3671 ], [ 1744681440000, 4359 ], [ 1744681470000, 4063 ], [ 1744681500000, 5231 ], [ 1744681530000, 3778 ], [ 1744681560000, 4684 ], [ 1744681590000, 4072 ], [ 1744681620000, 5029 ], [ 1744681650000, 3700 ], [ 1744681680000, 4670 ], [ 1744681710000, 3557 ], [ 1744681740000, 4590 ], [ 1744681770000, 3041 ], [ 1744681800000, 5043 ], [ 1744681830000, 3530 ], [ 1744681860000, 6807 ], [ 1744681890000, 4455 ], [ 1744681920000, 6841 ], [ 1744681950000, 4519 ], [ 1744681980000, 6617 ], [ 1744682010000, 4633 ], [ 1744682040000, 5997 ], [ 1744682070000, 4446 ], [ 1744682100000, 5569 ], [ 1744682130000, 4324 ], [ 1744682160000, 5354 ], [ 1744682190000, 7245 ], [ 1744682220000, 5258 ], [ 1744682250000, 4296 ], [ 1744682280000, 5349 ], [ 1744682310000, 4479 ], [ 1744682340000, 5127 ], [ 1744682370000, 4006 ], [ 1744682400000, 19058 ], [ 1744682430000, 14501 ], [ 1744682460000, 3810 ], [ 1744682490000, 12368 ], [ 1744682520000, 6976 ], [ 1744682550000, 4399 ], [ 1744682580000, 5482 ], [ 1744682610000, 4524 ], [ 1744682640000, 5478 ], [ 1744682670000, 4920 ], [ 1744682700000, 5347 ], [ 1744682730000, 4427 ], [ 1744682760000, 5102 ], [ 1744682790000, 4441 ], [ 1744682820000, 5596 ], [ 1744682850000, 4888 ], [ 1744682880000, 5306 ], [ 1744682910000, 4825 ], [ 1744682940000, 5897 ], [ 1744682970000, 4481 ], [ 1744683000000, 6086 ], [ 1744683030000, 4910 ], [ 1744683060000, 5676 ], [ 1744683090000, 3626 ], [ 1744683120000, 6929 ], [ 1744683150000, 4601 ], [ 1744683180000, 5525 ], [ 1744683210000, 4500 ], [ 1744683240000, 5617 ], [ 1744683270000, 4503 ], [ 1744683300000, 6328 ], [ 1744683330000, 4557 ], [ 1744683360000, 5356 ], [ 1744683390000, 4413 ], [ 1744683420000, 5335 ], [ 1744683450000, 4640 ], [ 1744683480000, 5399 ], [ 1744683510000, 4298 ], [ 1744683540000, 5415 ], [ 1744683570000, 4540 ], [ 1744683600000, 6949 ], [ 1744683630000, 4574 ], [ 1744683660000, 5757 ], [ 1744683690000, 4669 ], [ 1744683720000, 5706 ], [ 1744683750000, 4472 ], [ 1744683780000, 5386 ], [ 1744683810000, 4490 ], [ 1744683840000, 5104 ], [ 1744683870000, 4201 ], [ 1744683900000, 5979 ], [ 1744683930000, 4853 ], [ 1744683960000, 6691 ], [ 1744683990000, 4572 ], [ 1744684020000, 5554 ], [ 1744684050000, 5244 ], [ 1744684080000, 5392 ], [ 1744684110000, 4550 ] ]
  } ]
}`,
		},
	} {
		t.Run(fmt.Sprintf("%s", i), func(t *testing.T) {
			metadata.SetUser(ctx, "username:test", spaceUid, "true")

			res, err := queryTsWithPromEngine(ctx, c.queryTs)
			assert.Nil(t, err)
			excepted, err := json.Marshal(res)
			assert.Nil(t, err)
			assert.JSONEq(t, c.result, string(excepted))
		})
	}
}

func TestQueryTsWithEs(t *testing.T) {
	ctx := metadata.InitHashID(context.Background())

	viper.Set(bkapi.BkAPIAddressConfigPath, mock.EsUrlDomain)

	spaceUid := influxdb.SpaceUid
	tableID := influxdb.ResultTableEs

	mock.Init()
	promql.MockEngine()

	defaultStart := time.UnixMilli(1717027200000)
	defaultEnd := time.UnixMilli(1717027500000)

	for i, c := range map[string]struct {
		queryTs *structured.QueryTs
		result  string
	}{
		"查询 10 条原始数据，按照字段正向排序": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "gseIndex",
						Limit:         10,
						From:          0,
						ReferenceName: "a",
					},
				},
				OrderBy: structured.OrderBy{
					"_value",
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     false,
				SpaceUid:    spaceUid,
			},
		},
		"根据维度 __ext.container_name 进行 count 聚合，同时用值正向排序": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "gseIndex",
						Limit:         5,
						From:          0,
						ReferenceName: "a",
						TimeAggregation: structured.TimeAggregation{
							Function: "count_over_time",
							Window:   "30s",
						},
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method:     "sum",
								Dimensions: []string{"__ext.container_name"},
							},
							{
								Method: "topk",
								VArgsList: []interface{}{
									5,
								},
							},
						},
					},
				},
				OrderBy: structured.OrderBy{
					"gseIndex",
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     false,
				SpaceUid:    spaceUid,
				Step:        "30s",
			},
		},
	} {
		t.Run(fmt.Sprintf("%s", i), func(t *testing.T) {
			metadata.SetUser(ctx, "username:test", spaceUid, "true")

			res, err := queryTsWithPromEngine(ctx, c.queryTs)
			if err != nil {
				log.Errorf(ctx, err.Error())
				return
			}
			data := res.(*PromData)
			if data.Status != nil && data.Status.Code != "" {
				fmt.Println("code: ", data.Status.Code)
				fmt.Println("message: ", data.Status.Message)
				return
			}

			log.Infof(ctx, fmt.Sprintf("%+v", data.Tables))
		})
	}
}

func TestQueryReferenceWithEs(t *testing.T) {
	ctx := metadata.InitHashID(context.Background())

	spaceUid := influxdb.SpaceUid
	tableID := influxdb.ResultTableEs

	mock.Init()
	promql.MockEngine()
	influxdb.MockSpaceRouter(ctx)

	defaultStart := time.UnixMilli(1741154079123) // 2025-03-05 13:54:39
	defaultEnd := time.UnixMilli(1741155879987)   // 2025-03-05 14:24:39

	mock.Es.Set(map[string]any{
		`{"aggregations":{"_value":{"value_count":{"field":"gseIndex"}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_millis","from":1741154079123,"include_lower":true,"include_upper":true,"to":1741155879987}}}}},"size":0}`: `{"took":626,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"_value":{"value":182355}}}`,
		`{"aggregations":{"_value":{"value_count":{"field":"gseIndex"}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1741154079,"include_lower":true,"include_upper":true,"to":1741155879}}}}},"size":0}`:       `{"took":171,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"_value":{"value":182486}}}`,

		`{"aggregations":{"__ext.container_name":{"aggregations":{"_value":{"value_count":{"field":"gseIndex"}}},"terms":{"field":"__ext.container_name","missing":" ","order":[{"_value":"asc"}],"size":5}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_millis","from":1741154079123,"include_lower":true,"include_upper":true,"to":1741155879987}}}}},"size":0}`: `{"took":860,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"__ext.container_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"unify-query","doc_count":182355,"_value":{"value":182355}},{"key":" ","doc_count":182355,"_value":{"value":4325521}}]}}}`,

		`{"aggregations":{"__ext.container_name":{"aggregations":{"_value":{"value_count":{"field":"gseIndex"}}},"terms":{"field":"__ext.container_name","missing":" ","order":[{"_value":"desc"}],"size":5}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1741154079,"include_lower":true,"include_upper":true,"to":1741155879}}}}},"size":0}`:                                                                                       `{"took":885,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"__ext.container_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"unify-query","doc_count":182486,"_value":{"value":182486}}]}}}`,
		`{"aggregations":{"_value":{"value_count":{"field":"__ext.container_name"}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1741154079,"include_lower":true,"include_upper":true,"to":1741155879}}}}},"size":0}`:                                                                                                                                                                                                                 `{"took":283,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"_value":{"value":182486}}}`,
		`{"aggregations":{"_value":{"value_count":{"field":"__ext.io_kubernetes_pod"}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1741154079,"include_lower":true,"include_upper":true,"to":1741155879}}}}},"size":0}`:                                                                                                                                                                                                              `{"took":167,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"_value":{"value":182486}}}`,
		`{"aggregations":{"_value":{"cardinality":{"field":"__ext.io_kubernetes_pod"}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1741154079,"include_lower":true,"include_upper":true,"to":1741155879}}}}},"size":0}`:                                                                                                                                                                                                              `{"took":1595,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"_value":{"value":4}}}`,
		`{"aggregations":{"dtEventTimeStamp":{"aggregations":{"_value":{"value_count":{"field":"__ext.io_kubernetes_pod"}}},"date_histogram":{"extended_bounds":{"max":1741155879000,"min":1741154079000},"field":"dtEventTimeStamp","interval":"1m","min_doc_count":0}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1741154079,"include_lower":true,"include_upper":true,"to":1741155879}}}}},"size":0}`:                            `{"took":529,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"dtEventTimeStamp":{"buckets":[{"key_as_string":"1741154040000","key":1741154040000,"doc_count":3408,"_value":{"value":3408}},{"key_as_string":"1741154100000","key":1741154100000,"doc_count":4444,"_value":{"value":4444}},{"key_as_string":"1741154160000","key":1741154160000,"doc_count":4577,"_value":{"value":4577}},{"key_as_string":"1741154220000","key":1741154220000,"doc_count":4668,"_value":{"value":4668}},{"key_as_string":"1741154280000","key":1741154280000,"doc_count":5642,"_value":{"value":5642}},{"key_as_string":"1741154340000","key":1741154340000,"doc_count":4860,"_value":{"value":4860}},{"key_as_string":"1741154400000","key":1741154400000,"doc_count":35988,"_value":{"value":35988}},{"key_as_string":"1741154460000","key":1741154460000,"doc_count":7098,"_value":{"value":7098}},{"key_as_string":"1741154520000","key":1741154520000,"doc_count":5287,"_value":{"value":5287}},{"key_as_string":"1741154580000","key":1741154580000,"doc_count":5422,"_value":{"value":5422}},{"key_as_string":"1741154640000","key":1741154640000,"doc_count":4906,"_value":{"value":4906}},{"key_as_string":"1741154700000","key":1741154700000,"doc_count":4447,"_value":{"value":4447}},{"key_as_string":"1741154760000","key":1741154760000,"doc_count":4713,"_value":{"value":4713}},{"key_as_string":"1741154820000","key":1741154820000,"doc_count":4621,"_value":{"value":4621}},{"key_as_string":"1741154880000","key":1741154880000,"doc_count":4417,"_value":{"value":4417}},{"key_as_string":"1741154940000","key":1741154940000,"doc_count":5092,"_value":{"value":5092}},{"key_as_string":"1741155000000","key":1741155000000,"doc_count":4805,"_value":{"value":4805}},{"key_as_string":"1741155060000","key":1741155060000,"doc_count":5545,"_value":{"value":5545}},{"key_as_string":"1741155120000","key":1741155120000,"doc_count":4614,"_value":{"value":4614}},{"key_as_string":"1741155180000","key":1741155180000,"doc_count":5121,"_value":{"value":5121}},{"key_as_string":"1741155240000","key":1741155240000,"doc_count":4854,"_value":{"value":4854}},{"key_as_string":"1741155300000","key":1741155300000,"doc_count":5343,"_value":{"value":5343}},{"key_as_string":"1741155360000","key":1741155360000,"doc_count":4789,"_value":{"value":4789}},{"key_as_string":"1741155420000","key":1741155420000,"doc_count":4755,"_value":{"value":4755}},{"key_as_string":"1741155480000","key":1741155480000,"doc_count":5115,"_value":{"value":5115}},{"key_as_string":"1741155540000","key":1741155540000,"doc_count":4588,"_value":{"value":4588}},{"key_as_string":"1741155600000","key":1741155600000,"doc_count":6474,"_value":{"value":6474}},{"key_as_string":"1741155660000","key":1741155660000,"doc_count":5416,"_value":{"value":5416}},{"key_as_string":"1741155720000","key":1741155720000,"doc_count":5128,"_value":{"value":5128}},{"key_as_string":"1741155780000","key":1741155780000,"doc_count":5050,"_value":{"value":5050}},{"key_as_string":"1741155840000","key":1741155840000,"doc_count":1299,"_value":{"value":1299}}]}}}`,
		`{"aggregations":{"dtEventTimeStamp":{"aggregations":{"_value":{"value_count":{"field":"__ext.io_kubernetes_pod"}}},"date_histogram":{"extended_bounds":{"max":1741155879987,"min":1741154079123},"field":"dtEventTimeStamp","interval":"1m","min_doc_count":0}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_millis","from":1741154079123,"include_lower":true,"include_upper":true,"to":1741155879987}}}}},"size":0}`:                      `{"took":759,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"dtEventTimeStamp":{"buckets":[{"key_as_string":"1741154040000","key":1741154040000,"doc_count":3277,"_value":{"value":3277}},{"key_as_string":"1741154100000","key":1741154100000,"doc_count":4444,"_value":{"value":4444}},{"key_as_string":"1741154160000","key":1741154160000,"doc_count":4577,"_value":{"value":4577}},{"key_as_string":"1741154220000","key":1741154220000,"doc_count":4668,"_value":{"value":4668}},{"key_as_string":"1741154280000","key":1741154280000,"doc_count":5642,"_value":{"value":5642}},{"key_as_string":"1741154340000","key":1741154340000,"doc_count":4860,"_value":{"value":4860}},{"key_as_string":"1741154400000","key":1741154400000,"doc_count":35988,"_value":{"value":35988}},{"key_as_string":"1741154460000","key":1741154460000,"doc_count":7098,"_value":{"value":7098}},{"key_as_string":"1741154520000","key":1741154520000,"doc_count":5287,"_value":{"value":5287}},{"key_as_string":"1741154580000","key":1741154580000,"doc_count":5422,"_value":{"value":5422}},{"key_as_string":"1741154640000","key":1741154640000,"doc_count":4906,"_value":{"value":4906}},{"key_as_string":"1741154700000","key":1741154700000,"doc_count":4447,"_value":{"value":4447}},{"key_as_string":"1741154760000","key":1741154760000,"doc_count":4713,"_value":{"value":4713}},{"key_as_string":"1741154820000","key":1741154820000,"doc_count":4621,"_value":{"value":4621}},{"key_as_string":"1741154880000","key":1741154880000,"doc_count":4417,"_value":{"value":4417}},{"key_as_string":"1741154940000","key":1741154940000,"doc_count":5092,"_value":{"value":5092}},{"key_as_string":"1741155000000","key":1741155000000,"doc_count":4805,"_value":{"value":4805}},{"key_as_string":"1741155060000","key":1741155060000,"doc_count":5545,"_value":{"value":5545}},{"key_as_string":"1741155120000","key":1741155120000,"doc_count":4614,"_value":{"value":4614}},{"key_as_string":"1741155180000","key":1741155180000,"doc_count":5121,"_value":{"value":5121}},{"key_as_string":"1741155240000","key":1741155240000,"doc_count":4854,"_value":{"value":4854}},{"key_as_string":"1741155300000","key":1741155300000,"doc_count":5343,"_value":{"value":5343}},{"key_as_string":"1741155360000","key":1741155360000,"doc_count":4789,"_value":{"value":4789}},{"key_as_string":"1741155420000","key":1741155420000,"doc_count":4755,"_value":{"value":4755}},{"key_as_string":"1741155480000","key":1741155480000,"doc_count":5115,"_value":{"value":5115}},{"key_as_string":"1741155540000","key":1741155540000,"doc_count":4588,"_value":{"value":4588}},{"key_as_string":"1741155600000","key":1741155600000,"doc_count":6474,"_value":{"value":6474}},{"key_as_string":"1741155660000","key":1741155660000,"doc_count":5416,"_value":{"value":5416}},{"key_as_string":"1741155720000","key":1741155720000,"doc_count":5128,"_value":{"value":5128}},{"key_as_string":"1741155780000","key":1741155780000,"doc_count":5050,"_value":{"value":5050}},{"key_as_string":"1741155840000","key":1741155840000,"doc_count":1299,"_value":{"value":1299}}]}}}`,
		`{"aggregations":{"dtEventTimeStamp":{"aggregations":{"_value":{"value_count":{"field":"dtEventTimeStamp"}}},"date_histogram":{"extended_bounds":{"max":1741341600000,"min":1741320000000},"field":"dtEventTimeStamp","interval":"1d","min_doc_count":0,"time_zone":"Asia/Shanghai"}}},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_millis","from":1741320000000,"include_lower":true,"include_upper":true,"to":1741341600000}}}}},"size":0}`: `{"took":5,"timed_out":false,"_shards":{"total":68,"successful":68,"skipped":0,"failed":0},"hits":{"total":{"value":2367,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"dtEventTimeStamp":{"buckets":[{"key_as_string":"1741276800000","key":1741276800000,"doc_count":2367,"_value":{"value":2367}}]}}}`,
	})

	for i, c := range map[string]struct {
		queryTs *structured.QueryTs
		result  string
	}{
		"统计数量，毫秒查询": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "gseIndex",
						ReferenceName: "a",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method: "count",
							},
						},
					},
				},
				OrderBy: structured.OrderBy{
					"_value",
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.UnixMilli(), 10),
				End:         strconv.FormatInt(defaultEnd.UnixMilli(), 10),
				Instant:     true,
				SpaceUid:    spaceUid,
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":[],"group_values":[],"values":[[1741154079123,182355]]}]`,
		},
		"统计数量": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "gseIndex",
						ReferenceName: "a",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method: "count",
							},
						},
					},
				},
				OrderBy: structured.OrderBy{
					"_value",
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     true,
				SpaceUid:    spaceUid,
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":[],"group_values":[],"values":[[1741154079000,182486]]}]`,
		},
		"根据维度 __ext.container_name 进行 sum 聚合，同时用值正向排序": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "gseIndex",
						Limit:         5,
						From:          0,
						ReferenceName: "a",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method:     "count",
								Dimensions: []string{"__ext.container_name"},
							},
						},
					},
				},
				OrderBy: structured.OrderBy{
					"_value",
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.UnixMilli(), 10),
				End:         strconv.FormatInt(defaultEnd.UnixMilli(), 10),
				Instant:     true,
				SpaceUid:    spaceUid,
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":["__ext.container_name"],"group_values":["unify-query"],"values":[[1741154079123,182355]]},{"name":"_result1","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":["__ext.container_name"],"group_values":[""],"values":[[1741154079123,4325521]]}]`,
		},
		"根据维度 __ext.container_name 进行 count 聚合，同时用值倒序": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "gseIndex",
						Limit:         5,
						From:          0,
						ReferenceName: "a",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method:     "count",
								Dimensions: []string{"__ext.container_name"},
							},
						},
					},
				},
				OrderBy: structured.OrderBy{
					"-_value",
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     true,
				SpaceUid:    spaceUid,
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":["__ext.container_name"],"group_values":["unify-query"],"values":[[1741154079000,182486]]}]`,
		},
		"统计 __ext.container_name 和 __ext.io_kubernetes_pod 不为空的文档数量": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "__ext.container_name",
						ReferenceName: "a",
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "__ext.io_kubernetes_pod",
									Operator:      "ncontains",
									Value:         []string{""},
								},
								{
									DimensionName: "__ext.container_name",
									Operator:      "ncontains",
									Value:         []string{""},
								},
							},
							ConditionList: []string{
								"and",
							},
						},
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method: "count",
							},
						},
					},
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     true,
				SpaceUid:    spaceUid,
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":[],"group_values":[],"values":[[1741154079000,182486]]}]`,
		},
		"a + b": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "__ext.io_kubernetes_pod",
						ReferenceName: "a",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method: "count",
							},
						},
					},
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "__ext.io_kubernetes_pod",
						ReferenceName: "b",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method: "count",
							},
						},
					},
				},
				MetricMerge: "a + b",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     true,
				SpaceUid:    spaceUid,
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":[],"group_values":[],"values":[[1741154079000,364972]]}]`,
		},
		"__ext.io_kubernetes_pod 统计去重数量": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "__ext.io_kubernetes_pod",
						ReferenceName: "a",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method: "cardinality",
							},
						},
					},
				},
				MetricMerge: "a",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     true,
				SpaceUid:    spaceUid,
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":[],"group_values":[],"values":[[1741154079000,4]]}]`,
		},
		"__ext.io_kubernetes_pod 统计数量": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "__ext.io_kubernetes_pod",
						ReferenceName: "b",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method: "count",
							},
							{
								Method: "date_histogram",
								Window: "1m",
							},
						},
					},
				},
				MetricMerge: "b",
				Start:       strconv.FormatInt(defaultStart.Unix(), 10),
				End:         strconv.FormatInt(defaultEnd.Unix(), 10),
				Instant:     false,
				SpaceUid:    spaceUid,
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":[],"group_values":[],"values":[[1741154040000,3408],[1741154100000,4444],[1741154160000,4577],[1741154220000,4668],[1741154280000,5642],[1741154340000,4860],[1741154400000,35988],[1741154460000,7098],[1741154520000,5287],[1741154580000,5422],[1741154640000,4906],[1741154700000,4447],[1741154760000,4713],[1741154820000,4621],[1741154880000,4417],[1741154940000,5092],[1741155000000,4805],[1741155060000,5545],[1741155120000,4614],[1741155180000,5121],[1741155240000,4854],[1741155300000,5343],[1741155360000,4789],[1741155420000,4755],[1741155480000,5115],[1741155540000,4588],[1741155600000,6474],[1741155660000,5416],[1741155720000,5128],[1741155780000,5050],[1741155840000,1299]]}]`,
		},
		"__ext.io_kubernetes_pod 统计数量，毫秒": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "__ext.io_kubernetes_pod",
						ReferenceName: "b",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method: "count",
							},
							{
								Method: "date_histogram",
								Window: "1m",
							},
						},
					},
				},
				MetricMerge: "b",
				Start:       strconv.FormatInt(defaultStart.UnixMilli(), 10),
				End:         strconv.FormatInt(defaultEnd.UnixMilli(), 10),
				Instant:     false,
				SpaceUid:    spaceUid,
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":[],"group_values":[],"values":[[1741154040000,3277],[1741154100000,4444],[1741154160000,4577],[1741154220000,4668],[1741154280000,5642],[1741154340000,4860],[1741154400000,35988],[1741154460000,7098],[1741154520000,5287],[1741154580000,5422],[1741154640000,4906],[1741154700000,4447],[1741154760000,4713],[1741154820000,4621],[1741154880000,4417],[1741154940000,5092],[1741155000000,4805],[1741155060000,5545],[1741155120000,4614],[1741155180000,5121],[1741155240000,4854],[1741155300000,5343],[1741155360000,4789],[1741155420000,4755],[1741155480000,5115],[1741155540000,4588],[1741155600000,6474],[1741155660000,5416],[1741155720000,5128],[1741155780000,5050],[1741155840000,1299]]}]`,
		},
		"测试聚合周期大于查询周期": {
			queryTs: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(tableID),
						FieldName:     "dtEventTimeStamp",
						ReferenceName: "b",
						AggregateMethodList: structured.AggregateMethodList{
							{
								Method: "count",
							},
							{
								Method: "date_histogram",
								Window: "1d",
							},
						},
					},
				},
				MetricMerge: "b",
				Start:       "1741320000000",
				End:         "1741341600000",
				Instant:     false,
				SpaceUid:    spaceUid,
				Timezone:    "Asia/Shanghai",
				Step:        "1d",
			},
			result: `[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":[],"group_values":[],"values":[[1741276800000,2367]]}]`,
		},
	} {
		t.Run(fmt.Sprintf("%s", i), func(t *testing.T) {
			metadata.SetUser(ctx, "username:test", spaceUid, "true")

			data, err := queryReferenceWithPromEngine(ctx, c.queryTs)
			assert.Nil(t, err)

			if err != nil {
				return
			}

			if data.Status != nil && data.Status.Code != "" {
				fmt.Println("code: ", data.Status.Code)
				fmt.Println("message: ", data.Status.Message)
				return
			}

			actual, _ := json.Marshal(data.Tables)
			assert.Equal(t, c.result, string(actual))
		})
	}
}

func TestQueryTs(t *testing.T) {

	ctx := metadata.InitHashID(context.Background())
	mock.Init()
	influxdb.MockSpaceRouter(ctx)
	promql.MockEngine()

	mock.InfluxDB.Set(map[string]any{
		`SELECT mean("usage") AS _value, "time" AS _time FROM cpu_summary WHERE time > 1677081599999000000 and time < 1677085659999000000 AND (bk_biz_id='2') GROUP BY time(1m0s) LIMIT 100000005 SLIMIT 100005 TZ('UTC')`: &decoder.Response{
			Results: []decoder.Result{
				{
					Series: []*decoder.Row{
						{
							Name: "",
							Tags: map[string]string{},
							Columns: []string{
								influxdb.TimeColumnName,
								influxdb.ResultColumnName,
							},
							Values: [][]any{
								{
									1677081600000000000, 30,
								},
								{
									1677081660000000000, 21,
								},
								{
									1677081720000000000, 1,
								},
								{
									1677081780000000000, 7,
								},
								{
									1677081840000000000, 4,
								},
								{
									1677081900000000000, 2,
								},
								{
									1677081960000000000, 100,
								},
								{
									1677082020000000000, 94,
								},
								{
									1677082080000000000, 34,
								},
							},
						},
					},
				},
			},
		},
		`SELECT "usage" AS _value, *::tag, "time" AS _time FROM cpu_summary WHERE time > 1677081359999000000 and time < 1677085659999000000 AND ((notice_way='weixin' and status='failed') and bk_biz_id='2') LIMIT 100000005 SLIMIT 100005 TZ('UTC')`: &decoder.Response{
			Results: []decoder.Result{
				{
					Series: []*decoder.Row{
						{
							Name: "",
							Tags: map[string]string{},
							Columns: []string{
								influxdb.ResultColumnName,
								"job",
								"notice_way",
								"status",
								influxdb.TimeColumnName,
							},
							Values: [][]any{
								{
									30,
									"SLI",
									"weixin",
									"failed",
									1677081600000000000,
								},
								{
									21,
									"SLI",
									"weixin",
									"failed",
									1677081660000000000,
								},
								{
									1,
									"SLI",
									"weixin",
									"failed",
									1677081720000000000,
								},
								{
									7,
									"SLI",
									"weixin",
									"failed",
									1677081780000000000,
								},
								{
									4,
									"SLI",
									"weixin",
									"failed",
									1677081840000000000,
								},
								{
									2,
									"SLI",
									"weixin",
									"failed",
									1677081900000000000,
								},
								{
									100,
									"SLI",
									"weixin",
									"failed",
									1677081960000000000,
								},
								{
									94,
									"SLI",
									"weixin",
									"failed",
									1677082020000000000,
								},
								{
									34,
									"SLI",
									"weixin",
									"failed",
									1677082080000000000,
								},
							},
						},
					},
				},
			},
		},
		`SELECT count("usage") AS _value, "time" AS _time FROM cpu_summary WHERE time > 1677081599999000000 and time < 1677085659999000000 AND (bk_biz_id='2') GROUP BY "status", time(1m0s) LIMIT 100000005 SLIMIT 100005 TZ('UTC')`: &decoder.Response{
			Results: []decoder.Result{
				{
					Series: []*decoder.Row{
						{
							Name: "",
							Tags: map[string]string{
								"status": "failed",
							},
							Columns: []string{
								influxdb.TimeColumnName,
								influxdb.ResultColumnName,
							},
							Values: [][]any{
								{
									1677081600000000000, 30,
								},
								{
									1677081660000000000, 21,
								},
								{
									1677081720000000000, 1,
								},
								{
									1677081780000000000, 7,
								},
								{
									1677081840000000000, 4,
								},
								{
									1677081900000000000, 2,
								},
								{
									1677081960000000000, 100,
								},
								{
									1677082020000000000, 94,
								},
								{
									1677082080000000000, 34,
								},
							},
						},
					},
				},
			},
		},
	})

	testCases := map[string]struct {
		query  string
		result string
	}{
		"test query": {
			query:  `{"query_list":[{"data_source":"","table_id":"system.cpu_summary","field_name":"usage","field_list":null,"function":[{"method":"mean","without":false,"dimensions":[],"position":0,"args_list":null,"vargs_list":null}],"time_aggregation":{"function":"avg_over_time","window":"60s","position":0,"vargs_list":null},"reference_name":"a","dimensions":[],"limit":0,"timestamp":null,"start_or_end":0,"vector_offset":0,"offset":"","offset_forward":false,"slimit":0,"soffset":0,"conditions":{"field_list":[],"condition_list":[]},"keep_columns":["_time","a"]}],"metric_merge":"a","result_columns":null,"start_time":"1677081600","end_time":"1677085600","step":"60s"}`,
			result: `{"series":[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":[],"group_values":[],"values":[[1677081600000,30],[1677081660000,21],[1677081720000,1],[1677081780000,7],[1677081840000,4],[1677081900000,2],[1677081960000,100],[1677082020000,94],[1677082080000,34]]}]}`,
		},
		"test lost sample in increase": {
			query:  `{"query_list":[{"data_source":"bkmonitor","table_id":"system.cpu_summary","field_name":"usage","field_list":null,"function":null,"time_aggregation":{"function":"increase","window":"5m0s","position":0,"vargs_list":null},"reference_name":"a","dimensions":null,"limit":0,"timestamp":null,"start_or_end":0,"vector_offset":0,"offset":"","offset_forward":false,"slimit":0,"soffset":0,"conditions":{"field_list":[{"field_name":"notice_way","value":["weixin"],"op":"eq"},{"field_name":"status","value":["failed"],"op":"eq"}],"condition_list":["and"]},"keep_columns":null}],"metric_merge":"a","result_columns":null,"start_time":"1677081600","end_time":"1677085600","step":"60s"}`,
			result: `{"series":[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":["job","notice_way","status"],"group_values":["SLI","weixin","failed"],"values":[[1677081660000,52.499649999999995],[1677081720000,38.49981666666667],[1677081780000,46.66666666666667],[1677081840000,40],[1677081900000,16.25],[1677081960000,137.5],[1677082020000,247.5],[1677082080000,285],[1677082140000,263.6679222222223],[1677082200000,160.00106666666667],[1677082260000,51.00056666666667]]}]}`,
		},
		"test query support fuzzy __name__ with count": {
			query:  `{"query_list":[{"data_source":"","table_id":"system.cpu_summary","field_name":".*","is_regexp":true,"field_list":null,"function":[{"method":"sum","without":false,"dimensions":["status"],"position":0,"args_list":null,"vargs_list":null}],"time_aggregation":{"function":"count_over_time","window":"60s","position":0,"vargs_list":null},"reference_name":"a","dimensions":[],"limit":0,"timestamp":null,"start_or_end":0,"vector_offset":0,"offset":"","offset_forward":false,"slimit":0,"soffset":0,"conditions":{"field_list":[],"condition_list":[]},"keep_columns":["_time","a"]}],"metric_merge":"a","result_columns":null,"start_time":"1677081600","end_time":"1677085600","step":"60s"}`,
			result: `{"series":[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":["status"],"group_values":["failed"],"values":[[1677081600000,30],[1677081660000,21],[1677081720000,1],[1677081780000,7],[1677081840000,4],[1677081900000,2],[1677081960000,100],[1677082020000,94],[1677082080000,34]]}]}`,
		},
	}

	for name, c := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx = metadata.InitHashID(ctx)
			metadata.SetUser(ctx, "", influxdb.SpaceUid, "")

			body := []byte(c.query)
			query := &structured.QueryTs{}
			err := json.Unmarshal(body, query)
			assert.Nil(t, err)

			res, err := queryTsWithPromEngine(ctx, query)
			assert.Nil(t, err)
			out, err := json.Marshal(res)
			assert.Nil(t, err)
			actual := string(out)
			fmt.Printf("ActualResult: %v\n", actual)
			assert.Equal(t, c.result, actual)
		})
	}
}

func TestQueryRawWithInstance(t *testing.T) {
	ctx := metadata.InitHashID(context.Background())

	spaceUid := influxdb.SpaceUid

	mock.Init()
	influxdb.MockSpaceRouter(ctx)
	promql.MockEngine()

	start := "1723594000"
	end := "1723595000"

	mock.BkSQL.Set(map[string]any{
		"SHOW CREATE TABLE `2_bklog_bkunify_query_doris`.doris": `{"result":true,"message":"成功","code":"00","data":{"result_table_scan_range":{},"cluster":"doris-test","totalRecords":18,"external_api_call_time_mills":{"bkbase_auth_api":43,"bkbase_meta_api":0,"bkbase_apigw_api":33},"resource_use_summary":{"cpu_time_mills":0,"memory_bytes":0,"processed_bytes":0,"processed_rows":0},"source":"","list":[{"Field":"thedate","Type":"int","Null":"NO","Key":"YES","Default":null,"Extra":""},{"Field":"dteventtimestamp","Type":"bigint","Null":"NO","Key":"YES","Default":null,"Extra":""},{"Field":"dteventtime","Type":"varchar(32)","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"localtime","Type":"varchar(32)","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"__shard_key__","Type":"bigint","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"__ext","Type":"variant","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"cloudid","Type":"double","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"file","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"gseindex","Type":"double","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"iterationindex","Type":"double","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"level","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"log","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"message","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"path","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"report_time","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"serverip","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"time","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"},{"Field":"trace_id","Type":"text","Null":"YES","Key":"NO","Default":null,"Extra":"NONE"}],"stage_elapsed_time_mills":{"check_query_syntax":0,"query_db":5,"get_query_driver":0,"match_query_forbidden_config":0,"convert_query_statement":2,"connect_db":45,"match_query_routing_rule":0,"check_permission":43,"check_query_semantic":0,"pick_valid_storage":1},"select_fields_order":["Field","Type","Null","Key","Default","Extra"],"sql":"SHOW COLUMNS FROM mapleleaf_2.bklog_bkunify_query_doris_2","total_record_size":11776,"timetaken":0.096,"result_schema":[{"field_type":"string","field_name":"Field","field_alias":"Field","field_index":0},{"field_type":"string","field_name":"Type","field_alias":"Type","field_index":1},{"field_type":"string","field_name":"Null","field_alias":"Null","field_index":2},{"field_type":"string","field_name":"Key","field_alias":"Key","field_index":3},{"field_type":"string","field_name":"Default","field_alias":"Default","field_index":4},{"field_type":"string","field_name":"Extra","field_alias":"Extra","field_index":5}],"bksql_call_elapsed_time":0,"device":"doris","result_table_ids":["2_bklog_bkunify_query_doris"]},"errors":null,"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000"}`,
	})

	mock.Es.Set(map[string]any{
		`{"_source":{"includes":["__ext.container_id","dtEventTimeStamp"]},"from":0,"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1723594000,"include_lower":true,"include_upper":true,"to":1723595000}}}}},"size":20,"sort":[{"dtEventTimeStamp":{"order":"desc"}}]}`:      `{"took":301,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":0.0,"hits":[{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"c726c895a380ba1a9df04ba4a977b29b","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"fa209967d4a8c5d21b3e4f67d2cd579e","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"dc888e9a3789976aa11483626fc61a4f","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"c2dae031f095fa4b9deccf81964c7837","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"8a916e558c71d4226f1d7f3279cf0fdd","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"f6950fef394e813999d7316cdbf0de4d","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"328d487e284703b1d0bb8017dba46124","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"cb790ecb36bbaf02f6f0eb80ac2fd65c","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"bd8a8ef60e94ade63c55c8773170d458","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"c8401bb4ec021b038cb374593b8adce3","_score":0.0,"_source":{"dtEventTimeStamp":"1723594161000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}}]}}`,
		`{"_source":{"includes":["__ext.io_kubernetes_pod","dtEventTimeStamp"]},"from":0,"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1723594000,"include_lower":true,"include_upper":true,"to":1723595000}}}}},"size":20,"sort":[{"dtEventTimeStamp":{"order":"desc"}}]}`: `{"took":468,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":0.0,"hits":[{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"e058129ae18bff87c95e83f24584e654","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"c124dae69af9b86a7128ee4281820158","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"c7f73abf7e865a4b4d7fc608387d01cf","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"39c3ec662881e44bf26d2a6bfc0e35c3","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"58e03ce0b9754bf0657d49a5513adcb5","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"43a36f412886bf30b0746562513638d3","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"218ceafd04f89b39cda7954e51f4a48a","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"8d9abe9b782fe3a1272c93f0af6b39e1","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"0826407be7f04f19086774ed68eac8dd","_score":0.0,"_source":{"dtEventTimeStamp":"1723594224000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-llp94"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"d56b4120194eb37f53410780da777d43","_score":0.0,"_source":{"dtEventTimeStamp":"1723594224000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-llp94"}}}]}}`,
		`{"_source":{"includes":["__ext.container_id","dtEventTimeStamp"]},"from":1,"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1723594000,"include_lower":true,"include_upper":true,"to":1723595000}}}}},"size":1}`:                                                      `{"took":17,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":0.0,"hits":[{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"4f3a5e9c167097c9658e88b2f32364b2","_score":0.0,"_source":{"dtEventTimeStamp":"1723594209000","__ext":{"container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f"}}}]}}`,
		`{"_source":{"includes":["__ext.container_id","dtEventTimeStamp"]},"from":1,"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_millis","from":1723594000123,"include_lower":true,"include_upper":true,"to":1723595000234}}}}},"size":10}`:                                               `{"took":468,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":0.0,"hits":[{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"e058129ae18bff87c95e83f24584e654","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"c124dae69af9b86a7128ee4281820158","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"c7f73abf7e865a4b4d7fc608387d01cf","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"39c3ec662881e44bf26d2a6bfc0e35c3","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"58e03ce0b9754bf0657d49a5513adcb5","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"43a36f412886bf30b0746562513638d3","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"218ceafd04f89b39cda7954e51f4a48a","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"8d9abe9b782fe3a1272c93f0af6b39e1","_score":0.0,"_source":{"dtEventTimeStamp":"1723594211000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"0826407be7f04f19086774ed68eac8dd","_score":0.0,"_source":{"dtEventTimeStamp":"1723594224000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-llp94"}}},{"_index":"v2_2_bklog_bk_unify_query_20240814_0","_type":"_doc","_id":"d56b4120194eb37f53410780da777d43","_score":0.0,"_source":{"dtEventTimeStamp":"1723594224000","__ext":{"io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-llp94"}}}]}}`,

		// merge rt test mock data
		`{"_source":{"includes":["a","b"]},"from":0,"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1723594000,"include_lower":true,"include_upper":true,"to":1723595000}}}}},"size":5,"sort":[{"a":{"order":"asc"}},{"b":{"order":"asc"}}]}`:  `{"hits":{"total":{"value":123},"hits":[{"_index":"result_table_index","_id":"1","_source":{"a":"1","b":"1"},"sort":["1","1"]},{"_index":"result_table_index","_id":"2","_source":{"a":"1","b":"2"},"sort":["1","2"]},{"_index":"result_table_index","_id":"3","_source":{"a":"1","b":"3"},"sort":["1","3"]},{"_index":"result_table_index","_id":"4","_source":{"a":"1","b":"4"},"sort":["1","4"]},{"_index":"result_table_index","_id":"5","_source":{"a":"1","b":"5"},"sort":["1","5"]}]}}`,
		`{"_source":{"includes":["a","b"]},"from":5,"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1723594000,"include_lower":true,"include_upper":true,"to":1723595000}}}}},"size":5,"sort":[{"a":{"order":"asc"}},{"b":{"order":"asc"}}]}`:  `{"hits":{"total":{"value":123},"hits":[{"_index":"result_table_index","_id":"6","_source":{"a":"2","b":"1"},"sort":["2","1"]},{"_index":"result_table_index","_id":"7","_source":{"a":"2","b":"2"},"sort":["2","2"]},{"_index":"result_table_index","_id":"8","_source":{"a":"2","b":"3"},"sort":["2","3"]},{"_index":"result_table_index","_id":"9","_source":{"a":"2","b":"4"},"sort":["2","4"]},{"_index":"result_table_index","_id":"10","_source":{"a":"2","b":"5"},"sort":["2","5"]}]}}`,
		`{"_source":{"includes":["a","b"]},"from":0,"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1723594000,"include_lower":true,"include_upper":true,"to":1723595000}}}}},"size":10,"sort":[{"a":{"order":"asc"}},{"b":{"order":"asc"}}]}`: `{"hits":{"total":{"value":123},"hits":[{"_index":"result_table_index","_id":"1","_source":{"a":"1","b":"1"},"sort":["1","1"]},{"_index":"result_table_index","_id":"2","_source":{"a":"1","b":"2"},"sort":["1","2"]},{"_index":"result_table_index","_id":"3","_source":{"a":"1","b":"3"},"sort":["1","3"]},{"_index":"result_table_index","_id":"4","_source":{"a":"1","b":"4"},"sort":["1","4"]},{"_index":"result_table_index","_id":"5","_source":{"a":"1","b":"5"},"sort":["1","5"]},{"_index":"result_table_index","_id":"6","_source":{"a":"2","b":"1"},"sort":["2","1"]},{"_index":"result_table_index","_id":"7","_source":{"a":"2","b":"2"},"sort":["2","2"]},{"_index":"result_table_index","_id":"8","_source":{"a":"2","b":"3"},"sort":["2","3"]},{"_index":"result_table_index","_id":"9","_source":{"a":"2","b":"4"},"sort":["2","4"]},{"_index":"result_table_index","_id":"10","_source":{"a":"2","b":"5"},"sort":["2","5"]}]}}`,

		// scroll with 5m
		`{"_source":{"includes":["a","b"]},"query":{"bool":{"filter":{"range":{"dtEventTimeStamp":{"format":"epoch_second","from":1723594000,"include_lower":true,"include_upper":true,"to":1723595000}}}}},"size":5,"sort":[{"a":{"order":"asc"}},{"b":{"order":"asc"}}]}`: `{"_scroll_id":"one","hits":{"total":{"value":123},"hits":[{"_index":"result_table_index","_id":"1","_source":{"a":"1","b":"1"}},{"_index":"result_table_index","_id":"2","_source":{"a":"1","b":"2"}},{"_index":"result_table_index","_id":"3","_source":{"a":"1","b":"3"}},{"_index":"result_table_index","_id":"4","_source":{"a":"1","b":"4"}},{"_index":"result_table_index","_id":"5","_source":{"a":"1","b":"5"}}]}}`,

		// scroll id
		`{"scroll":"5m","scroll_id":"one"}`: `{"_scroll_id":"two","hits":{"total":{"value":123},"hits":[{"_index":"result_table_index","_id":"6","_source":{"a":"2","b":"1"}},{"_index":"result_table_index","_id":"7","_source":{"a":"2","b":"2"}},{"_index":"result_table_index","_id":"8","_source":{"a":"2","b":"3"}},{"_index":"result_table_index","_id":"9","_source":{"a":"2","b":"4"}},{"_index":"result_table_index","_id":"10","_source":{"a":"2","b":"5"}}]}}`,
	})

	tcs := map[string]struct {
		queryTs  *structured.QueryTs
		total    int64
		expected string
		options  string
	}{
		"query with EpochMillis": {
			queryTs: &structured.QueryTs{
				SpaceUid: spaceUid,
				QueryList: []*structured.Query{
					{
						DataSource:  structured.BkLog,
						TableID:     structured.TableID(influxdb.ResultTableBkBaseEs),
						KeepColumns: []string{"__ext.container_id", "dtEventTimeStamp"},
					},
				},
				From:  1,
				Limit: 10,
				Start: "1723594000123",
				End:   "1723595000234",
			},
			total:    1e4,
			expected: `[{"__data_label":"bkbase_es","__doc_id":"e058129ae18bff87c95e83f24584e654","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594211000","dtEventTimeStamp":"1723594211000"},{"__data_label":"bkbase_es","__doc_id":"c124dae69af9b86a7128ee4281820158","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594211000","dtEventTimeStamp":"1723594211000"},{"__data_label":"bkbase_es","__doc_id":"c7f73abf7e865a4b4d7fc608387d01cf","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594211000","dtEventTimeStamp":"1723594211000"},{"__data_label":"bkbase_es","__doc_id":"39c3ec662881e44bf26d2a6bfc0e35c3","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594211000","dtEventTimeStamp":"1723594211000"},{"__data_label":"bkbase_es","__doc_id":"58e03ce0b9754bf0657d49a5513adcb5","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594211000","dtEventTimeStamp":"1723594211000"},{"__data_label":"bkbase_es","__doc_id":"43a36f412886bf30b0746562513638d3","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594211000","dtEventTimeStamp":"1723594211000"},{"__data_label":"bkbase_es","__doc_id":"218ceafd04f89b39cda7954e51f4a48a","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594211000","dtEventTimeStamp":"1723594211000"},{"__data_label":"bkbase_es","__doc_id":"8d9abe9b782fe3a1272c93f0af6b39e1","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-599f9","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594211000","dtEventTimeStamp":"1723594211000"},{"__data_label":"bkbase_es","__doc_id":"0826407be7f04f19086774ed68eac8dd","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-llp94","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594224000","dtEventTimeStamp":"1723594224000"},{"__data_label":"bkbase_es","__doc_id":"d56b4120194eb37f53410780da777d43","__ext.io_kubernetes_pod":"bkmonitor-unify-query-64bd4f5df4-llp94","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.bk_base_es","_time":"1723594224000","dtEventTimeStamp":"1723594224000"}]`,
		},
		"query es with multi rt and multi from 0 - 5": {
			queryTs: &structured.QueryTs{
				SpaceUid: spaceUid,
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableBkBaseEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
				},
				OrderBy: structured.OrderBy{
					"a",
					"b",
				},
				Limit:       5,
				MetricMerge: "a",
				Start:       start,
				End:         end,
				IsMultiFrom: true,
				ResultTableOptions: map[string]*metadata.ResultTableOption{
					"result_table.es|http://127.0.0.1:93002": {
						From: function.IntPoint(0),
					},
					"result_table.bk_base_es|http://127.0.0.1:12001/bk_data/query_sync/es": {
						From: function.IntPoint(5),
					},
				},
			},
			total:    246,
			expected: `[{"__data_label":"es","__doc_id":"1","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"1"},{"__data_label":"es","__doc_id":"2","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"2"},{"__data_label":"es","__doc_id":"3","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"3"},{"__data_label":"es","__doc_id":"4","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"4"},{"__data_label":"es","__doc_id":"5","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"5"}]`,
			options:  `{"result_table.es|http://127.0.0.1:93002":{"from":5},"result_table.bk_base_es|http://127.0.0.1:12001/bk_data/query_sync/es":{"from":5}}`,
		},
		"query es with multi rt and multi from 5 - 10": {
			queryTs: &structured.QueryTs{
				SpaceUid: spaceUid,
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableBkBaseEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
				},
				OrderBy: structured.OrderBy{
					"a",
					"b",
					elasticsearch.KeyTableID,
				},
				Limit:       5,
				MetricMerge: "a",
				Start:       start,
				End:         end,
				IsMultiFrom: true,
				ResultTableOptions: map[string]*metadata.ResultTableOption{
					"result_table.es|http://127.0.0.1:93002": {
						From: function.IntPoint(5),
					},
					"result_table.bk_base_es|http://127.0.0.1:12001/bk_data/query_sync/es": {
						From: function.IntPoint(5),
					},
				},
			},
			total:    246,
			expected: `[{"__data_label":"bkbase_es","__doc_id":"6","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"2","b":"1"},{"__data_label":"es","__doc_id":"6","__index":"result_table_index","__result_table":"result_table.es","a":"2","b":"1"},{"__data_label":"bkbase_es","__doc_id":"7","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"2","b":"2"},{"__data_label":"es","__doc_id":"7","__index":"result_table_index","__result_table":"result_table.es","a":"2","b":"2"},{"__data_label":"bkbase_es","__doc_id":"8","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"2","b":"3"}]`,
			options:  `{"result_table.es|http://127.0.0.1:93002":{"from":7},"result_table.bk_base_es|http://127.0.0.1:12001/bk_data/query_sync/es":{"from":8}}`,
		},
		"query es with multi rt and one from 0 - 5": {
			queryTs: &structured.QueryTs{
				SpaceUid: spaceUid,
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableBkBaseEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
				},
				OrderBy: structured.OrderBy{
					"a",
					"b",
					elasticsearch.KeyTableID,
				},
				From:        0,
				Limit:       5,
				MetricMerge: "a",
				Start:       start,
				End:         end,
			},
			total:    246,
			expected: `[{"__data_label":"bkbase_es","__doc_id":"1","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"1"},{"__data_label":"es","__doc_id":"1","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"1"},{"__data_label":"bkbase_es","__doc_id":"2","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"2"},{"__data_label":"es","__doc_id":"2","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"2"},{"__data_label":"bkbase_es","__doc_id":"3","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"3"}]`,
			options:  `{"result_table.es|http://127.0.0.1:93002":{"search_after":["1","5"]},"result_table.bk_base_es|http://127.0.0.1:12001/bk_data/query_sync/es":{"search_after":["1","5"]}}`,
		},
		"query es with multi rt and one from 5 - 10": {
			queryTs: &structured.QueryTs{
				SpaceUid: spaceUid,
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableBkBaseEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
				},
				OrderBy: structured.OrderBy{
					"a",
					"b",
					elasticsearch.KeyTableID,
				},
				From:        5,
				Limit:       5,
				MetricMerge: "a",
				Start:       start,
				End:         end,
			},
			total:    246,
			expected: `[{"__data_label":"es","__doc_id":"3","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"3"},{"__data_label":"bkbase_es","__doc_id":"4","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"4"},{"__data_label":"es","__doc_id":"4","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"4"},{"__data_label":"bkbase_es","__doc_id":"5","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"5"},{"__data_label":"es","__doc_id":"5","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"5"}]`,
			options:  `{"result_table.bk_base_es|http://127.0.0.1:12001/bk_data/query_sync/es":{"search_after":["2","5"]},"result_table.es|http://127.0.0.1:93002":{"search_after":["2","5"]}}`,
		},
		"query_bk_base_es_1 to 1": {
			queryTs: &structured.QueryTs{
				SpaceUid: spaceUid,
				QueryList: []*structured.Query{
					{
						DataSource:  structured.BkLog,
						TableID:     structured.TableID(influxdb.ResultTableEs),
						From:        1,
						Limit:       1,
						KeepColumns: []string{"__ext.container_id", "dtEventTimeStamp"},
					},
				},
				Start: start,
				End:   end,
			},
			total:    1e4,
			expected: `[{"__data_label":"es","__doc_id":"4f3a5e9c167097c9658e88b2f32364b2","__ext.container_id":"77bd897e66402eb66ee97a1f832fb55b2114d83dc369f01e36ce4cec8483786f","__index":"v2_2_bklog_bk_unify_query_20240814_0","__result_table":"result_table.es","_time":"1723594209000","dtEventTimeStamp":"1723594209000"}]`,
		},
		"query with scroll - 1": {
			queryTs: &structured.QueryTs{
				SpaceUid: spaceUid,
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableBkBaseEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
				},
				OrderBy: structured.OrderBy{
					"a",
					"b",
					elasticsearch.KeyTableID,
				},
				From:        0,
				Limit:       5,
				MetricMerge: "a",
				Start:       start,
				End:         end,
				Scroll:      "5m",
			},
			expected: `[{"__data_label":"bkbase_es","__doc_id":"1","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"1"},{"__data_label":"es","__doc_id":"1","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"1"},{"__data_label":"bkbase_es","__doc_id":"2","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"2"},{"__data_label":"es","__doc_id":"2","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"2"},{"__data_label":"bkbase_es","__doc_id":"3","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"3"},{"__data_label":"es","__doc_id":"3","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"3"},{"__data_label":"bkbase_es","__doc_id":"4","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"4"},{"__data_label":"es","__doc_id":"4","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"4"},{"__data_label":"bkbase_es","__doc_id":"5","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"1","b":"5"},{"__data_label":"es","__doc_id":"5","__index":"result_table_index","__result_table":"result_table.es","a":"1","b":"5"}]`,
			total:    246,
			options:  `{"result_table.es|http://127.0.0.1:93002":{"scroll_id":"one"},"result_table.bk_base_es|http://127.0.0.1:12001/bk_data/query_sync/es":{"scroll_id":"one"}}`,
		},
		"query with scroll - 2": {
			queryTs: &structured.QueryTs{
				SpaceUid: spaceUid,
				QueryList: []*structured.Query{
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
					{
						DataSource:    structured.BkLog,
						TableID:       structured.TableID(influxdb.ResultTableBkBaseEs),
						KeepColumns:   []string{"a", "b"},
						ReferenceName: "a",
					},
				},
				OrderBy: structured.OrderBy{
					"a",
					"b",
					elasticsearch.KeyTableID,
				},
				From:        0,
				Limit:       5,
				MetricMerge: "a",
				Start:       start,
				End:         end,
				ResultTableOptions: map[string]*metadata.ResultTableOption{
					"result_table.es|http://127.0.0.1:93002": {
						ScrollID: "one",
					},
					"result_table.bk_base_es|http://127.0.0.1:12001/bk_data/query_sync/es": {
						ScrollID: "one",
					},
				},
				Scroll: "5m",
			},
			expected: `[{"__data_label":"bkbase_es","__doc_id":"6","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"2","b":"1"},{"__data_label":"es","__doc_id":"6","__index":"result_table_index","__result_table":"result_table.es","a":"2","b":"1"},{"__data_label":"bkbase_es","__doc_id":"7","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"2","b":"2"},{"__data_label":"es","__doc_id":"7","__index":"result_table_index","__result_table":"result_table.es","a":"2","b":"2"},{"__data_label":"bkbase_es","__doc_id":"8","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"2","b":"3"},{"__data_label":"es","__doc_id":"8","__index":"result_table_index","__result_table":"result_table.es","a":"2","b":"3"},{"__data_label":"bkbase_es","__doc_id":"9","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"2","b":"4"},{"__data_label":"es","__doc_id":"9","__index":"result_table_index","__result_table":"result_table.es","a":"2","b":"4"},{"__data_label":"bkbase_es","__doc_id":"10","__index":"result_table_index","__result_table":"result_table.bk_base_es","a":"2","b":"5"},{"__data_label":"es","__doc_id":"10","__index":"result_table_index","__result_table":"result_table.es","a":"2","b":"5"}]`,
			total:    246,
			options:  `{"result_table.es|http://127.0.0.1:93002":{"scroll_id":"two"},"result_table.bk_base_es|http://127.0.0.1:12001/bk_data/query_sync/es":{"scroll_id":"two"}}`,
		},
	}

	for name, c := range tcs {
		t.Run(name, func(t *testing.T) {
			total, list, options, err := queryRawWithInstance(ctx, c.queryTs)
			assert.Nil(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.total, total)

			actual := json.MarshalListMap(list)

			assert.Equal(t, c.expected, actual)

			if len(options) > 0 || c.options != "" {
				optActual, _ := json.Marshal(options)
				assert.JSONEq(t, c.options, string(optActual))
			}
		})
	}
}

// TestQueryExemplar comment lint rebel
func TestQueryExemplar(t *testing.T) {
	ctx := metadata.InitHashID(context.Background())

	mock.Init()
	promql.MockEngine()
	influxdb.MockSpaceRouter(ctx)

	body := []byte(`{"query_list":[{"data_source":"","table_id":"system.cpu_summary","field_name":"usage","field_list":["bk_trace_id","bk_span_id","bk_trace_value","bk_trace_timestamp"],"function":null,"time_aggregation":{"function":"","window":"","position":0,"vargs_list":null},"reference_name":"","dimensions":null,"limit":0,"timestamp":null,"start_or_end":0,"vector_offset":0,"offset":"","offset_forward":false,"slimit":0,"soffset":0,"conditions":{"field_list":[{"field_name":"bk_obj_id","value":["module"],"op":"contains"},{"field_name":"ip","value":["127.0.0.2"],"op":"contains"},{"field_name":"bk_inst_id","value":["14261"],"op":"contains"},{"field_name":"bk_biz_id","value":["7"],"op":"contains"}],"condition_list":["and","and","and"]},"keep_columns":null}],"metric_merge":"","result_columns":null,"start_time":"1677081600","end_time":"1677085600","step":"","down_sample_range":"1m"}`)

	query := &structured.QueryTs{}
	err := json.Unmarshal(body, query)
	assert.Nil(t, err)

	metadata.SetUser(ctx, "", influxdb.SpaceUid, "")

	mock.InfluxDB.Set(map[string]any{
		`select usage as _value, time as _time, bk_trace_id, bk_span_id, bk_trace_value, bk_trace_timestamp from cpu_summary where time > 1677081600000000000 and time < 1677085600000000000 and (bk_obj_id='module' and (ip='127.0.0.2' and (bk_inst_id='14261' and bk_biz_id='7'))) and bk_biz_id='2' and (bk_span_id != '' or bk_trace_id != '')  limit 100000005 slimit 100005`: &decoder.Response{
			Results: []decoder.Result{
				{
					Series: []*decoder.Row{
						{
							Name: "",
							Tags: map[string]string{},
							Columns: []string{
								influxdb.ResultColumnName,
								influxdb.TimeColumnName,
								"bk_trace_id",
								"bk_span_id",
								"bk_trace_value",
								"bk_trace_timestamp",
							},
							Values: [][]any{
								{
									30,
									1677081600000000000,
									"b9cc0e45d58a70b61e8db6fffb5e3376",
									"3d2a373cbeefa1f8",
									1,
									1680157900669,
								},
								{
									21,
									1677081660000000000,
									"fe45f0eccdce3e643a77504f6e6bd87a",
									"c72dcc8fac9bcead",
									1,
									1682121442937,
								},
								{
									1,
									1677081720000000000,
									"771073eb573336a6d3365022a512d6d8",
									"fca46f1c065452e8",
									1,
									1682150008969,
								},
							},
						},
					},
				},
			},
		},
	})

	res, err := queryExemplar(ctx, query)
	assert.Nil(t, err)
	out, err := json.Marshal(res)
	assert.Nil(t, err)
	actual := string(out)
	assert.Equal(t, `{"series":[{"name":"_result0","metric_name":"usage","columns":["_value","_time","bk_trace_id","bk_span_id","bk_trace_value","bk_trace_timestamp"],"types":["float","float","string","string","float","float"],"group_keys":[],"group_values":[],"values":[[30,1677081600000000000,"b9cc0e45d58a70b61e8db6fffb5e3376","3d2a373cbeefa1f8",1,1680157900669],[21,1677081660000000000,"fe45f0eccdce3e643a77504f6e6bd87a","c72dcc8fac9bcead",1,1682121442937],[1,1677081720000000000,"771073eb573336a6d3365022a512d6d8","fca46f1c065452e8",1,1682150008969]]}]}`, actual)
}

func TestVmQueryParams(t *testing.T) {
	ctx := metadata.InitHashID(context.Background())

	mock.Init()
	promql.MockEngine()

	testCases := []struct {
		username string
		spaceUid string
		query    string
		promql   string
		start    string
		end      string
		step     string
		params   string
		error    error
	}{
		{
			username: "vm-query",
			spaceUid: consul.VictoriaMetricsStorageType,
			query:    `{"query_list":[{"field_name":"bk_split_measurement","function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"increase","window":"1m0s"},"reference_name":"a","conditions":{"field_list":[{"field_name":"bcs_cluster_id","value":["cls-2"],"op":"req"},{"field_name":"bcs_cluster_id","value":["cls-2"],"op":"req"},{"field_name":"bk_biz_id","value":["100801"],"op":"eq"}],"condition_list":["and", "and"]}},{"field_name":"bk_split_measurement","function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"delta","window":"1m0s"},"reference_name":"b"}],"metric_merge":"a / b","start_time":"0","end_time":"600","step":"60s"}`,
			params:   `{"influx_compatible":true,"use_native_or":true,"api_type":"query_range","cluster_name":"","api_params":{"query":"sum by (bcs_cluster_id, namespace) (increase(a[1m] offset -59s999ms)) / sum by (bcs_cluster_id, namespace) (delta(b[1m] offset -59s999ms))","start":0,"end":600,"step":60},"result_table_list":["victoria_metrics"],"metric_filter_condition":{"a":"filter=\"bk_split_measurement\", bcs_cluster_id=~\"cls-2\", bcs_cluster_id=~\"cls-2\", bk_biz_id=\"100801\", result_table_id=\"victoria_metrics\", __name__=\"bk_split_measurement_value\"","b":"filter=\"bk_split_measurement\", result_table_id=\"victoria_metrics\", __name__=\"bk_split_measurement_value\""}}`,
		},
		{
			username: "vm-query-or",
			spaceUid: "vm-query",
			query:    `{"query_list":[{"field_name":"container_cpu_usage_seconds_total","field_list":null,"function":[{"method":"sum","without":false,"dimensions":[],"position":0,"args_list":null,"vargs_list":null}],"time_aggregation":{"function":"count_over_time","window":"60s","position":0,"vargs_list":null},"reference_name":"a","dimensions":[],"limit":0,"timestamp":null,"start_or_end":0,"vector_offset":0,"offset":"","offset_forward":false,"slimit":0,"soffset":0,"conditions":{"field_list":[{"field_name":"bk_biz_id","value":["7"],"op":"contains"},{"field_name":"ip","value":["127.0.0.1","127.0.0.2"],"op":"contains"},{"field_name":"ip","value":["[a-z]","[A-Z]"],"op":"req"},{"field_name":"api","value":["/metrics"],"op":"ncontains"},{"field_name":"bk_biz_id","value":["7"],"op":"contains"},{"field_name":"api","value":["/metrics"],"op":"contains"}],"condition_list":["and","and","and","or","and"]},"keep_columns":["_time","a"]}],"metric_merge":"a","result_columns":null,"start_time":"1697458200","end_time":"1697461800","step":"60s","down_sample_range":"3s","timezone":"Asia/Shanghai","look_back_delta":"","instant":false}`,
			params:   `{"influx_compatible":true,"use_native_or":true,"api_type":"query_range","cluster_name":"","api_params":{"query":"sum(count_over_time(a[1m] offset -59s999ms))","start":1697458200,"end":1697461800,"step":60},"result_table_list":["100147_bcs_prom_computation_result_table_25428","100147_bcs_prom_computation_result_table_25429"],"metric_filter_condition":{"a":"bcs_cluster_id=\"BCS-K8S-25428\", bk_biz_id=\"7\", ip=~\"^(127\\\\.0\\\\.0\\\\.1|127\\\\.0\\\\.0\\\\.2)$\", ip=~\"[a-z]|[A-Z]\", api!=\"/metrics\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=\"container_cpu_usage_seconds_total_value\" or bcs_cluster_id=\"BCS-K8S-25428\", bk_biz_id=\"7\", api=\"/metrics\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=\"container_cpu_usage_seconds_total_value\" or bcs_cluster_id=\"BCS-K8S-25430\", bk_biz_id=\"7\", ip=~\"^(127\\\\.0\\\\.0\\\\.1|127\\\\.0\\\\.0\\\\.2)$\", ip=~\"[a-z]|[A-Z]\", api!=\"/metrics\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=\"container_cpu_usage_seconds_total_value\" or bcs_cluster_id=\"BCS-K8S-25430\", bk_biz_id=\"7\", api=\"/metrics\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=\"container_cpu_usage_seconds_total_value\" or bcs_cluster_id=\"BCS-K8S-25429\", bk_biz_id=\"7\", ip=~\"^(127\\\\.0\\\\.0\\\\.1|127\\\\.0\\\\.0\\\\.2)$\", ip=~\"[a-z]|[A-Z]\", api!=\"/metrics\", result_table_id=\"100147_bcs_prom_computation_result_table_25429\", __name__=\"container_cpu_usage_seconds_total_value\" or bcs_cluster_id=\"BCS-K8S-25429\", bk_biz_id=\"7\", api=\"/metrics\", result_table_id=\"100147_bcs_prom_computation_result_table_25429\", __name__=\"container_cpu_usage_seconds_total_value\""}}`,
		},
		{
			username: "vm-query-or-for-internal",
			spaceUid: "vm-query",
			promql:   `{"promql":"sum by(job, metric_name) (delta(label_replace({__name__=~\"container_cpu_.+_total\", __name__ !~ \".+_size_count\", __name__ !~ \".+_process_time_count\", job=\"metric-social-friends-forever\"}, \"metric_name\", \"$1\", \"__name__\", \"ffs_rest_(.*)_count\")[2m:]))","start":"1698147600","end":"1698151200","step":"60s","bk_biz_ids":null,"timezone":"Asia/Shanghai","look_back_delta":"","instant":false}`,
			params:   `{"influx_compatible":true,"use_native_or":true,"api_type":"query_range","cluster_name":"","api_params":{"query":"sum by (job, metric_name) (delta(label_replace({__name__=~\"a\"} offset -59s999ms, \"metric_name\", \"$1\", \"__name__\", \"ffs_rest_(.*)_count_value\")[2m:]))","start":1698147600,"end":1698151200,"step":60},"result_table_list":["100147_bcs_prom_computation_result_table_25428","100147_bcs_prom_computation_result_table_25429"],"metric_filter_condition":{"a":"bcs_cluster_id=\"BCS-K8S-25428\", __name__!~\".+_size_count_value\", __name__!~\".+_process_time_count_value\", job=\"metric-social-friends-forever\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=~\"container_cpu_.+_total_value\" or bcs_cluster_id=\"BCS-K8S-25430\", __name__!~\".+_size_count_value\", __name__!~\".+_process_time_count_value\", job=\"metric-social-friends-forever\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=~\"container_cpu_.+_total_value\" or bcs_cluster_id=\"BCS-K8S-25429\", __name__!~\".+_size_count_value\", __name__!~\".+_process_time_count_value\", job=\"metric-social-friends-forever\", result_table_id=\"100147_bcs_prom_computation_result_table_25429\", __name__=~\"container_cpu_.+_total_value\""}}`,
		},
		{
			username: "vm-query",
			spaceUid: "vm-query",
			query:    `{"query_list":[{"field_name":"container_cpu_usage_seconds_total","function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"sum_over_time","window":"1m0s"},"reference_name":"a","conditions":{"field_list":[{"field_name":"bcs_cluster_id","value":["cls-2"],"op":"req"},{"field_name":"bcs_cluster_id","value":["cls-2"],"op":"req"},{"field_name":"bk_biz_id","value":["100801"],"op":"eq"}],"condition_list":["or", "and"]}},{"field_name":"container_cpu_usage_seconds_total","function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"count_over_time","window":"1m0s"},"reference_name":"b"}],"metric_merge":"a / b","start_time":"0","end_time":"600","step":"60s"}`,
			params:   `{"influx_compatible":true,"use_native_or":true,"api_type":"query_range","cluster_name":"","api_params":{"query":"sum by (bcs_cluster_id, namespace) (sum_over_time(a[1m] offset -59s999ms)) / sum by (bcs_cluster_id, namespace) (count_over_time(b[1m] offset -59s999ms))","start":0,"end":600,"step":60},"result_table_list":["100147_bcs_prom_computation_result_table_25428","100147_bcs_prom_computation_result_table_25429"],"metric_filter_condition":{"b":"bcs_cluster_id=\"BCS-K8S-25429\", result_table_id=\"100147_bcs_prom_computation_result_table_25429\", __name__=\"container_cpu_usage_seconds_total_value\" or bcs_cluster_id=\"BCS-K8S-25428\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=\"container_cpu_usage_seconds_total_value\" or bcs_cluster_id=\"BCS-K8S-25430\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=\"container_cpu_usage_seconds_total_value\""}}`,
		},
		{
			username: "vm-query",
			spaceUid: "vm-query",
			query:    `{"query_list":[{"field_name":"metric","function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"sum_over_time","window":"1m0s"},"reference_name":"a","conditions":{"field_list":[{"field_name":"bcs_cluster_id","value":["cls-2"],"op":"req"},{"field_name":"bcs_cluster_id","value":["cls-2"],"op":"req"},{"field_name":"bk_biz_id","value":["100801"],"op":"eq"}],"condition_list":["and","and"]}},{"field_name":"metric","function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"count_over_time","window":"1m0s"},"reference_name":"b"}],"metric_merge":"a / b","start_time":"0","end_time":"600","step":"60s"}`,
			params:   `{"influx_compatible":true,"use_native_or":true,"api_type":"query_range","cluster_name":"","api_params":{"query":"sum by (bcs_cluster_id, namespace) (sum_over_time(a[1m] offset -59s999ms)) / sum by (bcs_cluster_id, namespace) (count_over_time(b[1m] offset -59s999ms))","start":0,"end":600,"step":60},"result_table_list":["vm_rt"],"metric_filter_condition":{"a":"bcs_cluster_id=\"cls\", bcs_cluster_id=~\"cls-2\", bcs_cluster_id=~\"cls-2\", bk_biz_id=\"100801\", result_table_id=\"vm_rt\", __name__=\"metric_value\"","b":"bcs_cluster_id=\"cls\", result_table_id=\"vm_rt\", __name__=\"metric_value\""}}`,
		},
		{
			username: "vm-query",
			spaceUid: "vm-query",
			query:    `{"query_list":[{"field_name":"metric","function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"sum_over_time","window":"1m0s"},"reference_name":"a","conditions":{"field_list":[{"field_name":"namespace","value":["ns"],"op":"contains"}],"condition_list":[]}},{"field_name":"metric","function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"count_over_time","window":"1m0s"},"reference_name":"b"}],"metric_merge":"a / b","start_time":"0","end_time":"600","step":"60s"}`,
			params:   `{"influx_compatible":true,"use_native_or":true,"api_type":"query_range","cluster_name":"","api_params":{"query":"sum by (bcs_cluster_id, namespace) (sum_over_time(a[1m] offset -59s999ms)) / sum by (bcs_cluster_id, namespace) (count_over_time(b[1m] offset -59s999ms))","start":0,"end":600,"step":60},"result_table_list":["vm_rt"],"metric_filter_condition":{"a":"bcs_cluster_id=\"cls\", namespace=\"ns\", result_table_id=\"vm_rt\", __name__=\"metric_value\"","b":"bcs_cluster_id=\"cls\", result_table_id=\"vm_rt\", __name__=\"metric_value\""}}`,
		},
		{
			username: "vm-query-fuzzy-name",
			spaceUid: "vm-query",
			query:    `{"query_list":[{"field_name":"me.*","is_regexp":true,"function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"sum_over_time","window":"1m0s"},"reference_name":"a","conditions":{"field_list":[{"field_name":"namespace","value":["ns"],"op":"contains"}],"condition_list":[]}},{"field_name":"metric","function":[{"method":"sum","dimensions":["bcs_cluster_id","namespace"]}],"time_aggregation":{"function":"count_over_time","window":"1m0s"},"reference_name":"b"}],"metric_merge":"a / b","start_time":"0","end_time":"600","step":"60s"}`,
			params:   `{"influx_compatible":true,"use_native_or":true,"api_type":"query_range","cluster_name":"","api_params":{"query":"sum by (bcs_cluster_id, namespace) (sum_over_time({__name__=~\"a\"}[1m] offset -59s999ms)) / sum by (bcs_cluster_id, namespace) (count_over_time(b[1m] offset -59s999ms))","start":0,"end":600,"step":60},"result_table_list":["vm_rt"],"metric_filter_condition":{"a":"bcs_cluster_id=\"cls\", namespace=\"ns\", result_table_id=\"vm_rt\", __name__=~\"me.*_value\"","b":"bcs_cluster_id=\"cls\", result_table_id=\"vm_rt\", __name__=\"metric_value\""}}`,
		},
		{
			username: "vm-query",
			spaceUid: "vm-query",
			promql:   `{"promql":"max_over_time((increase(container_cpu_usage_seconds_total{}[10m]) \u003e 0)[1h:])","start":"1720765200","end":"1720786800","step":"10m","bk_biz_ids":null,"timezone":"Asia/Shanghai","look_back_delta":"","instant":false}`,
			params:   `{"influx_compatible":true,"use_native_or":true,"api_type":"query_range","cluster_name":"","api_params":{"query":"max_over_time((increase(a[10m] offset -9m59s999ms) \u003e 0)[1h:])","start":1720765200,"end":1720786800,"step":600},"result_table_list":["100147_bcs_prom_computation_result_table_25428","100147_bcs_prom_computation_result_table_25429"],"metric_filter_condition":{"a":"bcs_cluster_id=\"BCS-K8S-25428\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=\"container_cpu_usage_seconds_total_value\" or bcs_cluster_id=\"BCS-K8S-25430\", result_table_id=\"100147_bcs_prom_computation_result_table_25428\", __name__=\"container_cpu_usage_seconds_total_value\" or bcs_cluster_id=\"BCS-K8S-25429\", result_table_id=\"100147_bcs_prom_computation_result_table_25429\", __name__=\"container_cpu_usage_seconds_total_value\""}}`,
		},
	}

	for i, c := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var (
				query *structured.QueryTs
				err   error
			)
			ctx := metadata.InitHashID(ctx)
			metadata.SetUser(ctx, fmt.Sprintf("username:%s", c.username), c.spaceUid, "")

			if c.promql != "" {
				var queryPromQL *structured.QueryPromQL
				err = json.Unmarshal([]byte(c.promql), &queryPromQL)
				assert.Nil(t, err)
				query, err = promQLToStruct(ctx, queryPromQL)
			} else {
				err = json.Unmarshal([]byte(c.query), &query)
			}

			query.SpaceUid = c.spaceUid
			assert.Nil(t, err)
			_, err = queryTsWithPromEngine(ctx, query)
			if c.error != nil {
				assert.Contains(t, err.Error(), c.error.Error())
			} else {
				var vmParams map[string]string
				if vmParams != nil {
					assert.Equal(t, c.params, vmParams["sql"])
				}
			}
		})
	}
}

func TestStructAndPromQLConvert(t *testing.T) {
	ctx := metadata.InitHashID(context.Background())

	mock.Init()
	promql.MockEngine()

	testCase := map[string]struct {
		queryStruct bool
		query       *structured.QueryTs
		promql      *structured.QueryPromQL
		err         error
	}{
		"query struct with or": {
			queryStruct: true,
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "custom",
						TableID:    "dataLabel",
						FieldName:  "metric",
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function: "sum_over_time",
							Window:   "1m0s",
						},
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2",
									},
									Operator: "req",
								},
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2",
									},
									Operator: "req",
								},
							},
							ConditionList: []string{
								"or",
							},
						},
						ReferenceName: "a",
					},
				},
				MetricMerge: "a",
				Start:       "1691132705",
				End:         "1691136305",
				Step:        "1m",
			},
			err: fmt.Errorf("or 过滤条件无法直接转换为 promql 语句，请使用结构化查询"),
		},
		"query struct with and": {
			queryStruct: true,
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "custom",
						TableID:    "dataLabel",
						FieldName:  "metric",
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function:  "sum_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2",
									},
									Operator: "req",
								},
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2",
									},
									Operator: "req",
								},
							},
							ConditionList: []string{
								"and",
							},
						},
						ReferenceName: "a",
					},
				},
				MetricMerge: "a",
				Start:       `1691132705`,
				End:         `1691136305`,
				Step:        `1m`,
			},
			promql: &structured.QueryPromQL{
				PromQL: `sum by (bcs_cluster_id, result_table_id) (sum_over_time(custom:dataLabel:metric{bcs_cluster_id=~"cls-2",bcs_cluster_id=~"cls-2"}[1m]))`,
				Start:  `1691132705`,
				End:    `1691136305`,
				Step:   `1m`,
			},
		},
		"promql struct with and": {
			queryStruct: true,
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "custom",
						TableID:    "dataLabel",
						FieldName:  "metric",
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function:  "sum_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2",
									},
									Operator: "req",
								},
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2",
									},
									Operator: "req",
								},
							},
							ConditionList: []string{
								"and",
							},
						},
						ReferenceName: "a",
					},
				},
				MetricMerge: "a",
				Start:       `1691132705`,
				End:         `1691136305`,
				Step:        `1m`,
			},
			promql: &structured.QueryPromQL{
				PromQL: `sum by (bcs_cluster_id, result_table_id) (sum_over_time(custom:dataLabel:metric{bcs_cluster_id=~"cls-2",bcs_cluster_id=~"cls-2"}[1m]))`,
				Start:  `1691132705`,
				End:    `1691136305`,
				Step:   `1m`,
			},
		},
		"promql struct 1": {
			queryStruct: true,
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: structured.BkMonitor,
						FieldName:  "container_cpu_usage_seconds_total",
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function:  "sum_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2|cls-2",
									},
									Operator: "req",
								},
								{
									DimensionName: "bk_biz_id",
									Value: []string{
										"2",
									},
									Operator: "eq",
								},
							},
							ConditionList: []string{
								"and",
							},
						},
						ReferenceName: "a",
					},
					{
						DataSource: structured.BkMonitor,
						FieldName:  "container_cpu_usage_seconds_total",
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function:  "count_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						Conditions: structured.Conditions{
							FieldList:     []structured.ConditionField{},
							ConditionList: []string{},
						},
						ReferenceName: "b",
					},
				},
				MetricMerge: "a / on (bcs_cluster_id) group_left () b",
				Start:       `1691132705`,
				End:         `1691136305`,
				Step:        `1m`,
			},
			promql: &structured.QueryPromQL{
				PromQL: `sum by (bcs_cluster_id, result_table_id) (sum_over_time(bkmonitor:container_cpu_usage_seconds_total{bcs_cluster_id=~"cls-2|cls-2",bk_biz_id="2"}[1m])) / on (bcs_cluster_id) group_left () sum by (bcs_cluster_id, result_table_id) (count_over_time(bkmonitor:container_cpu_usage_seconds_total[1m]))`,
				Start:  `1691132705`,
				End:    `1691136305`,
				Step:   `1m`,
			},
		},
		"query struct 1": {
			queryStruct: true,
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: structured.BkMonitor,
						FieldName:  "container_cpu_usage_seconds_total",
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function:  "sum_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2|cls-2",
									},
									Operator: "req",
								},
								{
									DimensionName: "bk_biz_id",
									Value: []string{
										"2",
									},
									Operator: "eq",
								},
							},
							ConditionList: []string{
								"and",
							},
						},
						ReferenceName: "a",
					},
					{
						DataSource: structured.BkMonitor,
						FieldName:  "container_cpu_usage_seconds_total",
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function:  "count_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						Conditions: structured.Conditions{
							FieldList:     []structured.ConditionField{},
							ConditionList: []string{},
						},
						ReferenceName: "b",
					},
				},
				MetricMerge: "a / on (bcs_cluster_id) group_left () b",
				Start:       `1691132705`,
				End:         `1691136305`,
				Step:        `1m`,
			},
			promql: &structured.QueryPromQL{
				PromQL: `sum by (bcs_cluster_id, result_table_id) (sum_over_time(bkmonitor:container_cpu_usage_seconds_total{bcs_cluster_id=~"cls-2|cls-2",bk_biz_id="2"}[1m])) / on (bcs_cluster_id) group_left () sum by (bcs_cluster_id, result_table_id) (count_over_time(bkmonitor:container_cpu_usage_seconds_total[1m]))`,
				Start:  `1691132705`,
				End:    `1691136305`,
				Step:   `1m`,
			},
		},
		"query struct with __name__ ": {
			queryStruct: false,
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: structured.BkMonitor,
						TableID:    "table_id",
						FieldName:  ".*",
						IsRegexp:   true,
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function:  "sum_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						ReferenceName: "a",
						Dimensions:    nil,
						Limit:         0,
						Timestamp:     nil,
						StartOrEnd:    0,
						VectorOffset:  0,
						Offset:        "",
						OffsetForward: false,
						Slimit:        0,
						Soffset:       0,
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2|cls-2",
									},
									Operator: "req",
								},
								{
									DimensionName: "bk_biz_id",
									Value: []string{
										"2",
									},
									Operator: "eq",
								},
							},
							ConditionList: []string{
								"and",
							},
						},
						KeepColumns:         nil,
						AlignInfluxdbResult: false,
						Start:               "",
						End:                 "",
						Step:                "",
						Timezone:            "",
					},
					{
						DataSource: structured.BkMonitor,
						TableID:    "table_id",
						FieldName:  ".*",
						IsRegexp:   true,
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function:  "count_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						Conditions: structured.Conditions{
							FieldList:     []structured.ConditionField{},
							ConditionList: []string{},
						},
						ReferenceName: "b",
					},
				},
				MetricMerge: "a / on (bcs_cluster_id) group_left () b",
				Start:       `1691132705`,
				End:         `1691136305`,
				Step:        `1m`,
			},
			promql: &structured.QueryPromQL{
				PromQL: `sum by (bcs_cluster_id, result_table_id) (sum_over_time({__name__=~"bkmonitor:table_id:.*",bcs_cluster_id=~"cls-2|cls-2",bk_biz_id="2"}[1m])) / on (bcs_cluster_id) group_left () sum by (bcs_cluster_id, result_table_id) (count_over_time({__name__=~"bkmonitor:table_id:.*"}[1m]))`,
				Start:  `1691132705`,
				End:    `1691136305`,
				Step:   `1m`,
			},
		},
		"promql struct with __name__ ": {
			queryStruct: true,
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: structured.BkMonitor,
						TableID:    "table_id",
						FieldName:  ".*",
						IsRegexp:   true,
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function:  "sum_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						ReferenceName: "a",
						Dimensions:    nil,
						Limit:         0,
						Timestamp:     nil,
						StartOrEnd:    0,
						VectorOffset:  0,
						Offset:        "",
						OffsetForward: false,
						Slimit:        0,
						Soffset:       0,
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "bcs_cluster_id",
									Value: []string{
										"cls-2|cls-2",
									},
									Operator: "req",
								},
								{
									DimensionName: "bk_biz_id",
									Value: []string{
										"2",
									},
									Operator: "eq",
								},
							},
							ConditionList: []string{
								"and",
							},
						},
						KeepColumns:         nil,
						AlignInfluxdbResult: false,
						Start:               "",
						End:                 "",
						Step:                "",
						Timezone:            "",
					},
					{
						DataSource: structured.BkMonitor,
						TableID:    "table_id",
						FieldName:  ".*",
						IsRegexp:   true,
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"bcs_cluster_id",
									"result_table_id",
								},
							},
						},
						TimeAggregation: structured.TimeAggregation{
							Function: "count_over_time",
							Window:   "1m0s",
						},
						Conditions: structured.Conditions{
							FieldList:     []structured.ConditionField{},
							ConditionList: []string{},
						},
						ReferenceName: "b",
					},
				},
				MetricMerge: "a / on (bcs_cluster_id) group_left () b",
				Start:       `1691132705`,
				End:         `1691136305`,
				Step:        `1m`,
			},
			promql: &structured.QueryPromQL{
				PromQL: `sum by (bcs_cluster_id, result_table_id) (sum_over_time({__name__=~"bkmonitor:table_id:.*",bcs_cluster_id=~"cls-2|cls-2",bk_biz_id="2"}[1m])) / on (bcs_cluster_id) group_left () sum by (bcs_cluster_id, result_table_id) (count_over_time({__name__=~"bkmonitor:table_id:.*"}[1m]))`,
				Start:  `1691132705`,
				End:    `1691136305`,
				Step:   `1m`,
			},
		},
		"promql to struct with 1m": {
			queryStruct: true,
			promql: &structured.QueryPromQL{
				PromQL: `count_over_time(bkmonitor:metric[1m] @ start() offset -29s999ms)`,
				Start:  `1691132705`,
				End:    `1691136305`,
				Step:   `30s`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						AlignInfluxdbResult: true,
						DataSource:          `bkmonitor`,
						FieldName:           `metric`,
						StartOrEnd:          parser.START,
						//Offset:              "59s999ms",
						OffsetForward: true,
						TimeAggregation: structured.TimeAggregation{
							Function:  "count_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						Conditions: structured.Conditions{
							FieldList:     []structured.ConditionField{},
							ConditionList: []string{},
						},
						ReferenceName: `a`,
						Step:          `30s`,
					},
				},
				MetricMerge: "a",
				Start:       `1691132705`,
				End:         `1691136305`,
				Step:        `30s`,
			},
		},
		"promql to struct with delta label_replace 1m:2m": {
			queryStruct: true,
			promql: &structured.QueryPromQL{
				PromQL: `sum by (job, metric_name) (delta(label_replace({__name__=~"bkmonitor:container_cpu_.+_total",job="metric-social-friends-forever"} @ start() offset -29s999ms, "metric_name", "$1", "__name__", "ffs_rest_(.*)_count")[2m:]))`,
				Start:  `1691132705`,
				End:    `1691136305`,
				Step:   `30s`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:          `bkmonitor`,
						FieldName:           `container_cpu_.+_total`,
						IsRegexp:            true,
						StartOrEnd:          parser.START,
						AlignInfluxdbResult: true,
						TimeAggregation: structured.TimeAggregation{
							Function:   "delta",
							Window:     "2m0s",
							NodeIndex:  3,
							IsSubQuery: true,
							Step:       "0s",
						},
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "job",
									Operator:      "eq",
									Value: []string{
										"metric-social-friends-forever",
									},
								},
								//{
								//	DimensionName: "__name__",
								//	Operator:      "nreq",
								//	Value: []string{
								//		".+_size_count",
								//	},
								//},
								//{
								//	DimensionName: "__name__",
								//	Operator:      "nreq",
								//	Value: []string{
								//		".+_process_time_count",
								//	},
								//},
							},
							ConditionList: []string{},
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "label_replace",
								VArgsList: []interface{}{
									"metric_name",
									"$1",
									"__name__",
									"ffs_rest_(.*)_count",
								},
							},
							{
								Method: "sum",
								Dimensions: []string{
									"job",
									"metric_name",
								},
							},
						},
						ReferenceName: `a`,
						Offset:        "0s",
					},
				},
				MetricMerge: "a",
				Start:       `1691132705`,
				End:         `1691136305`,
				Step:        `30s`,
			},
		},
		"promql to struct with topk": {
			queryStruct: false,
			promql: &structured.QueryPromQL{
				PromQL: `topk(1, bkmonitor:metric)`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						FieldName:  "metric",
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "topk",
								VArgsList: []interface{}{
									1,
								},
							},
						},
						Conditions: structured.Conditions{
							FieldList:     []structured.ConditionField{},
							ConditionList: []string{},
						},
						ReferenceName: "a",
					},
				},
				MetricMerge: "a",
			},
		},
		"promql to struct with delta(metric[1m])`": {
			queryStruct: false,
			promql: &structured.QueryPromQL{
				PromQL: `delta(bkmonitor:metric[1m])`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						FieldName:  "metric",
						TimeAggregation: structured.TimeAggregation{
							Function:  "delta",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						Conditions: structured.Conditions{
							FieldList:     []structured.ConditionField{},
							ConditionList: []string{},
						},
						ReferenceName: "a",
					},
				},
				MetricMerge: "a",
			},
		},
		"promq to struct with metric @end()`": {
			queryStruct: false,
			promql: &structured.QueryPromQL{
				PromQL: `bkmonitor:metric @ end()`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						FieldName:  "metric",
						StartOrEnd: parser.END,
						Conditions: structured.Conditions{
							FieldList:     []structured.ConditionField{},
							ConditionList: []string{},
						},
						ReferenceName: "a",
					},
				},
				MetricMerge: "a",
			},
		},
		"promql to struct with condition contains`": {
			queryStruct: true,
			promql: &structured.QueryPromQL{
				PromQL: `bkmonitor:metric{dim_contains=~"^(val-1|val-2|val-3)$",dim_req=~"val-1|val-2|val-3"} @ end()`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						FieldName:  "metric",
						StartOrEnd: parser.END,
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "dim_contains",
									Value: []string{
										"val-1",
										"val-2",
										"val-3",
									},
									Operator: "contains",
								},
								{
									DimensionName: "dim_req",
									Value: []string{
										"val-1",
										"val-2",
										"val-3",
									},
									Operator: "req",
								},
							},
							ConditionList: []string{
								"and",
							},
						},
						ReferenceName: "a",
					},
				},
				MetricMerge: "a",
			},
		},
		"quantile and quantile_over_time": {
			queryStruct: true,
			promql: &structured.QueryPromQL{
				PromQL: `quantile(0.9, quantile_over_time(0.9, bkmonitor:metric[1m]))`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						FieldName:  "metric",
						Conditions: structured.Conditions{
							FieldList:     []structured.ConditionField{},
							ConditionList: []string{},
						},
						ReferenceName: "a",
						TimeAggregation: structured.TimeAggregation{
							Function:  "quantile_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
							VargsList: []interface{}{
								0.9,
							},
							Position: 1,
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "quantile",
								VArgsList: []interface{}{
									0.9,
								},
							},
						},
					},
				},
				MetricMerge: "a",
			},
		},
		"nodeIndex 3 with sum": {
			queryStruct: false,
			promql: &structured.QueryPromQL{
				PromQL: `increase(sum by (deployment_environment, result_table_id) (bkmonitor:5000575_bkapm_metric_tgf_server_gs_cn_idctest:__default__:trace_additional_duration_count{deployment_environment="g-5"})[2m:])`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						TableID:    "5000575_bkapm_metric_tgf_server_gs_cn_idctest.__default__",
						FieldName:  "trace_additional_duration_count",
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "deployment_environment",
									Value:         []string{"g-5"},
									Operator:      "eq",
								},
							},
							ConditionList: []string{},
						},
						ReferenceName: "a",
						TimeAggregation: structured.TimeAggregation{
							Function:   "increase",
							Window:     "2m0s",
							NodeIndex:  3,
							IsSubQuery: true,
							Step:       "0s",
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"deployment_environment", "result_table_id",
								},
							},
						},
						Offset: "0s",
					},
				},
				MetricMerge: "a",
			},
		},
		"nodeIndex 2 with sum": {
			queryStruct: false,
			promql: &structured.QueryPromQL{
				PromQL: `sum by (deployment_environment, result_table_id) (increase(bkmonitor:5000575_bkapm_metric_tgf_server_gs_cn_idctest:__default__:trace_additional_duration_count{deployment_environment="g-5"}[2m]))`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						TableID:    "5000575_bkapm_metric_tgf_server_gs_cn_idctest.__default__",
						FieldName:  "trace_additional_duration_count",
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "deployment_environment",
									Value:         []string{"g-5"},
									Operator:      "eq",
								},
							},
							ConditionList: []string{},
						},
						ReferenceName: "a",
						TimeAggregation: structured.TimeAggregation{
							Function:  "increase",
							Window:    "2m0s",
							NodeIndex: 2,
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "sum",
								Dimensions: []string{
									"deployment_environment", "result_table_id",
								},
							},
						},
					},
				},
				MetricMerge: "a",
			},
		},
		"predict_linear": {
			queryStruct: false,
			promql: &structured.QueryPromQL{
				PromQL: `predict_linear(bkmonitor:metric[1h], 4*3600)`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    "bkmonitor",
						TableID:       "",
						FieldName:     "metric",
						ReferenceName: "a",
						TimeAggregation: structured.TimeAggregation{
							Function:  "predict_linear",
							Window:    "1h0m0s",
							NodeIndex: 2,
							VargsList: []interface{}{4 * 3600},
						},
					},
				},
				MetricMerge: "a",
			},
		},
		"promql to struct with many time aggregate": {
			queryStruct: true,
			promql: &structured.QueryPromQL{
				PromQL: `min_over_time(increase(bkmonitor:metric[1m])[2m:])`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    "bkmonitor",
						TableID:       "",
						FieldName:     "metric",
						ReferenceName: "a",
						TimeAggregation: structured.TimeAggregation{
							Function:  "increase",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method:     "min_over_time",
								Window:     "2m0s",
								IsSubQuery: true,
								Step:       "0s",
							},
						},
						Offset: "0s",
					},
				},
				MetricMerge: "a",
			},
		},
		"promql to struct with many time aggregate and funciton": {
			queryStruct: true,
			promql: &structured.QueryPromQL{
				PromQL: `topk(5, floor(sum by (dim) (last_over_time(min_over_time(increase(label_replace(bkmonitor:metric, "name", "$0", "__name__", ".+")[1m:])[2m:])[3m:15s]))))`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource:    "bkmonitor",
						TableID:       "",
						FieldName:     "metric",
						ReferenceName: "a",
						TimeAggregation: structured.TimeAggregation{
							Function:   "increase",
							Window:     "1m0s",
							NodeIndex:  3,
							IsSubQuery: true,
							Step:       "0s",
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method: "label_replace",
								VArgsList: []interface{}{
									"name",
									"$0",
									"__name__",
									".+",
								},
							},
							{
								Method:     "min_over_time",
								Window:     "2m0s",
								IsSubQuery: true,
								Step:       "0s",
							},
							{
								Method:     "last_over_time",
								Window:     "3m0s",
								IsSubQuery: true,
								Step:       "15s",
							},
							{
								Method:     "sum",
								Dimensions: []string{"dim"},
							},
							{
								Method: "floor",
							},
							{
								Method: "topk",
								VArgsList: []interface{}{
									5,
								},
							},
						},
						Offset: "0s",
					},
				},
				MetricMerge: "a",
			},
		},
		"promql with match - 1": {
			queryStruct: false,
			promql: &structured.QueryPromQL{
				PromQL: `sum by (pod_name, bcs_cluster_id, namespace,instance) (rate(container_cpu_usage_seconds_total{namespace="ns-1"}[2m])) / on(bcs_cluster_id, namespace, pod_name) group_left() sum (sum_over_time(kube_pod_container_resource_limits_cpu_cores{namespace="ns-1"}[1m])) by (pod_name, bcs_cluster_id,namespace)`,
				Match:  `{pod_name="pod", bcs_cluster_id!="cls-1", namespace="ns-1", instance="ins-1"}`,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						FieldName:  "container_cpu_usage_seconds_total",
						TimeAggregation: structured.TimeAggregation{
							Function:  "rate",
							Window:    "2m0s",
							NodeIndex: 2,
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method:     "sum",
								Dimensions: []string{"pod_name", "bcs_cluster_id", "namespace", "instance"},
							},
						},
						ReferenceName: "a",
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "pod_name",
									Operator:      structured.ConditionEqual,
									Value:         []string{"pod"},
								},
								{
									DimensionName: "bcs_cluster_id",
									Operator:      structured.ConditionNotEqual,
									Value:         []string{"cls-1"},
								},
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "instance",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ins-1"},
								},
							},
							ConditionList: []string{"and", "and", "and", "and"},
						},
					},
					{
						DataSource: "bkmonitor",
						FieldName:  "kube_pod_container_resource_limits_cpu_cores",
						TimeAggregation: structured.TimeAggregation{
							Function:  "sum_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method:     "sum",
								Dimensions: []string{"pod_name", "bcs_cluster_id", "namespace"},
							},
						},
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "pod_name",
									Operator:      structured.ConditionEqual,
									Value:         []string{"pod"},
								},
								{
									DimensionName: "bcs_cluster_id",
									Operator:      structured.ConditionNotEqual,
									Value:         []string{"cls-1"},
								},
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "instance",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ins-1"},
								},
							},
							ConditionList: []string{"and", "and", "and", "and"},
						},
						ReferenceName: "b",
					},
				},
				MetricMerge: `a / on(bcs_cluster_id, namespace, pod_name) group_left() b`,
			},
		},
		"promql with match and verify - 1": {
			queryStruct: false,
			promql: &structured.QueryPromQL{
				PromQL:             `sum by (pod_name, bcs_cluster_id, namespace,instance) (rate(container_cpu_usage_seconds_total{namespace="ns-1"}[2m])) / on(bcs_cluster_id, namespace, pod_name) group_left() sum (sum_over_time(kube_pod_container_resource_limits_cpu_cores{namespace="ns-1"}[1m])) by (pod_name, bcs_cluster_id,namespace)`,
				Match:              `{pod_name="pod", bcs_cluster_id!="cls-1", namespace="ns-1", instance="ins-1"}`,
				IsVerifyDimensions: true,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						FieldName:  "container_cpu_usage_seconds_total",
						TimeAggregation: structured.TimeAggregation{
							Function:  "rate",
							Window:    "2m0s",
							NodeIndex: 2,
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method:     "sum",
								Dimensions: []string{"pod_name", "bcs_cluster_id", "namespace", "instance"},
							},
						},
						ReferenceName: "a",
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "pod_name",
									Operator:      structured.ConditionEqual,
									Value:         []string{"pod"},
								},
								{
									DimensionName: "bcs_cluster_id",
									Operator:      structured.ConditionNotEqual,
									Value:         []string{"cls-1"},
								},
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "instance",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ins-1"},
								},
							},
							ConditionList: []string{"and", "and", "and", "and"},
						},
					},
					{
						DataSource: "bkmonitor",
						FieldName:  "kube_pod_container_resource_limits_cpu_cores",
						TimeAggregation: structured.TimeAggregation{
							Function:  "sum_over_time",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method:     "sum",
								Dimensions: []string{"pod_name", "bcs_cluster_id", "namespace"},
							},
						},
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "pod_name",
									Operator:      structured.ConditionEqual,
									Value:         []string{"pod"},
								},
								{
									DimensionName: "bcs_cluster_id",
									Operator:      structured.ConditionNotEqual,
									Value:         []string{"cls-1"},
								},
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
							},
							ConditionList: []string{"and", "and", "and"},
						},
						ReferenceName: "b",
					},
				},
				MetricMerge: `a / on(bcs_cluster_id, namespace, pod_name) group_left() b`,
			},
		},
		"promql with match and verify - 2": {
			queryStruct: false,
			promql: &structured.QueryPromQL{
				PromQL:             `sum by (pod_name) (rate(container_cpu_usage_seconds_total{namespace="ns-1"}[2m])) / on(bcs_cluster_id, namespace, pod_name) group_left() kube_pod_container_resource_limits_cpu_cores{namespace="ns-1"} or sum by (bcs_cluster_id, namespace, pod_name, instance) (rate(container_cpu_usage_seconds_total{namespace="ns-1"}[1m]))`,
				Match:              `{pod_name="pod", bcs_cluster_id!="cls-1", namespace="ns-1", instance="ins-1"}`,
				IsVerifyDimensions: true,
			},
			query: &structured.QueryTs{
				QueryList: []*structured.Query{
					{
						DataSource: "bkmonitor",
						FieldName:  "container_cpu_usage_seconds_total",
						TimeAggregation: structured.TimeAggregation{
							Function:  "rate",
							Window:    "2m0s",
							NodeIndex: 2,
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method:     "sum",
								Dimensions: []string{"pod_name"},
							},
						},
						ReferenceName: "a",
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "pod_name",
									Operator:      structured.ConditionEqual,
									Value:         []string{"pod"},
								},
							},
							ConditionList: []string{"and"},
						},
					},
					{
						DataSource: "bkmonitor",
						FieldName:  "kube_pod_container_resource_limits_cpu_cores",
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
							},
						},
						ReferenceName: "b",
					},
					{
						DataSource: "bkmonitor",
						FieldName:  "container_cpu_usage_seconds_total",
						TimeAggregation: structured.TimeAggregation{
							Function:  "rate",
							Window:    "1m0s",
							NodeIndex: 2,
						},
						AggregateMethodList: []structured.AggregateMethod{
							{
								Method:     "sum",
								Dimensions: []string{"bcs_cluster_id", "namespace", "pod_name", "instance"},
							},
						},
						Conditions: structured.Conditions{
							FieldList: []structured.ConditionField{
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "pod_name",
									Operator:      structured.ConditionEqual,
									Value:         []string{"pod"},
								},
								{
									DimensionName: "bcs_cluster_id",
									Operator:      structured.ConditionNotEqual,
									Value:         []string{"cls-1"},
								},
								{
									DimensionName: "namespace",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ns-1"},
								},
								{
									DimensionName: "instance",
									Operator:      structured.ConditionEqual,
									Value:         []string{"ins-1"},
								},
							},
							ConditionList: []string{"and", "and", "and", "and"},
						},
						ReferenceName: "c",
					},
				},
				MetricMerge: `a / on(bcs_cluster_id, namespace, pod_name) group_left() b or c`,
			},
		},
	}

	for n, c := range testCase {
		t.Run(n, func(t *testing.T) {
			ctx, _ = context.WithCancel(ctx)
			if c.queryStruct {
				promql, err := structToPromQL(ctx, c.query)
				if c.err != nil {
					assert.Equal(t, c.err, err)
				} else {
					assert.Nil(t, err)
					if err == nil {
						equalWithJson(t, c.promql, promql)
					}
				}
			} else {
				query, err := promQLToStruct(ctx, c.promql)
				if c.err != nil {
					assert.Equal(t, c.err, err)
				} else {
					assert.Nil(t, err)
					if err == nil {
						equalWithJson(t, c.query, query)
					}
				}
			}
		})
	}
}

func equalWithJson(t *testing.T, a, b interface{}) {
	a1, a1Err := json.Marshal(a)
	assert.Nil(t, a1Err)

	b1, b1Err := json.Marshal(b)
	assert.Nil(t, b1Err)
	if a1Err == nil && b1Err == nil {
		assert.Equal(t, string(a1), string(b1))
	}
}

func TestQueryTs_ToQueryReference(t *testing.T) {
	ctx := metadata.InitHashID(context.Background())

	mock.Init()
	influxdb.MockSpaceRouter(ctx)

	metadata.SetUser(ctx, "", influxdb.SpaceUid, "")
	jsonData := `{"query_list":[{"data_source":"","table_id":"","field_name":"container_cpu_usage_seconds_total","is_regexp":false,"field_list":null,"function":[{"method":"sum","without":false,"dimensions":["namespace"],"position":0,"args_list":null,"vargs_list":null}],"time_aggregation":{"function":"rate","window":"5m","node_index":0,"position":0,"vargs_list":[],"is_sub_query":false,"step":""},"reference_name":"a","dimensions":["namespace"],"limit":0,"timestamp":null,"start_or_end":0,"vector_offset":0,"offset":"","offset_forward":false,"slimit":0,"soffset":0,"conditions":{"field_list":[{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["ieg-blueking-gse-data-common"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["ieg-blueking-gse"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["flux-cd-deploy"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["kube-system"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bkmonitor-operator-bkop"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bkmonitor-operator"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["ieg-blueking-gse-data-jk"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["kyverno"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["ieg-bscp-prod"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["ieg-bkce-bcs-k8s-40980"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["ieg-costops-grey"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["ieg-bscp-test"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bcs-system"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bkop-system"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bk-system"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bcs-k8s-25186"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bcs-k8s-25451"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bcs-k8s-25326"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bcs-k8s-25182"],"op":"contains"},{"field_name":"job","value":["kubelet"],"op":"contains"},{"field_name":"image","value":[""],"op":"ncontains"},{"field_name":"container_name","value":["POD"],"op":"ncontains"},{"field_name":"bcs_cluster_id","value":["BCS-K8S-00000"],"op":"contains"},{"field_name":"namespace","value":["bcs-k8s-25037"],"op":"contains"}],"condition_list":["and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and","or","and","and","and","and"]},"keep_columns":["_time","a","namespace"],"step":""}],"metric_merge":"a","result_columns":null,"start_time":"1702266900","end_time":"1702871700","step":"150s","down_sample_range":"5m","timezone":"Asia/Shanghai","look_back_delta":"","instant":false}`
	var query *structured.QueryTs
	err := json.Unmarshal([]byte(jsonData), &query)
	assert.Nil(t, err)

	queryReference, err := query.ToQueryReference(ctx)
	assert.Nil(t, err)

	vmExpand := queryReference.ToVmExpand(ctx)
	expectData := `job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="ieg-blueking-gse-data-common", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="ieg-blueking-gse", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="flux-cd-deploy", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="kube-system", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bkmonitor-operator-bkop", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bkmonitor-operator", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="ieg-blueking-gse-data-jk", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="kyverno", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="ieg-bscp-prod", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="ieg-bkce-bcs-k8s-40980", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="ieg-costops-grey", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="ieg-bscp-test", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bcs-system", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bkop-system", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bk-system", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bcs-k8s-25186", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bcs-k8s-25451", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bcs-k8s-25326", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bcs-k8s-25182", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value" or job="kubelet", image!="", container_name!="POD", bcs_cluster_id="BCS-K8S-00000", namespace="bcs-k8s-25037", result_table_id="2_bcs_prom_computation_result_table", __name__="container_cpu_usage_seconds_total_value"`
	assert.Equal(t, expectData, vmExpand.MetricFilterCondition["a"])
	assert.Nil(t, err)
	assert.True(t, metadata.GetQueryParams(ctx).IsDirectQuery())
}

func TestQueryTsClusterMetrics(t *testing.T) {
	ctx := metadata.InitHashID(context.Background())

	mock.Init()
	promql.MockEngine()
	influxdb.MockSpaceRouter(ctx)

	var (
		key string
		err error
	)

	key = fmt.Sprintf("%s:%s", ClusterMetricQueryPrefix, redis.ClusterMetricMetaKey)
	_, err = redisUtil.HSet(ctx, key, "influxdb_shard_write_points_ok", `{"metric_name":"influxdb_shard_write_points_ok","tags":["bkm_cluster","database","engine","hostname","id","index_type","path","retention_policy","wal_path"]}`)
	if err != nil {
		return
	}

	key = fmt.Sprintf("%s:%s", ClusterMetricQueryPrefix, redis.ClusterMetricKey)
	_, err = redisUtil.HSet(ctx, key, "influxdb_shard_write_points_ok|bkm_cluster=default", `[{"bkm_cluster":"default","database":"_internal","engine":"tsm1","hostname":"influxdb-0","id":"43","index_type":"inmem","path":"/var/lib/influxdb/data/_internal/monitor/43","retention_policy":"monitor","wal_path":"/var/lib/influxdb/wal/_internal/monitor/43","time":1700903220,"value":1498687},{"bkm_cluster":"default","database":"_internal","engine":"tsm1","hostname":"influxdb-0","id":"44","index_type":"inmem","path":"/var/lib/influxdb/data/_internal/monitor/44","retention_policy":"monitor","wal_path":"/var/lib/influxdb/wal/_internal/monitor/44","time":1700903340,"value":1499039.5}]`)
	if err != nil {
		return
	}

	testCases := map[string]struct {
		query  string
		result string
	}{
		"rangeCase": {
			query: `
                {
                    "space_uid": "influxdb",
                    "query_list": [
                        {
                            "data_source": "",
                            "table_id": "",
                            "field_name": "influxdb_shard_write_points_ok",
                            "field_list": null,
                            "function": [
                                {
                                    "method": "sum",
                                    "without": false,
                                    "dimensions": ["bkm_cluster"],
                                    "position": 0,
                                    "args_list": null,
                                    "vargs_list": null
                                }
                            ],
                            "time_aggregation": {
                                "function": "avg_over_time",
                                "window": "60s",
                                "position": 0,
                                "vargs_list": null
                            },
                            "reference_name": "a",
                            "dimensions": [],
                            "limit": 0,
                            "timestamp": null,
                            "start_or_end": 0,
                            "vector_offset": 0,
                            "offset": "",
                            "offset_forward": false,
                            "slimit": 0,
                            "soffset": 0,
                            "conditions": {
                                "field_list": [{"field_name": "bkm_cluster", "value": ["default"], "op": "eq"}],
                                "condition_list": []
                            },
                            "keep_columns": [
                                "_time",
                                "a"
                            ]
                        }
                    ],
                    "metric_merge": "a",
                    "result_columns": null,
                    "start_time": "1700901370",
                    "end_time": "1700905370",
                    "step": "60s",
					"instant": false
                }
			`,
			result: `{"series":[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":["bkm_cluster"],"group_values":["default"],"values":[[1700903220000,1498687],[1700903340000,1499039.5]]}]}`,
		},
		"instanceCase": {
			query: `
                {
                    "space_uid": "influxdb",
                    "query_list": [
                        {
                            "data_source": "",
                            "table_id": "",
                            "field_name": "influxdb_shard_write_points_ok",
                            "field_list": null,
                            "reference_name": "a",
                            "dimensions": [],
                            "limit": 0,
                            "timestamp": null,
                            "start_or_end": 0,
                            "vector_offset": 0,
                            "offset": "",
                            "offset_forward": false,
                            "slimit": 0,
                            "soffset": 0,
                            "conditions": {
                                "field_list": [
									{"field_name": "bkm_cluster", "value": ["default"], "op": "eq"},
									{"field_name": "id", "value": ["43"], "op": "eq"},
									{"field_name": "database", "value": ["_internal"], "op": "eq"},
									{"field_name": "bkm_cluster", "value": ["default"], "op": "eq"},
									{"field_name": "id", "value": ["44"], "op": "eq"}
								],
                                "condition_list": ["and", "or", "and", "and"]
                            },
                            "keep_columns": [
                                "_time",
                                "a"
                            ]
                        }
                    ],
                    "metric_merge": "a",
                    "result_columns": null,
                    "end_time": "1700905370",
					"instant": true
                }
			`,
			result: `{"series":[{"name":"_result0","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":["bkm_cluster","database","engine","hostname","id","index_type","path","retention_policy","wal_path"],"group_values":["default","_internal","tsm1","influxdb-0","43","inmem","/var/lib/influxdb/data/_internal/monitor/43","monitor","/var/lib/influxdb/wal/_internal/monitor/43"],"values":[[1700903220000,1498687]]},{"name":"_result1","metric_name":"","columns":["_time","_value"],"types":["float","float"],"group_keys":["bkm_cluster","database","engine","hostname","id","index_type","path","retention_policy","wal_path"],"group_values":["default","_internal","tsm1","influxdb-0","44","inmem","/var/lib/influxdb/data/_internal/monitor/44","monitor","/var/lib/influxdb/wal/_internal/monitor/44"],"values":[[1700903340000,1499039.5]]}]}`,
		},
	}
	for name, c := range testCases {
		t.Run(name, func(t *testing.T) {
			body := []byte(c.query)
			query := &structured.QueryTs{}
			err := json.Unmarshal(body, query)
			assert.Nil(t, err)

			res, err := QueryTsClusterMetrics(ctx, query)
			t.Logf("QueryTsClusterMetrics error: %+v", err)
			assert.Nil(t, err)
			out, err := json.Marshal(res)
			actual := string(out)
			assert.Nil(t, err)
			fmt.Printf("ActualResult: %v\n", actual)
			assert.JSONEq(t, c.result, actual)
		})
	}
}

func TestQueryTsToInstanceAndStmt(t *testing.T) {

	ctx := metadata.InitHashID(context.Background())

	spaceUid := influxdb.SpaceUid

	mock.Init()
	promql.MockEngine()
	influxdb.MockSpaceRouter(ctx)

	testCases := map[string]struct {
		query        *structured.QueryTs
		promql       string
		stmt         string
		instanceType string
	}{
		"test_matcher_with_vm": {
			promql:       `datasource:result_table:vm:container_cpu_usage_seconds_total{}`,
			stmt:         `a`,
			instanceType: consul.VictoriaMetricsStorageType,
		},
		"test_matcher_with_influxdb": {
			promql:       `datasource:result_table:influxdb:cpu_summary{}`,
			stmt:         `a`,
			instanceType: consul.PrometheusStorageType,
		},
		"test_group_with_vm": {
			promql:       `sum(count_over_time(datasource:result_table:vm:container_cpu_usage_seconds_total{}[1m]))`,
			stmt:         `sum(count_over_time(a[1m] offset -59s999ms))`,
			instanceType: consul.VictoriaMetricsStorageType,
		},
		"test_group_with_influxdb": {
			promql:       `sum(count_over_time(datasource:result_table:influxdb:cpu_summary{}[1m]))`,
			stmt:         `sum(last_over_time(a[1m] offset -59s999ms))`,
			instanceType: consul.PrometheusStorageType,
		},
	}

	err := featureFlag.MockFeatureFlag(ctx, `{
	  	"must-vm-query": {
	  		"variations": {
	  			"true": true,
	  			"false": false
	  		},
	  		"targeting": [{
	  			"query": "tableID in [\"result_table.vm\"]",
	  			"percentage": {
	  				"true": 100,
	  				"false":0 
	  			}
	  		}],
	  		"defaultRule": {
	  			"variation": "false"
	  		}
	  	}
	  }`)
	if err != nil {
		log.Fatalf(ctx, err.Error())
	}

	for name, c := range testCases {
		t.Run(name, func(t *testing.T) {
			if c.promql != "" {
				query, err := promQLToStruct(ctx, &structured.QueryPromQL{PromQL: c.promql})
				if err != nil {
					log.Fatalf(ctx, err.Error())
				}
				c.query = query
			}
			c.query.SpaceUid = spaceUid

			instance, stmt, err := queryTsToInstanceAndStmt(metadata.InitHashID(ctx), c.query)
			if err != nil {
				log.Fatalf(ctx, err.Error())
			}

			assert.Equal(t, c.stmt, stmt)
			if instance != nil {
				assert.Equal(t, c.instanceType, instance.InstanceType())
			}
		})
	}
}
