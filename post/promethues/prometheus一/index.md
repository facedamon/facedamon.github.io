
# 摘要
&emsp;&emsp;**由于该系统涉及较多的概念和技术，为了快速上手，第一篇文档主要讲解环境搭建。后期系列文档会逐一解释每一个关键的结点。**
</br>
# 系统架构
</br>

![avator](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/architecture.png)
</br>

# 简略解释
</br>

&emsp;&emsp;**Prometheus**是一个开放性的监控解决方案，用户可以非常方便的安装和使用Prometheus并且能够非常方便的对其进行扩展。为了能够更加直观的了解Prometheus Server，接下来我们将在本地部署并运行一个Prometheus Server实例，通过**Node Exporter**采集当前主机的系统资源使用情况。 并通过**Grafana**创建一个简单的可视化仪表盘。
&emsp;&emsp;**Prometheus Server**是Prometheus组件中的核心部分，负责实现对监控数据的获取，存储以及查询。 Prometheus Server可以通过静态配置管理监控目标，也可以配合使用**Service Discovery**的方式动态管理监控目标，并从这些监控目标中获取数据。其次Prometheus Server需要对采集到的监控数据进行存储，Prometheus Server本身就是一个时序数据库，将采集到的监控数据按照时间序列的方式存储在本地磁盘当中。最后Prometheus Server对外提供了自定义的**PromQL**语言，实现对数据的查询以及分析。
&emsp;&emsp;Prometheus Server内置的**Express Browser UI**，通过这个UI可以直接通过PromQL实现数据的查询以及可视化。
&emsp;&emsp;Prometheus Server的集群能力可以使其从其他的Prometheus Server实例中获取数据，因此在大规模监控的情况下，可以通过集群以及功能分区的方式对Prometheus Server进行扩展。
&emsp;&emsp;**Exporter**将监控数据采集的端点通过HTTP服务的形式暴露给Prometheus Server，Prometheus Server通过访问该Exporter提供的Endpoint端点，即可获取到需要采集的监控数据。
&emsp;&emsp;一般来说可以将Exporter分为2类：
- 直接采集：这一类Exporter直接内置了对Prometheus监控的支持，比如cAdvisor，Kubernetes，**Etcd**(最好用)，Gokit等，都直接内置了用于向Prometheus暴露监控数据的端点。
- 间接采集：间接采集，原有监控目标并不直接支持Prometheus，因此我们需要通过Prometheus提供的Client Library编写该监控目标的监控采集程序。例如： **Mysqld Exporter，JMX Exporter，Nginx Exporter**等。
&emsp;&emsp;在Prometheus Server中支持基于PromQL创建告警规则，如果满足PromQL定义的规则，则会产生一条告警，而告警的后续处理流程则由AlertManager进行管理。在AlertManager中我们可以与**邮件**，Slack等等内置的通知方式进行集成，也可以通过**Webhook**自定义告警处理方式。AlertManager即Prometheus体系中的告警处理中心。
</br>

# 环境准备

</br>

		# useradd -d /home/ddd -s /bin/bash ddd
		# su - ddd
		$ wget https://github.com/prometheus/prometheus/releases/download/v2.19.2/prometheus-2.19.2.linux-amd64.tar.gz
		$ wget https://github.com/prometheus/alertmanager/releases/download/v0.21.0/alertmanager-0.21.0.linux-amd64.tar.gz
		$ wget https://github.com/prometheus/node_exporter/releases/download/v1.0.1/node_exporter-1.0.1.linux-amd64.tar.gz
		$ for i in (ls *.tar.gz);do tar -zxvf $i;done
		$ wget https://dl.grafana.com/oss/release/grafana-7.0.4-1.x86_64.rpm
		$ sudo yum install grafana-7.0.4-1.x86_64.rpm

</br>

