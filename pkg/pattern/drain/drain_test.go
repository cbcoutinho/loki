package drain

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/loki/v3/pkg/logql/log/pattern"
)

func TestDrain_TrainExtractsPatterns(t *testing.T) {
	printUpdatedPatterns := false
	t.Parallel()
	tests := []struct {
		name      string
		tokenizer PatternTokenizer
		inputFile string
		patterns  []string
	}{
		{
			// High variation leads to many patterns including some that are too generic (many tokens matched) and some that are too specific (too few matchers)
			name:      "Generate patterns on high variation logfmt logs",
			tokenizer: NewExpLogfmtTokenizer(),
			inputFile: "testdata/agent-logfmt.txt",
			patterns: []string{
				`ts=<_> caller=filetarget.go:192 level=info component=logs logs_config=default msg="filetarget: watcher closed, tailer stopped, positions saved" path=<_>`,
				`ts=<_> caller=filetarget.go:313 level=info component=logs logs_config=default msg="watching new directory" directory=<_>`,
				`ts=<_> caller=filetarget.go:326 level=info component=logs logs_config=default msg="removing directory from watcher" directory=<_>`,
				`ts=<_> caller=filetargetmanager.go:181 level=info component=logs logs_config=default msg="received file watcher event" name=<_> op=CREATE`,
				`ts=<_> caller=filetargetmanager.go:361 level=info component=logs logs_config=default msg="Adding target" key="/var/log/pods/*19a1cce8-5f04-46e0-a124-292b0dd9b34<_> batch_kubernetes_io_job_name="testcoordinator-job-<_> container="testcoordinator", controller_uid="25ec5edf-f78e-468b-b6f3-3b9685f0cc<_> job="k6-cloud/testcoordinator", job_name="testcoordinator-job-2665838", name="testcoordinator", namespace="k6-cloud", pod="testcoordinator-job-2665838-9g8ds"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Adding target" key="/var/log/pods/*1954fd4ff7221e619e2d202bfb2c4ab9/ku<_> container="kube-proxy", job="kube-system/kube-proxy", namespace="kube-system", pod="kube-proxy-gke-ops-us-east-0-main-n2s32-1-1dd<_> tier="node"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Adding target" key="/var/log/pods/*6e95d7c0-d863-461c-a6d1-68653e438e3<_> container="kube-proxy", job="kube-system/gke-ops-us-east-0-main-n2s32-1-1d<_> namespace="kube-system", pod="kube-proxy-gke-ops-us-east-0-main-n2s32-1-1dd<_> tier="node"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Adding target" key="/var/log/pods/*b92ee988-5c26-4c64-bba3-ff6a0172375<_> conprof="true", <_> instanceId="i1111", job="hosted-grafana/grafana", name="grafana", namespace="hosted-grafana", org="orgnamehere", plan="free", pod="orgnamehere-grafana-7c65678f86-9zhlb", pod_template_hash="7c65678f86", resource_version="143638246", slug="orgnamehere", stackId="866772"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Removing target" key="/var/log/pods/*043372d7-9411-443f-ba7f-80988f77d8b<_> conprof="true", <_> instanceId="i3333", job="hosted-grafana/grafana", name="grafana", namespace="hosted-grafana", org="org4", plan="free", pod="org4-grafana-b4f87fcc5-fflgn", pod_template_hash="b4f87fcc5", resource_version="167289888", slug="org4", stackId="333333"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Removing target" key="/var/log/pods/*0ecafb81-c168-4bc4-99e3-e8b2315a09b<_> conprof="true", <_> instanceId="i5555", job="hosted-grafana/grafana", name="grafana", namespace="hosted-grafana", org="org6", plan="free", pod="org6-grafana-fdfdc64bb-srrx6", pod_template_hash="fdfdc64bb", resource_version="95745089", slug="org6", stackId="666666"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Removing target" key="/var/log/pods/*35649bfd-52ff-4281-9294-5f65fd5a89f<_> job="grafana-com/marketplaces-api", name="marketplaces-api", namespace="grafana-com", pod="marketplaces-api-f67ff7567-gqrvb", pod_template_hash="f67ff7567"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Removing target" key="/var/log/pods/*37ae8d4e-1a76-40f2-be88-2251a3528a0<_> conprof="true", <_> instanceId="i2222", job="hosted-grafana/grafana", name="grafana", namespace="hosted-grafana", org="someorg", plan="free", pod="someorg-grafana-666bd48cf9-7zrtv", pod_template_hash="666bd48cf9", resource_version="167212086", slug="someorg", stackId="444444"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Removing target" key="/var/log/pods/*404e8595-2e9f4fcf-9495-925f6d245e20<_> conprof="true", <_> instanceId="222222", job="hosted-grafana/grafana", name="grafana", namespace="hosted-grafana", org="org3", plan="free", pod="org3-grafana-7fd6786f4b-242cb", pod_template_hash="7fd6786f4b", resource_version="167282051", slug="org3", stackId="1111111"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Removing target" key="/var/log/pods/*64872ae4-62eb-4757-b148-72bab4a9e88<_> conprof="true", <_> instanceId="i6666", job="hosted-grafana/grafana", name="grafana", namespace="hosted-grafana", org="org7", plan="free", pod="org7-grafana-647dc5b44f-pmz8j", pod_template_hash="647dc5b44f", resource_version="167297262", slug="org7", stackId="777777"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Removing target" key="/var/log/pods/*9afcdf84-163b402e-bdd2-cfb711593385<_> conprof="true", <_> instanceId="i4444", job="hosted-grafana/grafana", name="grafana", namespace="hosted-grafana", org="org5", plan="free", pod="org5-grafana-ddbf649cc-zgtf6", pod_template_hash="ddbf649cc", resource_version="95783554", slug="org5", stackId="555555"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Removing target" key="/var/log/pods/*c3c249a2-c8ff-40f4-a66d-9d746b39110<_> conprof="true", <_> instanceId="i7777", job="hosted-grafana/grafana", name="grafana", namespace="hosted-grafana", org="org8", plan="free", pod="org8-grafana-6c66686686-mqtcr", pod_template_hash="6c66686686", resource_version="95723800", slug="org8", stackId="888888"}"`,
				`ts=<_> caller=filetargetmanager.go:<_> level=info component=logs logs_config=default msg="Removing target" key="<_> conprof="true", <_> <_> job="hosted-grafana/grafana", name="grafana", namespace="hosted-grafana", <_> plan="free", <_> <_> <_> <_> <_>"`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Re-opening moved/deleted file /var/log/pods/grafana-ruler_grafana-ruler-7cb758bf<_> ..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Re-opening moved/deleted file /var/log/pods/hosted-grafana_.something-grafana-7c<_> ..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Re-opening moved/deleted file /var/log/pods/hosted-grafana_.something-grafana-d8<_> ..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Re-opening moved/deleted file /var/log/pods/hosted-grafana_.something-grafana-ga<_> ..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Re-opening moved/deleted file /var/log/pods/insight-logs_promtail-insight-logs-p<_> ..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Re-opening moved/deleted file /var/log/pods/mimir-dedicated-48_ingester-zone-b-7<_> ..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Re-opening moved/deleted file /var/log/pods/pyroscope-ebpf_profiler-w87sj_04264a<_> ..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Re-opening moved/deleted file <_> ..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Seeked /var/log/pods/hosted-grafana_.something-grafana-5f<_> - &{Offset:0 Whence:0}"`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Seeked /var/log/pods/hosted-grafana_.something-grafana-6b<_> - &{Offset:0 Whence:0}"`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Seeked /var/log/pods/hosted-grafana_.something-grafana-6f<_> - &{Offset:0 Whence:0}"`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Seeked /var/log/pods/hosted-grafana_.something-grafana-<_><_> - &{Offset:0 Whence:0}"`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Seeked /var/log/pods/hosted-grafana_.something-grafana-bc<_> - &{Offset:0 Whence:0}"`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Seeked /var/log/pods/hosted-grafana_.something-grafana-ff<_> - &{Offset:0 Whence:0}"`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Seeked <_> - &{Offset:0 Whence:0}"`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Successfully reopened <_>"`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Waiting for /var/log/pods/hosted-grafana_.something-grafana-<_><_> to appear..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Waiting for /var/log/pods/hosted-grafana_.something-grafana-b4<_> to appear..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Waiting for /var/log/pods/hosted-grafana_.something-grafana-b9<_> to appear..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Waiting for /var/log/pods/hosted-grafana_.something-grafana-d8<_> to appear..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Waiting for /var/log/pods/kube-system_calico-node-7zxvh_a1ad89<_> to appear..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Waiting for /var/log/pods/kube-system_calico-node-jwmqj_894554<_> to appear..."`,
				`ts=<_> caller=log.go:168 component=logs logs_config=default level=info msg="Waiting for <_> to appear..."`,
				`ts=<_> caller=logfmt.go:139 level=error component=logs logs_config=default component=file_pipeline component=stage type=logfmt msg="failed to decode logfmt" err="bufio.Scanner: token too long"`,
				`ts=<_> caller=logfmt.go:139 level=error component=logs logs_config=default component=file_pipeline component=stage type=logfmt msg="failed to decode logfmt" err="logfmt syntax error at pos <_> on line 1: unexpected '"'"`,
				`ts=<_> caller=tailer.go:164 level=info component=logs logs_config=default component=tailer msg="tail routine: tail channel closed, stopping tailer" path=<_> reason=null`,
				`ts=<_> caller=tailer.go:207 level=info component=logs logs_config=default component=tailer msg="skipping update of position for a file which does not currently exist" path=<_>`,
				`ts=<_> caller=tailer.go:<_> level=info component=logs logs_config=default component=tailer msg="position timer: exited" path=<_>`,
				`ts=<_> caller=tailer.go:<_> level=info component=logs logs_config=default component=tailer msg="stopped tailing file" path=<_>`,
				`ts=<_> caller=tailer.go:<_> level=info component=logs logs_config=default component=tailer msg="tail routine: <_>" path=<_>`,
				`ts=<_> level=info msg="finished node evaluation" controller_id=module.http.cloudwatch_pipelines node_id=prometheus.scrape.stack_378175_cloudwatch_notags duration=<_>`,
				`ts=<_> level=info msg="finished node evaluation" controller_id=module.http.cloudwatch_pipelines node_id=prometheus.scrape.stack_<_>_cloudwatch_tags duration=<_>`,
			},
		},
		{
			// Lower variation leads to fewer patterns including some with limited value (single lines, no matchers)
			name:      "Generate patterns on low variation logfmt logs",
			tokenizer: NewExpLogfmtTokenizer(),
			inputFile: "testdata/ingester-logfmt.txt",
			patterns: []string{
				`ts=<_> caller=head.go:216 level=debug tenant=987678 msg="profile is empty after delta computation" metricName=memory`,
				`ts=<_> caller=http.go:194 level=debug traceID=1b48f5156a61ca69 msg="GET /debug/pprof/delta_mutex (200) <_>"`,
				`ts=<_> caller=http.go:194 level=debug traceID=<_> orgID=<_> msg="POST /ingester.v1.IngesterService/Push (200) <_>"`,
			},
		},
		{
			// Lower variation logs in json leads to a high number of patterns with very few matchers
			name:      "Generate patterns on json formatted logs",
			tokenizer: &AdaptiveTokenizer{},
			inputFile: "testdata/drone-json.txt",
			patterns: []string{
				`{"duration":<_>,"level":"debug","method":"GET","msg":"request completed","referer":"","remote":"<_>:52702","request":"/metrics","status":200,"time":"<_>","user-agent":"GrafanaAgent/v0.40.3 (flow; linux; helm)"}`,
				`{"id":"15eSzaEG0enf86Kl","level":"debug","max-pool":4,"min-pool":0,"msg":"check capacity","pending-builds":0,"running-builds":0,"server-buffer":0,"server-capacity":0,"server-count":0,"time":"<_>"}`,
				`{"id":"15eSzaEG0enf86Kl","level":"debug","msg":"calculate server capacity","time":"<_>"}`,
				`{"id":"15eSzaEG0enf86Kl","level":"debug","msg":"calculate unfinished jobs","time":"<_>"}`,
				`{"id":"15eSzaEG0enf86Kl","level":"debug","msg":"check capacity complete","time":"<_>"}`,
				`{"id":"15eSzaEG0enf86Kl","level":"debug","msg":"no capacity changes required","time":"<_>"}`,
				`{"id":"9eA72xOtx8kzMhXn","level":"debug","max-pool":4,"min-pool":0,"msg":"check capacity","pending-builds":0,"running-builds":0,"server-buffer":0,"server-capacity":0,"server-count":0,"time":"<_>"}`,
				`{"id":"9eA72xOtx8kzMhXn","level":"debug","msg":"calculate server capacity","time":"<_>"}`,
				`{"id":"9eA72xOtx8kzMhXn","level":"debug","msg":"calculate unfinished jobs","time":"<_>"}`,
				`{"id":"9eA72xOtx8kzMhXn","level":"debug","msg":"check capacity complete","time":"<_>"}`,
				`{"id":"9eA72xOtx8kzMhXn","level":"debug","msg":"no capacity changes required","time":"<_>"}`,
				`{"id":"JO1OT5ADoNA8NYqr","level":"debug","max-pool":4,"min-pool":0,"msg":"check capacity","pending-builds":0,"running-builds":0,"server-buffer":0,"server-capacity":0,"server-count":0,"time":"<_>"}`,
				`{"id":"JO1OT5ADoNA8NYqr","level":"debug","msg":"calculate server capacity","time":"<_>"}`,
				`{"id":"JO1OT5ADoNA8NYqr","level":"debug","msg":"calculate unfinished jobs","time":"<_>"}`,
				`{"id":"JO1OT5ADoNA8NYqr","level":"debug","msg":"check capacity complete","time":"<_>"}`,
				`{"id":"JO1OT5ADoNA8NYqr","level":"debug","msg":"no capacity changes required","time":"<_>"}`,
				`{"id":"T0I8Dsnw3uSi3Gal","level":"debug","max-pool":4,"min-pool":0,"msg":"check capacity","pending-builds":0,"running-builds":0,"server-buffer":0,"server-capacity":0,"server-count":0,"time":"<_>"}`,
				`{"id":"T0I8Dsnw3uSi3Gal","level":"debug","msg":"calculate server capacity","time":"<_>"}`,
				`{"id":"T0I8Dsnw3uSi3Gal","level":"debug","msg":"calculate unfinished jobs","time":"<_>"}`,
				`{"id":"T0I8Dsnw3uSi3Gal","level":"debug","msg":"check capacity complete","time":"<_>"}`,
				`{"id":"T0I8Dsnw3uSi3Gal","level":"debug","msg":"no capacity changes required","time":"<_>"}`,
				`{"id":"m6SpYHzdXrDAFqDR","level":"debug","max-pool":4,"min-pool":0,"msg":"check capacity","pending-builds":0,"running-builds":0,"server-buffer":0,"server-capacity":0,"server-count":0,"time":"<_>"}`,
				`{"id":"m6SpYHzdXrDAFqDR","level":"debug","msg":"calculate server capacity","time":"<_>"}`,
				`{"id":"m6SpYHzdXrDAFqDR","level":"debug","msg":"calculate unfinished jobs","time":"<_>"}`,
				`{"id":"m6SpYHzdXrDAFqDR","level":"debug","msg":"check capacity complete","time":"<_>"}`,
				`{"id":"m6SpYHzdXrDAFqDR","level":"debug","msg":"no capacity changes required","time":"<_>"}`,
				`{"id":"pet7QVfO1yE8fk56","level":"debug","max-pool":4,"min-pool":0,"msg":"check capacity","pending-builds":0,"running-builds":0,"server-buffer":0,"server-capacity":0,"server-count":0,"time":"<_>"}`,
				`{"id":"pet7QVfO1yE8fk56","level":"debug","msg":"calculate server capacity","time":"<_>"}`,
				`{"id":"pet7QVfO1yE8fk56","level":"debug","msg":"calculate unfinished jobs","time":"<_>"}`,
				`{"id":"pet7QVfO1yE8fk56","level":"debug","msg":"check capacity complete","time":"<_>"}`,
				`{"id":"pet7QVfO1yE8fk56","level":"debug","msg":"no capacity changes required","time":"<_>"}`,
				`{"id":"q62wCcIkEOueqFKF","level":"debug","max-pool":4,"min-pool":0,"msg":"check capacity","pending-builds":0,"running-builds":0,"server-buffer":0,"server-capacity":0,"server-count":0,"time":"<_>"}`,
				`{"id":"q62wCcIkEOueqFKF","level":"debug","msg":"calculate server capacity","time":"<_>"}`,
				`{"id":"q62wCcIkEOueqFKF","level":"debug","msg":"calculate unfinished jobs","time":"<_>"}`,
				`{"id":"q62wCcIkEOueqFKF","level":"debug","msg":"check capacity complete","time":"<_>"}`,
				`{"id":"q62wCcIkEOueqFKF","level":"debug","msg":"no capacity changes required","time":"<_>"}`,
				`{"id":<_>,"level":"debug","max-pool":4,"min-pool":0,"msg":"check capacity","pending-builds":0,"running-builds":0,"server-buffer":0,"server-capacity":0,"server-count":0,"time":"<_>"}`,
				`{"id":<_>,"level":"debug","msg":"calculate server capacity","time":"<_>"}`,
				`{"id":<_>,"level":"debug","msg":"calculate unfinished jobs","time":"<_>"}`,
				`{"id":<_>,"level":"debug","msg":"check capacity complete","time":"<_>"}`,
				`{"id":<_>,"level":"debug","msg":"no capacity changes required","time":"<_>"}`,
			},
		},
		{
			name:      "Patterns for distributor logs",
			tokenizer: NewExpLogfmtTokenizer(),
			inputFile: "testdata/distributor-logfmt.txt",
			patterns: []string{
				`ts=<_> caller=http.go:194 level=debug traceID=<_> orgID=<_> msg="POST /ingest?aggregationType=&from=1714652227232613927&<_> (200) <_>"`,
				`ts=<_> caller=http.go:194 level=debug traceID=<_> orgID=<_> msg="POST /ingest?aggregationType=average&from=1714652227232<_> (200) <_>"`,
				`ts=<_> caller=http.go:194 level=debug traceID=<_> orgID=<_> msg="POST /ingest?aggregationType=sum&from=17146522271076410<_> (200) <_>"`,
				`ts=<_> caller=http.go:194 level=debug traceID=<_> orgID=<_> msg="POST /push.v1.PusherService/Push (200) <_>"`,
				`ts=<_> caller=http.go:194 level=debug traceID=<_> orgID=<_> msg="POST /push.v1.PusherService/Push (<_>) <_>"`,
				`ts=<_> caller=http.go:194 level=debug traceID=<_> orgID=<_> msg="POST /pyroscope/ingest?aggregationType=sum&from=1714652<_> (200) <_>"`,
			},
		},
		{
			name:      "Patterns for journald logs",
			tokenizer: &AdaptiveTokenizer{},
			inputFile: "testdata/journald.txt",
			patterns: []string{
				`						exec /bin/hgrun -log.level=debug launch -bundledPluginsManifest /proc/$(pidof plugins-pause)/root/manifest.json -bundledPluginsDir /proc/$(pidof plugins-pause)/root/plugins],WorkingDir:,Ports:[]C<_>-profile-port=6060 -profile-addr=<_>,ValueFrom:nil,},EnvVar{Name:HG_<_>{{536870912 0} {<nil>}  BinarySI},},Requests:ResourceList{cpu: {{26 3} {<nil>} 26m DecimalSI},memory: {{293601280 0} {<nil>}  BinarySI},},Claims:[]ResourceClaim{},},VolumeMount<_>80 },Host:,Scheme:HTTP,HTTPHeaders:[]HTTPHeader{},},T<_>check],},HTTPGet:nil,TCPSocket:nil,GRPC:nil,},Init<_>drain -timeout 1m0s -waitTime 55s],},HTTPGet:nil,TCPSocket:nil,},},TerminationMe<_>start failed in pod <_>ErrImagePull: [rpc error: code = NotFound desc = failed to pull and unpack image "us.gcr.io/hosted-grafana/hosted-grafana-pro:<_><_>failed to resolve reference "us.gcr.io/hosted-grafana/hosted-grafana-pro:<_><_>us.gcr.io/hosted-grafana/hosted-grafana-pro:<_>.<_>not found, failed to pull and unpack image "us.gcr.io/hosted-grafana/hosted-grafana-pro:<_><_>failed to resolve reference "us.gcr.io/hosted-grafana/hosted-grafana-pro:<_><_>unexpected status from HEAD request to https://us.gcr.io/v2/hosted-grafana/hosted-grafana<_>403 Forbidden]`,
				`						ln --force -s /proc/$(pidof hgrun-pause)/root/bin/hgrun /bin/hgrun;`,
				`						while [ "$(pidof plugins-pause)" = "" ]; do sleep 0.5; done;`,
				`	ts=<_> level=error caller=http_client.go:56 app=hgrun hgrun_version=<_>.<_>.<_> msg="request failed" error="Get \"http://<_>:3000/api/health\": dial t<_>method=GET url=http://<_>:3000/api/health`,
				`	ts=<_> level=error caller=http_client.go:56 app=hgrun hgrun_version=<_>.<_>.<_>-59-gf3f63162a msg="request`,
				`	ts=<_> level=error caller=http_client.go:56 app=hgrun hgrun_version=<_>.<_>.<_>-59-gf3f63162a msg="request failed" error="Get \"http://<_>:3000/api/health\": dial t<_>method=GET url=http://<_>:3000/api/health`,
				`	ts=<_> level=error caller=http_client.go:56 app=hgrun hgrun_version=<_>.<_>.<_>-62-g2605e8595 msg="request failed" error="Get \"http://<_>:3000/api/health\": dial t<_>method=GET url=http://<_>:3000/api/health`,
				` >`,
				`<_> INFO ExtHandler ExtHandler Downloading agent manifest`,
				`<_> INFO TelemetryEventsCollector ExtHandler Collected 2 events for extension: Microsoft.Azure.Extensions.CustomScript`,
				`AVC apparmor="DENIED" operation="ptrace" profile="cri-containerd.apparmor.d" pid=<_> comm="pidof" requested_mask="read" denied_mask="read" peer="unconfined"`,
				`E0507 11:59:34.923938    3027 kuberuntime_manager.go:1261] container &Container{Name:mysqld-exporter,Image:prom/mysqld-<_>start failed in pod testcrossplane-exporter-c67cfc58f-vbzl4_crossplane<_>CreateContainerConfigError: secret "testcrossplane-user-exporter" not found`,
				`E0507 11:59:41.375655    4736 kuberuntime_manager.go:1256] container &Container{Name:ruler,Image:grafana/enterprise-met<_>-config.expand-env=true -config.file=/etc/mimir/mimir.yaml -distributor.remote-timeout=10s],WorkingDir:,Ports<_>{{100 3} {<nil>} 100m DecimalSI},memory: {{134217728 0} {<nil>}  BinarySI},},Claims:[]ResourceClaim{},},VolumeMount<_>0 http-metrics},Host:,Scheme:HTTP,HTTPHeaders:[]HTTP<_>start failed in pod gem-mimir-ruler-5f56f7846b-fgxdm_ge-metrics-federa<_>CreateContainerConfigError: secret "ruler-alertmanager-token" not found`,
				`E0507 <_>:<_>:<_>    <_> kuberuntime_manager.go:1256] container &Container{Name:grafana,Image:us.gcr.io/hosted-gra<_>set -e; while [ "$(pidof hgrun-pause)" = "" ]; do sleep 0.5; done;`,
				`E0507 <_>:<_>:<_>    <_> kuberuntime_manager.go:1256] container &Container{Name:pdc,Image:us.gcr.io/hosted-grafana<_>-proxy.socks-server.addr=:10443 -proxy.ssh-server.addr=:2222 -proxy.use-socks-username-for-routing -proxy.api.http-address=:9182 -proxy.check-connpool-address-in-ring -memberlist.join=dns+gossip-ring.pdc.svc.cluster.l<_>-api.http-address=:11443 -distributor.enabled=true -distributor.addr=:10444 -distributor.use-socks-username-for-routing -gateway.enabled=true -gateway.addr=:2244 -log.level=debug -certs.ca-private-key-file=/var/run/secrets/pdc-ce<_>-certs.ca-cert-file=/var/run/secrets/pdc-certs/ca.<_>-certs.ca-pub-file=/var/run/secrets/pdc-certs/ca.p<_>-certs.cluster=local-k8s -shard-size=3 -graceful-shutdown-period=30s -enable-multiple-networks],WorkingDir:,Ports:[]Con<_>{{500 3} {<nil>} 500m DecimalSI},memory: {{134217728 0} {<nil>}  BinarySI},},Requests:ResourceList{cpu: {{250 3} {<nil>} 250m DecimalSI},memory: {{67108864 0} {<nil>}  BinarySI},},Claims:[]ResourceClaim{},},VolumeMount<_>11443 },Host:,Scheme:HTTP,HTTPHeaders:[]HTTPHeader{},},T<_>5],},HTTPGet:nil,TCPSocket:nil,},},TerminationMess<_>start failed in pod <_>ErrImageNeverPull: Container image "us.gcr.io/hosted-grafana/pdc:<_>.<_>.<_>" is not present with pull policy of Never`,
				`E0507 <_>:<_>:<_>    <_> kuberuntime_manager.go:1256] container &Container{Name:ruler,Image:grafana/enterprise-met<_>-config.expand-env=true -config.file=/etc/mimir/mimir.yaml],WorkingDir:,Po<_>{{100 3} {<nil>} 100m DecimalSI},memory: {{134217728 0} {<nil>}  BinarySI},},Claims:[]ResourceClaim{},},VolumeMount<_>0 http-metrics},Host:,Scheme:HTTP,HTTPHeaders:[]HTTP<_>start failed in pod <_>CreateContainerConfigError: secret "ruler-alertmanager-token" not found`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" <_><_><_>`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"agent\" wi<_>pod="jaeger/jaeger-agent-856f67c6d7-6xj9z" podUID="1a240429-7c6f4c4c8c4e-d2579a6e737e"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"agent\" wi<_>pod="jaeger/jaeger-agent-856f67c6d7-tcsmd" podUID="9121c1a3-6d79-4411-be8e-41406c88944a"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"cluster-ag<_>pod="integration/appdynamics-cluster-agent-appdyna<_>podUID="69bc5e6c-0451-443e-af8a-c831871afbb8"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"cortex-gw\<_>pod="faro/cortex-gw-6f7f764f94-rgtw8" podUID="d6bf8bcc-35b9-4c1f-ab69-f857a2328d11"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"cortex-gw\<_>pod="faro/cortex-gw-74f78948ff-9pcl6" podUID="643043e2-707a4a3f-adf3-08beab1d1ea7"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"cortex-gw\<_>pod="faro/cortex-gw-78bc9b5ccc-8hkmp" podUID="44b54226-b4bd-46e0-a3f0-257cb44d9ea8"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"gcom-sync\<_>pod="faro/update-usage-28487080-9sqzn" podUID="2cc85139-2f31-44ae-a308-3dc0df893592"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"gcom-sync\<_>pod="faro/update-usage-28487090-xg5bt" podUID="6e8f7589-7d91-47e6-9128-7ec922779773"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"grafana\" <_><_><_>`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"grafana\" <_>pod="hosted-grafana/benchloadtestingxxl2-grafana-5<_><_>`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"grafana\" <_>pod="hosted-grafana/ephemeral1180076306267marefr-g<_><_>`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"grafana\" <_>pod="hosted-grafana/ephemeral1511182185282svenner-<_><_>`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"grafana\" <_>pod="hosted-grafana/johan6-grafana-85546bbbf5-xbkr<_>podUID="a1ca81cd-1fd3-4f14-b6a5-a129930ba761"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"grafana\" <_>pod="hosted-grafana/johangrafana10-grafana-69c6449<_>podUID="bb953c26-c201-4082-9b56-85ab12c1d0e1"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"grafana\" <_>pod="hosted-grafana/k6testslow2-grafana-7b64f97bd7<_>podUID="9890650a-e338-4648-be7a-bb7f9726aa46"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"grafana\" <_>pod="hosted-grafana/oncalldev-grafana-7b88d9459-fv<_>podUID="fc7753d0-4067-4626-b539-5fd27ded163b"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"ksm\" with<_>pod="integration/new-relic-nri-bundle-nrk8s-ksm-6c<_>podUID="f7cc3cca-2ffb-4fde-a73e-a4ba8b0f6b3c"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"overrides-<_>pod="loki-dev-010/overrides-exporter-98c77fd66-6zj<_>podUID="1ff5bf3e-5856-4f6f-ae04-273f2dee170b"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"pdc\" with<_>pod="pdc/private-datasource-connect-564fb6cfbb-fd2<_>podUID="ac6bc6d0-43a4-4885-9ee4-ba3441b0b527"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"pdc\" with<_>pod="pdc/private-datasource-connect-564fb6cfbb-l8p<_>podUID="57e4a0cb-5e77-47bd-b277-70f4b1512c44"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"prometheus<_>pod="bryan-prometheus/bryan-prometheus-0" podUID="6dadfe71-eb19-4231-a96e-c64bb5499a1e"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"ruler\" wi<_>pod="ge-metrics-federation/gem-mimir-ruler-5f56f78<_>podUID="07c06e21-137b4fdd-b7d3-703f0a567720"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"ruler\" wi<_>pod="ge-metrics-federation/gem-mimir-ruler-8c54cd6<_>podUID="0a159d8c5540-44c2-a592-f43db7a1aae6"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"ruler\" wi<_>pod="ge-metrics-federation/gem-mimir-ruler-bd7cbc8<_>podUID="f39fa140-2a71-4cba-bcb7-b37b2fafa343"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"support-ag<_>pod="support-agent/support-agent-557dff8b77-c6f8b"<_>_>dUID="ede5a224-96fb-45d0-b452-1eb2de73cf19"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"support-ag<_>pod="support-agent/support-agent-557dff8b77-sx6hb"<_>_>dUID="f7b72dbb-4f3a45b1-88c0-62337a3e8d3d"`,
				`E0507 <_>:<_>:<_>    <_> pod_workers.go:1300] "Error syncing pod, skipping" err="unmounted volumes=[custom-grafana-agent], una<_>pod="loki-dev-010/custom-grafana-agent-856948968f6<_>podUID="17b244cc-ecb9-4fbc-beaa-8fa47fafe013"`,
				`E0507 <_>:<_>:<_>    <_> prober.go:104] "Probe errored" err="rpc error: code = NotFound desc = failed to e<_>probeType="Readiness" pod="hosted-grafana/benchloadtestingxxl2-grafana-5<_>podUID="<_>" containerName="grafana"`,
				`E0507 <_>:<_>:<_>    <_> prober.go:104] "Probe errored" err="rpc error: code = NotFound desc = failed to e<_>probeType="Readiness" pod="hosted-grafana/dafdeveuwest2-grafana-7845d969<_>podUID="14ac9939-b36a-40d7-9ca9-a0367aab99d8" containerName="grafana"`,
				`E0507 <_>:<_>:<_>    <_> prober.go:104] "Probe errored" err="rpc error: code = NotFound desc = failed to e<_>probeType="Readiness" pod="hosted-grafana/k6teststeady3-grafana-659d5ff5<_>podUID="85274c17-190e4275-a8f3-6e111cd833bf" containerName="grafana"`,
				`E0507 <_>:<_>:<_>    <_> prober.go:239] "Unable to write all bytes from execInContainer" err="short write" expectedBytes=<_> actualBytes=10240`,
				`E0507 <_>:<_>:<_>    <_> remote_image.go:180] "PullImage from image service failed" err="rpc error: code = NotFound desc = failed to p<_>image="us.gcr.io/hosted-grafana/hosted-grafana-pro<_>`,
				`E0507 <_>:<_>:<_>    <_> remote_image.go:180] "PullImage from image service failed" err="rpc error: code = Unknown desc = failed to pu<_>image="us.gcr.io/hosted-grafana/hosted-grafana-pro<_>`,
				`E0507 <_>:<_>:<_>    <_> remote_runtime.go:432] "ContainerStatus from runtime service failed" err="rpc error: code = NotFound desc = an error oc<_><_>`,
				`E0507 <_>:<_>:<_>    <_> remote_runtime.go:496] "ExecSync cmd from runtime service failed" err="rpc error: code = NotFound desc = failed to e<_>containerID="<_><_>cmd=["/bin/hgrun","check"]`,
				`I0507 11:59:33.422254 1537502 kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-gke-dev-us-central-0-m<_>status="Running"`,
				`I0507 11:59:34.518822    3224 kuberuntime_container.go:745] "Killing container with a grace period" pod="hosted-grafana/hosted-grafana-api-7b6bd9b949-<_>podUID="25cb986c-3d6c4ed0-abf3-ee59ed6175f9" containerName="hgapi" containerID="containerd://c91436db00920ec961b9d5d6<_>gracePeriod=30`,
				`I0507 <_>:<_>:<_>    3224 operation_generator.go:888] UnmountVolume.TearDown succeeded for volume "kubernetes.io/projected/25cb986c-3d6c4ed0-abf3-ee<_>(OuterVolumeSpecName: "kube-api-access-95j2t") pod "25cb986c-3d6c4ed0-abf3-ee59ed6175f9" (UID: "25cb986c-3d6c4ed0-abf3-ee59ed6175f9"). InnerVolumeSpecName "kube-api-access-95j2t". PluginName "kubernetes.io/projected", VolumeGidValue ""`,
				`I0507 <_>:<_>:<_>    3224 operation_generator.go:888] UnmountVolume.TearDown succeeded for volume "kubernetes.io/secret/25cb986c-3d6c4ed0-abf3-ee59e<_>(OuterVolumeSpecName: "gcs-serviceaccount") pod "25cb986c-3d6c4ed0-abf3-ee59ed6175f9" (UID: "25cb986c-3d6c4ed0-abf3-ee59ed6175f9"). InnerVolumeSpecName "gcs-serviceaccount". PluginName "kubernetes.io/secret", VolumeGidValue ""`,
				`I0507 <_>:<_>:<_>    3224 operation_generator.go:888] UnmountVolume.TearDown succeeded for volume "kubernetes.io/secret/25cb986c-3d6c4ed0-abf3-ee59e<_>(OuterVolumeSpecName: "pdc-certs") pod "25cb986c-3d6c4ed0-abf3-ee59ed6175f9" (UID: "25cb986c-3d6c4ed0-abf3-ee59ed6175f9"). InnerVolumeSpecName "pdc-certs". PluginName "kubernetes.io/secret", VolumeGidValue ""`,
				`I0507 <_>:<_>:<_>    3224 reconciler_common.go:<_>] <_>`,
				`I0507 <_>:<_>:<_>    <_> azure_credentials.go:220] image(us.gcr.io/hosted-grafana/hg-plugins) is not from ACR, return empty authentication`,
				`I0507 <_>:<_>:<_>    <_> azure_credentials.go:220] image(us.gcr.io/hosted-grafana/hgrun) is not from ACR, return empty authentication`,
				`I0507 <_>:<_>:<_>    <_> azure_credentials.go:220] image(us.gcr.io/hosted-grafana/hosted-grafana-pro)<_>_> not from ACR, return empty authentication`,
				`I0507 <_>:<_>:<_>    <_> generic.go:334] "Generic (PLEG): container finished" podID="25cb986c-3d6c4ed0-abf3-ee59ed6175f9" containerID="c91436db00920ec961b9d5d6b4859d80a912e<_>exitCode=1`,
				`I0507 <_>:<_>:<_>    <_> generic.go:334] "Generic (PLEG): container finished" podID="85274c17-190e4275-a8f3-6e111cd833bf" containerID="fc7a558bca122d6b5fb9aa81e62a87053c8a6<_>exitCode=1`,
				`I0507 <_>:<_>:<_>    <_> generic.go:334] "Generic (PLEG): container finished" podID="<_>" containerID="<_><_>exitCode=1`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (PLEG): event for pod" pod="hosted-grafana/benchloadtestingxxl2-grafana-5<_>event={"ID":"<_>","Type":"ContainerDied","Data":<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (PLEG): event for pod" pod="hosted-grafana/dafdeveuwest2-grafana-546fbd78<_>event={"ID":"fc6ba4ea-9950-4999-8ad2-bdc9a577fb34","Type":"ContainerStarted","Data":<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (PLEG): event for pod" pod="hosted-grafana/dafdeveuwest2-grafana-7845d969<_>event={"ID":"14ac9939-b36a-40d7-9ca9-a0367aab99d8","Type":"ContainerDied","Data":"eeccb21da13bfae40b1a01984522c7a8f8dcb65dba3cc1cc2<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (PLEG): event for pod" pod="hosted-grafana/dafdeveuwest2-grafana-7845d969<_>event={"ID":"14ac9939-b36a-40d7-9ca9-a0367aab99d8","Type":"ContainerStarted","Data":"eeccb21da13bfae40b1a01984522c7a8f8dcb65dba3cc1cc2<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (PLEG): event for pod" pod="hosted-grafana/hosted-grafana-api-7b6bd9b949-<_>event={"ID":"25cb986c-3d6c4ed0-abf3-ee59ed6175f9","Type":"ContainerDied","Data":<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (PLEG): event for pod" pod="hosted-grafana/k6teststeady3-grafana-659d5ff5<_>event={"ID":"85274c17-190e4275-a8f3-6e111cd833bf","Type":"ContainerDied","Data":"fc7a558bca122d6b5fb9aa81e62a87053c8a6a84945fd7a5f<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (PLEG): event for pod" pod="hosted-grafana/k6teststeady4-grafana-5c4f6cd5<_>event={"ID":"a95be6bc-a7bc-48cb-8935-f7040f91f7f9","Type":"ContainerDied","Data":"c6da2382101cc3ca3a9a6de7b86f62dfd7b344559c7e17cec<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (PLEG): event for pod" pod="hosted-grafana/victor-grafana-7b7bb568cc-grfl<_>event={"ID":"1803645b5526-41b4-bf88-271be4827277","Type":"ContainerStarted","Data":<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (PLEG): event for pod" pod="otel-demo/otel-demo-dev-checkoutservice-6ddf9<_>event={"ID":"f263b787-926e459a95a0-f9ef8e4e9bc2","Type":"ContainerStarted","Data":"95bf586cd79d43120ff44582d4dbd2476de61744411f8515b<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (probe)" probe="liveness" status="unhealthy" <_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (probe)" probe="readiness" status="" pod="hosted-grafana/dafdeveuwest2-grafana-7845d969<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop (probe)" probe="readiness" status="ready" <_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop DELETE" source="api" <_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet.go:<_>] "SyncLoop REMOVE" source="api" <_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-gke-dev-eu-west-3-main<_>status="Running"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-gke-dev-us-central-0-c<_>status="Running"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-gke-dev-us-central-0-d<_>status="Running"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-gke-dev-us-central-0-h<_>status="Running"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-gke-dev-us-central-0-o<_>status="Running"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-gke-dev-us-central-0-p<_>status="Running"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-gke-dev-us-central-0-s<_>status="Running"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-gke-dev-us-central-<_>-m<_>status="Running"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_pods.go:906] "Unable to retrieve pull secret, the image pull ma<_>pod="grafana-apps/bigquery-datasource-grafana-app-<_>secret="" err="secret \"dockerhub\" not found"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_pods.go:906] "Unable to retrieve pull secret, the image pull ma<_>pod="grafana-apps/loki-datasource-grafana-app-fast<_>secret="" err="secret \"dockerhub\" not found"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_pods.go:906] "Unable to retrieve pull secret, the image pull ma<_>pod="grafana-apps/query-grafana-app-fast-7d6dfcc78<_>secret="" err="secret \"dockerhub\" not found"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_pods.go:906] "Unable to retrieve pull secret, the image pull ma<_>pod="integration/grafana-render-service-cbff479fc-<_>secret="" err="secret \"us-gcr-io-hosted-grafana\" not found<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet_pods.go:906] "Unable to retrieve pull secret, the image pull ma<_>pod="kafka/kafka-broker-1" secret="" err="secret \"gcr\" not found"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_pods.go:906] "Unable to retrieve pull secret, the image pull ma<_>pod="kafka/kafka-controller-2" secret="" err="secret \"gcr\" not found"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_pods.go:906] "Unable to retrieve pull secret, the image pull ma<_>pod="logs-endpoint-dev-005/kafka-exporter-766c6757<_>secret="" err="secret \"not-needed\" not found"`,
				`I0507 <_>:<_>:<_>    <_> kubelet_volumes.go:<_>] "Cleaned up orphaned pod volumes dir" podUID="10bdda8a-7f0b466e9c81-045fb5150dc4" path="/var/lib/kubelet/pods/10bdda8a-7f0b466e9c81-<_>`,
				`I0507 <_>:<_>:<_>    <_> kubelet_volumes.go:<_>] "Cleaned up orphaned pod volumes dir" podUID="25cb986c-3d6c4ed0-abf3-ee59ed6175f9" path="/var/lib/kubelet/pods/25cb986c-3d6c4ed0-abf3<_>`,
				`I0507 <_>:<_>:<_>    <_> pod_container_deletor.go:53] "DeleteContainer returned error" containerID={"Type":"containerd","ID":"c50338fdb99<_>err="failed to get container status \"c50338fdb990<_>`,
				`I0507 <_>:<_>:<_>    <_> pod_container_deletor.go:53] "DeleteContainer returned error" containerID={"Type":"containerd","ID":"c8a30401d2a<_>err="failed to get container status \"c8a30401d2ac<_>`,
				`I0507 <_>:<_>:<_>    <_> pod_container_deletor.go:53] "DeleteContainer returned error" containerID={"Type":"containerd","ID":"c91436db009<_>err="failed to get container status \"c91436db0092<_>`,
				`I0507 <_>:<_>:<_>    <_> pod_container_deletor.go:53] "DeleteContainer returned error" containerID={"Type":"containerd","ID":"ea8c181d2a9<_>err="failed to get container status \"ea8c181d2a9b<_>`,
				`I0507 <_>:<_>:<_>    <_> prober.go:107] "Probe failed" probeType="Readiness" pod="agent-management-dev-002/agent-management-api<_>podUID="9893f9ac-f3e4-41fb-8da7-592061d2386c" containerName="agent-management-api" probeResult="failure" output="HTTP probe failed with statuscode: 400"`,
				`I0507 <_>:<_>:<_>    <_> prober.go:107] "Probe failed" probeType="Readiness" pod="grafana-agent/grafana-agent-helm-4" podUID="c36c5200-1cd6-4093-893c-c022f91af996" containerName="grafana-agent" probeResult="failure" output="Get \"http://<_>:3090/-/ready\": dial tcp<_>`,
				`I0507 <_>:<_>:<_>    <_> prober.go:107] "Probe failed" probeType="Readiness" pod="hosted-grafana/benchloadtestingxxl2-grafana-5<_>podUID="<_>" containerName="grafana" probeResult="failure" output=<`,
				`I0507 <_>:<_>:<_>    <_> prober.go:107] "Probe failed" probeType="Readiness" pod="hosted-grafana/dafdeveuwest2-grafana-7845d969<_>podUID="14ac9939-b36a-40d7-9ca9-a0367aab99d8" containerName="grafana" probeResult="failure" output=<`,
				`I0507 <_>:<_>:<_>    <_> prober.go:107] "Probe failed" probeType="Readiness" pod="hosted-grafana/k6teststeady3-grafana-659d5ff5<_>podUID="85274c17-190e4275-a8f3-6e111cd833bf" containerName="grafana" probeResult="failure" output=<`,
				`I0507 <_>:<_>:<_>    <_> prober.go:107] "Probe failed" probeType="Readiness" pod="loki-dev-014/loki-dev-014-rollout-operator-58<_>podUID="e6504036-2514-4ecc-b78c-c47061f60c9f" containerName="rollout-operator" probeResult="failure" output="HTTP probe failed with statuscode: 500"`,
				`I0507 <_>:<_>:<_>    <_> scope.go:117] "RemoveContainer" <_>`,
				`I0507 <_>:<_>:<_>  581823 cache.go:40] re-using cached key and certificate`,
				`I0507 <_>:<_>:<_> 1537502 kubelet_pods.go:906] "Unable to retrieve pull secret, the image pull ma<_>pod="logs-endpoint-dev-005/kafka-controller-0" secret="" err="secret \"not-needed\" not found"`,
				`I0507 <_>:<_>:<_> <_> cache.go:40] re-using cached key and certificate`,
				`IPv4: martian source <_> from <_>, on dev eth0`,
				`PRC: Renewing lease on eth0.`,
				`RCV: Reply message on eth0 from fe80::e9:7eff:fedf:3d37.`,
				`Removed slice libcontainer container kubepods-burstable-pod25cb986c_3d6c_4ed0_abf3_ee59<_>`,
				`Started cri-containerd-95bf586cd79d43120ff44582d4dbd2476de<_>`,
				`Started libcontainer container <_>`,
				`XMT: Renew on eth0, interval 9700ms.`,
				`XMT: Solicit on eth0, interval <_>`,
				`audit: type=1400 audit(<_>:<_>): apparmor="DENIED" operation="ptrace" profile="cri-containerd.apparmor.d" pid=<_> comm="pidof" requested_mask="read" denied_mask="read" peer="unconfined"`,
				`cri-containerd-<_><_>Consumed <_> CPU time.`,
				`cri-containerd-<_><_>Deactivated successfully.`,
				`kauditd_printk_skb: <_> callbacks suppressed`,
				`ll header: 00000000: 42 01 0a 80 00 17 42 01 0a 80 00 01 08 00`,
				`ll header: 00000000: 42 01 0a 80 00 7c 42 01 0a 80 00 01 08 00`,
				`ll header: 00000000: 42 01 0a 80 00 8f 42 01 0a 80 00 01 08 00`,
				`net_ratelimit: 2 callbacks suppressed`,
				`run-containerd-io.containerd.runtime.v2.task-k8s.i<_>Deactivated successfully.`,
				`run-containerd-runc-k8s.io-e5f17d69eee483ec8d43b26<_>Deactivated successfully.`,
				`time="<_>" level=error msg="ContainerStatus for \"<_><_>error="rpc error: code = NotFound desc = an error <_>`,
				`time="<_>" level=error msg="ExecSync for \"<_><_>error="rpc error: code = NotFound desc = failed to<_>`,
				`time="<_>" level=error msg="Failed to delete exec process \"d9e0a1867ce73<_>error="ttrpc: closed: unknown"`,
				`time="<_>" level=error msg="PullImage \"us.gcr.io/hosted-grafana/hosted-g<_><_>`,
				`time="<_>" level=info <_>`,
				`time="<_>" level=info msg="cleaning up dead shim" namespace=k8s.io`,
				`time="<_>" level=info msg="shim disconnected" id=<_><_>namespace=k8s.io`,
				`time="<_>" level=info msg="trying next host - response was http.StatusNo<_>host=us.gcr.io`,
				`time="<_>" level=warning msg="cleaning up after shim disconnected" id=<_><_>namespace=k8s.io`,
				`var-lib-containerd-tmpmounts-containerd\x2dmount40<_>Deactivated successfully.`,
				`var-lib-containerd-tmpmounts-containerd\x2dmount77<_>Deactivated successfully.`},
		},
		{
			name:      "Patterns for kafka logs",
			tokenizer: &AdaptiveTokenizer{},
			inputFile: "testdata/kafka.txt",
			patterns: []string{
				`[<_>,<_>] INFO Deleted log /bitnami/kafka/data/ingest-<_>/<_>.<_>(kafka.log.LogSegment)`,
				`[<_>,<_>] INFO Deleted log /bitnami/kafka/data/mimir-dev-09-aggregations-offs<_>(kafka.log.LogSegment)`,
				`[<_>,<_>] INFO Deleted offset index /bitnami/kafka/data/ingest-<_>/<_>.<_>(kafka.log.LogSegment)`,
				`[<_>,<_>] INFO Deleted offset index /bitnami/kafka/data/mimir-dev-09-aggregations-offs<_>(kafka.log.LogSegment)`,
				`[<_>,<_>] INFO Deleted producer state snapshot /bitnami/kafka/data/ingest-<_>/<_>.<_>(kafka.log.SnapshotFile)`,
				`[<_>,<_>] INFO Deleted producer state snapshot /bitnami/kafka/data/mimir-dev-09-aggregations-offs<_>(kafka.log.SnapshotFile)`,
				`[<_>,<_>] INFO Deleted time index /bitnami/kafka/data/ingest-<_>/<_>.<_>(kafka.log.LogSegment)`,
				`[<_>,<_>] INFO Deleted time index /bitnami/kafka/data/mimir-dev-09-aggregations-offs<_>(kafka.log.LogSegment)`,
				`[<_>,<_>] INFO [LocalLog partition=cortex-dev-01-aggregations-offsets-1, dir=/bitnami/kafka/data] Rolled new log segment at offset 2142125 in 0 ms. (kafka.log.LocalLog)`,
				`[<_>,<_>] INFO [LocalLog partition=ingest-6, dir=/bitnami/kafka/data] Deleting segment files LogSegment(baseOffset=180391157, size=16991045, lastModifiedTime=1715075754780, largestRecordTimestamp=Some(1715075754774)),LogSeg<_>size=16997692, lastModifiedTime=1715075760206, largestRecordTimestamp=Some(1715075760186)),LogSeg<_>size=16998200, lastModifiedTime=1715075765542, largestRecordTimestamp=Some(1715075765526)),LogSeg<_>size=16977347, lastModifiedTime=1715075770515, largestRecordTimestamp=Some(1715075770504)) (kafka.log.LocalLog$)`,
				`[<_>,<_>] INFO [LocalLog partition=ingest-<_>, dir=/bitnami/kafka/data] Deleting segment files LogSegment(baseOffset=<_>, size=<_>, lastModifiedTime=<_>, largestRecordTimestamp=Some(<_>)),LogSeg<_>size=<_>, lastModifiedTime=<_>, largestRecordTimestamp=Some(<_>)),LogSeg<_>size=16989895, lastModifiedTime=1715075786205, largestRecordTimestamp=Some(1715075786174)),LogSeg<_>size=16998698, lastModifiedTime=1715075791681, largestRecordTimestamp=Some(1715075791673)),LogSeg<_>size=16995676, lastModifiedTime=1715075796438, largestRecordTimestamp=Some(1715075796430)),LogSeg<_>size=16963278, lastModifiedTime=1715075800534, largestRecordTimestamp=Some(1715075800511)),LogSeg<_>size=16984328, lastModifiedTime=1715075805272, largestRecordTimestamp=Some(1715075805230)),LogSeg<_>size=16989109, lastModifiedTime=1715075810381, largestRecordTimestamp=Some(1715075810372)),LogSeg<_>size=16996871, lastModifiedTime=1715075815153, largestRecordTimestamp=Some(1715075815125)),LogSeg<_>size=16988558, lastModifiedTime=1715075819785, largestRecordTimestamp=Some(1715075819763)),LogSeg<_>size=16999292, lastModifiedTime=1715075825336, largestRecordTimestamp=Some(1715075825303)),LogSeg<_>size=16990595, lastModifiedTime=1715075830839, largestRecordTimestamp=Some(1715075830827)),LogSeg<_>size=16995859, lastModifiedTime=1715075835942, largestRecordTimestamp=Some(1715075835904)),LogSeg<_>size=16992294, lastModifiedTime=1715075841219, largestRecordTimestamp=Some(1715075841214)),LogSeg<_>size=16966736, lastModifiedTime=1715075846443, largestRecordTimestamp=Some(1715075846401)),LogSeg<_>size=16894731, lastModifiedTime=1715075853273, largestRecordTimestamp=Some(1715075853244)),LogSeg<_>size=16983529, lastModifiedTime=1715075858911, largestRecordTimestamp=Some(1715075858891)),LogSeg<_>size=16996933, lastModifiedTime=1715075863566, largestRecordTimestamp=Some(1715075863554)),LogSeg<_>size=16999841, lastModifiedTime=1715075866199, largestRecordTimestamp=Some(1715075866185)),LogSeg<_>size=16992471, lastModifiedTime=1715075870385, largestRecordTimestamp=Some(1715075870347)),LogSeg<_>size=16999996, lastModifiedTime=1715075875102, largestRecordTimestamp=Some(1715075875091)),LogSeg<_>size=16994426, lastModifiedTime=1715075879927, largestRecordTimestamp=Some(1715075879926)),LogSeg<_>size=16998020, lastModifiedTime=1715075885293, largestRecordTimestamp=Some(1715075885263)),LogSeg<_>size=16992231, lastModifiedTime=1715075890424, largestRecordTimestamp=Some(1715075890409)),LogSeg<_>size=16970315, lastModifiedTime=1715075895719, largestRecordTimestamp=Some(1715075895690)),LogSeg<_>size=16990785, lastModifiedTime=1715075900996, largestRecordTimestamp=Some(1715075900985)),LogSeg<_>size=16996655, lastModifiedTime=1715075905847, largestRecordTimestamp=Some(1715075905841)),LogSeg<_>size=16982181, lastModifiedTime=1715075911052, largestRecordTimestamp=Some(1715075911028)),LogSeg<_>size=16997630, lastModifiedTime=1715075915962, largestRecordTimestamp=Some(1715075915953)),LogSeg<_>size=16995723, lastModifiedTime=1715075920325, largestRecordTimestamp=Some(1715075920308)),LogSeg<_><...>`,
				`[<_>,<_>] INFO [LocalLog partition=ingest-<_>, dir=/bitnami/kafka/data] Deleting segment files LogSegment(baseOffset=<_>, size=<_>, lastModifiedTime=<_>, largestRecordTimestamp=Some(<_>)),LogSeg<_>size=<_>, lastModifiedTime=<_>, largestRecordTimestamp=Some(<_>)),LogSeg<_>size=16994485, lastModifiedTime=1715075770425, largestRecordTimestamp=Some(1715075770404)),LogSeg<_>size=16996810, lastModifiedTime=1715075775622, largestRecordTimestamp=Some(1715075775619)),LogSeg<_>size=16998520, lastModifiedTime=1715075780912, largestRecordTimestamp=Some(1715075780889)),LogSeg<_>size=16988474, lastModifiedTime=1715075786051, largestRecordTimestamp=Some(1715075786030)),LogSeg<_>size=16956099, lastModifiedTime=1715075791514, largestRecordTimestamp=Some(1715075791486)),LogSeg<_>size=16995476, lastModifiedTime=1715075796360, largestRecordTimestamp=Some(1715075796329)),LogSeg<_>size=16993313, lastModifiedTime=1715075800440, largestRecordTimestamp=Some(1715075800430)),LogSeg<_>size=16992142, lastModifiedTime=1715075805147, largestRecordTimestamp=Some(1715075805135)),LogSeg<_>size=16999919, lastModifiedTime=1715075810155, largestRecordTimestamp=Some(1715075810153)),LogSeg<_>size=16995021, lastModifiedTime=1715075815018, largestRecordTimestamp=Some(1715075815016)),LogSeg<_>size=16966526, lastModifiedTime=1715075819528, largestRecordTimestamp=Some(1715075819521)),LogSeg<_>size=16990848, lastModifiedTime=1715075825066, largestRecordTimestamp=Some(1715075825042)),LogSeg<_>size=16997833, lastModifiedTime=1715075830662, largestRecordTimestamp=Some(1715075830656)),LogSeg<_>size=16992619, lastModifiedTime=1715075835771, largestRecordTimestamp=Some(1715075835741)),LogSeg<_>size=16999091, lastModifiedTime=1715075841031, largestRecordTimestamp=Some(1715075841022)),LogSeg<_>size=16993953, lastModifiedTime=1715075846197, largestRecordTimestamp=Some(1715075846181)),LogSeg<_>size=16997479, lastModifiedTime=1715075853192, largestRecordTimestamp=Some(1715075853172)),LogSeg<_>size=16997174, lastModifiedTime=1715075858693, largestRecordTimestamp=Some(1715075858682)),LogSeg<_>size=16986004, lastModifiedTime=1715075863400, largestRecordTimestamp=Some(1715075863396)),LogSeg<_>size=16995316, lastModifiedTime=1715075866123, largestRecordTimestamp=Some(1715075866112)),LogSeg<_>size=16990492, lastModifiedTime=1715075870154, largestRecordTimestamp=Some(1715075870146)),LogSeg<_>size=16999541, lastModifiedTime=1715075874980, largestRecordTimestamp=Some(1715075874961)),LogSeg<_>size=16987383, lastModifiedTime=1715075879670, largestRecordTimestamp=Some(1715075879639)),LogSeg<_>size=16991701, lastModifiedTime=1715075885010, largestRecordTimestamp=Some(1715075884995)),LogSeg<_>size=16989109, lastModifiedTime=1715075890220, largestRecordTimestamp=Some(1715075890208)),LogSeg<_>size=16962782, lastModifiedTime=1715075895466, largestRecordTimestamp=Some(1715075895456)),LogSeg<_>size=16974715, lastModifiedTime=1715075900757, largestRecordTimestamp=Some(1715075900746)),LogSeg<_>size=16993973, lastModifiedTime=1715075905639, largestRecordTimestamp=Some(1715075905638)),LogSeg<_><...>`,
				`[<_>,<_>] INFO [LocalLog partition=ingest-<_>, dir=/bitnami/kafka/data] Deleting segment files LogSegment(baseOffset=<_>, size=<_>, lastModifiedTime=<_>, largestRecordTimestamp=Some(<_>)),LogSeg<_>size=<_>, lastModifiedTime=<_>, largestRecordTimestamp=Some(<_>)),LogSeg<_>size=16994792, lastModifiedTime=1715075768711, largestRecordTimestamp=Some(1715075768697)),LogSeg<_>size=16987578, lastModifiedTime=1715075773552, largestRecordTimestamp=Some(1715075773536)),LogSeg<_>size=16987705, lastModifiedTime=1715075779055, largestRecordTimestamp=Some(1715075779046)),LogSeg<_>size=16997466, lastModifiedTime=1715075784005, largestRecordTimestamp=Some(1715075784004)),LogSeg<_>size=16981250, lastModifiedTime=1715075789523, largestRecordTimestamp=Some(1715075789487)),LogSeg<_>size=16980484, lastModifiedTime=1715075794637, largestRecordTimestamp=Some(1715075794632)),LogSeg<_>size=16999738, lastModifiedTime=1715075799008, largestRecordTimestamp=Some(1715075799000)),LogSeg<_>size=16872695, lastModifiedTime=1715075803273, largestRecordTimestamp=Some(1715075803251)),LogSeg<_>size=16999890, lastModifiedTime=1715075808368, largestRecordTimestamp=Some(1715075808355)),LogSeg<_>size=16959982, lastModifiedTime=1715075813294, largestRecordTimestamp=Some(1715075813293)),LogSeg<_>size=16988073, lastModifiedTime=1715075817816, largestRecordTimestamp=Some(1715075817783)),LogSeg<_>size=16974731, lastModifiedTime=1715075823018, largestRecordTimestamp=Some(1715075823016)),LogSeg<_>size=16996090, lastModifiedTime=1715075828672, largestRecordTimestamp=Some(1715075828632)),LogSeg<_>size=16999327, lastModifiedTime=1715075833742, largestRecordTimestamp=Some(1715075833709)),LogSeg<_>size=16992947, lastModifiedTime=1715075839121, largestRecordTimestamp=Some(1715075839114)),LogSeg<_>size=16982572, lastModifiedTime=1715075844268, largestRecordTimestamp=Some(1715075844254)),LogSeg<_>size=16994786, lastModifiedTime=1715075850659, largestRecordTimestamp=Some(1715075850642)),LogSeg<_>size=16998391, lastModifiedTime=1715075856704, largestRecordTimestamp=Some(1715075856684)),LogSeg<_>size=16994403, lastModifiedTime=1715075861956, largestRecordTimestamp=Some(1715075861922)),LogSeg<_>size=16984546, lastModifiedTime=1715075865194, largestRecordTimestamp=Some(1715075865180)),LogSeg<_>size=16987846, lastModifiedTime=1715075868470, largestRecordTimestamp=Some(1715075868460)),LogSeg<_>size=16958237, lastModifiedTime=1715075873168, largestRecordTimestamp=Some(1715075873151)),LogSeg<_>size=16999432, lastModifiedTime=1715075877858, largestRecordTimestamp=Some(1715075877850)),LogSeg<_>size=16938567, lastModifiedTime=1715075882952, largestRecordTimestamp=Some(1715075882938)),LogSeg<_>size=16998214, lastModifiedTime=1715075888306, largestRecordTimestamp=Some(1715075888285)),LogSeg<_>size=16996264, lastModifiedTime=1715075893370, largestRecordTimestamp=Some(1715075893365)),LogSeg<_>size=16991650, lastModifiedTime=1715075898806, largestRecordTimestamp=Some(1715075898802)),LogSeg<_>size=16998234, lastModifiedTime=1715075903737, largestRecordTimestamp=Some(1715075903733)),LogSeg<_><...>`,
				`[<_>,<_>] INFO [LocalLog partition=ingest-<_>, dir=/bitnami/kafka/data] Rolled new log segment at offset <_> in <_> ms. (kafka.log.LocalLog)`,
				`[<_>,<_>] INFO [LocalLog partition=mimir-dev-09-aggregations-offsets-0, dir=/bitnami/kafka/data] Deleting segment files LogSegment(baseOffset=<_>, size=948, lastModifiedTime=<_>, largestRecordTimestamp=Some(<_>)) (kafka.log.LocalLog$)`,
				`[<_>,<_>] INFO [LocalLog partition=mimir-dev-<_>-aggregations-offsets-<_>, dir=/bitnami/kafka/data] Deleting segment files LogSegment(baseOffset=447957, size=948, lastModifiedTime=1715059232052, largestRecordTimestamp=Some(1715059232002)),LogSeg<_>size=948, lastModifiedTime=1715059424352, largestRecordTimestamp=Some(1715059424301)) (kafka.log.LocalLog$)`,
				`[<_>,<_>] INFO [LocalLog partition=mimir-dev-<_>-aggregations-offsets-<_>, dir=/bitnami/kafka/data] Rolled new log segment at offset 27664 in 0 ms. (kafka.log.LocalLog)`,
				`[<_>,<_>] INFO [ProducerStateManager partition=cortex-dev-01-aggregations-offsets-1] Wrote producer snapshot at offset 2142125 with 0 producer ids in 6 ms. (kafka.log.ProducerStateManager)`,
				`[<_>,<_>] INFO [ProducerStateManager partition=ingest-<_>] Wrote producer snapshot at offset <_> with 0 producer ids in <_>ms. (kafka.log.ProducerStateManager)`,
				`[<_>,<_>] INFO [ProducerStateManager partition=mimir-dev-14-aggregations-offsets-3] Wrote producer snapshot at offset 27664 with 0 producer ids in 43 ms. (kafka.log.ProducerStateManager)`,
				`[<_>,<_>] INFO [UnifiedLog partition=cortex-dev-<_>-aggregations-offsets-<_>, dir=/bitnami/kafka/data] Incremented log start offset to <_> due to leader offset increment (kafka.log.UnifiedLog)`,
				`[<_>,<_>] INFO [UnifiedLog partition=mimir-dev-11-aggregations-offsets-0, dir=/bitnami/kafka/data] Deleting segments due to log start offset 1452491 breach: LogSegment(baseOffset=1452479, size=972, lastModifiedTime=1715059950760, largestRecordTimestamp=Some(1715059950710)) (kafka.log.UnifiedLog)`,
				`[<_>,<_>] INFO [UnifiedLog partition=mimir-dev-<_>-aggregations-<_>, dir=/bitnami/kafka/data] Deleting segment LogSegment(baseOffset=<_>, size=<_>, lastModifiedTime=<_>, largestRecordTimestamp=Some(<_>)) due to retention size 38386270208 breach. Log size after deletion will be <_>(kafka.log.UnifiedLog)`,
				`[<_>,<_>] INFO [UnifiedLog partition=mimir-dev-<_>-aggregations-<_>, dir=/bitnami/kafka/data] Deleting segments LogSegment(baseOffset=<_>, to log <_> offset <_> breach: LogSegment(baseOffset=<_>, <_><_><_>(kafka.log.UnifiedLog)`,
				`[<_>,<_>] INFO [UnifiedLog partition=mimir-dev-<_>-aggregations-<_>, dir=/bitnami/kafka/data] Incremented log start offset to <_> due to segment deletion (kafka.log.UnifiedLog)`,
				`[<_>,<_>] INFO [UnifiedLog partition=mimir-dev-<_>-aggregations-offsets-<_>, dir=/bitnami/kafka/data] Deleting segment LogSegment(baseOffset=<_>, size=<_>, lastModifiedTime=<_>, largestRecordTimestamp=Some(<_>)) due to retention size 102400 breach. Log size after deletion will be <_>(kafka.log.UnifiedLog)`,
				`[<_>,<_>] INFO [UnifiedLog partition=mimir-dev-<_>-aggregations-offsets-<_>, dir=/bitnami/kafka/data] Deleting segments due size=<_>, log start offset <_> breach: LogSegment(baseOffset=<_>, size=948, <_><_>size=948, <_><_>(kafka.log.UnifiedLog)`,
				`[<_>,<_>] INFO [UnifiedLog partition=mimir-dev-<_>-aggregations-offsets-<_>, dir=/bitnami/kafka/data] Incremented log start offset to <_> <_> to segment deletion (kafka.log.UnifiedLog)`,
			},
		},
		{
			name:      "Patterns for kubernetes logs",
			tokenizer: &AdaptiveTokenizer{},
			inputFile: "testdata/kubernetes.txt",
			patterns: []string{
				`I0507 12:02:27.947830       1 nodeutilization.go:274] "Evicting pods based on priority, if they have sam<_>`,
				`I0507 <_>:<_>:<_>       1 defaultevictor.go:163] "pod does not fit on any other node because of nod<_><_>`,
				`I0507 <_>:<_>:<_>       1 defaultevictor.go:202] "Pod fails the following checks" <_><_>`,
				`I0507 <_>:<_>:<_>       1 defaultevictor.go:202] "Pod fails the following checks" pod="kube-system/calico-node-cnc6m" checks="[pod is a DaemonSet pod, pod has system cr<_>`,
				`I0507 <_>:<_>:<_>       1 defaultevictor.go:202] "Pod fails the following checks" pod="kube-system/calico-typha-7cc4789bc8-qhw5r" checks="[pod has system critical priority, pod has<_>`,
				`I0507 <_>:<_>:<_>       1 defaultevictor.go:202] "Pod fails the following checks" pod="kube-system/konnectivity-agent-6f8f85c4fb-7bh<_>checks="[pod has system critical priority, pod has<_>`,
				`I0507 <_>:<_>:<_>       1 defaultevictor.go:202] "Pod fails the following checks" pod="kube-system/pdcsi-node-7khn6" checks="[pod is a DaemonSet pod, pod has system cr<_>`,
				`I0507 <_>:<_>:<_>       1 defaultevictor.go:202] "Pod fails the following checks" pod="netfilter-exporter/netfilter-exporter-jkrhn" checks="[pod is a DaemonSet pod, pod has higher pr<_>`,
				`I0507 <_>:<_>:<_>       1 defaultevictor.go:202] "Pod fails the following checks" pod="node-exporter/node-exporter-h82wd" checks="[pod is a DaemonSet pod, pod has higher pr<_>`,
				`I0507 <_>:<_>:<_>       1 defaultevictor.go:202] "Pod fails the following checks" pod="promtail-ops/loki-canary-n5p56" checks="[pod is a DaemonSet pod, pod has higher pr<_>`,
				`I0507 <_>:<_>:<_>       1 descheduler.go:155] Building a pod evictor`,
				`I0507 <_>:<_>:<_>       1 descheduler.go:<_>] "Number of evicted pods" <_>`,
				`I0507 <_>:<_>:<_>       1 highnodeutilization.go:107] "Criteria for a node below target utilization" CPU=50 Mem=50 Pods=100`,
				`I0507 <_>:<_>:<_>       1 highnodeutilization.go:108] "Number of underutilized nodes" totalNumber=1`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="gel-sbdev/gel-4" <_><_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="gel-sbdev/gel-4" node:="gke-dev-us-central-0-cache-n2hc8-1-1d61155f<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="gel-sbdev/gel-4" node:="gke-dev-us-central-0-databenchloki-n2-8c6b6<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="gel-sbdev/gel-4" node:="gke-dev-us-central-0-hg-n2s4-7-1dd39c-6f2ad<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="gel-sbdev/gel-4" node:="gke-dev-us-central-0-hg-n2s8-6-1dd39c-3bfd0<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="gel-sbdev/gel-4" node:="gke-dev-us-central-0-main-n2s16-3-1dd-9b502<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="gel-sbdev/gel-4" node:="gke-dev-us-central-0-otel-n2s4-0-1dd3-b196a<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="gel-sbdev/gel-4" node:="gke-dev-us-central-0-spot-n2s8-0-1dd3-f8133<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="loki-dev-005/querier-burst-6b5f6db455-5zvkm"<_><_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="loki-dev-005/querier-burst-6b5f6db455-5zvkm"<_>_>de:="gke-dev-us-central-0-cache-n2hc8-1-1d61155f<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="loki-dev-005/querier-burst-6b5f6db455-5zvkm"<_>_>de:="gke-dev-us-central-0-databenchloki-n2-8c6b6<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="loki-dev-005/querier-burst-6b5f6db455-5zvkm"<_>_>de:="gke-dev-us-central-0-hg-n2s8-6-1dd39c-3bfd0<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="loki-dev-005/querier-burst-6b5f6db455-5zvkm"<_>_>de:="gke-dev-us-central-0-main-n2s16-3-1dd-9b502<_><_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="loki-dev-005/querier-burst-6b5f6db455-5zvkm"<_>_>de:="gke-dev-us-central-0-otel-alt-n2s4-0-3cf760<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="loki-dev-005/querier-burst-6b5f6db455-5zvkm"<_>_>de:="gke-dev-us-central-0-perf-n2s8-0-1dd3-91689<_>error:="[pod node selector does not match the node<_>`,
				`I0507 <_>:<_>:<_>       1 node.go:157] "Pod does not fit on any other node" pod:="loki-dev-005/querier-burst-6b5f6db455-5zvkm"<_>_>de:="gke-dev-us-central-0-spot-n2s8-0-1dd3-f8133<_><_>`,
				`I0507 <_>:<_>:<_>       1 node.go:339] "no Pod antiaffinity rule found" <_>`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:260] "Total capacity to be moved" CPU=5060 Mem=112216292800 Pods=163`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:<_>] "Evicting pods from node" node="gke-dev-eu-west-3-main-n2s8-1-1dd39c-d1c9206<_>usage={"cpu":"984m","memory":"611Mi","pods":"16"}`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:<_>] "Evicting pods from node" node="gke-dev-us-central-0-main-n2s16-3-1dd-9b502d<_>usage={"cpu":"<_>","memory":<_>,"pods":"64"}`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:<_>] "Evicting pods from node" node="gke-dev-us-central-0-spot-n2s8-0-1dd3-f81338<_>usage={"cpu":"6826m","memory":"16564Mi","pods":"20"}`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:<_>] "No removable pods on node, try next node" <_>`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:<_>] "Node is overutilized" node="gke-dev-eu-west-3-main-n2s8-1-1dd39c-d1c9206<_>usage={"cpu":"<_>","memory":<_>,"pods":<_>usagePercentage={"cpu":<_>,"memory":<_>,"pods":<_>`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:<_>] "Node is underutilized" node="gke-dev-eu-west-3-main-n2s8-1-1dd39c-d1c9206<_>usage={"cpu":"984m","memory":"611Mi","pods":"16"} usagePercentage={"cpu":12.44,"memory":2.15,"pods":25}`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:<_>] "Pods on node" node="gke-dev-eu-west-3-main-n2s8-1-1dd39c-d1c9206<_>allPods=16 nonRemovablePods=16 removablePods=0`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:<_>] "Pods on node" node="gke-dev-us-central-0-main-n2s16-3-1dd-9b502d<_>allPods=64 nonRemovablePods=64 removablePods=0`,
				`I0507 <_>:<_>:<_>       1 nodeutilization.go:<_>] "Pods on node" node="gke-dev-us-central-0-spot-n2s8-0-1dd3-f81338<_>allPods=20 nonRemovablePods=19 removablePods=1`,
				`I0507 <_>:<_>:<_>       1 profile.go:<_>] "Total number of pods evicted" extension point="Balance" <_>`,
				`I0507 <_>:<_>:<_>       1 reflector.go:<_>] k8s.io/client-go/informers/factory.go:<_>: Watch close - *v1.Node total 71 items received`,
				`I0507 <_>:<_>:<_>       1 reflector.go:<_>] k8s.io/client-go/informers/factory.go:<_>: Watch close - *v1.PriorityClass total 7 items received`,
			},
		},
		{
			name:      "Patterns for vault logs",
			tokenizer: &AdaptiveTokenizer{},
			inputFile: "testdata/vault.txt",
			patterns: []string{
				"<_> [INFO]  expiration: revoked lease: <_>",
			},
		},
		{
			name:      "Patterns for calico logs",
			tokenizer: &AdaptiveTokenizer{},
			inputFile: "testdata/calico.txt",
			patterns: []string{
				`<_> [DEBUG][<_>] felix/endpoint_mgr.go 443: Reporting endpoint status. dirtyEndpoints=set.Set{}`,
				`<_> [DEBUG][<_>] felix/feature_detect.go <_>: Parsed iptables version version=<_>.<_>.<_>`,
				`<_> [DEBUG][<_>] felix/feature_detect.go <_>: Ran iptables --version rawVersion="iptables v1.8.4 (legacy)\n"`,
				`<_> [DEBUG][<_>] felix/feature_detect.go <_>: Refreshing detected iptables features`,
				`<_> [DEBUG][<_>] felix/health.go 196: Checking state of reporter reporter=&health.reporterState{name:"async_calc_gr<_>reports:health.HealthReport{Live:true, Ready:true, Detail:""}, timeout:20000000000, latest:health.HealthReport{Live:true, Ready:true, Detail:""}, <_><_>loc:(*time.Location)(0x4ce3aa0)}}`,
				`<_> [DEBUG][<_>] felix/health.go 196: Checking state of reporter reporter=&health.reporterState{name:"felix-startup<_>reports:health.HealthReport{Live:true, Ready:true, Detail:""}, timeout:0, latest:health.HealthReport{Live:true, Ready:true, Detail:""}, <_><_>loc:(*time.Location)(0x4ce3aa0)}}`,
				`<_> [DEBUG][<_>] felix/health.go 196: Checking state of reporter reporter=&health.reporterState{name:"int_dataplane<_>reports:health.HealthReport{Live:true, Ready:true, Detail:""}, timeout:90000000000, latest:health.HealthReport{Live:true, Ready:true, Detail:""}, <_><_>loc:(*time.Location)(0x4ce3aa0)}}`,
				`<_> [DEBUG][<_>] felix/health.go 245: Calculated health summary healthResult=&health.HealthReport{Live:true, Ready:true, Detail:"+------------------+---------+------------<_>`,
				`<_> [DEBUG][<_>] felix/health.go <_>: GET <_>`,
				`<_> [DEBUG][<_>] felix/health.go <_>: Health: <_>`,
				`<_> [DEBUG][<_>] felix/int_dataplane.go 1777: Refreshing routes`,
				`<_> [DEBUG][<_>] felix/int_dataplane.go <_>: Applying dataplane updates`,
				`<_> [DEBUG][<_>] felix/int_dataplane.go <_>: Asked to reschedule. delay=<_>`,
				`<_> [DEBUG][<_>] felix/int_dataplane.go <_>: Examining link for MTU calculation mtu=1500 name="eth0"`,
				`<_> [DEBUG][<_>] felix/int_dataplane.go <_>: Refreshing IP sets state`,
				`<_> [DEBUG][<_>] felix/int_dataplane.go <_>: Reschedule kick received`,
				`<_> [DEBUG][<_>] felix/int_dataplane.go <_>: Skipping interface for MTU detection mtu=<_> <_>`,
				`<_> [DEBUG][<_>] felix/ipsets.go 234: Asked to resync with the dataplane on next update. family="inet"`,
				`<_> [DEBUG][<_>] felix/ipsets.go 467: Found member in dataplane canon=<_> family="inet" member="<_>" setID="this-host"`,
				`<_> [DEBUG][<_>] felix/ipsets.go 607: Skipping expected Calico IP set. family="inet" <_>`,
				`<_> [DEBUG][<_>] felix/ipsets.go <_>: Finished IPSets resync family="inet" numInconsistenciesFound=0 resyncDuration=<_>`,
				`<_> [DEBUG][<_>] felix/ipsets.go <_>: No dirty IP sets. family="inet"`,
				`<_> [DEBUG][<_>] felix/ipsets.go <_>: Parsing IP set. family="inet" <_>`,
				`<_> [DEBUG][<_>] felix/ipsets.go <_>: Resyncing ipsets with dataplane. family="inet"`,
				`<_> [DEBUG][<_>] felix/ipsets.go <_>: Whitelisting IP sets. ID="all-ipam-pools" family="inet" mainName="cali40all-ipam-pools"`,
				`<_> [DEBUG][<_>] felix/ipsets.go <_>: Whitelisting IP sets. ID="masq-ipam-pools" family="inet" mainName="cali40masq-ipam-pools"`,
				`<_> [DEBUG][<_>] felix/ipsets.go <_>: Whitelisting IP sets. ID="this-host" family="inet" mainName="cali40this-host"`,
				`<_> [DEBUG][<_>] felix/route_rule.go 179: Queueing a resync of routing rules. ipVersion=4`,
				`<_> [DEBUG][<_>] felix/route_table.go 533: Check interfaces matching regex`,
				`<_> [DEBUG][<_>] felix/route_table.go 584: Flag no OIF for full re-sync`,
				`<_> [DEBUG][<_>] felix/route_table.go 661: Syncing interface routes <_>ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go 661: Syncing interface routes ifaceName="*NoOIF*" ifaceRegex="^wireguard.cali$" ipVersion=0x4 tableIndex=1`,
				`<_> [DEBUG][<_>] felix/route_table.go 661: Syncing interface routes ifaceName="azv1e0e3e8aac0" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go 661: Syncing interface routes ifaceName="azv24bd4f90868" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go 661: Syncing interface routes ifaceName="azv6767b9519e3" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go 661: Syncing interface routes ifaceName="azv7209a4b4cbc" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go 661: Syncing interface routes ifaceName="azvd32f7c1c18e" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go 661: Syncing interface routes ifaceName="azvd9f11c4f109" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go 661: Syncing interface routes ifaceName="azvddd03b40b4a" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Processing route: 254 <_> <_>/32 <_>ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Processing route: 254 <_> <_>/32 ifaceName="azv1e0e3e8aac0" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Processing route: 254 <_> <_>/32 ifaceName="azv24bd4f90868" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Processing route: 254 <_> <_>/32 ifaceName="azv6767b9519e3" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Processing route: 254 <_> <_>/32 ifaceName="azv7209a4b4cbc" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Processing route: 254 <_> <_>/32 ifaceName="azvd32f7c1c18e" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Processing route: 254 <_> <_>/32 ifaceName="azvd9f11c4f109" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Processing route: 254 <_> <_>/32 ifaceName="azvddd03b40b4a" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Queueing a resync of routing table. ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Queueing a resync of routing table. ifaceRegex="^wireguard.cali$" ipVersion=0x4 tableIndex=1`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Reconcile against kernel programming <_>ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Reconcile against kernel programming ifaceName="*NoOIF*" ifaceRegex="^wireguard.cali$" ipVersion=0x4 tableIndex=1`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Reconcile against kernel programming ifaceName="azv1e0e3e8aac0" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Reconcile against kernel programming ifaceName="azv24bd4f90868" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Reconcile against kernel programming ifaceName="azv6767b9519e3" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Reconcile against kernel programming ifaceName="azv7209a4b4cbc" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Reconcile against kernel programming ifaceName="azvd32f7c1c18e" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Reconcile against kernel programming ifaceName="azvd9f11c4f109" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Reconcile against kernel programming ifaceName="azvddd03b40b4a" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Resync: found calico-owned interface <_>ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Resync: found calico-owned interface ifaceName="azv24bd4f90868" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Resync: found calico-owned interface ifaceName="azv443ad95a1ab" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Resync: found calico-owned interface ifaceName="azv6767b9519e3" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Resync: found calico-owned interface ifaceName="azv7209a4b4cbc" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Resync: found calico-owned interface ifaceName="azvd9f11c4f109" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Resync: found calico-owned interface ifaceName="azvddd03b40b4a" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Resync: found calico-owned interface ifaceName="azve1df6b75675" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Route is correct dest=<_>/32 <_>ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Route is correct dest=<_>/32 ifaceName="azv1e0e3e8aac0" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Route is correct dest=<_>/32 ifaceName="azv24bd4f90868" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Route is correct dest=<_>/32 ifaceName="azv6767b9519e3" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Route is correct dest=<_>/32 ifaceName="azv7209a4b4cbc" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Route is correct dest=<_>/32 ifaceName="azvd32f7c1c18e" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Route is correct dest=<_>/32 ifaceName="azvd9f11c4f109" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Route is correct dest=<_>/32 ifaceName="azvddd03b40b4a" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Synchronised routes on interface <_>ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Synchronised routes on interface ifaceName="*NoOIF*" ifaceRegex="^wireguard.cali$" ipVersion=0x4 tableIndex=1`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Synchronised routes on interface ifaceName="azv1e0e3e8aac0" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Synchronised routes on interface ifaceName="azv24bd4f90868" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Synchronised routes on interface ifaceName="azv6767b9519e3" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Synchronised routes on interface ifaceName="azv7209a4b4cbc" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Synchronised routes on interface ifaceName="azvd32f7c1c18e" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Synchronised routes on interface ifaceName="azvd9f11c4f109" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/route_table.go <_>: Synchronised routes on interface ifaceName="azvddd03b40b4a" ifaceRegex="^azv.*" ipVersion=0x4 tableIndex=0`,
				`<_> [DEBUG][<_>] felix/sync_client.go 434: New message from Typha. connID=0x0 connection=&discovery.Typha{Addr:"", IP:"", NodeName:(*string)(nil)} envelope=syncproto.Envelope{Message:syncproto.MsgP<_>time.May, 8, 15, 23, <_><_>time.Local)}} type=""`,
				`<_> [DEBUG][<_>] felix/sync_client.go <_>: Ping received from Typha connID=0x0 connection=&discovery.Typha{Addr:"", IP:"", NodeName:(*string)(nil)} type=""`,
				`<_> [DEBUG][<_>] felix/sync_client.go <_>: Pong sent to Typha connID=0x0 connection=&discovery.Typha{Addr:"", IP:"", NodeName:(*string)(nil)} type=""`,
				`<_> [DEBUG][<_>] felix/table.go 851: Parsing line ipVersion=0x4 <_>table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go 851: Parsing line ipVersion=0x4 line=":KUBE-SEP-527TDH7QDHLCYDTX - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go 851: Parsing line ipVersion=0x4 line=":KUBE-SEP-5KVHYONDUWXKZLCF - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go 851: Parsing line ipVersion=0x4 line=":KUBE-SEP-DZCXKX63Q3ZRE2XB - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go 851: Parsing line ipVersion=0x4 line=":KUBE-SEP-IOPUYNOJID4CYL5S - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go 851: Parsing line ipVersion=0x4 line=":KUBE-SEP-P53KRBBAHF7EH6MF - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go 851: Parsing line ipVersion=0x4 line=":KUBE-SEP-RK34UV6XMAMZC6JG - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go 851: Parsing line ipVersion=0x4 line=":KUBE-SEP-YADSGSG25SR3HQ6W - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go 881: Not an append, skipping ipVersion=0x4 line="# Generated by iptables-nft-save v1.8.4 on <<_>table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go 881: Not an append, skipping ipVersion=0x4 line="*nat" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Finished loading iptables state ipVersion=0x4 table="filter"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference <_>ipVersion=0x4 <_>table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-527TDH7QDHLCYDTX" ipVersion=0x4 line=":KUBE-SEP-527TDH7QDHLCYDTX - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-5KVHYONDUWXKZLCF" ipVersion=0x4 line=":KUBE-SEP-5KVHYONDUWXKZLCF - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-DZCXKX63Q3ZRE2XB" ipVersion=0x4 line=":KUBE-SEP-DZCXKX63Q3ZRE2XB - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-E2GHKVOJHBYBPZ3C" ipVersion=0x4 line=":KUBE-SEP-E2GHKVOJHBYBPZ3C - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-IOPUYNOJID4CYL5S" ipVersion=0x4 line=":KUBE-SEP-IOPUYNOJID4CYL5S - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-JLFGSS5Y56HOFTOX" ipVersion=0x4 line=":KUBE-SEP-JLFGSS5Y56HOFTOX - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-RK34UV6XMAMZC6JG" ipVersion=0x4 line=":KUBE-SEP-RK34UV6XMAMZC6JG - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-TX2E3S6G3BZ6VCYU" ipVersion=0x4 line=":KUBE-SEP-TX2E3S6G3BZ6VCYU - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-U4I77HBN3HVEYELA" ipVersion=0x4 line=":KUBE-SEP-U4I77HBN3HVEYELA - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-VQOBWT5QN7AMFSUO" ipVersion=0x4 line=":KUBE-SEP-VQOBWT5QN7AMFSUO - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-WKQFW72ZLNYTB4P7" ipVersion=0x4 line=":KUBE-SEP-WKQFW72ZLNYTB4P7 - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-YADSGSG25SR3HQ6W" ipVersion=0x4 line=":KUBE-SEP-YADSGSG25SR3HQ6W - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SEP-YCN2JZZKB3DRPNC4" ipVersion=0x4 line=":KUBE-SEP-YCN2JZZKB3DRPNC4 - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Found forward-reference chainName="KUBE-SVC-TS6C4FBECULI2LCC" ipVersion=0x4 line=":KUBE-SVC-TS6C4FBECULI2LCC - [0:0]" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: In nftables mode, restarting transaction between updates and deletions. ipVersion=0x4 <_>`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Invalidating dataplane cache ipVersion=0x4 reason="refresh timer" table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Loading current iptables state and checking it is correct. ipVersion=0x4 table="nat"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Skipping expected chain <_>ipVersion=0x4 table="filter"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Skipping expected chain chainName="cali-FORWARD" ipVersion=0x4 table="filter"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Skipping expected chain chainName="cali-from-wl-dispatch-b" ipVersion=0x4 table="filter"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Skipping expected chain chainName="cali-pri-_78B28-fZujIjQTQ2aI" ipVersion=0x4 table="filter"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Skipping expected chain chainName="cali-pri-_qr-cFgKHOI4CiiUEEX" ipVersion=0x4 table="filter"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Skipping expected chain chainName="cali-pro-_8C_MHVnZxZL2yzVTdL" ipVersion=0x4 table="filter"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Skipping expected chain chainName="cali-pro-_qr-cFgKHOI4CiiUEEX" ipVersion=0x4 table="filter"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Skipping expected chain chainName="cali-pro-ksa.startup.default" ipVersion=0x4 table="filter"`,
				`<_> [DEBUG][<_>] felix/table.go <_>: Update ended up being no-op, skipping call to ip(6)tables-restore. ipVersion=0x4 <_>`,
				`<_> [DEBUG][<_>] felix/versionparse.go <_>: Parsed kernel version version=<_>.<_>.<_>-1057`,
				`<_> [DEBUG][<_>] felix/versionparse.go <_>: Raw kernel version rawVersion="Linux version <_>.<_>.<_>-1057-a<_>`,
				`<_> [DEBUG][<_>] felix/wireguard.go 1503: Wireguard is disabled and does not exist ifaceName="wireguard.cali" ipVersion=0x4`,
				`<_> [DEBUG][<_>] felix/wireguard.go 654: Wireguard is not in-sync - verifying wireguard configuration is removed ipVersion=0x4`,
				`<_> [DEBUG][<_>] felix/wireguard.go <_>: Queueing a resync of wireguard configuration ipVersion=0x4`,
				`<_> [DEBUG][<_>] felix/wireguard.go <_>: Wireguard is not enabled, skipping sync ipVersion=0x4`,
				`<_> [DEBUG][<_>] felix/xdp_state.go 1043: Processing pending diff state. cs=&intdataplane.xdpSystemState{IfaceNameToData:ma<_>XDPEligiblePolicies:map[proto.PolicyID]intdataplan<_>family=4`,
				`<_> [DEBUG][<_>] felix/xdp_state.go 1270: Finished processing pending diff state. bpfActions=intdataplane.xdpBPFActions{CreateMap:se<_>RemoveMap:set.Typed[string]{}, AddToMap:map[string]map[string]uint32{}, RemoveFromMap:map[string]map[string]uint32{}, InstallXDP:set.Typed[string]{}, UninstallXDP:set.Typed[string]{}, MembersToDrop:map[string]map[string]uint32{}, MembersToAdd:map[string]map[string]uint32{}} family=4 newCS=&intdataplane.xdpSystemState{IfaceNameToData<_>XDPEligiblePolicies:map[proto.PolicyID]intdataplan<_>`,
				`<_> [DEBUG][<_>] felix/xdp_state.go <_>: Finished processing BPF actions. family="ipv4"`,
				`<_> [DEBUG][<_>] felix/xdp_state.go <_>: Getting member changes. family=4 oldMembers=map[string]set.Set[string]{}`,
				`<_> [DEBUG][<_>] felix/xdp_state.go <_>: Processing BPF actions. family="ipv4"`,
				`<_> [DEBUG][<_>] felix/xdp_state.go <_>: Processing member updates. family=4`,
				`<_> [DEBUG][<_>] felix/xdp_state.go <_>: Updating ipsetIDsToMembers cache. family=4`,
				`<_> [INFO][<_>] felix/summary.go 100: Summarising <_> dataplane reconciliation loops over 1m<_>: avg=<_> longest=<_> <_>`,
				`<_> [INFO][<_>] felix/summary.go 100: Summarising <_> dataplane reconciliation loops over <_>: avg=<_> longest=<_> <_>`,
				`<_> [WARNING][56] felix/table.go 654: Detected out-of-sync inserts, marking for resync actualRuleIDs=[]string{"", "", "", "", "", "", "", "", "", "", "", "", "tVnHkvAo15HuiPy0", "", ""} chainName="OUTPUT" expectedRuleIDs=[]string{"tVnHkvAo15HuiPy0", "", "", "", "", "", "", "", "", "", "", "", "", "", ""} ipVersion=0x4 table="raw"`,
				`<_> [WARNING][56] felix/table.go 654: Detected out-of-sync inserts, marking for resync actualRuleIDs=[]string{"", "", "", "", "6gwbT8clXdHdC1b1"} chainName="PREROUTING" expectedRuleIDs=[]string{"6gwbT8clXdHdC1b1", "", "", "", ""} ipVersion=0x4 table="raw"`,
				`<_> [WARNING][56] felix/table.go 654: Detected out-of-sync inserts, marking for resync actualRuleIDs=[]string{"", "", "", "", "Cz_u1IQiXIMmKD4c", "", "", "", "", "", "", "", "", "", "", "", ""} chainName="INPUT" expectedRuleIDs=[]string{"Cz_u1IQiXIMmKD4c", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""} ipVersion=0x4 table="filter"`,
				`<_> [WARNING][56] felix/table.go 654: Detected out-of-sync inserts, marking for resync actualRuleIDs=[]string{"", "", "", "", "tVnHkvAo15HuiPy0", "", "", "", "", ""} chainName="OUTPUT" expectedRuleIDs=[]string{"tVnHkvAo15HuiPy0", "", "", "", "", "", "", "", "", ""} ipVersion=0x4 table="filter"`,
				`bird: Netlink: No route to host`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.inputFile)
			require.NoError(t, err)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			drain := NewWithTokenizer(DefaultConfig(), tt.tokenizer)
			for _, line := range lines {
				drain.Train(line, 0)
			}

			var output []string
			clusters := drain.Clusters()
			for _, cluster := range clusters {
				output = append(output, cluster.String())
			}
			sort.Slice(output, func(i, j int) bool {
				return output[i] < output[j]
			})

			if printUpdatedPatterns {
				for _, out := range output {
					fmt.Printf("`%s`,\n", out)
				}
			}

			require.Equal(t, tt.patterns, output)
		})
	}
}

