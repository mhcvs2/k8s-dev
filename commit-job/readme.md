<pre><code>
# helm install commit-job --set targetHostname=k8s-n1,podName=tiller-deploy-8694f8fddc-c2rql,imageName=registry.bst-1.cns.bstjpc.com:5000/test-name,imageTag=test-tag,author="mhc haha",message="hello world"  --name=test-job-id
NAME:   test-job-id
LAST DEPLOYED: Sat May 12 17:07:07 2018
NAMESPACE: default
STATUS: DEPLOYED

RESOURCES:
==> v1/Job
NAME                    DESIRED  SUCCESSFUL  AGE
test-job-id-commit-job  1        0           2s

==> v1/Pod(related)
NAME                          READY  STATUS             RESTARTS  AGE
test-job-id-commit-job-zm9f8  0/1    ContainerCreating  0         1s
</code></pre>
<pre><code>
# helm status test-job-id
LAST DEPLOYED: Sat May 12 17:07:07 2018
NAMESPACE: default
STATUS: DEPLOYED

RESOURCES:
==> v1/Job
NAME                    DESIRED  SUCCESSFUL  AGE
test-job-id-commit-job  1        1           19s

==> v1/Pod(related)
NAME                          READY  STATUS     RESTARTS  AGE
test-job-id-commit-job-zm9f8  0/1    Completed  0         18s
</code></pre>

<pre><code>
# helm delete test-job-id --purge
release "test-job-id" deleted

</code></pre>