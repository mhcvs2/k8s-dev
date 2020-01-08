<pre><code>

# etcdctl set /confd/dbclud/dns/serial 111
# etcdctl set /confd/dbclud/dns/nodes [\"109.105.1.253\",\"109.105.1.254\",\"109.105.1.246\",\"109.105.1.208\",\"109.105.1.180\",\"109.105.1.176\"]

# cat /etc/confd/confd.toml 
backend = "etcd"
confdir = "/etc/confd"
log-level = "debug"
interval = 600
nodes = [
"http://10.111.0.90:2379"
]
noop = false
prefix = "/confd"
scheme = "http"
watch = true

# cat /etc/confd/conf.d/dbcloud.ksyun.com.zone.toml 
[template]
src = "dbcloud.ksyun.com.zone.tmpl"
dest = "/var/named/dbcloud.ksyun.com.zone"
keys = [
"/dbclud/dns/serial",
"/dbclud/dns/names",
"/dbclud/dns/nodes",
]
check_cmd = "chown named:named {{.src}} && named-checkzone dbcloud.ksyun.com.zone {{.src}}"
reload_cmd = "rndc reload"

# cat /etc/confd/templates/dbcloud.ksyun.com.zone.tmpl 
$TTL 600
@       IN SOA  @ rname.invalid. (
                                        {{ getv "/dbclud/dns/serial" }}       ; serial
                                        1D      ; refresh
                                        1H      ; retry
                                        1W      ; expire
                                        3H )    ; minimum
        NS      @
        A       10.111.0.90
        AAAA    ::1

         IN   NS  ns1
ns1      IN   A   10.111.0.90
{{ $nodes := jsonArray (getv "/dbclud/dns/nodes") }}
{{ range jsonArray (getv "/dbclud/dns/names") }}
{{ . }} {{ range $nodes }} IN A {{ . }}
{{ end }}{{ end }}


</code></pre>