# prometheus 配置
</br>

		$ cd /home/ddd/prometheus-2.19.1.linux-amd64
		$ vim prometheus.yml
		# my global config
		global:
		  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
		  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
		  # scrape_timeout is set to the global default (10s).
		# Alertmanager configuration
		alerting:
		  alertmanagers:
		  - static_configs:
		    - targets: ['127.0.0.1:9093']
		      # - alertmanager:9093
		# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
		rule_files:
		  # - "first_rules.yml"
		  # - "second_rules.yml"
		  - /home/ddd/prometheus-2.19.1.linux-amd64/rules/*.rules
		# A scrape configuration containing exactly one endpoint to scrape:
		# Here it's Prometheus itself.
		scrape_configs:
		  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
		  - job_name: 'prometheus'
		    # metrics_path defaults to '/metrics'
		    # scheme defaults to 'http'.
		    scrape_interval: 5s
		    static_configs:
		    - targets: ['127.0.0.1:9090']
		  # operating system
		  - job_name: 'OS'
		    scrape_interval: 3s
		    static_configs:
		        - targets: ['127.0.0.1:9100']
		  # mysqld
		  # - job_name: 'mysqld'
		  #  static_configs:
		  #      - targets: ['127.0.0.1:9104']
		$ mkdir rules && cd rules
		$ vim hoststats-alert.rules
		 groups:
		 - name: memoryrules
		   rules:
		   - alert: MemoryUsageTooHigh
		     expr: (1-node_memory_MemAvailable_bytes/node_memory_MemTotal_bytes) * 100 > 85
		     for: 5m
		     labels:
		       #team: node
		       severity: warning
		       action: text
		       kind: node_export
		     annotations:
		       summary: "{{$labels.instance}}: Memory usage is too high"
		       description: "{{$labels.instance}}: job {{$labels.job}} memory usage is more than 85% last 5min "
		       value: "{{$value}}"
		 - name: cpurules
		   rules:
		   - alert: CpuUsageTooHigh
		     expr: 100-avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by(instance)*100 > 80
		     for: 1m
		     labels:
		       severity: warning
		     annotations:
		       summary: "Instance {{ $labels.instance }} cpu使用率过高"
		       description: "{{ $labels.instance }} of job {{ $labels.job }}cpu使用率超过80%,当前使用率[{{ $value }}]."
</br>

# AlertManager 配置
</br>

		$ cd /home/ddd/alertmanager-0.21.0.linux-amd64
		$ vim alertmanager.yml
		global:
		  resolve_timeout: 5m
		  smtp_smarthost: smtp.263.net:25
		  smtp_from: th@ddddian.com
		  smtp_auth_username: th@ddddian.com
		  smtp_auth_identity: th@ddddian.com
		  smtp_auth_password: ******
		templates:
		  - 'template/*.tmpl'
		route:
		  group_by: ['alertname']
		  group_wait: 10s
		  group_interval: 10s
		  repeat_interval: 1m
		  receiver: 'default-receiver'
		  routes:
		  - receiver: 'os-pager'
		    group_wait: 10s
		    match_re:
		      service: cpu|mem|network
		  - receiver: 'database-pager'
		    group_wait: 10s
		    match_re:
		      service: mysql|mariadb
		 # - receiver: 'frontend-pager'
		 #   group_by: [product, environment]
		 #   match:
		 #     team: frontend
		receivers:
		  - name: 'default-receiver'
		    email_configs:
		    - to: 'tsq@ddddian.com'
		      require_tls: false
		      html: '{{ template "test.html" . }}'
		      headers: { Subject: "[WARNING] Alertmanager 告警提醒" }
		  - name: 'os-pager'
		    email_configs:
		    - to: 'tsq@ddddian.com'
		      require_tls: false
		   #   html:
		   #   headers:
		  - name: 'database-pager'
		    email_configs:
		    - to: 'tsq@ddddian.com'
		      require_tls: false
		   #   html:
		   #   headers:
		inhibit_rules:
		  - source_match:
		       severity: 'critical'
		    target_match:
		       severity: 'warning'
		    equal: ['alertname', 'dev', 'instance']
		$ mkdir template && cd template
		$ vim hoststats-alert.tmpl
		{{ define "test.html" }}
		 <table border="1">
		         <tr>
		                 <td>报警项</td>
		                 <td>实例IP</td>
		                 <td>报警阀值</td>
		                 <td>开始时间</td>
		         </tr>
		         {{ range $i, $alert := .Alerts }}
		                 <tr>
		                         <td>{{ index $alert.Labels "alertname" }}</td>
		                         <td>{{ index $alert.Labels "instance" }}</td>
		                         <td>{{ index $alert.Annotations "description" }}</td>
		                         <td>{{ ($alert.StartsAt.Add 28800e9).Format "2006-02-02 15:04:05" }}</td>
		                 </tr>
		         {{ end }}
		 </table>
		 {{ end }}
</br>

# 启动服务
</br>

		$ cd /home/ddd/prometheus-2.19.1.linux-amd64
		$ nohup ./prometheus &
		$ cd /home/ddd/alertmanager-0.21.0.linux-amd64
		$ nohup ./alertmanager &
		$ cd /home/ddd/node_exporter-1.0.1.linux-amd64
		$ nohup ./node_exporter
		$ sudo systemctl status grafana-server.service
		$ sudo systemctl restart grafana-server.service
		$ sudo firewall-cmd --zone=public --add-port=3000/tcp --permanent
		$ sudo firewall-cmd --zone=public --add-port=9090/tcp --permanent
		$ sudo firewall-cmd --reload
</br>

# 面板配置
&emsp;&emsp;登入前台系统`http://xxx:3000`，默认用户/密码：admin/123456
![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/4668b31c615a7536fcfac05f19a2a3a.png)
</br>

## 配置数据源
</br>

![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/072978f6b21acaa7cca66a179df2bb5.png)
</br>

![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/1d4fe0c1ea5e8eec1038a6ac79ac511.png)
</br>

![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/1342122957529bf8fb5406129ea2adc.png)
</br>

## 配置面板
</br>

&emsp;&emsp;这里的**8919**编号是grafana官网给出的面板主题编号，你可以在官网下载自己喜欢的面板样式。
</br>

![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/874251548bbd2e1e6ae3cf71df720d5.png)
</br>

## 主页配置
</br>

![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/fcc22e5936aded76abe5da97b00853d.png)
</br>
&emsp;&emsp;**点击收藏**

![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/71094e1f412c50c774960f47cc608ee.png)
</br>
&emsp;&emsp;**回到个人设置页面**

![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/f22efb03f5f090c2d7bb02fc2f2a26d.png)

> **点击左上角主页按钮**

![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/9ea46edb06ea3970640df7daa30ee79.png)

&emsp;&emsp;**grafana最新版也自带了Alert&Rules组件，但是功能太过单一，最主要的是没有通知策略，所以我们还是选择prometheus的AlertManager组件作为通知载体。**

# 测试预警

&emsp;&emsp;登入该主机，手动执行以下shell来拉高CPU占用率。（该语句的含义是遍历主机的每一个CPU，让其都循环加载磁盘。请勿长时间运行）
		$ for i in `seq 1 $(cat /proc/cpuinfo |grep "physical id" |wc -l)`; do dd if=/dev/zero of=/dev/null & done
</br>

&emsp;&emsp;检查一下配置的接收邮箱。

</br>
![avaror](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/prometheus/20200629/5a7aec9ad3df19ad19d91735f05d479.png)

&emsp;&emsp;也许你发现了问题，开始时间不对，开始时间采用的是UTC，差了8个小时。我们后面会修改这一问题。还有，你可能觉得不太美观，没关系，后面可以自己编写html来美化。下一期我们将继续讲解各个组件的配置，以此来解决通知的**沉默、抑制、合并**，Alertmanager的抑制机制可以避免当某种问题告警产生之后用户接收到大量由此问题导致的一系列的其它告警通知。例如当集群不可用时，用户可能只希望接收到一条告警，告诉他这时候集群出现了问题，而不是大量的如集群中的应用异常、中间件服务异常的告警通知。

> 别忘了关闭拉高CPU的线程

		$ killall dd
