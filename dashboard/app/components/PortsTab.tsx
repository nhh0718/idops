"use client";

import {
  ArrowUpDown,
  Crosshair,
  ExternalLink,
  Filter,
  Globe,
  RefreshCw,
  Search,
  Wifi,
} from "lucide-react";
import { useEffect, useState } from "react";
import { portsApi } from "../lib/api";
import { useI18n } from "../lib/i18n";
import type { PortEntry } from "../types";

type SortField = "port" | "pid" | "process" | "protocol";

export default function PortsTab({
  ports: initialPorts,
}: {
  ports: PortEntry[];
}) {
  const { t } = useI18n();
  const [ports, setPorts] = useState<PortEntry[]>(initialPorts);
  const [filter, setFilter] = useState("");
  const [protocolFilter, setProtocolFilter] = useState<"" | "tcp" | "udp">("");
  const [portRange, setPortRange] = useState({ min: "", max: "" });
  const [sortField, setSortField] = useState<SortField>("port");
  const [sortAsc, setSortAsc] = useState(true);
  const [confirmKill, setConfirmKill] = useState<PortEntry | null>(null);
  const [statusMsg, setStatusMsg] = useState<{
    text: string;
    isError: boolean;
  } | null>(null);
  const [watchMode, setWatchMode] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadPorts();
  }, []);

  async function loadPorts() {
    setIsLoading(true);
    try {
      const range =
        portRange.min && portRange.max
          ? `${portRange.min}-${portRange.max}`
          : "";
      const data = await portsApi.scan(
        protocolFilter || undefined,
        range || undefined,
      );
      setPorts(data);
    } catch {
      showStatus("Failed to scan ports", true);
    } finally {
      setIsLoading(false);
    }
  }

  async function handleRefresh() {
    await loadPorts();
    showStatus("Port scan refreshed");
  }

  async function handleKill(port: PortEntry) {
    setConfirmKill(null);
    try {
      const result = await portsApi.kill(port.port);
      if (result.error) {
        showStatus(result.error, true);
      } else {
        showStatus(result.message || `Killed process on port ${port.port}`);
        await loadPorts();
      }
    } catch {
      showStatus("Failed to kill process", true);
    }
  }

  function showStatus(text: string, isError = false) {
    setStatusMsg({ text, isError });
    setTimeout(() => setStatusMsg(null), 3000);
  }

  const filtered = ports
    .filter((p) => {
      if (protocolFilter && p.protocol !== protocolFilter) return false;
      if (portRange.min && p.port < parseInt(portRange.min)) return false;
      if (portRange.max && p.port > parseInt(portRange.max)) return false;
      if (filter) {
        const q = filter.toLowerCase();
        return (
          p.process.toLowerCase().includes(q) ||
          p.port.toString().includes(q) ||
          p.pid.toString().includes(q) ||
          p.address.toLowerCase().includes(q)
        );
      }
      return true;
    })
    .sort((a, b) => {
      let cmp = 0;
      switch (sortField) {
        case "port":
          cmp = a.port - b.port;
          break;
        case "pid":
          cmp = a.pid - b.pid;
          break;
        case "process":
          cmp = a.process.localeCompare(b.process);
          break;
        case "protocol":
          cmp = a.protocol.localeCompare(b.protocol);
          break;
      }
      return sortAsc ? cmp : -cmp;
    });

  function toggleSort(field: SortField) {
    if (sortField === field) setSortAsc(!sortAsc);
    else {
      setSortField(field);
      setSortAsc(true);
    }
  }

  // Watch mode: auto-refresh
  useEffect(() => {
    if (!watchMode) return;
    const interval = setInterval(loadPorts, 2000);
    return () => clearInterval(interval);
  }, [watchMode, protocolFilter, portRange]);

  function openInBrowser(port: number) {
    window.open(`http://localhost:${port}`, "_blank");
  }

  const tcpCount = ports.filter((p) => p.protocol === "tcp").length;
  const udpCount = ports.filter((p) => p.protocol === "udp").length;

  return (
    <div className="space-y-4 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-[var(--color-foreground)] flex items-center gap-3">
            <Search size={24} className="text-[var(--info)]" />
            {t("ports.title")}
          </h2>
          <p className="text-sm text-[var(--color-muted)] mt-1">
            {t("ports.subtitle")}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setWatchMode(!watchMode)}
            className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-all border ${
              watchMode
                ? "bg-[var(--info)]/10 border-[var(--info)]/30 text-[var(--info)]"
                : "bg-[var(--color-card)] border-[var(--color-border)] text-[var(--color-muted)] hover:text-[var(--color-foreground)]"
            }`}
          >
            <Wifi size={14} />
            {t("ports.watchMode")} {watchMode ? "ON" : "OFF"}
          </button>
          <button
            onClick={handleRefresh}
            disabled={isLoading}
            className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] hover:border-[var(--color-primary)] transition-all disabled:opacity-50"
          >
            <RefreshCw size={14} className={isLoading ? "animate-spin" : ""} />
            {t("common.refresh")}
          </button>
        </div>
      </div>

      {/* Stats */}
      <div className="flex gap-4">
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <Globe size={12} className="text-[var(--info)]" />
          <span className="text-xs text-[var(--color-muted)]">
            {t("ports.stats.total")}:
          </span>
          <span className="text-xs font-bold text-[var(--color-foreground)]">
            {ports.length}
          </span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <span className="text-xs px-1.5 py-0.5 rounded bg-[var(--info)]/10 text-[var(--info)] font-mono">
            TCP
          </span>
          <span className="text-xs font-bold text-[var(--info)]">
            {tcpCount}
          </span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <span className="text-xs px-1.5 py-0.5 rounded bg-[var(--warning)]/10 text-[var(--warning)] font-mono">
            UDP
          </span>
          <span className="text-xs font-bold text-[var(--warning)]">
            {udpCount}
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

      {/* Filters */}
      <div className="flex gap-3 flex-wrap">
        <div className="relative flex-1 min-w-[200px]">
          <Search
            size={16}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-muted)]"
          />
          <input
            type="text"
            placeholder={t("common.search")}
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
          />
        </div>
        <select
          value={protocolFilter}
          onChange={(e) =>
            setProtocolFilter(e.target.value as "" | "tcp" | "udp")
          }
          className="px-3 py-2.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] focus:outline-none focus:border-[var(--color-primary)]"
        >
          <option value="">{t("ports.protocol.both")}</option>
          <option value="tcp">{t("ports.protocol.tcp")}</option>
          <option value="udp">{t("ports.protocol.udp")}</option>
        </select>
        <div className="flex items-center gap-2">
          <Filter size={14} className="text-[var(--color-muted)]" />
          <input
            type="text"
            placeholder="Min"
            value={portRange.min}
            onChange={(e) =>
              setPortRange((p) => ({ ...p, min: e.target.value }))
            }
            className="w-20 px-3 py-2.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
          />
          <span className="text-[var(--color-muted)]">-</span>
          <input
            type="text"
            placeholder="Max"
            value={portRange.max}
            onChange={(e) =>
              setPortRange((p) => ({ ...p, max: e.target.value }))
            }
            className="w-20 px-3 py-2.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
          />
        </div>
      </div>

      {/* Table */}
      <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-[var(--color-border)]">
                {[
                  {
                    key: "protocol" as SortField,
                    label: t("ports.table.protocol"),
                  },
                  { key: "port" as SortField, label: "Address" },
                  { key: "port" as SortField, label: t("ports.table.port") },
                  { key: "pid" as SortField, label: t("ports.table.pid") },
                  {
                    key: "process" as SortField,
                    label: t("ports.table.process"),
                  },
                ].map((col) => (
                  <th key={col.label} className="px-4 py-3 text-left">
                    <button
                      onClick={() => toggleSort(col.key)}
                      className="flex items-center gap-1 text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider hover:text-[var(--color-foreground)] transition-colors"
                    >
                      {col.label}
                      <ArrowUpDown size={10} />
                    </button>
                  </th>
                ))}
                <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  {t("common.status")}
                </th>
                <th className="px-4 py-3 text-right text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  {t("common.actions")}
                </th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((p, i) => (
                <tr
                  key={i}
                  className="border-b border-[var(--color-border)] hover:bg-[var(--color-card-hover)] transition-colors"
                >
                  <td className="px-4 py-3">
                    <span
                      className={`text-xs px-2 py-0.5 rounded font-mono font-medium ${
                        p.protocol === "tcp"
                          ? "bg-[var(--info)]/10 text-[var(--info)]"
                          : "bg-[var(--warning)]/10 text-[var(--warning)]"
                      }`}
                    >
                      {(p.protocol || "tcp").toUpperCase()}
                    </span>
                  </td>
                  <td className="px-4 py-3 font-mono text-xs text-[var(--color-muted)]">
                    {p.address || "0.0.0.0"}
                  </td>
                  <td className="px-4 py-3 font-mono text-sm font-bold text-[var(--color-foreground)]">
                    {p.port}
                  </td>
                  <td className="px-4 py-3 font-mono text-xs text-[var(--color-muted)]">
                    {p.pid}
                  </td>
                  <td className="px-4 py-3 text-xs text-[var(--color-foreground)]">
                    {p.process}
                  </td>
                  <td className="px-4 py-3">
                    <span className="text-xs px-2 py-0.5 rounded bg-[var(--success)]/10 text-[var(--success)]">
                      {p.status}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        onClick={() => setConfirmKill(p)}
                        title={t("ports.killProcess")}
                        className="p-1.5 rounded hover:bg-[var(--danger)]/10 text-[var(--color-muted)] hover:text-[var(--danger)] transition-colors"
                      >
                        <Crosshair size={14} />
                      </button>
                      <button
                        onClick={() => openInBrowser(p.port)}
                        title="Open in Browser"
                        className="p-1.5 rounded hover:bg-[var(--info)]/10 text-[var(--color-muted)] hover:text-[var(--info)] transition-colors"
                      >
                        <ExternalLink size={14} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td
                    colSpan={7}
                    className="px-4 py-8 text-center text-[var(--color-muted)]"
                  >
                    {t("ports.noPorts")}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Confirm Kill */}
      {confirmKill && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => setConfirmKill(null)}
        >
          <div
            className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-6 max-w-sm w-full mx-4"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="text-lg font-bold text-[var(--color-foreground)] mb-2">
              {t("ports.killProcess")}
            </h3>
            <p className="text-sm text-[var(--color-muted)] mb-1">
              Kill{" "}
              <strong className="text-[var(--color-foreground)]">
                {confirmKill.process}
              </strong>{" "}
              (PID {confirmKill.pid})?
            </p>
            <p className="text-xs text-[var(--color-muted)] mb-4">
              This will terminate the process listening on port{" "}
              {confirmKill.port}. Requires admin privileges.
            </p>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setConfirmKill(null)}
                className="px-4 py-2 rounded-lg bg-[var(--color-background)] border border-[var(--color-border)] text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] transition-colors"
              >
                {t("common.cancel")}
              </button>
              <button
                onClick={() => handleKill(confirmKill)}
                className="px-4 py-2 rounded-lg bg-[var(--danger)]/20 border border-[var(--danger)]/30 text-sm text-[var(--danger)] hover:bg-[var(--danger)]/30 transition-colors"
              >
                {t("ports.killProcess")}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