func TestDrain_TrainGeneratesMatchablePatterns(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		drain      *Drain
		inputLines []string
	}{
		{
			name:  `should match each line against a pattern`,
			drain: New(DefaultConfig()),
			inputLines: []string{
				"test test test",
				"test test test",
				"test test test",
				"test test test",
			},
		},
		{
			name:  `should also match newlines`,
			drain: New(DefaultConfig()),
			inputLines: []string{
				`test test test
`,
				`test test test
`,
				`test test test
`,
				`test test test
`,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for _, line := range tt.inputLines {
				tt.drain.Train(line, 0)
			}
			t.Log(`Learned clusters`, tt.drain.Clusters())

			for _, line := range tt.inputLines {
				match := tt.drain.Match(line)
				require.NotNil(t, match, `Line should match a cluster`)
			}
		})
	}

}

func TestDrain_TrainGeneratesPatternsMatchableByLokiPatternFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		drain      *Drain
		inputLines []string
	}{
		{
			name:  `should extract patterns that all lines match`,
			drain: New(DefaultConfig()),
			inputLines: []string{
				"test 1 test",
				"test 2 test",
				"test 3 test",
				"test 4 test",
			},
		},
		{
			name:  `should extract patterns that match if line ends with newlines`,
			drain: New(DefaultConfig()),
			inputLines: []string{
				`test 1 test
`,
				`test 2 test
`,
				`test 3 test
`,
				`test 4 test
`,
			},
		},
		{
			name:  `should extract patterns that match if line ends with empty space`,
			drain: New(DefaultConfig()),
			inputLines: []string{
				`test 1 test			`,
				`test 2 test			`,
				`test 3 test			`,
				`test 4 test			`,
			},
		},
		{
			name:  `should extract patterns that match if line starts with empty space`,
			drain: New(DefaultConfig()),
			inputLines: []string{
				`			test 1 test`,
				`			test 2 test`,
				`			test 3 test`,
				`			test 4 test`,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for _, line := range tt.inputLines {
				tt.drain.Train(line, 0)
			}
			require.Equal(t, 1, len(tt.drain.Clusters()))
			cluster := tt.drain.Clusters()[0]
			t.Log(`Extracted cluster: `, cluster)

			matcher, err := pattern.ParseLineFilter([]byte(cluster.String()))
			require.NoError(t, err)

			for _, line := range tt.inputLines {
				passes := matcher.Test([]byte(line))
				require.Truef(t, passes, `Line %q should match extracted pattern`, line)
			}
		})
	}

}
