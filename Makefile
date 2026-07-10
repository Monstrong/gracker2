.PHONY: proto proto-clean $(addprefix proto-,$(SERVICES))

# да, этот файл я сгенерил в нейронке

MODULE := github.com/monstrong/gracker2
PROTO_OPTS := --go_out=. --go_opt=module=$(MODULE) --go-grpc_out=. --go-grpc_opt=module=$(MODULE)

# Список сервисов
SERVICES := torrent auth

# Версия по умолчанию - последняя (находим максимальную)
LATEST_VERSION := $(shell ls -d proto/*/v* 2>/dev/null | sed 's/.*\/v//' | sort -n | tail -1)
VERSION ?= v$(LATEST_VERSION)

# Если VERSION не задана и нет папок, ставим v1
ifeq ($(VERSION),)
VERSION := v1
endif

# Все proto файлы
PROTO_FILES := $(foreach s,$(SERVICES),proto/$(s)/$(VERSION)/$(s).proto)

proto: $(PROTO_FILES)
	@echo "Генерация для версии $(VERSION)"
	@protoc -I. $(PROTO_OPTS) $^

# Генерация для конкретного сервиса
$(addprefix proto-,$(SERVICES)): proto-%:
	@echo "Генерация $* для версии $(VERSION)"
	@protoc -I. $(PROTO_OPTS) proto/$*/$(VERSION)/$*.proto

proto-clean:
	@rm -rf proto/gen/go

proto-versions:
	@echo "Доступные версии:"
	@ls -d proto/*/v* 2>/dev/null | sed 's/.*\/v//' | sort -n | uniq