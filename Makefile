NAMESPACE=rhdh

.PHONY: deploy-operator
deploy-operator:
	oc apply -k ./profile

.PHONY: delete-operator
delete-operator:
	oc delete -k ./presets/rhdh-complete
	oc delete -k ./profile/namespace
	oc delete -k ./profile/default-config

.PHONY: deploy-rhdh
deploy-rhdh:
	oc apply -n $(NAMESPACE) -k ./presets/rhdh-complete

.PHONY: delete-rhdh
delete-rhdh:
	oc delete -n $(NAMESPACE) -k ./presets/rhdh-complete

.PHONY: deploy-all
deploy-all:
	$(MAKE) deploy-operator
	@echo "Waiting for operator to be ready ..."
	sleep 10
	$(MAKE) deploy-rhdh

.PHONY: delete-all
delete-all:
	delete-operator
	delete-rhdh

.PHONY: create-namespace
create-namespace:
	oc create namespace $(NAMESPACE)

.PHONY: delete-namespace
delete-namespace:
	oc delete namespace $(NAMESPACE)