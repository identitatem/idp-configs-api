#!/bin/bash

echo "os: $OSTYPE"
echo "shell: $SHELL"
export PATH=$PATH:$PWD

# --------------------------------------------
# Options that must be configured by app owner
# --------------------------------------------
APP_NAME="idp-configs"  # name of app-sre "application" folder this component lives in
COMPONENT_NAME="idp-configs-api"  # name of app-sre "resourceTemplate" in deploy.yaml for this component
IMAGE="quay.io/cloudservices/idp-configs-api"

# IQE_PLUGINS="idp-configs-api"
# IQE_MARKER_EXPRESSION="smoke"
# IQE_FILTER_EXPRESSION=""

echo "LABEL quay.expires-after=3d" >> ./Dockerfile # tag expire in 3 days

# Install bonfire repo/initialize
CICD_URL=https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd
curl -s $CICD_URL/bootstrap.sh > .cicd_bootstrap.sh && source .cicd_bootstrap.sh

source $CICD_ROOT/build.sh
source $CICD_ROOT/deploy_ephemeral_env.sh
#source $CICD_ROOT/smoke_test.sh

mkdir -p $WORKSPACE/artifacts
cat << EOF > ${WORKSPACE}/artifacts/junit-dummy.xml
<testsuite tests="1">
    <testcase classname="dummy" name="dummytest"/>
</testsuite>
EOF

echo "** pr_check.sh - services"
oc get service -n $NAMESPACE
# echo "** port forward"
# oc port-forward svc/idp-configs-api-service 8000:8000 -n $NAMESPACE
echo "** curl"
curl -v http://idp-configs-api-service:8000/
echo "** end pr_check.sh"
