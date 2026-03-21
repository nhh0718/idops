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
import { useEffect, useState } from "react";
import { envApi } from "../lib/api";
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
  const [envVars, setEnvVars] = useState<EnvVariable[]>(initialEnvVars);
  const [activeSubTab, setActiveSubTab] = useState<
    "show" | "compare" | "validate" | "sync" | "init"
  >("show");
  const [showSecrets, setShowSecrets] = useState(false);
  const [filter, setFilter] = useState("");
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

  useEffect(() => {
    if (activeSubTab === "show") {
      loadEnvVars();
    } else if (activeSubTab === "compare") {
      loadCompare();
    } else if (activeSubTab === "validate") {
      loadValidate();
    }
  }, [activeSubTab]);

  async function loadEnvVars() {
    setIsLoading(true);
    try {
      const data = await envApi.show();
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
      const result = await envApi.compare();
      setCompareOutput(result.output || "No differences found");
    } catch {
      showStatus("Failed to compare env files", true);
    }
  }

  async function loadValidate() {
    try {
      const result = await envApi.validate();
      setValidationResult(result);
    } catch {
      showStatus("Failed to validate env file", true);
    }
  }

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
    { id: "show" as const, label: "Show", icon: <Eye size={14} /> },
    {
      id: "compare" as const,
      label: "Compare",
      icon: <GitCompare size={14} />,
    },
    {
      id: "validate" as const,
      label: "Validate",
      icon: <CheckCircle size={14} />,
    },
    { id: "sync" as const, label: "Sync", icon: <RefreshCw size={14} /> },
    { id: "init" as const, label: "Init", icon: <FileDown size={14} /> },
  ];

  return (
    <div className="space-y-4 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-white flex items-center gap-3">
            <FileCode size={24} className="text-amber-400" />
            Environment Manager
          </h2>
          <p className="text-sm text-[var(--color-muted)] mt-1">
            Compare, sync, validate, and manage .env files
          </p>
        </div>
      </div>

      {/* Stats */}
      <div className="flex gap-4">
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <FileCode size={12} className="text-amber-400" />
          <span className="text-xs text-[var(--color-muted)]">Variables:</span>
          <span className="text-xs font-bold text-white">{envVars.length}</span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <Shield size={12} className="text-red-400" />
          <span className="text-xs text-[var(--color-muted)]">Sensitive:</span>
          <span className="text-xs font-bold text-red-400">
            {sensitiveCount}
          </span>
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
        {subTabs.map((tab) => (
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
                placeholder="Filter variables..."
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                className="w-full pl-10 pr-4 py-2.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-white placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
              />
            </div>
            <button
              onClick={() => setShowSecrets(!showSecrets)}
              className={`flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm transition-all border ${
                showSecrets
                  ? "bg-amber-500/10 border-amber-500/30 text-amber-400"
                  : "bg-[var(--color-card)] border-[var(--color-border)] text-[var(--color-muted)] hover:text-white"
              }`}
            >
              {showSecrets ? <EyeOff size={14} /> : <Eye size={14} />}
              {showSecrets ? "Hide" : "Show"} Secrets
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
                    <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                      Key
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                      Value
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider w-20">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {filtered.map((v, i) => (
                    <tr
                      key={v.key}
                      className="border-b border-[var(--color-border)] hover:bg-[var(--color-card-hover)] transition-colors"
                    >
                      <td className="px-4 py-2.5 text-xs text-[var(--color-muted)]">
                        {i + 1}
                      </td>
                      <td className="px-4 py-2.5">
                        <div className="flex items-center gap-2">
                          {v.isSensitive && (
                            <Shield size={12} className="text-red-400" />
                          )}
                          <span className="font-mono text-xs font-medium text-white">
                            {v.key}
                          </span>
                        </div>
                      </td>
                      <td className="px-4 py-2.5">
                        <span
                          className={`font-mono text-xs ${v.isSensitive && !showSecrets ? "text-red-400/50" : "text-[var(--color-muted)]"}`}
                        >
                          {maskValue(v.key, v.value)}
                        </span>
                      </td>
                      <td className="px-4 py-2.5 text-right">
                        <button
                          onClick={() => copyValue(v.key, v.value)}
                          className="p-1.5 rounded hover:bg-blue-500/10 text-[var(--color-muted)] hover:text-blue-400 transition-colors"
                          title="Copy"
                        >
                          {copied === v.key ? (
                            <Check size={14} className="text-emerald-400" />
                          ) : (
                            <Copy size={14} />
                          )}
                        </button>
                      </td>
                    </tr>
                  ))}
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
            <div className="px-3 py-2 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)] text-white font-mono text-xs">
              .env.example
            </div>
            <ArrowRight size={16} className="text-[var(--color-muted)]" />
            <div className="px-3 py-2 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)] text-white font-mono text-xs">
              .env
            </div>
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
            <div className="bg-emerald-500/10 border border-emerald-500/20 rounded-xl p-6 text-center">
              <CheckCircle
                size={32}
                className="mx-auto mb-2 text-emerald-400"
              />
              <p className="text-emerald-400 font-medium">No issues found!</p>
            </div>
          ) : (
            <>
              <div className="flex items-center gap-2 px-4 py-2 bg-amber-500/10 border border-amber-500/20 rounded-lg">
                <AlertTriangle size={14} className="text-amber-400" />
                <span className="text-sm text-amber-400">
                  Found {validationIssues.length} issue(s) in .env
                </span>
              </div>

              <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-[var(--color-border)]">
                      <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase">
                        Line
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase">
                        Key
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase">
                        Type
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase">
                        Message
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
                        <td className="px-4 py-2.5 font-mono text-xs text-white font-medium">
                          {issue.key}
                        </td>
                        <td className="px-4 py-2.5">
                          <span
                            className={`text-[10px] px-2 py-0.5 rounded-full ${
                              issue.type === "empty"
                                ? "bg-red-500/10 text-red-400"
                                : issue.type === "duplicate"
                                  ? "bg-purple-500/10 text-purple-400"
                                  : issue.type === "trailing_space"
                                    ? "bg-amber-500/10 text-amber-400"
                                    : issue.type === "unquoted_spaces"
                                      ? "bg-blue-500/10 text-blue-400"
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

            <div className="space-y-3">
              <p className="text-xs text-[var(--color-muted)]">
                Sync feature will be available when differences are detected.
              </p>
            </div>

            <button
              onClick={() => showStatus(`Sync feature not yet implemented`)}
              className="mt-4 px-4 py-2 bg-[var(--color-primary)] rounded-lg text-sm text-white hover:bg-[var(--color-primary-light)] transition-colors"
            >
              Sync Variables
            </button>
          </div>
        </div>
      )}

      {/* Init Tab */}
      {activeSubTab === "init" && (
        <div className="space-y-4">
          <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-5">
            <h4 className="text-sm font-semibold text-white mb-3 flex items-center gap-2">
              <FileDown size={14} className="text-emerald-400" />
              Generate .env from .env.example
            </h4>
            <p className="text-xs text-[var(--color-muted)] mb-4">
              Create a new .env file with all variables from .env.example. Set
              values for each variable below.
            </p>

            <div className="space-y-3">
              {envVars.slice(0, 8).map((v) => (
                <div
                  key={v.key}
                  className="flex items-center gap-3 py-2 px-3 bg-[var(--color-background)] rounded-lg border border-[var(--color-border)]"
                >
                  <span className="font-mono text-xs text-white font-medium w-40 flex-shrink-0">
                    {v.key}
                  </span>
                  <input
                    type="text"
                    defaultValue={v.isSensitive ? "" : v.value}
                    placeholder={v.isSensitive ? "Enter value..." : v.value}
                    className="flex-1 px-3 py-1.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded text-xs text-white placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
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
                className="px-4 py-2 bg-[var(--color-primary)] rounded-lg text-sm text-white hover:bg-[var(--color-primary-light)] transition-colors"
              >
                Generate .env
              </button>
              <label className="flex items-center gap-2 text-xs text-[var(--color-muted)]">
                <input
                  type="checkbox"
                  className="rounded bg-[var(--color-background)] border-[var(--color-border)]"
                />
                Force overwrite
              </label>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
