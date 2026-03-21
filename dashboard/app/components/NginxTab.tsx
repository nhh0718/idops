"use client";

import {
  Check,
  CheckCircle,
  Copy,
  Download,
  Eye,
  FileCode,
  List,
  Plus,
  Server,
  Shield,
  Trash2,
} from "lucide-react";
import { useEffect, useState } from "react";
import { nginxApi } from "../lib/api";
import type { NginxConfig } from "../types";

const templates = [
  {
    value: "reverse-proxy",
    label: "Reverse Proxy",
    desc: "Upstream proxy sites",
  },
  { value: "static-site", label: "Static Site", desc: "Serve static files" },
  { value: "php-fpm", label: "PHP-FPM", desc: "PHP applications" },
  {
    value: "load-balancer",
    label: "Load Balancer",
    desc: "Upstream load balancing",
  },
  {
    value: "websocket",
    label: "WebSocket Proxy",
    desc: "WebSocket connections",
  },
];

function generateNginxConfig(config: NginxConfig): string {
  const lines: string[] = [];

  if (
    config.template === "load-balancer" &&
    config.backends &&
    config.backends.length > 0
  ) {
    lines.push(`upstream ${config.upstreamName || "backend"} {`);
    if (config.method === "least_conn") lines.push("    least_conn;");
    if (config.method === "ip_hash") lines.push("    ip_hash;");
    for (const b of config.backends) {
      lines.push(
        `    server ${b.host}:${b.port}${b.weight > 1 ? ` weight=${b.weight}` : ""};`,
      );
    }
    lines.push("}\n");
  }

  lines.push("server {");
  lines.push(
    `    listen ${config.listenPort}${config.sslEnabled ? " ssl" : ""};`,
  );
  lines.push(`    server_name ${config.serverName};`);
  lines.push("");

  if (config.sslEnabled) {
    lines.push(
      `    ssl_certificate ${config.sslCertPath || "/etc/ssl/certs/cert.pem"};`,
    );
    lines.push(
      `    ssl_certificate_key ${config.sslKeyPath || "/etc/ssl/private/key.pem"};`,
    );
    lines.push("    ssl_protocols TLSv1.2 TLSv1.3;");
    lines.push("    ssl_ciphers HIGH:!aNULL:!MD5;");
    lines.push("");
  }

  switch (config.template) {
    case "reverse-proxy":
      lines.push("    location / {");
      lines.push(
        `        proxy_pass http://${config.upstreamHost || "127.0.0.1"}:${config.upstreamPort || 3000};`,
      );
      lines.push("        proxy_set_header Host $host;");
      lines.push("        proxy_set_header X-Real-IP $remote_addr;");
      lines.push(
        "        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;",
      );
      lines.push("        proxy_set_header X-Forwarded-Proto $scheme;");
      if (config.webSocket) {
        lines.push("        proxy_http_version 1.1;");
        lines.push("        proxy_set_header Upgrade $http_upgrade;");
        lines.push('        proxy_set_header Connection "upgrade";');
      }
      lines.push("    }");
      break;

    case "static-site":
      lines.push(`    root ${config.rootPath || "/var/www/html"};`);
      lines.push("    index index.html index.htm;");
      lines.push("");
      if (config.enableGzip) {
        lines.push("    gzip on;");
        lines.push(
          "    gzip_types text/plain text/css application/json application/javascript;",
        );
        lines.push("");
      }
      lines.push("    location / {");
      lines.push("        try_files $uri $uri/ =404;");
      lines.push("    }");
      if (config.cacheMaxAge) {
        lines.push("");
        lines.push("    location ~* \\.(js|css|png|jpg|jpeg|gif|ico|svg)$ {");
        lines.push(`        expires ${config.cacheMaxAge}d;`);
        lines.push('        add_header Cache-Control "public, immutable";');
        lines.push("    }");
      }
      break;

    case "php-fpm":
      lines.push(`    root ${config.rootPath || "/var/www/html"};`);
      lines.push("    index index.php index.html;");
      lines.push("");
      lines.push("    location / {");
      lines.push("        try_files $uri $uri/ /index.php?$query_string;");
      lines.push("    }");
      lines.push("");
      lines.push("    location ~ \\.php$ {");
      lines.push(
        `        fastcgi_pass unix:${config.phpSocket || "/run/php/php8.1-fpm.sock"};`,
      );
      lines.push(
        "        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;",
      );
      lines.push("        include fastcgi_params;");
      lines.push("    }");
      break;

    case "load-balancer":
      lines.push("    location / {");
      lines.push(
        `        proxy_pass http://${config.upstreamName || "backend"};`,
      );
      lines.push("        proxy_set_header Host $host;");
      lines.push("        proxy_set_header X-Real-IP $remote_addr;");
      lines.push("    }");
      break;

    case "websocket":
      lines.push("    location / {");
      lines.push(
        `        proxy_pass http://${config.upstreamHost || "127.0.0.1"}:${config.upstreamPort || 3000};`,
      );
      lines.push("        proxy_http_version 1.1;");
      lines.push("        proxy_set_header Upgrade $http_upgrade;");
      lines.push('        proxy_set_header Connection "upgrade";');
      lines.push("        proxy_set_header Host $host;");
      lines.push("        proxy_read_timeout 86400;");
      lines.push("    }");
      break;
  }

  lines.push("}");
  return lines.join("\n");
}

