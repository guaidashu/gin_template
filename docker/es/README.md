# 配置ES 和 Kibana密码

1. 注释es配置文件 yml 文件内的

        xpack.security.transport.ssl.keystore.path: /usr/share/elasticsearch/config/elastic-certificates.p12
        xpack.security.transport.ssl.truststore.path: /usr/share/elasticsearch/config/elastic-certificates.p12
        xpack.security.transport.ssl.keystore.password: password
        xpack.security.transport.ssl.truststore.password: password

2. 进入容器

        进入容器
        docker exec -it es-master(或容器ID) /bin/bash

3. 生成证书

        ./bin/elasticsearch-certutil cert -out config/elastic-certificates.p12 -pass "证书密码"

4. 拷贝证书到所有容器中

        docker cp es-master:/usr/share/elasticsearch/config/elastic-certificates.p12 ./config

5. 证书改为 777 权限

        每个容器都要搞
        chomod 777 /usr/share/elasticsearch/config/elastic-certificates.p12

6. 打开1中的注释

        xpack.security.transport.ssl.keystore.password: password
        xpack.security.transport.ssl.truststore.password: password
        password为3生成证书时候的密码

7. 进入master容器并设置密码

        docker exec -it es-master /bin/bash
        自动创建密码(这个生成的密码一定要记下来)：
        ./bin/elasticsearch-setup-passwords auto
        ./bin/elasticsearch-setup-passwords interactive

8. kibana.yml设置密码为刚刚设置的密码即可