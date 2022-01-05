[comment]: # ( Copyright Red Hat )
# idp-configs-api

Service to store and retrieve Identity Provider configurations. 

---
### Basic setup:

1. Install Golang 1.16
2. Run `make run-migrate`
3. Run `make run`
4. Access the API at http://localhost:3000/api/idp-configs-api/v0/openapi.json
   (All other routes are authenticated and require the X-RH-Identity request header to be set)

---
### Setup with ephemeral clusters
1. See instructions [here](https://clouddot.pages.redhat.com/docs/dev/getting-started/ephemeral/onboarding.html) for onboarding an ephemeral cluster.
2. Using the instructions in the step above, reserve a namespace on the ephemeral cluster.
3. If you're running this from within the cloned `idp-configs-api` directory, run `make bonfire-config-github` to generate the bonfire config.yml. But if you're in a separate directory, edit the bonfire config (`~/.config/bonfire/config.yaml`) with the contents of `default_config.yaml.github.example`.
4. After successfully reserving a namespace (eg. ephemeral-10), use this namespace to deploy the app:
    * If you're within the `idp-configs-api` directory:
      ```
      make deploy-app NAMESPACE=<reserved_ephemeral_namespace>
      ```
      Otherwise run:
      ```
      bonfire deploy idp-configs -n <reserved_ephemeral_namespace>
      ```
    * Now the application should be running. You can check the pod with `oc get pods -n <reserved_ephemeral_namespace>` You can test by port-forwarding the app in one terminal and running a curl command in another 
      ```
        oc port-forward service/idp-configs-api-service 8000:8000 -n <reserved_ephemeral_namespace>
      ```
      In a separate terminal:
      ```
      curl -v http://localhost:8000/api/idp-configs-api/v0/openapi.json
      ```
      Note: To test authenticated endpoints, an X-RH-Identity header needs to be set on the request
    * When you're done testing, release your ephemeral namespace:
      ```
      bonfire namespace release <reserved_ephemeral_namespace>
      ```

---

### Setup with Minikube 
*Note: Testing on minikube can be a bit flaky and less reliable than testing on ephemeral clusters due to resource constraints.*

This setup utilizes the following tools. Follow the steps in the links for installation:
- [minikube](https://minikube.sigs.k8s.io/docs/)
- [Clowder](https://github.com/RedHatInsights/clowder)
- [Bonfire](https://github.com/RedHatInsights/bonfire) (see installation instructions for setting up the virtual environment for bonfire)

Following the information above you should have Docker, a minikube cluster running with Clowder installed, and a Python environment with `bonfire` installed. Now move on to running the `idp-configs-api` application.

1. Clone the project.
```
git clone git@github.com:identitatem/idp-configs-api.git
```
2. Change directories to the project.
```
cd idp-configs-api
```
3. Setup access to the Docker enviroment within minikube, so you can build images directly to the cluster's registry.
```
eval $(minikube -p minikube docker-env)
```
4. Build the container image.
```
make docker-build
```
5. Create Bonfire configuration. To deploy from your local repository run the following:
```
make bonfire-config-local
```
The above command will create a file named `config.yaml` pointing to your local repository. At times you may need to update the branch which is referred to with the `ref` parameter (defaults to `main`).

Bonfire can also deploy from GitHub. Running the following command will setup the GitHub based configuration:
```
make bonfire-config-github
```
6. Setup a test namespace for deployment.
```
make create-ns NAMESPACE=test
```
7. Deploy a Clowder environment (*ClowdEnviroment*) to the namespace with bonfire.
```
make deploy-env NAMESPACE=test
```
8. Deploy the application to the namespace.
```
make deploy-app NAMESPACE=test
```

Now the application should be running. You can test this by port-forwarding the app in one terminal and running a curl command in another as follows:

**Terminal 1**
```
kubectl -n test port-forward service/idp-configs-api-service 8000:8000
```
**Terminal 2**
```
curl -v http://localhost:8000/
```

You should get a 200 response.

---
### Testing and Linting:

Unit tests:
```
make test
```
Lint:
```
go get -u honnef.co/go/tools/cmd/staticcheck@latest
make lint
```

---
#### Viewing the OpenAPI 3.0 spec:

* Run `make generate-docs`

The API will serve the docs under a `/docs` endpoint.
test