export default function NginxTab({
  configs: initialConfigs,
}: {
  configs: string[];
}) {
  const [config, setConfig] = useState<NginxConfig>({
    template: "reverse-proxy",
    serverName: "example.com",
    listenPort: 80,
    sslEnabled: false,
    sslCertPath: "/etc/ssl/certs/cert.pem",
    sslKeyPath: "/etc/ssl/private/key.pem",
    upstreamHost: "127.0.0.1",
    upstreamPort: 3000,
    webSocket: false,
    rootPath: "/var/www/html",
    phpSocket: "/run/php/php8.1-fpm.sock",
    backends: [{ host: "127.0.0.1", port: 8001, weight: 1 }],
    method: "round-robin",
    upstreamName: "backend",
    enableGzip: true,
    cacheMaxAge: 30,
  });

  const [copied, setCopied] = useState(false);
  const [statusMsg, setStatusMsg] = useState<{
    text: string;
    isError: boolean;
  } | null>(null);
  const [activeSubTab, setActiveSubTab] = useState<
    "generate" | "validate" | "list"
  >("generate");
  const [existingConfigs, setExistingConfigs] =
    useState<string[]>(initialConfigs);
  const [isLoading, setIsLoading] = useState(false);
  const [validationResult, setValidationResult] = useState<{
    valid: boolean;
    message?: string;
    error?: string;
  } | null>(null);

  useEffect(() => {
    if (activeSubTab === "list") {
      loadConfigs();
    } else if (activeSubTab === "validate") {
      loadValidation();
    }
  }, [activeSubTab]);

  async function loadConfigs() {
    setIsLoading(true);
    try {
      const data = await nginxApi.list();
      setExistingConfigs(data);
    } catch {
      showStatus("Failed to load nginx configs", true);
    } finally {
      setIsLoading(false);
    }
  }

  async function loadValidation() {
    try {
      const result = await nginxApi.validate();
      setValidationResult(result);
    } catch {
      showStatus("Failed to validate nginx config", true);
    }
  }

  function showStatus(text: string, isError = false) {
    setStatusMsg({ text, isError });
    setTimeout(() => setStatusMsg(null), 3000);
  }

  function copyConfig() {
    navigator.clipboard.writeText(generateNginxConfig(config));
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  function updateConfig<K extends keyof NginxConfig>(
    key: K,
    value: NginxConfig[K],
  ) {
    setConfig((prev) => ({ ...prev, [key]: value }));
  }

  function addBackend() {
    setConfig((prev) => ({
      ...prev,
      backends: [
        ...(prev.backends || []),
        { host: "127.0.0.1", port: 8000, weight: 1 },
      ],
    }));
  }

  function removeBackend(idx: number) {
    setConfig((prev) => ({
      ...prev,
      backends: (prev.backends || []).filter((_, i) => i !== idx),
    }));
  }

  function updateBackend(idx: number, field: string, value: string | number) {
    setConfig((prev) => ({
      ...prev,
      backends: (prev.backends || []).map((b, i) =>
        i === idx ? { ...b, [field]: value } : b,
      ),
    }));
  }

  const generatedConfig = generateNginxConfig(config);

  return (
    <div className="space-y-4 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-white flex items-center gap-3">
            <Server size={24} className="text-rose-400" />
            Nginx Config Generator
          </h2>
          <p className="text-sm text-[var(--color-muted)] mt-1">
            Generate, validate, and manage nginx configurations from templates
          </p>
        </div>
      </div>

      {/* Status */}
      {statusMsg && (
        <div
          className={`px-4 py-2 rounded-lg text-sm ${statusMsg.isError ? "bg-red-500/10 text-red-400 border border-red-500/20" : "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"}`}
        >
          {statusMsg.text}
        </div>
      )}

      {/* Sub Tabs */}
      <div className="flex gap-1 bg-[var(--color-card)] p-1 rounded-lg border border-[var(--color-border)]">
        {[
          {
            id: "generate" as const,
            label: "Generate",
            icon: <FileCode size={14} />,
          },
          {
            id: "validate" as const,
            label: "Validate",
            icon: <Shield size={14} />,
          },
          {
            id: "list" as const,
            label: "List Configs",
            icon: <List size={14} />,
          },
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveSubTab(tab.id)}
            className={`flex items-center gap-2 px-4 py-2 rounded-md text-xs font-medium transition-all ${
              activeSubTab === tab.id
                ? "bg-[var(--color-primary)] text-white"
                : "text-[var(--color-muted)] hover:text-white hover:bg-[var(--color-card-hover)]"
            }`}
          >
            {tab.icon}
            {tab.label}
          </button>
        ))}
      </div>

      {/* Generate Tab */}
      {activeSubTab === "generate" && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {/* Config Form */}
          <div className="space-y-4">
            {/* Template Selection */}
            <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-4">
              <h4 className="text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider mb-3">
                Template
              </h4>
              <div className="grid grid-cols-1 gap-2">
                {templates.map((t) => (
                  <button
                    key={t.value}
                    onClick={() => updateConfig("template", t.value)}
                    className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-all border ${
                      config.template === t.value
                        ? "bg-[var(--color-primary)]/10 border-[var(--color-primary)]/30 text-white"
                        : "bg-[var(--color-background)] border-[var(--color-border)] text-[var(--color-muted)] hover:text-white hover:border-[var(--color-border)]"
                    }`}
                  >
                    <Server
                      size={14}
                      className={
                        config.template === t.value
                          ? "text-[var(--color-primary)]"
                          : ""
                      }
                    />
                    <div>
                      <p className="text-xs font-medium">{t.label}</p>
                      <p className="text-[10px] text-[var(--color-muted)]">
                        {t.desc}
                      </p>
                    </div>
                  </button>
                ))}
              </div>
            </div>

            {/* Common Settings */}
            <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-4">
              <h4 className="text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider mb-3">
                Common Settings
              </h4>
              <div className="space-y-3">
                <div>
                  <label className="block text-xs text-[var(--color-muted)] mb-1">
                    Server Name (domain)
                  </label>
                  <input
                    type="text"
                    value={config.serverName}
                    onChange={(e) => updateConfig("serverName", e.target.value)}
                    className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                  />
                </div>
                <div>
                  <label className="block text-xs text-[var(--color-muted)] mb-1">
                    Listen Port
                  </label>
                  <input
                    type="number"
                    value={config.listenPort}
                    onChange={(e) =>
                      updateConfig("listenPort", parseInt(e.target.value) || 80)
                    }
                    className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                  />
                </div>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={config.sslEnabled}
                    onChange={(e) =>
                      updateConfig("sslEnabled", e.target.checked)
                    }
                    className="rounded bg-[var(--color-background)] border-[var(--color-border)]"
                  />
                  <span className="text-xs text-white">Enable SSL</span>
                </label>
                {config.sslEnabled && (
                  <>
                    <div>
                      <label className="block text-xs text-[var(--color-muted)] mb-1">
                        SSL Cert Path
                      </label>
                      <input
                        type="text"
                        value={config.sslCertPath}
                        onChange={(e) =>
                          updateConfig("sslCertPath", e.target.value)
                        }
                        className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                      />
                    </div>
                    <div>
                      <label className="block text-xs text-[var(--color-muted)] mb-1">
                        SSL Key Path
                      </label>
                      <input
                        type="text"
                        value={config.sslKeyPath}
                        onChange={(e) =>
                          updateConfig("sslKeyPath", e.target.value)
                        }
                        className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                      />
                    </div>
                  </>
                )}
              </div>
            </div>

            {/* Template-specific fields */}
            {(config.template === "reverse-proxy" ||
              config.template === "websocket") && (
              <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-4">
                <h4 className="text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider mb-3">
                  Proxy Settings
                </h4>
                <div className="space-y-3">
                  <div>
                    <label className="block text-xs text-[var(--color-muted)] mb-1">
                      Upstream Host
                    </label>
                    <input
                      type="text"
                      value={config.upstreamHost}
                      onChange={(e) =>
                        updateConfig("upstreamHost", e.target.value)
                      }
                      className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-[var(--color-muted)] mb-1">
                      Upstream Port
                    </label>
                    <input
                      type="number"
                      value={config.upstreamPort}
                      onChange={(e) =>
                        updateConfig(
                          "upstreamPort",
                          parseInt(e.target.value) || 3000,
                        )
                      }
                      className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                    />
                  </div>
                  {config.template === "reverse-proxy" && (
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={config.webSocket}
                        onChange={(e) =>
                          updateConfig("webSocket", e.target.checked)
                        }
                        className="rounded bg-[var(--color-background)] border-[var(--color-border)]"
                      />
                      <span className="text-xs text-white">
                        Enable WebSocket support
                      </span>
                    </label>
                  )}
                </div>
              </div>
            )}

            {config.template === "static-site" && (
              <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-4">
                <h4 className="text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider mb-3">
                  Static Site Settings
                </h4>
                <div className="space-y-3">
                  <div>
                    <label className="block text-xs text-[var(--color-muted)] mb-1">
                      Root Path
                    </label>
                    <input
                      type="text"
                      value={config.rootPath}
                      onChange={(e) => updateConfig("rootPath", e.target.value)}
                      className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                    />
                  </div>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={config.enableGzip}
                      onChange={(e) =>
                        updateConfig("enableGzip", e.target.checked)
                      }
                      className="rounded bg-[var(--color-background)] border-[var(--color-border)]"
                    />
                    <span className="text-xs text-white">Enable Gzip</span>
                  </label>
                  <div>
                    <label className="block text-xs text-[var(--color-muted)] mb-1">
                      Cache Max Age (days)
                    </label>
                    <input
                      type="number"
                      value={config.cacheMaxAge}
                      onChange={(e) =>
                        updateConfig(
                          "cacheMaxAge",
                          parseInt(e.target.value) || 0,
                        )
                      }
                      className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                    />
                  </div>
                </div>
              </div>
            )}

            {config.template === "php-fpm" && (
              <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-4">
                <h4 className="text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider mb-3">
                  PHP-FPM Settings
                </h4>
                <div className="space-y-3">
                  <div>
                    <label className="block text-xs text-[var(--color-muted)] mb-1">
                      Root Path
                    </label>
                    <input
                      type="text"
                      value={config.rootPath}
                      onChange={(e) => updateConfig("rootPath", e.target.value)}
                      className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-[var(--color-muted)] mb-1">
                      PHP Socket
                    </label>
                    <input
                      type="text"
                      value={config.phpSocket}
                      onChange={(e) =>
                        updateConfig("phpSocket", e.target.value)
                      }
                      className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                    />
                  </div>
                </div>
              </div>
            )}

            {config.template === "load-balancer" && (
              <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-4">
                <h4 className="text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider mb-3">
                  Load Balancer Settings
                </h4>
                <div className="space-y-3">
                  <div>
                    <label className="block text-xs text-[var(--color-muted)] mb-1">
                      Upstream Name
                    </label>
                    <input
                      type="text"
                      value={config.upstreamName}
                      onChange={(e) =>
                        updateConfig("upstreamName", e.target.value)
                      }
                      className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-[var(--color-muted)] mb-1">
                      Method
                    </label>
                    <select
                      value={config.method}
                      onChange={(e) => updateConfig("method", e.target.value)}
                      className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white focus:outline-none focus:border-[var(--color-primary)]"
                    >
                      <option value="round-robin">Round Robin</option>
                      <option value="least_conn">Least Connections</option>
                      <option value="ip_hash">IP Hash</option>
                    </select>
                  </div>
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <label className="text-xs text-[var(--color-muted)]">
                        Backends
                      </label>
                      <button
                        onClick={addBackend}
                        className="flex items-center gap-1 text-xs text-[var(--color-primary)] hover:text-[var(--color-primary-light)]"
                      >
                        <Plus size={12} /> Add
                      </button>
                    </div>
                    <div className="space-y-2">
                      {(config.backends || []).map((b, i) => (
                        <div key={i} className="flex items-center gap-2">
                          <input
                            type="text"
                            value={b.host}
                            onChange={(e) =>
                              updateBackend(i, "host", e.target.value)
                            }
                            placeholder="Host"
                            className="flex-1 px-2 py-1.5 bg-[var(--color-background)] border border-[var(--color-border)] rounded text-xs text-white focus:outline-none focus:border-[var(--color-primary)]"
                          />
                          <input
                            type="number"
                            value={b.port}
                            onChange={(e) =>
                              updateBackend(
                                i,
                                "port",
                                parseInt(e.target.value) || 80,
                              )
                            }
                            placeholder="Port"
                            className="w-20 px-2 py-1.5 bg-[var(--color-background)] border border-[var(--color-border)] rounded text-xs text-white focus:outline-none focus:border-[var(--color-primary)]"
                          />
                          <input
                            type="number"
                            value={b.weight}
                            onChange={(e) =>
                              updateBackend(
                                i,
                                "weight",
                                parseInt(e.target.value) || 1,
                              )
                            }
                            placeholder="Weight"
                            className="w-16 px-2 py-1.5 bg-[var(--color-background)] border border-[var(--color-border)] rounded text-xs text-white focus:outline-none focus:border-[var(--color-primary)]"
                          />
                          <button
                            onClick={() => removeBackend(i)}
                            className="p-1 text-red-400 hover:bg-red-500/10 rounded"
                          >
                            <Trash2 size={12} />
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Preview Panel */}
          <div className="space-y-4">
            <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl overflow-hidden sticky top-4">
              <div className="flex items-center justify-between px-4 py-3 border-b border-[var(--color-border)]">
                <h4 className="text-xs font-semibold text-white flex items-center gap-2">
                  <Eye size={14} className="text-blue-400" />
                  Preview
                </h4>
                <div className="flex items-center gap-2">
                  <button
                    onClick={copyConfig}
                    className="flex items-center gap-1 px-2 py-1 rounded text-xs text-[var(--color-muted)] hover:text-white hover:bg-[var(--color-card-hover)] transition-colors"
                  >
                    {copied ? (
                      <Check size={12} className="text-emerald-400" />
                    ) : (
                      <Copy size={12} />
                    )}
                    {copied ? "Copied!" : "Copy"}
                  </button>
                  <button
                    onClick={() =>
                      showStatus(`Saved to ${config.serverName}.conf`)
                    }
                    className="flex items-center gap-1 px-2 py-1 rounded text-xs text-[var(--color-muted)] hover:text-white hover:bg-[var(--color-card-hover)] transition-colors"
                  >
                    <Download size={12} />
                    Save
                  </button>
                </div>
              </div>
              <pre className="p-4 overflow-auto max-h-[60vh] text-xs font-mono text-emerald-300 leading-relaxed bg-[#0d1117]">
                <code>{generatedConfig}</code>
              </pre>
            </div>
          </div>
        </div>
      )}

      {/* Validate Tab */}
      {activeSubTab === "validate" && (
        <div className="space-y-4">
          <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-5">
            <h4 className="text-sm font-semibold text-white mb-3 flex items-center gap-2">
              <Shield size={14} className="text-emerald-400" />
              Validate Nginx Configuration
            </h4>
            <p className="text-xs text-[var(--color-muted)] mb-4">
              Run{" "}
              <code className="px-1 py-0.5 bg-[var(--color-background)] rounded text-white">
                nginx -t
              </code>{" "}
              to check configuration syntax.
            </p>
            <button
              onClick={() => showStatus("Nginx configuration is valid")}
              className="px-4 py-2 bg-emerald-500/20 border border-emerald-500/30 rounded-lg text-sm text-emerald-400 hover:bg-emerald-500/30 transition-colors"
            >
              Run Validation
            </button>

            <div className="mt-4 p-3 bg-[var(--color-background)] rounded-lg font-mono text-xs">
              <p className="text-emerald-400">$ nginx -t</p>
              <p className="text-[var(--color-muted)]">
                nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
              </p>
              <p className="text-[var(--color-muted)]">
                nginx: configuration file /etc/nginx/nginx.conf test is
                successful
              </p>
              <p className="text-emerald-400 mt-2 flex items-center gap-1">
                <CheckCircle size={12} /> Configuration valid
              </p>
            </div>
          </div>
        </div>
      )}

      {/* List Configs Tab */}
      {activeSubTab === "list" && (
        <div className="space-y-4">
          <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-5">
            <h4 className="text-sm font-semibold text-white mb-3 flex items-center gap-2">
              <List size={14} className="text-blue-400" />
              Nginx Configs — /etc/nginx/sites-available
            </h4>
            <div className="space-y-2">
              {existingConfigs.map((conf) => (
                <div
                  key={conf}
                  className="flex items-center justify-between py-2.5 px-3 bg-[var(--color-background)] rounded-lg border border-[var(--color-border)]"
                >
                  <div className="flex items-center gap-2">
                    <FileCode size={14} className="text-[var(--color-muted)]" />
                    <span className="font-mono text-xs text-white">{conf}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() =>
                        showStatus(`Applied ${conf} and reloaded nginx`)
                      }
                      className="text-xs px-2 py-1 rounded bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20 transition-colors"
                    >
                      Apply
                    </button>
                    <button className="text-xs px-2 py-1 rounded bg-blue-500/10 text-blue-400 hover:bg-blue-500/20 transition-colors">
                      View
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
