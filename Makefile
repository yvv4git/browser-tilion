SHELL := /bin/bash
.DEFAULT_GOAL := help

ENV_FILE := .env
URL ?= https://example.com
Q ?= tilion stealth browser

-include $(ENV_FILE)
export

COMPOSE := docker compose

.PHONY: help env build rebuild up down restart logs ps config shell cdp clean \
	install install-mcp install-vision activate fetch search crawl screenshot

help:
	@echo "browser-tilion"
	@echo ""
	@echo "  make env                 create .env from .env.example if missing"
	@echo "  make build               build the image"
	@echo "  make rebuild             rebuild the image without cache"
	@echo "  make up                  start the browser"
	@echo "  make down                stop the browser"
	@echo "  make restart             restart the browser"
	@echo "  make logs                follow logs"
	@echo "  make ps                  show containers"
	@echo "  make config              render compose config"
	@echo "  make shell               open a shell in the container"
	@echo "  make cdp                 print the CDP endpoint"
	@echo "  make clean               stop and remove the image"
	@echo ""
	@echo "  make install             pip install tilion (core)"
	@echo "  make install-mcp         pip install tilion[mcp]"
	@echo "  make install-vision      pip install tilion[vision]"
	@echo "  make activate            activate the license in the container"
	@echo ""
	@echo "  make fetch URL=...       fetch a page"
	@echo "  make search Q=...        web search"
	@echo "  make crawl URL=...       crawl a site"
	@echo "  make screenshot URL=...  screenshot a page"

env:
	@test -f $(ENV_FILE) || cp .env.example $(ENV_FILE)

build:
	$(COMPOSE) build

rebuild:
	$(COMPOSE) build --no-cache

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

restart:
	$(COMPOSE) restart

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

config:
	$(COMPOSE) config

shell:
	$(COMPOSE) exec tilion bash

cdp:
	curl -fsS http://127.0.0.1:$(or $(TILION_PORT),9222)/json/version

clean:
	$(COMPOSE) down --rmi local --remove-orphans

install:
	pip install tilion

install-mcp:
	pip install "tilion[mcp]"

install-vision:
	pip install "tilion[vision]"

activate:
	$(COMPOSE) exec tilion tilion activate --token "$(TILION_API_KEY)"

fetch:
	$(COMPOSE) exec -T tilion python3 - fetch "$(URL)" < docker/tilion/smoke.py

search:
	$(COMPOSE) exec -T tilion python3 - search "$(Q)" < docker/tilion/smoke.py

crawl:
	$(COMPOSE) exec -T tilion python3 - crawl "$(URL)" < docker/tilion/smoke.py

screenshot:
	$(COMPOSE) exec -T tilion python3 - screenshot "$(URL)" < docker/tilion/smoke.py
