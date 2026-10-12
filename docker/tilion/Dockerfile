# syntax=docker/dockerfile:1

ARG TILION_ENGINE_IMAGE=tilion/fortress:latest

FROM ${TILION_ENGINE_IMAGE} AS builder
USER root
ARG TILION_EXTRAS=mcp,vision
RUN apt-get update \
    && apt-get install -y --no-install-recommends python3-venv \
    && rm -rf /var/lib/apt/lists/* \
    && python3 -m venv /opt/venv \
    && /opt/venv/bin/pip install --no-cache-dir --upgrade pip \
    && /opt/venv/bin/pip install --no-cache-dir "tilion[${TILION_EXTRAS}]"

FROM ${TILION_ENGINE_IMAGE} AS runtime
USER root
COPY --from=builder /opt/venv /opt/venv
COPY docker/tilion/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh \
    && SP="$(/opt/venv/bin/python -c 'import site; print(site.getsitepackages()[0])')" \
    && if [ -d "$SP/tillion_fortress" ] && [ ! -e "$SP/tilion_fortress" ]; then ln -s tillion_fortress "$SP/tilion_fortress"; fi \
    && VER="$(/opt/venv/bin/python -c 'import tilion_fortress as t; print(t.__version__)')" \
    && PLAT="$(/opt/venv/bin/python -c 'import tilion_fortress as t; print(t.resolve_platform())')" \
    && mkdir -p "/home/tilion/.cache/tillion-fortress/${VER}/${PLAT}" \
    && ln -sf /opt/tilion/tilion "/home/tilion/.cache/tillion-fortress/${VER}/${PLAT}/tilion" \
    && chown -R tilion:tilion /home/tilion
ENV PATH="/opt/venv/bin:${PATH}"
EXPOSE 9222
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
    CMD curl -fsS http://127.0.0.1:9222/json/version >/dev/null || exit 1
USER tilion
ENTRYPOINT ["dumb-init", "--", "/usr/local/bin/entrypoint.sh"]
