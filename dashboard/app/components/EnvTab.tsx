"use client";

import {
  AlertTriangle,
  ArrowRight,
  Check,
  CheckCircle,
  Copy,
  Eye,
  EyeOff,
  FileCode,
  FileDown,
  GitCompare,
  RefreshCw,
  Search,
  Shield,
} from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { envApi } from "../lib/api";
import { useI18n } from "../lib/i18n";
import type { EnvValidationIssue, EnvVariable } from "../types";

interface EnvTabProps {
  envVars: EnvVariable[];
  setEnvVars?: React.Dispatch<React.SetStateAction<EnvVariable[]>>;
}

export default function EnvTab({
  envVars: initialEnvVars,
}: {
  envVars: EnvVariable[];
}) {
  const { t } = useI18n();
  const [envVars, setEnvVars] = useState<EnvVariable[]>(initialEnvVars);
  const [activeSubTab, setActiveSubTab] = useState<
    "show" | "compare" | "validate" | "sync" | "init"
  >("show");
  const [showSecrets, setShowSecrets] = useState(false);
  const [filter, setFilter] = useState("");
  const [envPath, setEnvPath] = useState("");
  const [statusMsg, setStatusMsg] = useState<{
    text: string;
    isError: boolean;
  } | null>(null);
  const [copied, setCopied] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [compareOutput, setCompareOutput] = useState("");
  const [validationResult, setValidationResult] = useState<{
    valid: boolean;
    output?: string;
    error?: string;
  } | null>(null);

  // Validation mock data
  const [validationIssues] = useState<EnvValidationIssue[]>([
    {
      line: 5,
      key: "APP_NAME",
      type: "unquoted_spaces",
      message: "value has spaces but is not quoted",
    },
    { line: 12, key: "API_KEY", type: "empty", message: "empty value" },
    {
      line: 18,
      key: "DB_HOST",
      type: "duplicate",
      message: "duplicate key (first at line 3)",
    },
    {
      line: 25,
      key: "SECRET_KEY",
      type: "trailing_space",
      message: "trailing whitespace in value",
    },
  ]);

  const loadEnvVarsCb = useCallback(loadEnvVars, [envPath]);
  const loadCompareCb = useCallback(loadCompare, [envPath]);
  const loadValidateCb = useCallback(loadValidate, [envPath]);

  useEffect(() => {
    if (activeSubTab === "show") {
      loadEnvVarsCb();
    } else if (activeSubTab === "compare") {
      loadCompareCb();
    } else if (activeSubTab === "validate") {
      loadValidateCb();
    }
  }, [activeSubTab, loadEnvVarsCb, loadCompareCb, loadValidateCb]);

  async function loadEnvVars() {
    setIsLoading(true);
    try {
      const data = await envApi.show(envPath || ".env");
      const vars: EnvVariable[] = Object.entries(data).map(([key, value]) => ({
        key,
        value: value as string,
        isSensitive:
          key.toLowerCase().includes("key") ||
          key.toLowerCase().includes("secret") ||
          key.toLowerCase().includes("password"),
      }));
      setEnvVars(vars);
    } catch {
      showStatus("Failed to load env variables", true);
    } finally {
      setIsLoading(false);
    }
  }

  async function loadCompare() {
    try {
      // Pass the selected env file path instead of just the default
      // Assuming envPath points to a specific .env file, we compare against .env.example in the same dir
      const target = envPath || ".env";
      let source = ".env.example";

      // Basic heuristic to find .env.example relative to the given .env path
      if (target.includes("/") || target.includes("\\")) {
        const parts = target.replace(/\\/g, "/").split("/");
        parts.pop(); // remove filename
        source = parts.join("/") + "/.env.example";
      }

      const result = await envApi.compare(source, target);
      setCompareOutput(result.output || "No differences found");
    } catch {
      showStatus("Failed to compare env files", true);
    }
  }

  async function loadValidate() {
    try {
      const result = await envApi.validate(envPath || ".env");
      setValidationResult(result);
    } catch {
      showStatus("Failed to validate env file", true);
    }
  }

  function showStatus(text: string, isError = false) {
    setStatusMsg({ text, isError });
    setTimeout(() => setStatusMsg(null), 3000);
  }

  function copyValue(key: string, value: string) {
    navigator.clipboard.writeText(`${key}=${value}`);
    setCopied(key);
    setTimeout(() => setCopied(null), 2000);
  }

  function maskValue(key: string, value: string): string {
    if (!showSecrets) {
      const sensitive = [
        "PASSWORD",
        "SECRET",
        "KEY",
        "TOKEN",
        "API_KEY",
        "PRIVATE",
        "CREDENTIAL",
      ];
      const upper = key.toUpperCase();
      if (sensitive.some((p) => upper.includes(p))) return "****";
    }
    return value;
  }

  const filtered = envVars.filter(
    (v) =>
      v.key.toLowerCase().includes(filter.toLowerCase()) ||
      v.value.toLowerCase().includes(filter.toLowerCase()),
  );
  const sensitiveCount = envVars.filter((v) => v.isSensitive).length;

  const subTabs = [
    { id: "show" as const, label: t("env.tabs.show"), icon: <Eye size={14} /> },
    {
      id: "compare" as const,
      label: t("env.tabs.compare"),
      icon: <GitCompare size={14} />,
    },
    {
      id: "validate" as const,
      label: t("env.tabs.validate"),
      icon: <Shield size={14} />,
    },
    {
      id: "sync" as const,
      label: t("env.tabs.sync"),
      icon: <RefreshCw size={14} />,
    },
    {
      id: "init" as const,
      label: t("env.tabs.init"),
      icon: <FileDown size={14} />,
    },
  ];

  return (
    <div className="space-y-4 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-[var(--color-foreground)] flex items-center gap-3">
            <FileCode size={24} className="text-[var(--warning)]" />
            {t("env.title")}
          </h2>
          <p className="text-sm text-[var(--color-muted)] mt-1">
            {t("env.subtitle")}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <div className="relative">
            <input
              type="text"
              placeholder={t("env.pathPlaceholder")}
              title={t("env.pathHelp")}
              value={envPath}
              onChange={(e) => setEnvPath(e.target.value)}
              className="w-48 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
            />
          </div>
          <button
            onClick={() => {
              if (activeSubTab === "show") loadEnvVars();
              if (activeSubTab === "compare") loadCompare();
              if (activeSubTab === "validate") loadValidate();
            }}
            disabled={isLoading}
            className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] hover:border-[var(--color-primary)] transition-all disabled:opacity-50"
          >
            <RefreshCw size={14} className={isLoading ? "animate-spin" : ""} />
            {t("common.refresh")}
          </button>
        </div>
      </div>

      <div className="bg-[var(--info)]/10 border border-[var(--info)]/20 rounded-xl p-4 flex gap-3 text-sm text-[var(--info)]">
        <FileCode size={18} className="flex-shrink-0 mt-0.5" />
        <p>{t("env.description")}</p>
      </div>

      {/* Stats */}
      <div className="flex gap-4">
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <FileCode size={12} className="text-[var(--warning)]" />
          <span className="text-xs text-[var(--color-muted)]">
            {t("env.stats.total")}:
          </span>
          <span className="text-xs font-bold text-[var(--color-foreground)]">
            {envVars.length}
          </span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <Shield size={12} className="text-[var(--danger)]" />
          <span className="text-xs text-[var(--color-muted)]">
            {t("env.stats.sensitive")}:
          </span>
          <span className="text-xs font-bold text-[var(--danger)]">
            {sensitiveCount}
          </span>
        </div>
      </div>

      {/* Status */}
      {statusMsg && (
        <div
          className={`px-4 py-2 rounded-lg text-sm ${statusMsg.isError ? "bg-[var(--danger)]/10 text-[var(--danger)] border border-[var(--danger)]/20" : "bg-[var(--success)]/10 text-[var(--success)] border border-[var(--success)]/20"}`}
        >
          {statusMsg.text}
        </div>
      )}

      {/* Sub Tabs */}
      <div className="flex gap-1 bg-[var(--color-card)] p-1 rounded-lg border border-[var(--color-border)]">
        {subTabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveSubTab(tab.id)}
            className={`flex items-center gap-2 px-4 py-2 rounded-md text-xs font-medium transition-all ${
              activeSubTab === tab.id
                ? "bg-[var(--color-primary)] text-white"
                : "text-[var(--color-muted)] hover:text-[var(--color-foreground)] hover:bg-[var(--color-card-hover)]"
            }`}
          >
            {tab.icon}
            {tab.label}
          </button>
        ))}
      </div>

      {/* Show Tab */}
      {activeSubTab === "show" && (
        <div className="space-y-3">
          <div className="flex gap-3">
            <div className="relative flex-1">
              <Search
                size={16}
                className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-muted)]"
              />
              <input
                type="text"
                placeholder={t("env.filterPlaceholder")}
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                className="w-full pl-10 pr-4 py-2.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
              />
            </div>
            <button
              onClick={() => setShowSecrets(!showSecrets)}
              className={`flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm transition-all border ${
                showSecrets
                  ? "bg-[var(--danger)]/10 border-[var(--danger)]/30 text-[var(--danger)]"
                  : "bg-[var(--color-card)] border-[var(--color-border)] text-[var(--color-muted)] hover:text-[var(--color-foreground)]"
              }`}
            >
              {showSecrets ? <EyeOff size={14} /> : <Eye size={14} />}
              {showSecrets ? t("env.hideSecrets") : t("env.showSecrets")}
            </button>
          </div>

          <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-[var(--color-border)]">
                    <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider w-8">
                      #
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider w-1/3">
                      {t("env.table.key")}
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                      {t("env.table.value")}
                    </th>
                    <th className="px-4 py-3 w-16"></th>
                  </tr>
                </thead>
                <tbody>
                  {filtered.map((v, i) => (
                    <tr
                      key={v.key}
                      className="border-b border-[var(--color-border)] hover:bg-[var(--color-card-hover)] transition-colors group"
                    >
                      <td className="px-4 py-3 text-xs text-[var(--color-muted)] font-mono">
                        {i + 1}
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <span className="font-mono text-sm font-medium text-[var(--color-foreground)]">
                            {v.key}
                          </span>
                          {v.isSensitive && (
                            <div title="Sensitive variable">
                              <Shield
                                size={12}
                                className="text-[var(--danger)]"
                              />
                            </div>
                          )}
                        </div>
                      </td>
                      <td className="px-4 py-3 font-mono text-xs text-[var(--color-muted)]">
                        <div className="flex items-center gap-2 max-w-[400px] xl:max-w-[600px]">
                          <span className="truncate flex-1">
                            {v.isSensitive && !showSecrets
                              ? "••••••••"
                              : v.value}
                          </span>
                          {v.comment && (
                            <span className="text-[10px] text-[var(--color-muted)] italic truncate">
                              {"// "}
                              {v.comment}
                            </span>
                          )}
                        </div>
                      </td>
                      <td className="px-4 py-3">
                        <button
                          onClick={() => {
                            navigator.clipboard.writeText(v.value);
                            setCopied(v.key);
                            setTimeout(() => setCopied(null), 2000);
                          }}
                          className="p-1.5 rounded-lg text-[var(--color-muted)] hover:text-[var(--color-foreground)] hover:bg-[var(--color-background)] transition-all opacity-0 group-hover:opacity-100"
                          title={t("common.copy")}
                        >
                          {copied === v.key ? (
                            <Check
                              size={14}
                              className="text-[var(--success)]"
                            />
                          ) : (
                            <Copy size={14} />
                          )}
                        </button>
                      </td>
                    </tr>
                  ))}
                  {filtered.length === 0 && (
                    <tr>
                      <td
                        colSpan={4}
                        className="px-4 py-8 text-center text-[var(--color-muted)]"
                      >
                        {t("env.noVariables")}
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* Compare Tab */}
      {activeSubTab === "compare" && (
        <div className="space-y-4">
          <div className="flex items-center gap-3 text-sm">
            <div className="px-3 py-2 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)] text-[var(--color-foreground)] font-mono text-xs">
              .env.example
            </div>
            <ArrowRight size={16} className="text-[var(--color-muted)]" />
            <div className="px-3 py-2 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)] text-[var(--color-foreground)] font-mono text-xs">
              .env
            </div>
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg p-3">
              <p className="text-xs text-[var(--color-muted)] mb-1">
                {t("env.compare.status")}
              </p>
              <p
                className={`text-lg font-bold ${
                  compareOutput.includes("MISSING") ||
                  compareOutput.includes("EXTRA")
                    ? "text-[var(--danger)]"
                    : "text-[var(--success)]"
                }`}
              >
                {compareOutput.includes("MISSING") ||
                compareOutput.includes("EXTRA")
                  ? t("env.compare.different")
                  : t("env.compare.inSync")}
              </p>
            </div>
          </div>
          <div className="bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg p-4 font-mono text-xs whitespace-pre-wrap text-[var(--color-foreground)]">
            {compareOutput || t("env.compare.noDiff")}
          </div>
        </div>
      )}

      {/* Sync Tab */}
      {activeSubTab === "sync" && (
        <div className="space-y-4">
          <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-5">
            <h4 className="text-sm font-semibold text-white mb-3 flex items-center gap-2">
              <RefreshCw size={14} className="text-blue-400" />
              Interactive Sync
            </h4>
            <p className="text-xs text-[var(--color-muted)] mb-4">
              Sync missing variables from .env.example to .env. You can set
              values for each missing variable.
            </p>
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg p-3">
              <p className="text-xs text-[var(--color-muted)] mb-1">Status</p>
              <p className="text-lg font-bold text-white">
                {compareOutput.includes("MISSING") ||
                compareOutput.includes("EXTRA")
                  ? "Different"
                  : "In Sync"}
              </p>
            </div>
          </div>
          <div className="bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg p-4 font-mono text-xs whitespace-pre-wrap">
            {compareOutput ||
              "No differences found between .env and .env.example"}
          </div>
        </div>
      )}

      {/* Validate Tab */}
      {activeSubTab === "validate" && (
        <div className="space-y-4">
          {validationIssues.length === 0 ? (
            <div className="bg-[var(--success)]/10 border border-[var(--success)]/20 rounded-xl p-6 text-center">
              <CheckCircle
                size={32}
                className="mx-auto mb-2 text-[var(--success)]"
              />
              <p className="text-[var(--success)] font-medium">
                {t("env.validate.noIssues")}
              </p>
            </div>
          ) : (
            <>
              <div className="flex items-center gap-2 px-4 py-2 bg-[var(--warning)]/10 border border-[var(--warning)]/20 rounded-lg">
                <AlertTriangle size={14} className="text-[var(--warning)]" />
                <span className="text-sm text-[var(--warning)]">
                  {t("env.validate.foundIssues").replace(
                    "{count}",
                    validationIssues.length.toString(),
                  )}
                </span>
              </div>

              <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-[var(--color-border)]">
                      <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase">
                        {t("env.validate.table.line")}
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase">
                        {t("env.validate.table.key")}
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase">
                        {t("env.validate.table.type")}
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase">
                        {t("env.validate.table.message")}
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {validationIssues.map((issue, i) => (
                      <tr
                        key={i}
                        className="border-b border-[var(--color-border)] hover:bg-[var(--color-card-hover)]"
                      >
                        <td className="px-4 py-2.5 font-mono text-xs text-[var(--color-muted)]">
                          {issue.line}
                        </td>
                        <td className="px-4 py-2.5 font-mono text-xs text-[var(--color-foreground)] font-medium">
                          {issue.key}
                        </td>
                        <td className="px-4 py-2.5">
                          <span
                            className={`text-[10px] px-2 py-0.5 rounded-full ${
                              issue.type === "empty"
                                ? "bg-[var(--danger)]/10 text-[var(--danger)]"
                                : issue.type === "duplicate"
                                  ? "bg-[var(--primary)]/10 text-[var(--primary)]"
                                  : issue.type === "trailing_space"
                                    ? "bg-[var(--warning)]/10 text-[var(--warning)]"
                                    : issue.type === "unquoted_spaces"
                                      ? "bg-[var(--info)]/10 text-[var(--info)]"
                                      : "bg-[var(--color-muted)]/10 text-[var(--color-muted)]"
                            }`}
                          >
                            {issue.type.replace("_", " ")}
                          </span>
                        </td>
                        <td className="px-4 py-2.5 text-xs text-[var(--color-muted)]">
                          {issue.message}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </>
          )}
        </div>
      )}

      {/* Init Tab */}
      {activeSubTab === "init" && (
        <div className="space-y-4">
          <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-5">
            <h4 className="text-sm font-semibold text-[var(--color-foreground)] mb-3 flex items-center gap-2">
              <FileDown size={14} className="text-[var(--success)]" />
              {t("env.init.title")}
            </h4>
            <p className="text-xs text-[var(--color-muted)] mb-4">
              {t("env.init.desc")}
            </p>

            <div className="space-y-3">
              {envVars.slice(0, 8).map((v) => (
                <div
                  key={v.key}
                  className="flex items-center gap-3 py-2 px-3 bg-[var(--color-background)] rounded-lg border border-[var(--color-border)]"
                >
                  <span className="font-mono text-xs text-[var(--color-foreground)] font-medium w-40 flex-shrink-0">
                    {v.key}
                  </span>
                  <input
                    type="text"
                    defaultValue={v.isSensitive ? "" : v.value}
                    placeholder={
                      v.isSensitive ? t("env.init.enterValue") : v.value
                    }
                    className="flex-1 px-3 py-1.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded text-xs text-[var(--color-foreground)] placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
                  />
                </div>
              ))}
            </div>

            <div className="flex gap-3 mt-4">
              <button
                onClick={() =>
                  showStatus(
                    `Generated .env with ${envVars.length} variable(s)`,
                  )
                }
                className="px-4 py-2 bg-[var(--color-primary)] rounded-lg text-sm text-[var(--color-foreground)] hover:bg-[var(--color-primary-light)] transition-colors"
              >
                {t("env.init.button")}
              </button>
              <label className="flex items-center gap-2 text-xs text-[var(--color-muted)]">
                <input
                  type="checkbox"
                  className="rounded bg-[var(--color-background)] border-[var(--color-border)]"
                />
                {t("env.init.forceOverwrite")}
              </label>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
